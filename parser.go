package useragent

import "sort"

// Parser holds an immutable set of custom rules. The zero value is a valid
// parser with only the built-in detectors. Use New to add custom rules.
type Parser struct {
	browserRules []BrowserRule
	osRules      []OSRule
	botRules     []BotRule
	engineRules  []EngineRule
	deviceRules  []DeviceRule
}

// New creates a Parser and applies opts. Custom rules are copied, sorted by
// descending priority, and never mutated afterwards.
func New(opts ...Option) *Parser {
	p := &Parser{}
	for _, opt := range opts {
		if opt != nil {
			opt(p)
		}
	}
	sort.SliceStable(p.browserRules, func(i, j int) bool {
		return p.browserRules[i].Priority > p.browserRules[j].Priority
	})
	sort.SliceStable(p.osRules, func(i, j int) bool {
		return p.osRules[i].Priority > p.osRules[j].Priority
	})
	sort.SliceStable(p.botRules, func(i, j int) bool {
		return p.botRules[i].Priority > p.botRules[j].Priority
	})
	sort.SliceStable(p.engineRules, func(i, j int) bool {
		return p.engineRules[i].Priority > p.engineRules[j].Priority
	})
	sort.SliceStable(p.deviceRules, func(i, j int) bool {
		return p.deviceRules[i].Priority > p.deviceRules[j].Priority
	})
	return p
}

// DefaultParser is used by the package-level Parse function. It contains no
// custom rules and is safe for concurrent use.
var DefaultParser = New()

// Parse parses raw using DefaultParser.
func Parse(raw string) UserAgent { return DefaultParser.Parse(raw) }

// Parse parses one user-agent string. It is safe to call concurrently.
func (p *Parser) Parse(raw string) UserAgent {
	ts := tokenize(raw)
	ua := UserAgent{Raw: raw}

	p.applyCustomRules(&ua, &ts)

	detectOS(&ua, &ts)
	detectDevice(&ua, &ts)
	detectCPU(&ua, &ts)
	detectBot(&ua, &ts)
	detectBrowser(&ua, &ts)
	detectApp(&ua, &ts)
	detectEngine(&ua, &ts)
	finalize(&ua, &ts)
	return ua
}

func (p *Parser) applyCustomRules(ua *UserAgent, ts *tokenSet) {
	if p == nil {
		return
	}
	for _, r := range p.browserRules {
		if r.Name == "" || !r.Match.matches(ts) {
			continue
		}
		version := tokenValue(ts, r.VersionKeys...)
		if version == "" && r.Product != "" {
			version = productVersion(ua.Raw, r.Product)
		}
		setBrowser(ua, r.Name, version)
		break
	}
	for _, r := range p.osRules {
		if r.Name == "" || !r.Match.matches(ts) {
			continue
		}
		version := tokenValue(ts, r.VersionKeys...)
		if version == "" && r.Product != "" {
			version = productVersion(ua.Raw, r.Product)
		}
		setOS(ua, r.Name, version)
		break
	}
	for _, r := range p.botRules {
		if r.Name == "" || !r.Match.matches(ts) {
			continue
		}
		version := tokenValue(ts, r.VersionKeys...)
		if version == "" && r.Product != "" {
			version = productVersion(ua.Raw, r.Product)
		}
		url := r.URL
		if url == "" {
			url = ts.url
		}
		setBot(ua, r.Name, version, url, r.Category)
		break
	}
	for _, r := range p.engineRules {
		if r.Name == "" || !r.Match.matches(ts) {
			continue
		}
		version := tokenValue(ts, r.VersionKeys...)
		if version == "" && r.Product != "" {
			version = productVersion(ua.Raw, r.Product)
		}
		setEngine(ua, r.Name, version)
		break
	}
	for _, r := range p.deviceRules {
		if r.Type == DeviceUnknown || !r.Match.matches(ts) {
			continue
		}
		ua.Device.Type = r.Type
		if r.Vendor != "" {
			ua.Device.Vendor = r.Vendor
		}
		if r.Model != "" {
			ua.Device.Model = r.Model
		}
		break
	}
}

