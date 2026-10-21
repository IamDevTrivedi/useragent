# Changelog

All notable changes to this project are documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/),
and this project follows [Semantic Versioning](https://semver.org/spec/v2.0.0.html)
after the public API stabilizes.

## [Unreleased]

### Planned

- Expand the conformance corpus for less common browsers and TV platforms.
- Stabilize the public API for a future `v1.0.0`.

## [0.1.0] - 2026-10-21

Initial release of the parser.

### Added

- Requires Go 1.24 or newer and has no runtime dependencies.
- Nested, typed result model:
  - `Browser`
  - `OS`
  - `Device` and `DeviceType`
  - `Engine`
  - `CPU`
  - `Bot` and `BotCategory`
  - `App`
  - `WebView`
- Zero-dependency tokenizer and parser with case-insensitive product matching.
- Browser detection for mainstream desktop, mobile, embedded, and vendor browsers.
- OS detection for Windows, Windows Phone, macOS, iOS, Android, Linux, ChromeOS,
  Tizen, webOS, HarmonyOS, KaiOS, BlackBerry, FreeBSD, OpenBSD, NetBSD, and Solaris.
- Device detection for desktop, mobile, tablet, TV, console, and wearable devices,
  including Android vendor/model extraction and normalized CPU architectures.
- Bot detection with categories for search, social preview, AI crawlers, SEO tools,
  monitoring, security scanners, scrapers, HTTP clients, headless browsers, and ads.
- Embedded application and WebView detection.
- Immutable custom rules for browsers, OSes, bots, engines, and devices, ordered by
  descending priority.
- User-Agent Client Hints parsing and application through `ClientHints`,
  `ParseBrandList`, and `ApplyClientHints`.
- `uahttp` integration for `*http.Request` and `Sec-CH-UA-*` headers.
- `cmd/ua` CLI with table, JSON, pretty JSON, summary, file, stdin, and combined-log
  input modes.
- Typed table tests, behavioral regressions, a JSON conformance corpus, fuzz tests,
  benchmarks, and runnable examples.
- CI, CodeQL, nightly fuzz/vulnerability/benchmark, and release workflow definitions
  prepared locally for the repository owner.

### Changed

- Unknown input keeps `Browser.Name` empty instead of using an arbitrary token.
- Device information is modeled as one `DeviceType` plus vendor and model instead of
  independent `Mobile`, `Tablet`, and `Desktop` booleans.
- Bot information is a structured value instead of a single boolean.
- Parsing is safe for concurrent use with a shared `Parser`.

### Fixed

- QQ Browser is no longer reported as `Core`.
- Facebook and iMessage preview requests are classified as social bots without losing
  the underlying Safari browser.
- Splash and Scrapy requests are detected as scrapers with version and URL data.
- QXSEC Android app requests now report Android, version, vendor, model, app, and
  WebView state.
- UC Browser and QQ app requests expose their Android device model and app metadata.
- Bare AppleWebKit requests are reported as `AppleWebKit`, not a parenthetical
  operating-system string.
- HeyTapBrowser is recognized and normalized to Android.
- GCP proxy prefixes no longer become the browser name; `darwin/arm64` is parsed as
  CPU metadata and Chrome wins.
- Smart TV requests are classified as TVs with Tizen/webOS/HbbTV-family OS data.
- Flutter iOS `CFNetwork` requests are parsed as iOS WebViews.
- GoogleOther is recognized as a search bot.
- Coc Coc and a broader set of browsers are detected; custom rules cover the long tail.
- CPU architecture is exposed for common x86, x64, ARM, and ARM64 User-Agent forms.
- ChromeOS platform version extraction now uses the CrOS platform token.
- Regression coverage includes lowercase input, hostile input, Brave, desktop Safari,
  Pixel models, Miui Browser, HeadlessChrome, Android WebView, Google Ads bots,
  Seafile, NetFront/FreeBSD, Samsung Internet, Bytespider, and facebookcatalog.

### Performance

- Common desktop Chrome parsing is allocation-free on the benchmark path.
- Mobile parsing dropped from seven allocations to one by removing repeated lowercase
  copies in Android model cleanup.
- Current benchmark results and methodology are documented in
  [docs/benchmarks.md](docs/benchmarks.md).

### Security

- Fuzz tests exercise arbitrary parsing, tokenization, version parsing, brand lists,
  and HTTP request headers.
- The parser performs no network access, DNS lookup, file access, or command execution.
- Parsed User-Agent values are documented as untrusted, spoofable input.

[Unreleased]: https://github.com/IamDevTrivedi/useragent/compare/v0.1.0...HEAD
[0.1.0]: https://github.com/IamDevTrivedi/useragent/releases/tag/v0.1.0
