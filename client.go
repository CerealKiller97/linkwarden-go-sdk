// Package linkwarden is a Go client for the Linkwarden API.
// See https://docs.linkwarden.app/api/api-introduction for the upstream API reference.
package linkwarden

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"mime/multipart"
	"net/http"
	"net/url"
	"strings"
)

// Client is a Linkwarden API client. Construct one with New.
type Client struct {
	baseURL    string
	token      string
	httpClient *http.Client
	logger     *slog.Logger

	Links       *LinksService
	Collections *CollectionsService
	Tags        *TagsService
	Tokens      *TokensService
	RSS         *RSSService
	Search      *SearchService
	Dashboard   *DashboardService
	Highlights  *HighlightsService
	Archives    *ArchivesService
	Config      *ConfigService
	Favicon     *FaviconService
	Migrations  *MigrationsService
}

// Option configures a Client.
type Option func(*Client)

// WithHTTPClient overrides the default http.Client used for requests.
func WithHTTPClient(hc *http.Client) Option {
	return func(c *Client) {
		c.httpClient = hc
	}
}

// WithLogger overrides the default (no-op) logger used for internal diagnostics.
// A nil logger is ignored.
func WithLogger(logger *slog.Logger) Option {
	return func(c *Client) {
		if logger != nil {
			c.logger = logger
		}
	}
}

// New creates a Linkwarden API client for the given instance base URL, authenticated
// with the given bearer token (an API token or session token).
func New(baseURL, token string, opts ...Option) *Client {
	c := &Client{
		baseURL:    strings.TrimRight(baseURL, "/"),
		token:      token,
		httpClient: http.DefaultClient,
		logger:     slog.New(slog.DiscardHandler),
	}

	for _, opt := range opts {
		opt(c)
	}

	c.Links = &LinksService{c: c}
	c.Collections = &CollectionsService{c: c}
	c.Tags = &TagsService{c: c}
	c.Tokens = &TokensService{c: c}
	c.RSS = &RSSService{c: c}
	c.Search = &SearchService{c: c}
	c.Dashboard = &DashboardService{c: c}
	c.Highlights = &HighlightsService{c: c}
	c.Archives = &ArchivesService{c: c}
	c.Config = &ConfigService{c: c}
	c.Favicon = &FaviconService{c: c}
	c.Migrations = &MigrationsService{c: c}

	return c
}

type requestOptions struct {
	method      string
	path        string
	query       url.Values
	jsonBody    any
	rawBody     io.Reader
	contentType string
}

// request performs an HTTP request against the Linkwarden API, returning the raw
// response and its fully-read body. On a non-2xx status it returns a *APIError.
func (c *Client) request(ctx context.Context, opt requestOptions) (*http.Response, []byte, error) {
	reqURL := c.baseURL + opt.path
	if len(opt.query) > 0 {
		reqURL += "?" + opt.query.Encode()
	}

	var body io.Reader
	contentType := opt.contentType

	switch {
	case opt.rawBody != nil:
		body = opt.rawBody
	case opt.jsonBody != nil:
		b, err := json.Marshal(opt.jsonBody)
		if err != nil {
			return nil, nil, fmt.Errorf("marshal request body: %w", err)
		}
		body = bytes.NewReader(b)
		contentType = "application/json"
	}

	req, err := http.NewRequestWithContext(ctx, opt.method, reqURL, body)
	if err != nil {
		return nil, nil, fmt.Errorf("create request: %w", err)
	}

	if contentType != "" {
		req.Header.Set("Content-Type", contentType)
	}
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Authorization", "Bearer "+c.token)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, nil, fmt.Errorf("do request: %w", err)
	}
	defer func() {
		if cerr := resp.Body.Close(); cerr != nil {
			c.logger.Error("failed to close response body", slog.Any("error", cerr))
		}
	}()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return resp, nil, fmt.Errorf("read response body: %w", err)
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return resp, respBody, newAPIError(resp.StatusCode, respBody)
	}

	return resp, respBody, nil
}

// buildMultipartBody encodes fields and a file part (under fileField, with the given
// filename) as a multipart/form-data body, returning the body reader and the
// Content-Type header (including boundary) to send alongside it.
func buildMultipartBody(fields map[string]string, fileField, filename string, r io.Reader) (io.Reader, string, error) {
	buf := &bytes.Buffer{}
	w := multipart.NewWriter(buf)

	for k, v := range fields {
		if err := w.WriteField(k, v); err != nil {
			return nil, "", fmt.Errorf("write field %q: %w", k, err)
		}
	}

	part, err := w.CreateFormFile(fileField, filename)
	if err != nil {
		return nil, "", fmt.Errorf("create form file: %w", err)
	}

	if _, err := io.Copy(part, r); err != nil {
		return nil, "", fmt.Errorf("copy file contents: %w", err)
	}

	if err := w.Close(); err != nil {
		return nil, "", fmt.Errorf("close multipart writer: %w", err)
	}

	return buf, w.FormDataContentType(), nil
}

// decodeResponseEnvelope unwraps the common {"response": T} envelope shape.
func decodeResponseEnvelope[T any](body []byte) (T, error) {
	var envelope struct {
		Response T `json:"response"`
	}

	var zero T
	if err := json.Unmarshal(body, &envelope); err != nil {
		return zero, fmt.Errorf("decode response envelope: %w", err)
	}

	return envelope.Response, nil
}

// decodeDataEnvelope unwraps the {"data": T, "message": ..., "success": ...} envelope shape.
func decodeDataEnvelope[T any](body []byte) (T, error) {
	var envelope struct {
		Data    T      `json:"data"`
		Message string `json:"message"`
		Success bool   `json:"success"`
	}

	var zero T
	if err := json.Unmarshal(body, &envelope); err != nil {
		return zero, fmt.Errorf("decode data envelope: %w", err)
	}

	return envelope.Data, nil
}

// decodeJSON unmarshals a response body with no enclosing envelope.
func decodeJSON(body []byte, v any) error {
	if err := json.Unmarshal(body, v); err != nil {
		return fmt.Errorf("decode response body: %w", err)
	}

	return nil
}