// finalize fills derived fields and applies the few documented cross-field
// overrides that cannot be expressed as independent detectors.
func finalize(ua *UserAgent, ts *tokenSet) {
	if ua.Browser.Version != "" {
		ua.Browser.VersionNo = ParseVersion(ua.Browser.Version)
	}
	if ua.OS.Version != "" {
		ua.OS.VersionNo = ParseVersion(ua.OS.Version)
	}

	// HeyTapBrowser is Android-only. Some vendor UAs carry a bogus
	// "(X11; Windows)" platform comment.
	if ua.Browser.Name == HeyTapBrowser {
		if ua.OS.Name == "" || ua.OS.Name == Windows {
			setOS(ua, Android, "")
		}
		ua.Device.Type = DeviceMobile
	}

	// These browsers are Android-only in practice. Some vendor UAs report an
	// X11/Linux platform comment, which should not turn the device into a
	// desktop.
	switch ua.Browser.Name {
	case MiuiBrowser, SamsungBrowser, HuaweiBrowser, UCBrowser:
		if ua.OS.Name == "" || ua.OS.Name == Linux {
			setOS(ua, Android, androidVersion(ts, ua.Raw))
		}
		if ua.Device.Type == DeviceUnknown || ua.Device.Type == DeviceDesktop {
			ua.Device.Type = DeviceMobile
		}
	}

	if ua.Browser.Name == "QXSEC Android App" {
		setOS(ua, Android, ua.OS.Version)
		ua.Device.Type = DeviceMobile
		if ua.App.Name == "" {
			ua.App.Name = "QxsecApp"
			ua.App.Version = ua.Browser.Version
			ua.App.InApp = true
		}
		ua.WebView = true
	}

	if ua.Browser.Name == AppleWebKit || ua.Browser.Name == CFNetwork {
		ua.WebView = true
	}
	if ua.App.Name == CFNetwork {
		ua.WebView = true
	}

	// A primary bot replaces an incidental browser token, but social preview
	// bots such as facebookexternalhit keep the underlying browser.
	if ua.Bot.Detected {
		if ua.Browser.Name == "" || botOverridesBrowser(ua.Bot.Name) {
			ua.Browser.Name = ua.Bot.Name
			ua.Browser.Version = ua.Bot.Version
			ua.Browser.VersionNo = ParseVersion(ua.Bot.Version)
		}
	}

	if ua.OS.Name == Android && ua.Device.Type == DeviceUnknown {
		ua.Device.Type = DeviceMobile
	}
	if ua.OS.Name == IOS && ua.Device.Type == DeviceUnknown {
		ua.Device.Type = DeviceMobile
	}
	if ua.OS.Name == WindowsPhone || ua.OS.Name == KaiOS ||
		ua.OS.Name == BlackBerry || ua.OS.Name == HarmonyOS {
		if ua.Device.Type == DeviceUnknown || ua.Device.Type == DeviceDesktop {
			ua.Device.Type = DeviceMobile
		}
	}
	if ua.Device.Type == DeviceUnknown {
		switch ua.OS.Name {
		case Windows, MacOS, Linux, ChromeOS, FreeBSD, OpenBSD, NetBSD, Solaris:
			ua.Device.Type = DeviceDesktop
		}
	}

	if ua.Device.Type == DeviceUnknown && ua.App.InApp {
		ua.Device.Type = DeviceMobile
	}

	// An app container is treated as a web view.
	if ua.App.InApp {
		ua.WebView = true
	}
}

func botOverridesBrowser(name string) bool {
	switch name {
	case FacebookExternal, "Facebot", Twitterbot, LinkedInBot, Slackbot, Discordbot,
		TelegramBot, "Pinterestbot", "redditbot", "facebookcatalog", "Google-InspectionTool":
		return false
	default:
		return name != ""
	}
}
