# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

Each release also names the Linkwarden version its tests passed against: the one pinned in
[`.linkwarden-version`](.linkwarden-version) when it was tagged. Other Linkwarden versions
may work, but they weren't tested.

<!--
How to write an entry:
- Add changes under [Unreleased] as you go. On release, rename it to
  "## [x.y.z] - YYYY-MM-DD", start a new empty [Unreleased], and update the links at the bottom.
- Keep the "🧪 Tested against Linkwarden: vX.Y.Z" line in every section. The release
  workflow checks it against .linkwarden-version.
- Group changes under these headings, in this order, and leave out empty ones:
    ### ✨ Added       new features
    ### 🔄 Changed     changes to existing behaviour
    ### ⚠️ Deprecated  features that will be removed in a future release
    ### 🗑️ Removed     features removed in this release
    ### 🐛 Fixed       bug fixes
    ### 🔒 Security    vulnerability fixes
- Mark breaking changes with "💥 **Breaking:**" at the start of the entry.
- A pulled release gets " [YANKED]" after its date.
-->

## [Unreleased]

🧪 Tested against Linkwarden: v2.16.3

## [1.0.1] - 2026-11-08

🧪 Tested against Linkwarden: v2.16.3

### 🐛 Fixed

- `Links.List` calls `GET /api/v1/search` instead of `GET /api/v1/links`, which newer
  Linkwarden versions (including Linkwarden Cloud) reject with `400 This route is
  deprecated`. Filters, sorting and the cursor work as before, and a page with no links
  is now `nil`.

## [1.0.0] - 2026-11-03

🧪 Tested against Linkwarden: v2.16.3

🎉 First stable release.

### ✨ Added

- `New` client with `WithHTTPClient` and `WithLogger` options. Every call takes a
  `context.Context`. Log output is discarded unless you pass a logger.
- Services for links, collections, tags, tokens, RSS, search, dashboard, highlights,
  archives, config, favicon, migrations, preserved files, avatars, logins, public
  collections, session, auth, users and the background worker.
- `*APIError` with the status code, message and raw body, plus `IsNotFound` and
  `IsConflict` helpers.
- Doc comments on every exported method and type.

[unreleased]: https://github.com/CerealKiller97/linkwarden-go-sdk/compare/v1.0.1...HEAD
[1.0.1]: https://github.com/CerealKiller97/linkwarden-go-sdk/compare/v1.0.0...v1.0.1
[1.0.0]: https://github.com/CerealKiller97/linkwarden-go-sdk/releases/tag/v1.0.0
