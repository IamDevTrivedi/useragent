# Architecture

This document describes how `useragent` turns a raw User-Agent string into a
`UserAgent` value.

## Design goals

- deterministic, allocation-conscious parsing
- no third-party runtime dependencies
- no regular expressions in the hot path
- no arbitrary "best match" browser names
- immutable parsers that are safe for concurrent use
- a small public API with explicit extension points

The implementation is split into a tokenizer, a result model, detector files, a rule
layer, and a thin HTTP/client-hints integration.

## Parse pipeline

`Parser.Parse` performs these steps in order:

1. `tokenize(raw)` builds a `tokenSet`.
2. Custom browser/OS/bot/engine/device rules are applied in priority order.
3. `detectOS` fills the operating system.
4. `detectDevice` fills device type, vendor, and model.
5. `detectCPU` normalizes the reported CPU architecture.
6. `detectBot` fills bot name, version, URL, and category.
7. `detectBrowser` fills browser name and version.
8. `detectApp` detects embedded application containers.
9. `detectEngine` fills engine name and version.
10. `finalize` derives `VersionNo`, WebView state, and documented cross-field
    overrides.

The order matters. For example, device detection can use the detected OS to skip TV
and wearable marker scans on desktop platforms, and finalization can normalize a
platform-specific browser such as HeyTapBrowser to Android.

```text
raw string
   │
   ▼
tokenize ──► tokenSet (products, comments, URL, architecture tokens)
   │
   ├── custom rules (priority, first match wins)
   ▼
OS ──► device ──► CPU ──► bot ──► browser ──► app ──► engine ──► finalize
                                                               │
                                                               ▼
                                                         UserAgent
```

## Tokenizer

The tokenizer is intentionally small and bounded. It does not use `strings.Split`,
regular expressions, or `bytes.Buffer` in the hot path.

### Token storage

- At most 64 product tokens are stored in a fixed-size array.
- At most 32 parenthesized comment fields are retained.
- Token keys are matched case-insensitively with ASCII folding.

A token has:

```go
type token struct {
	key     string
	value   string
	comment bool
}
```

Product tokens come from `Product/version` syntax, bare product names, nested vendor
products, and a small set of multi-word products. Comment fields are retained for
device-model extraction and are also scanned for product-like fields such as
`Opera Mini/28.0`.

Nested product syntax is exposed one level deep, so a raw token such as
`XiaoMi/MiuiBrowser/12.11.5-gn` makes both `XiaoMi` and `MiuiBrowser` available.
Architecture words found in products or comments are promoted to tokens so CPU
detection can avoid full-string scans.

### Matching helpers

All case-insensitive comparisons are ASCII-oriented and allocation-free:

- `equalFold` compares short strings.
- `hasPrefixFold` and `hasSuffixFold` provide case-insensitive prefix/suffix checks.
- `containsFold` finds a substring without allocating.
- `indexFold` performs the allocation-free substring search used by version and URL
  extraction.

The token set also carries a small three-bit `keyMask` filter based on the first byte,
last byte, and length of a key. It rejects most missing tokens before the linear token
scan. This filter is a performance optimization, not a security boundary.

For detectors that need a substring of the raw User-Agent, `tokenSet.contains` builds a
lowercased copy lazily, once, and only when needed. The common desktop path never
requests it, which keeps that path allocation-free.

### URL extraction

A URL is extracted only when an `http://` or `https://` scheme is present. It is stored
on the token set for bot parsing, and trailing punctuation and common delimiters are
trimmed.

## Detectors

### Operating systems

`detect_os.go` checks, in order:

1. iOS
2. Windows Phone
3. HarmonyOS
4. KaiOS
5. BlackBerry
6. Android
7. ChromeOS
8. Tizen
9. webOS
10. macOS
11. Windows
12. FreeBSD
13. OpenBSD
14. NetBSD
15. Solaris
16. Linux

Each detector is token-first. Version helpers use a token value when available and
fall back to a product scan only when necessary. The order prevents Windows Phone and
Android from being swallowed by the more general Windows or Linux checks.

The OS detection intentionally does not inflate Windows NT `10.0` to Windows 11. The
token is preserved; client hints can provide a more precise platform version.

### Device

`detect_device.go` classifies one of:

- desktop
- mobile
- tablet
- TV
- console
- wearable
- embedded (available to custom rules; built-in detection leaves it mostly unset)

The order is:

1. short-circuit desktop operating systems
2. consoles
3. TVs
4. wearables
5. iOS device families
6. Android device and model extraction
7. mobile/tablet markers
8. desktop fallback for desktop OSes

Android model extraction scans parenthesized comments after the Android platform field.
It strips known noise (locale, build, resolution, memory, browser products), recognizes
vendor prefixes, and infers vendor from model prefixes when possible. The result is
normalized into `Device.Vendor` and `Device.Model` rather than overloading a single
device string.

The parser deliberately does not treat a model string as a unique device identifier.

### CPU

`detect_cpu.go` normalizes common architecture markers:

