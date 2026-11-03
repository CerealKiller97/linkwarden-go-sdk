# Changelog

All notable changes to this SDK are documented here. The format follows
[Keep a Changelog](https://keepachangelog.com/en/1.1.0/), and the project uses
[Semantic Versioning](https://semver.org/).

Each release names the Linkwarden version its test suite passed against. That's the
version pinned in [`.linkwarden-version`](.linkwarden-version) when the release was
tagged. Other Linkwarden versions may work, but they weren't tested.

## [Unreleased]

Tested against Linkwarden: v2.16.3

## [1.0.0] - 2026-11-03

Tested against Linkwarden: v2.16.3

First stable release.

### Added

- `New` client with `WithHTTPClient` and `WithLogger` options. Every call takes a
  `context.Context`. Log output is discarded unless you pass a logger.
- Services for links, collections, tags, tokens, RSS, search, dashboard, highlights,
  archives, config, favicon, migrations, preserved files, avatars, logins, public
  collections, session, auth, users and the background worker.
- `*APIError` with the status code, message and raw body, plus `IsNotFound` and
  `IsConflict` helpers.
- Doc comments on every exported method and type.

[Unreleased]: https://github.com/CerealKiller97/linkwarden-go-sdk/compare/v1.0.0...HEAD
[1.0.0]: https://github.com/CerealKiller97/linkwarden-go-sdk/releases/tag/v1.0.0
