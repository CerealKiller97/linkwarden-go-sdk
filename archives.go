package linkwarden

import (
	"context"
	"io"
	"mime"
	"net/http"
	"net/url"
	"strconv"
)

// ArchiveFormat identifies the format of an archived file, as used by
// ArchivesService.Get and ArchivesService.Upload.
type ArchiveFormat int

const (
	ArchiveFormatPNG             ArchiveFormat = 0
	ArchiveFormatJPEG            ArchiveFormat = 1
	ArchiveFormatPDF             ArchiveFormat = 2
	ArchiveFormatReadabilityJSON ArchiveFormat = 3
	ArchiveFormatMonolithHTML    ArchiveFormat = 4
)

// ArchiveFile is a binary archive file retrieved from ArchivesService.Get. Filename
// is best-effort, populated only if the server sends a Content-Disposition header.
type ArchiveFile struct {
	Data        []byte
	ContentType string
	Filename    string
}

// ArchivesService accesses the /api/v1/archives endpoints.
type ArchivesService struct {
	c *Client
}

func (s *ArchivesService) Get(ctx context.Context, linkID int, format ArchiveFormat, preview bool) (*ArchiveFile, error) {
	q := url.Values{}
	q.Set("format", strconv.Itoa(int(format)))
	if preview {
		q.Set("preview", "true")
	}

	resp, body, err := s.c.request(ctx, requestOptions{
		method: http.MethodGet,
		path:   "/api/v1/archives/" + strconv.Itoa(linkID),
		query:  q,
	})
	if err != nil {
		return nil, err
	}

	file := &ArchiveFile{
		Data:        body,
		ContentType: resp.Header.Get("Content-Type"),
	}

	if _, params, perr := mime.ParseMediaType(resp.Header.Get("Content-Disposition")); perr == nil {
		file.Filename = params["filename"]
	}

	return file, nil
}

// Upload attaches a local file as the archive (or preview, if preview is true) for
// an existing link, returning the link with its updated archive fields.
func (s *ArchivesService) Upload(ctx context.Context, linkID int, format ArchiveFormat, preview bool, filename string, r io.Reader) (*Link, error) {
	body, contentType, err := buildMultipartBody(nil, "file", filename, r)
	if err != nil {
		return nil, err
	}

	q := url.Values{}
	q.Set("format", strconv.Itoa(int(format)))
	if preview {
		q.Set("preview", "true")
	}

	_, respBody, err := s.c.request(ctx, requestOptions{
		method:      http.MethodPost,
		path:        "/api/v1/archives/" + strconv.Itoa(linkID),
		query:       q,
		rawBody:     body,
		contentType: contentType,
	})
	if err != nil {
		return nil, err
	}

	link, err := decodeResponseEnvelope[Link](respBody)
	if err != nil {
		return nil, err
	}

	return &link, nil
}

// UploadStandalone uploads a local file as a new link's archive, creating the link
// in the process.
func (s *ArchivesService) UploadStandalone(ctx context.Context, format ArchiveFormat, preview bool, linkURL, filename string, r io.Reader) (*Link, error) {
	body, contentType, err := buildMultipartBody(map[string]string{"url": linkURL}, "file", filename, r)
	if err != nil {
		return nil, err
	}

	q := url.Values{}
	q.Set("format", strconv.Itoa(int(format)))
	if preview {
		q.Set("preview", "true")
	}

	_, respBody, err := s.c.request(ctx, requestOptions{
		method:      http.MethodPost,
		path:        "/api/v1/archives",
		query:       q,
		rawBody:     body,
		contentType: contentType,
	})
	if err != nil {
		return nil, err
	}

	link, err := decodeResponseEnvelope[Link](respBody)
	if err != nil {
		return nil, err
	}

	return &link, nil
}
