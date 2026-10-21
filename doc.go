// Package useragent parses User-Agent strings into a typed, nested result.
//
// The package is dependency-free and safe for concurrent use. Parse is the
// zero-configuration entry point:
//
//	ua := useragent.Parse(raw)
//	fmt.Println(ua.Browser.Name, ua.OS.Name, ua.Device.Type)
//
// Use New with WithBrowserRules, WithOSRules, WithBotRules, WithEngineRules,
// or WithDeviceRules to add application-specific detection. Custom rules are
// immutable after construction.
//
// The parser is intentionally conservative: unknown or hostile input leaves
// Browser.Name empty instead of selecting an arbitrary token. Client hints can
// be applied with ApplyClientHints or through the uahttp subpackage.
package useragent
