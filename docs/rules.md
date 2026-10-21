# Custom rules

Rules extend or override detection without modifying the package. This guide covers
the public rule API and the conventions used when adding built-in rules.

## Basic behavior

- Custom rules run before built-in detectors.
- Rule families are evaluated separately: browser, OS, bot, engine, device.
- Within a family, rules are stable-sorted by descending `Priority`.
- The first matching rule wins.
- A zero `Match` never matches.
- Matching is case-insensitive for ASCII.
- Rule slices are copied during `New`; a `Parser` is immutable afterward.

`New` returns a parser. The package-level `Parse` function uses `DefaultParser`, which
has no custom rules.

```go
parser := useragent.New(
	useragent.WithBrowserRules(browserRules...),
	useragent.WithOSRules(osRules...),
	useragent.WithBotRules(botRules...),
	useragent.WithEngineRules(engineRules...),
	useragent.WithDeviceRules(deviceRules...),
)

ua := parser.Parse(raw)
```

## Match

`Match` selects tokens. Unspecified fields are ignored, but at least one field must be
set. `None` alone is valid and matches when none of the listed tokens exists.

```go
type Match struct {
	Token    string   // exact token key
	Prefix   string   // token key starts with this
	Suffix   string   // token key ends with this
	Contains string   // token key or value contains this
	All      []string // every token must exist
	Any      []string // at least one token must exist
	None     []string // none of these tokens may exist
}
```

Examples:

```go
// Exact product token.
useragent.Match{Token: "AcmeClient"}

// Vendor-prefixed product.
useragent.Match{Prefix: "Acme"}

// A token that must be present, with a guard token that must not be.
useragent.Match{
	All:  []string{"AcmeShell", "Mobile"},
	None: []string{"AcmeBot"},
}
```

`Contains` checks the token keys and values collected by the tokenizer, not every
arbitrary substring of the raw User-Agent. If a client embeds a marker only in a
parenthesized field, verify that the tokenizer retains it before relying on
`Contains`.

## Browser rules

```go
browserRule := useragent.BrowserRule{
	Name:        "Acme Desktop",
	Match:       useragent.Match{Token: "AcmeDesktop"},
	VersionKeys: []string{"AcmeDesktop", "AcmeVersion"},
	Product:     "AcmeDesktop",
	Priority:    100,
}
```

- `Name` is the result written to `Browser.Name`. An empty name disables the rule.
- `VersionKeys` are token keys tried in order.
- `Product` is an additional raw `Product/version` lookup used only if `VersionKeys`
  produce nothing.
- Use a high priority for an exact client, a lower priority for a family, and negative
  priorities only when intentionally deferring to other custom rules.

A vendor browser that must win over Chromium should use a higher priority than a
generic Chromium rule. Built-in Chromium detection still runs after all custom rules if
nothing matches.

## OS rules

```go
useragent.OSRule{
	Name:        "AcmeOS",
	Match:       useragent.Match{Prefix: "AcmeOS"},
	VersionKeys: []string{"AcmeOS"},
	Priority:    50,
}
```

OS rules use the same version lookup semantics as browser rules. A rule can set an OS
name with an empty version when the version is unavailable.

## Bot rules

```go
useragent.BotRule{
	Name:        "AcmeBot",
	Match:       useragent.Match{Prefix: "AcmeBot"},
	VersionKeys: []string{"AcmeBot"},
	Category:    useragent.BotScraper,
	URL:         "https://example.com/acmebot",
	Priority:    200,
}
```

- `Category` should use one of the exported `BotCategory` values.
- `URL` overrides the URL extracted from the User-Agent. Leave it empty to use the
  extracted `+https://...` URL when present.
- A primary bot rule can replace the browser name during finalization. Social preview
  bots are intentionally classified as secondary, so they do not erase a detected
  Safari browser.

## Engine rules

```go
useragent.EngineRule{
	Name:        "AcmeEngine",
	Match:       useragent.Match{Token: "AcmeEngine"},
	VersionKeys: []string{"AcmeEngine"},
	Priority:    20,
}
```

Engine rules are useful for embedded runtimes whose engine name is not the same as the
browser name.

## Device rules

```go
useragent.DeviceRule{
	Type:     useragent.DeviceTablet,
	Vendor:   "Acme",
	Model:    "Acme Slate",
	Match:    useragent.Match{Token: "AcmeSlate"},
	Priority: 100,
}
```

A device rule sets `Device.Type` and optionally overrides vendor and model. It runs
before built-in device detection, which then leaves the custom type in place. OS
detection and finalization still run and may add metadata.

## Built-in rule development

When a client is broadly useful, add it to the built-in tables instead of asking every
user to write a rule.

### Browser

Edit `detect_browser.go`. Put the most specific product before generic Chromium,
Firefox, Safari, and WebKit cases. Vendor browsers must be checked before Chromium.
Keep the unknown fallback empty; do not use a first-token or best-match fallback.

### Bot

Add an entry to `botDefs` in `detect_bot.go`:

```go
{name: "AcmeBot", keys: []string{"AcmeBot"}, category: BotScraper},
```

Use `versionKeys` when the version is in a different token. Use `contains` only for
markers that cannot be represented as token keys. The first-byte index is built at
initialization, so no manual index maintenance is required.

### OS

Add a token-first detector in `detect_os.go`, then place the case in the detection
switch according to specificity. For example, a platform-specific OS must be checked
before a generic Linux check.

### Device

Add markers or model handling in `detect_device.go`. Keep desktop short-circuits in
mind: a new TV/console/wearable platform that also reports Windows or macOS needs an
explicit exception and tests.

### Engine

Add the product/version mapping in `detect_engine.go`, and ensure browser detection
does not mask the engine's identifying token.

## Tests for a new rule

Every behavior change needs:

1. a typed table or named regression test with the exact User-Agent
2. an assertion for the fields that changed
3. a negative assertion where the old bug produced a wrong value
4. a conformance corpus entry for broad coverage
5. a benchmark check if the change touches the hot path

Example:

```go
func TestAcmeBrowser(t *testing.T) {
	ua := Parse("AcmeClient/2.5 (Windows NT 10.0; Win64; x64)")
	if ua.Browser.Name != "Acme Desktop" {
		t.Fatalf("Browser.Name = %q", ua.Browser.Name)
	}
	if ua.Browser.Version != "2.5" {
		t.Fatalf("Browser.Version = %q", ua.Browser.Version)
	}
}
```

For a behavioral regression, name the test after the client or behavior under test and
preserve the exact reported fixture. See `parse_test.go` for examples.

## Common mistakes

- **Catch-all custom rule:** a zero `Match` is ignored, so `Match{}` will not act as a
  fallback.
- **Wrong precedence:** a generic Chromium rule at a higher priority can mask a vendor
  browser. Put the most specific match first.
- **Using raw substring matching as a first choice:** token matching is cheaper and less
  likely to match unrelated comments.
- **Treating client hints as verified:** hints are spoofable. They improve parsing but
  are not identity.
- **Weakening regression tests:** fix the detector or document the intended behavior
  change; do not remove assertions to make a fixture pass.
