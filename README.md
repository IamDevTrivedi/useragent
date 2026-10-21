# useragent

A zero-dependency, typed User-Agent parser for Go.

`useragent` turns a User-Agent string into a structured result covering the browser,
operating system, device, engine, CPU architecture, bot, embedded app, and WebView
state. It is designed for analytics, request classification, and observability rather
than authentication or device identification.

[![Go Reference](https://pkg.go.dev/badge/github.com/IamDevTrivedi/useragent.svg)](https://pkg.go.dev/github.com/IamDevTrivedi/useragent)
[![CI](https://github.com/IamDevTrivedi/useragent/actions/workflows/ci.yml/badge.svg)](https://github.com/IamDevTrivedi/useragent/actions/workflows/ci.yml)
[![License: MIT](https://img.shields.io/badge/license-MIT-blue.svg)](LICENSE)

## Why this parser

- **Typed, nested API.** Use `ua.Browser.Name`, `ua.OS.Version`, `ua.Device.Type`,
  `ua.Bot.Category`, and `ua.CPU.Architecture` instead of a flat set of booleans.
- **No arbitrary names.** Unknown or hostile input leaves `Browser.Name` empty rather
  than guessing from an unrelated token.
- **Zero runtime dependencies.** The core package uses only the Go standard library.
- **Safe for concurrent use.** `DefaultParser` and every `Parser` are immutable after
  construction.
- **Extensible.** Add browser, OS, bot, engine, and device rules without editing the
  package.
- **Client hints supported.** Merge `Sec-CH-UA-*` values through `uahttp` or
  `ApplyClientHints`.
- **Tested as a parser.** The repository includes a conformance corpus, behavioral
  regressions, fuzz tests, benchmarks, and runnable examples.

## Install

```sh
go get github.com/IamDevTrivedi/useragent@v0.1.0
```

The module requires Go 1.24 or newer.

## Quick start

```go
package main

import (
	"fmt"

	"github.com/IamDevTrivedi/useragent"
)

func main() {
	ua := useragent.Parse("Mozilla/5.0 (Windows NT 10.0; Win64; x64) " +
		"AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.6099.109 Safari/537.36")

	fmt.Println(ua.Browser.Name)              // Chrome
	fmt.Println(ua.Browser.Version)           // 120.0.6099.109
	fmt.Println(ua.OS.Name, ua.OS.Version)    // Windows 10.0
	fmt.Println(ua.Device.Type)               // desktop
	fmt.Println(ua.Engine.Name)               // Blink
	fmt.Println(ua.CPU.Architecture)          // amd64
	fmt.Println(ua.Bot.Detected)              // false
}
```

A parsed result has this shape:

```go
type UserAgent struct {
	Raw     string
	Browser Browser
	OS      OS
	Device  Device
	Engine  Engine
	CPU     CPU
	Bot     Bot
	App     App
	WebView bool
}
```

All nested types are JSON-ready. `DeviceType` and `BotCategory` marshal as readable
strings and also accept numeric values when unmarshaling.

## Bot and application detection

```go
ua := useragent.Parse("Mozilla/5.0 (compatible; Googlebot/2.1; +http://www.google.com/bot.html)")

if ua.IsBot() {
	fmt.Println(ua.Bot.Name, ua.Bot.Category, ua.Bot.URL)
	// Googlebot search http://www.google.com/bot.html
}
```

Bot categories include search crawlers, social previewers, AI crawlers, SEO tools,
monitoring services, security scanners, scrapers, HTTP clients, headless browsers, and
advertising crawlers. Social preview bots can retain the underlying browser; for
example, a Facebook preview request can report both `Bot.Name = facebookexternalhit`
and `Browser.Name = Safari`.

Embedded app containers such as Facebook, Instagram, TikTok, WeChat, QQ, Seafile, and
CFNetwork are exposed through `ua.App` and `ua.WebView`.

## Client hints

Use `uahttp` with an `*http.Request`:

```go
import (
	"net/http"

	"github.com/IamDevTrivedi/useragent/uahttp"
)

func handler(w http.ResponseWriter, r *http.Request) {
	ua := uahttp.Parse(r)
	_ = ua
}
```

Or apply hints directly:

```go
ua := useragent.Parse(frozenAndroidUA)
ua.ApplyClientHints(useragent.ClientHints{
	FullVersionList: []useragent.Brand{{Name: "Google Chrome", Version: "120.0.6099.109"}},
	Mobile:          true,
	Platform:        "Android",
	PlatformVersion: "13.0.0",
	Model:           "Pixel 7",
	Architecture:    "arm64",
})
```

Client hints are treated as higher-quality data for the browser version, platform,
CPU architecture, and device model. They do not replace an already-detected Safari or
Firefox browser with a generic Chromium brand.

## Custom rules

Custom rules run before the built-in detectors. Higher `Priority` values run first,
and the first matching rule wins.

```go
parser := useragent.New(
	useragent.WithBrowserRules(useragent.BrowserRule{
		Name:        "Acme Desktop",
		Match:       useragent.Match{Token: "AcmeClient"},
		VersionKeys: []string{"AcmeClient"},
		Priority:    100,
	}),
	useragent.WithBotRules(useragent.BotRule{
		Name:     "AcmeBot",
		Match:    useragent.Match{Prefix: "AcmeBot"},
		Category: useragent.BotScraper,
	}),
)

ua := parser.Parse("AcmeClient/2.5")
fmt.Println(ua.Browser.Name, ua.Browser.Version) // Acme Desktop 2.5
```

`Match` supports exact tokens, prefixes, suffixes, contained text, and all/any/none
token combinations. A zero `Match` never matches. See [docs/rules.md](docs/rules.md)
for all rule types and ordering details.

## CLI

Install the `ua` command:

```sh
go install github.com/IamDevTrivedi/useragent/cmd/ua@v0.1.0
```

Parse one or more User-Agent strings:

```sh
ua "Mozilla/5.0 (Windows NT 10.0; Win64; x64) Chrome/120.0.0.0 Safari/537.36"

cat user-agents.txt | ua -summary
ua -json -file user-agents.txt
ua -log combined -file access.log
```

Available output modes:

- default table
- `-json` for newline-delimited JSON
- `-json -pretty` for indented JSON
- `-summary` for aggregate counts

## Detection behavior

The parser is deliberately conservative:

- Unknown input keeps `Browser.Name` empty. It does not return a random comment token.
- Product tokens are matched case-insensitively, so lowercased User-Agent strings work.
- A bare `AppleWebKit` product is reported as `AppleWebKit` and marked as a WebView.
- Contradictory vendor UAs are normalized only where the browser is platform-specific,
  such as HeyTapBrowser being Android-only.
- Windows 11 continues to report Windows NT `10.0`; browsers have not changed the token.
- Client hints are supplemental and spoofable, like the User-Agent itself.

No User-Agent parser can identify every client. The built-in rules cover common desktop,
mobile, tablet, TV, console, wearable, bot, app, and WebView clients. Custom rules cover
the long tail.

## Performance

Measured on a 12th Gen Intel Core i5-12450HX with Go 1.27.2.

| Benchmark | Time | Bytes | Allocations |
|---|---:|---:|---:|
| Desktop Chrome | ~1.75 µs | 0 B | 0 |
| Android mobile | ~4.77 µs | 176 B | 1 |
| Googlebot | ~2.41 µs | 80 B | 1 |
| Desktop, parallel | ~0.34 µs | 0 B | 0 |

The common desktop path is allocation-free. Methodology and full results are in
[docs/benchmarks.md](docs/benchmarks.md).

## Documentation

- [Architecture](docs/architecture.md): tokenizer, detector pipeline, and precedence.
- [Custom rules](docs/rules.md): add or override detection.
- [Detection matrix](docs/detection-matrix.md): supported families and test coverage.
- [Benchmarks](docs/benchmarks.md): methodology and results.
- [GitHub setup](docs/github-setup.md): owner runbook for publishing and repository
  configuration.

## FAQ

### Can I use the User-Agent to uniquely identify a device?

No. A User-Agent is **not a unique device identifier**. It can be spoofed, omitted,
frozen, shared by many devices, and changed by browser updates. `Device.Vendor` and
`Device.Model` are useful for coarse analytics, but they are not stable identity. Use
an application-level session or device identifier for identity, and treat every parsed
field as untrusted input.

### Why does Windows 11 report Windows 10.0?

The User-Agent still contains `Windows NT 10.0` on Windows 11. The parser preserves the
reported token rather than guessing the marketing version. Use
`Sec-CH-UA-Platform-Version` when the client sends client hints.

### Why is my custom browser not detected?

Add a `BrowserRule` with a token that appears in its User-Agent. See
[docs/rules.md](docs/rules.md). If the client is a generic Chromium shell, the best
result may be `Chromium` rather than a product-specific name.

### Why is `Browser.Name` empty?

That is intentional for unknown or hostile input. The parser does not choose an
arbitrary token as a browser name. If the client is important to you, add a custom rule
and a regression test.

## Contributing

Contributions are welcome. Read [CONTRIBUTING.md](CONTRIBUTING.md), add a regression
fixture for every behavior change, and run `make check` before opening a pull request.

Security reports should follow [SECURITY.md](SECURITY.md). Community participation is
covered by [CODE_OF_CONDUCT.md](CODE_OF_CONDUCT.md).

## License

MIT. See [LICENSE](LICENSE) for the copyright and license notice.
