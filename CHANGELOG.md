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

🧪 Tested against Linkwarden: v2.17.0

## [1.1.0] - 2026-11-11

🧪 Tested against Linkwarden: v2.17.0

🎉 Support for Linkwarden v2.17.0.

### ✨ Added

- `Links.Reorder` moves a link within its collection's custom order
  (`PUT /api/v1/links/{id}/order`), and `Link.CollectionOrder` holds its position.
- `SortDateNewestFirst`, `SortDateOldestFirst`, `SortNameAZ`, `SortNameZA` and the new
  `SortCustom` constants for the `Sort` option of link lists and searches.
- `Session.Delete` signs out by revoking the client's session token
  (`DELETE /api/v1/session`).
- `CreateSessionRequest.Code` and `Auth.RequestVerificationEmail` for verifying a new
  account's email with the 6-digit code Linkwarden now emails.
- `APIError.Code` holds the API's machine-readable error code (`ErrorCodeEmailNotVerified`,
  `ErrorCodeInvalidVerificationCode`, `ErrorCodeTooManyAttempts`), `APIError.RetryAfter`
  the `Retry-After` of a 429, and `IsTooManyRequests` checks for one.
- `ShareUsageAnalytics` on `User` and `UserPreferences`.

### ⚠️ Deprecated

- `UpdateUserRequest.IsPrivate`: Linkwarden v2.17 ignores it.

### 🔄 Changed

- Tested against Linkwarden v2.17.0. `Auth.ForgotPassword` now succeeds for any email,
  and changing or resetting a password signs the user out elsewhere and revokes their
  access tokens.

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

[unreleased]: https://github.com/CerealKiller97/linkwarden-go-sdk/compare/v1.1.0...HEAD
[1.1.0]: https://github.com/CerealKiller97/linkwarden-go-sdk/compare/v1.0.1...v1.1.0
[1.0.1]: https://github.com/CerealKiller97/linkwarden-go-sdk/compare/v1.0.0...v1.0.1
[1.0.0]: https://github.com/CerealKiller97/linkwarden-go-sdk/releases/tag/v1.0.0