| User-Agent marker | Normalized value |
|---|---|
| `Win64`, `x64`, `x86_64`, `amd64` | `amd64` |
| `WOW64`, `i386`, `i686`, `x86` | `386` |
| `arm64`, `aarch64` | `arm64` |
| `arm`, `armv7`, `armv7l`, `armv8l` | `arm` |
| `ppc64` | `ppc64` |
| `riscv64` | `riscv64` |
| `mips`, `mips64` | `mips` |

iOS UAs do not expose a CPU architecture through the UA token, so the field remains
empty unless client hints provide one.

### Bots

`detect_bot.go` stores bot definitions in an ordered table. Each definition contains
token keys, optional version keys, a category, and optional raw contains markers.
At initialization, keys are indexed by their first byte so detection does not scan
every definition for every token.

The first matching definition in table order wins. Social preview bots are treated as
secondary markers: they set `Bot.Detected` but do not overwrite an already-detected
Safari browser. More primary bots, such as `Googlebot`, replace the browser with the bot
name during finalization.

Bot categories are intentionally coarse:

`search`, `social`, `ai-crawler`, `seo-tool`, `monitoring`, `security-scanner`,
`scraper`, `http-client`, `headless`, `ads`, and `other`.

### Browsers

`detect_browser.go` uses ordered cases rather than a generic best-match fallback. Vendor
and embedded browsers are checked before generic Chromium, Firefox, and WebKit products.
Examples include QQ Browser, UC Browser, Samsung Internet, Miui Browser, HeyTapBrowser,
Opera variants, Electron, QtWebEngine, and CEF.

If no browser rule matches, `Browser.Name` stays empty. Parenthesized OS fields and
unrelated comment tokens are never promoted to browser names. A bare AppleWebKit UA is
an explicit exception: it returns `AppleWebKit` and is marked as a WebView.

`finalize` applies a small number of documented platform-specific corrections, such as
forcing HeyTapBrowser to Android/mobile. These overrides are tested and listed in the
detection matrix.

### Apps and WebView

`detect_app.go` recognizes common embedded containers and app-only User-Agent markers.
An app container implies `WebView = true`. Android `; wv)`, Flutter `CFNetwork`,
QXSEC, Facebook, Instagram, TikTok, WeChat, QQ, Seafile, and similar containers are
covered.

### Engines

`detect_engine.go` detects Blink, WebKit, Gecko, Trident, EdgeHTML, and Presto from
browser and product tokens. Engine version data is best-effort; browsers often freeze or
omit it.

## Custom rules

Custom rules run before built-ins. Each rule family is stable-sorted by descending
`Priority`, and the first match wins. A zero `Match` matches nothing, preventing
catch-all rules from being added accidentally.

The rule layer is described in [rules.md](rules.md).

## Client hints

`ClientHints` and `ApplyClientHints` live in `clienthints.go`. `uahttp` converts an
`*http.Request` into the same value.

Hints are supplemental:

- full-version lists can replace a reduced/frozen UA version
- platform, platform version, architecture, model, and mobile flags can fill gaps
- an existing Safari or Firefox browser is not replaced by a generic Chromium brand
- a generic Chromium result can be upgraded to Edge, Opera, Chrome, Brave, Vivaldi,
  Samsung Internet, or another known Chromium brand

Hints are client-controlled and must not be treated as verified facts.

## Concurrency and immutability

`New` copies and sorts custom rule slices. The resulting `Parser` is immutable.
`DefaultParser` is a package-level parser with no custom rules. Parsing allocates local
state only; no global mutable state is touched.

## Performance choices

- fixed-capacity stack token arrays
- small key filter for negative token lookups
- allocation-free ASCII folding and substring search
- lazy lowercase raw copy only for detectors that need it
- architecture tokens promoted during tokenization
- bot keys indexed by first byte
- desktop device short-circuit to avoid TV/wearable marker scans
- no regular expressions, no `strings.Split`, no `bytes.Buffer`

The measured performance and remaining gaps are documented in
[benchmarks.md](benchmarks.md).

## Bounds and failure behavior

- The parser never panics on arbitrary input.
- Long comment fields are discarded by Android model cleanup.
- Token and comment arrays are bounded.
- Unknown input produces empty fields instead of invented names.
- The parser performs no I/O and no network access.
- Returned strings are slices of, or normalized copies derived from, the input.

Fuzz tests cover arbitrary parsing, tokenization, version parsing, brand-list parsing,
and HTTP request/header integration. The conformance corpus covers representative
clients and regression fixtures.

## Where to make changes

- `tokens.go`: token storage, ASCII matching helpers, product version extraction.
- `tokenizer.go`: product/comment scanning and URL extraction.
- `detect_os.go`: OS detection and OS versions.
- `detect_device.go`: device classification and Android model extraction.
- `detect_cpu.go`: architecture normalization.
- `detect_bot.go`: bot table and bot matching.
- `detect_browser.go`: browser precedence.
- `detect_app.go`: app containers and WebView state.
- `detect_engine.go`: engine detection.
- `rules.go`: custom rule types and option constructors.
- `parser.go`: pipeline and finalization.
- `clienthints.go`: client hint parsing and application.
- `uahttp/`: HTTP integration.
- `cmd/ua/`: CLI.

When changing a detector, add a typed regression test and, for broad behavior, a
conformance case. Keep the detection matrix, fuzz targets, and changelog current.
