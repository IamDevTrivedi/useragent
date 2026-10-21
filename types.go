package useragent

import (
	"encoding/json"
	"fmt"
	"strings"
)

// UserAgent is the result of parsing one user-agent string.
// The zero value is valid and represents an unknown client.
type UserAgent struct {
	Raw     string  `json:"raw"`
	Browser Browser `json:"browser"`
	OS      OS      `json:"os"`
	Device  Device  `json:"device"`
	Engine  Engine  `json:"engine"`
	CPU     CPU     `json:"cpu"`
	Bot     Bot     `json:"bot"`
	App     App     `json:"app"`
	WebView bool    `json:"webview"`
}

// Browser describes the detected browser or client.
type Browser struct {
	Name      string  `json:"name,omitempty"`
	Version   string  `json:"version,omitempty"`
	VersionNo Version `json:"version_no"`
}

// Is reports whether the browser name matches name.
func (b Browser) Is(name string) bool { return name != "" && b.Name == name }

// OS describes the detected operating system.
type OS struct {
	Name      string  `json:"name,omitempty"`
	Version   string  `json:"version,omitempty"`
	VersionNo Version `json:"version_no"`
}

// Is reports whether the OS name matches name.
func (o OS) Is(name string) bool { return name != "" && o.Name == name }

// DeviceType classifies the device.
type DeviceType uint8

// Device types returned by Device.Type.
const (
	DeviceUnknown DeviceType = iota
	DeviceDesktop
	DeviceMobile
	DeviceTablet
	DeviceTV
	DeviceConsole
	DeviceWearable
	DeviceEmbedded
)

// String returns the lower-case device type name.
func (d DeviceType) String() string {
	switch d {
	case DeviceDesktop:
		return "desktop"
	case DeviceMobile:
		return "mobile"
	case DeviceTablet:
		return "tablet"
	case DeviceTV:
		return "tv"
	case DeviceConsole:
		return "console"
	case DeviceWearable:
		return "wearable"
	case DeviceEmbedded:
		return "embedded"
	default:
		return "unknown"
	}
}

// MarshalJSON encodes the device type as a string.
func (d DeviceType) MarshalJSON() ([]byte, error) { return json.Marshal(d.String()) }

// UnmarshalJSON accepts the string form produced by MarshalJSON. Numeric
// values are accepted for backwards compatibility with simple integer output.
func (d *DeviceType) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err == nil {
		switch strings.ToLower(s) {
		case "unknown":
			*d = DeviceUnknown
		case "desktop":
			*d = DeviceDesktop
		case "mobile":
			*d = DeviceMobile
		case "tablet":
			*d = DeviceTablet
		case "tv":
			*d = DeviceTV
		case "console":
			*d = DeviceConsole
		case "wearable":
			*d = DeviceWearable
		case "embedded":
			*d = DeviceEmbedded
		default:
			return fmt.Errorf("useragent: unknown device type %q", s)
		}
		return nil
	}
	var n uint8
	if err := json.Unmarshal(data, &n); err == nil {
		*d = DeviceType(n)
		return nil
	}
	return fmt.Errorf("useragent: invalid device type %s", data)
}

// Device describes the device family, vendor, and model.
type Device struct {
	Type   DeviceType `json:"type"`
	Vendor string     `json:"vendor,omitempty"`
	Model  string     `json:"model,omitempty"`
}

// IsMobile reports whether the device is a phone.
func (d Device) IsMobile() bool { return d.Type == DeviceMobile }

// Engine describes the browser engine.
type Engine struct {
	Name    string `json:"name,omitempty"`
	Version string `json:"version,omitempty"`
}

// CPU describes the CPU architecture reported by the user agent.
type CPU struct {
	// Architecture is normalized to Go-style names where possible:
	// "amd64", "386", "arm64", "arm", "ppc64", "riscv64", "mips", or "".
	Architecture string `json:"architecture,omitempty"`
}

// BotCategory classifies a detected bot.
type BotCategory uint8

// Bot categories returned by Bot.Category.
const (
	BotUnknown BotCategory = iota
	BotSearch
	BotSocial
	BotAICrawler
	BotSEOTool
	BotMonitoring
	BotSecurityScanner
	BotScraper
	BotHTTPClient
	BotHeadless
	BotAds
	BotOther
)

// String returns the lower-case bot category name.
func (c BotCategory) String() string {
	switch c {
	case BotSearch:
		return "search"
	case BotSocial:
		return "social"
	case BotAICrawler:
		return "ai-crawler"
	case BotSEOTool:
		return "seo-tool"
	case BotMonitoring:
		return "monitoring"
	case BotSecurityScanner:
		return "security-scanner"
	case BotScraper:
		return "scraper"
	case BotHTTPClient:
		return "http-client"
	case BotHeadless:
		return "headless"
	case BotAds:
		return "ads"
	case BotOther:
		return "other"
	default:
		return "unknown"
	}
}

// MarshalJSON encodes the bot category as a string.
func (c BotCategory) MarshalJSON() ([]byte, error) { return json.Marshal(c.String()) }

// UnmarshalJSON accepts the string form produced by MarshalJSON. Numeric
// values are accepted for backwards compatibility with simple integer output.
func (c *BotCategory) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err == nil {
		switch strings.ToLower(s) {
		case "unknown":
			*c = BotUnknown
		case "search":
			*c = BotSearch
		case "social":
			*c = BotSocial
		case "ai-crawler":
			*c = BotAICrawler
		case "seo-tool":
			*c = BotSEOTool
		case "monitoring":
			*c = BotMonitoring
		case "security-scanner":
			*c = BotSecurityScanner
		case "scraper":
			*c = BotScraper
		case "http-client":
			*c = BotHTTPClient
		case "headless":
			*c = BotHeadless
		case "ads":
			*c = BotAds
		case "other":
			*c = BotOther
		default:
			return fmt.Errorf("useragent: unknown bot category %q", s)
		}
		return nil
	}
	var n uint8
	if err := json.Unmarshal(data, &n); err == nil {
		*c = BotCategory(n)
		return nil
	}
	return fmt.Errorf("useragent: invalid bot category %s", data)
}

// Bot describes a detected bot or automated client.
type Bot struct {
	Detected bool        `json:"detected"`
	Name     string      `json:"name,omitempty"`
	Version  string      `json:"version,omitempty"`
	URL      string      `json:"url,omitempty"`
	Category BotCategory `json:"category"`
}

// App describes an embedded application when one is detectable.
type App struct {
	Name    string `json:"name,omitempty"`
	Version string `json:"version,omitempty"`
	InApp   bool   `json:"in_app"`
}

// Browser and client names used by the detector and exposed for comparisons.
const (
	Chrome           = "Chrome"
	Chromium         = "Chromium"
	Firefox          = "Firefox"
	Safari           = "Safari"
	MobileSafari     = "Mobile Safari"
	Edge             = "Edge"
	Opera            = "Opera"
	OperaMini        = "Opera Mini"
	OperaTouch       = "Opera Touch"
	InternetExplorer = "Internet Explorer"
	Vivaldi          = "Vivaldi"
	Brave            = "Brave"
	Arc              = "Arc"
	DuckDuckGo       = "DuckDuckGo"
	CocCoc           = "Coc Coc"
	YandexBrowser    = "Yandex Browser"
	Whale            = "Whale"
	QQBrowser        = "QQBrowser"
	UCBrowser        = "UCBrowser"
	HeyTapBrowser    = "HeyTapBrowser"
	SamsungBrowser   = "Samsung Browser"
	MiuiBrowser      = "Miui Browser"
	HuaweiBrowser    = "Huawei Browser"
	BaiduBrowser     = "Baidu Browser"
	SogouBrowser     = "Sogou Browser"
	QuarkBrowser     = "Quark"
	Browser360       = "360 Browser"
	NetFront         = "NetFront"
	Konqueror        = "Konqueror"
	HeadlessChrome   = "HeadlessChrome"
	AppleWebKit      = "AppleWebKit"
	CFNetwork        = "CFNetwork"
	Splash           = "Splash"
	Scrapy           = "Scrapy"
	GoogleOther      = "GoogleOther"
	Googlebot        = "Googlebot"
	Bingbot          = "Bingbot"
	YandexBot        = "YandexBot"
	Baiduspider      = "Baiduspider"
	DuckDuckBot      = "DuckDuckBot"
	Applebot         = "Applebot"
	Bytespider       = "Bytespider"
	AdsBotGoogle     = "AdsBot-Google"
	YahooAdMonitor   = "Yahoo Ad monitoring"
	FacebookExternal = "facebookexternalhit"
	FacebookCatalog  = "facebookcatalog"
	Twitterbot       = "Twitterbot"
	LinkedInBot      = "LinkedInBot"
	Slackbot         = "Slackbot"
	Discordbot       = "Discordbot"
	TelegramBot      = "TelegramBot"
	AndroidBrowser   = "Android browser"
	UnknownBrowser   = ""
)

// OS names used by the detector and exposed for comparisons.
const (
	Windows      = "Windows"
	WindowsPhone = "Windows Phone"
	Android      = "Android"
	IOS          = "iOS"
	MacOS        = "macOS"
	Linux        = "Linux"
	ChromeOS     = "ChromeOS"
	Tizen        = "Tizen"
	WebOS        = "webOS"
	HarmonyOS    = "HarmonyOS"
	KaiOS        = "KaiOS"
	BlackBerry   = "BlackBerry"
	FreeBSD      = "FreeBSD"
	OpenBSD      = "OpenBSD"
	NetBSD       = "NetBSD"
	Solaris      = "Solaris"
)

// Bot names used by the detector and exposed for comparisons.
const (
	GPTBot         = "GPTBot"
	ChatGPTUser    = "ChatGPT-User"
	OAI            = "OAI-SearchBot"
	ClaudeBot      = "ClaudeBot"
	ClaudeWeb      = "Claude-Web"
	PerplexityBot  = "PerplexityBot"
	CCBot          = "CCBot"
	Amazonbot      = "Amazonbot"
	SemrushBot     = "SemrushBot"
	AhrefsBot      = "AhrefsBot"
	MJ12bot        = "MJ12bot"
	DotBot         = "DotBot"
	UptimeRobot    = "UptimeRobot"
	Pingdom        = "Pingdom"
	StatusCake     = "StatusCake"
	Nuclei         = "Nuclei"
	Censys         = "Censys"
	Shodan         = "Shodan"
	ZGrab          = "ZGrab"
	Masscan        = "masscan"
	Nmap           = "Nmap"
	SQLmap         = "sqlmap"
	Nikto          = "Nikto"
	Headless       = "HeadlessChrome"
	PhantomJS      = "PhantomJS"
	Selenium       = "Selenium"
	Playwright     = "Playwright"
	Puppeteer      = "Puppeteer"
	HTTrack        = "HTTrack"
	Curl           = "curl"
	Wget           = "Wget"
	GoHTTPClient   = "Go-http-client"
	PythonRequests = "python-requests"
	OkHTTP         = "okhttp"
	JavaHTTP       = "Java"
	Axios          = "axios"
	NodeFetch      = "node-fetch"
	Undici         = "undici"
	Postman        = "PostmanRuntime"
	Insomnia       = "Insomnia"
	HTTPie         = "HTTPie"
)

// String returns a compact human-readable summary suitable for logs.
func (u UserAgent) String() string {
	var b strings.Builder
	switch {
	case u.Bot.Detected:
		b.WriteString("bot ")
		if u.Bot.Name != "" {
			b.WriteString(u.Bot.Name)
			if u.Bot.Version != "" {
				b.WriteByte(' ')
				b.WriteString(u.Bot.Version)
			}
		} else {
			b.WriteString("unknown")
		}
	case u.Browser.Name != "":
		b.WriteString(u.Browser.Name)
		if u.Browser.Version != "" {
			b.WriteByte(' ')
			b.WriteString(u.Browser.Version)
		}
	default:
		b.WriteString("unknown")
	}

	if u.OS.Name != "" {
		b.WriteString(" on ")
		b.WriteString(u.OS.Name)
		if u.OS.Version != "" {
			b.WriteByte(' ')
			b.WriteString(u.OS.Version)
		}
	}
	if u.Device.Type != DeviceUnknown {
		b.WriteString(" (")
		b.WriteString(u.Device.Type.String())
		if u.Device.Model != "" {
			b.WriteString(": ")
			b.WriteString(u.Device.Model)
		}
		b.WriteByte(')')
	}
	if u.App.InApp && u.App.Name != "" {
		b.WriteString(" app=")
		b.WriteString(u.App.Name)
	}
	return b.String()
}

// IsBot reports whether the user agent was identified as a bot.
func (u UserAgent) IsBot() bool { return u.Bot.Detected }

// IsMobile reports whether the device is a phone.
func (u UserAgent) IsMobile() bool { return u.Device.Type == DeviceMobile }

// IsTablet reports whether the device is a tablet.
func (u UserAgent) IsTablet() bool { return u.Device.Type == DeviceTablet }

// IsDesktop reports whether the device is a desktop or laptop.
func (u UserAgent) IsDesktop() bool { return u.Device.Type == DeviceDesktop }

// IsTV reports whether the device is a TV or set-top box.
func (u UserAgent) IsTV() bool { return u.Device.Type == DeviceTV }

// IsConsole reports whether the device is a game console.
func (u UserAgent) IsConsole() bool { return u.Device.Type == DeviceConsole }

// IsWebView reports whether the request came from an embedded web view.
func (u UserAgent) IsWebView() bool { return u.WebView }

// IsUnknown reports whether nothing useful could be determined.
func (u UserAgent) IsUnknown() bool {
	return u.Browser.Name == "" &&
		u.OS.Name == "" &&
		u.Device.Type == DeviceUnknown &&
		u.App.Name == "" &&
		!u.Bot.Detected
}

// IsChrome reports whether the browser is Chrome.
func (u UserAgent) IsChrome() bool { return u.Browser.Name == Chrome }

// IsFirefox reports whether the browser is Firefox.
func (u UserAgent) IsFirefox() bool { return u.Browser.Name == Firefox }

// IsSafari reports whether the browser is Safari.
func (u UserAgent) IsSafari() bool { return u.Browser.Name == Safari }

// IsEdge reports whether the browser is Edge.
func (u UserAgent) IsEdge() bool { return u.Browser.Name == Edge }

// IsOpera reports whether the browser is Opera.
func (u UserAgent) IsOpera() bool { return u.Browser.Name == Opera }

// IsAndroid reports whether the OS is Android.
func (u UserAgent) IsAndroid() bool { return u.OS.Name == Android }

// IsWindows reports whether the OS is Windows.
func (u UserAgent) IsWindows() bool { return u.OS.Name == Windows }

// IsMacOS reports whether the OS is macOS.
func (u UserAgent) IsMacOS() bool { return u.OS.Name == MacOS }

// IsIOS reports whether the OS is iOS.
func (u UserAgent) IsIOS() bool { return u.OS.Name == IOS }

// IsLinux reports whether the OS is Linux.
func (u UserAgent) IsLinux() bool { return u.OS.Name == Linux }

// IsChromeOS reports whether the OS is ChromeOS.
func (u UserAgent) IsChromeOS() bool { return u.OS.Name == ChromeOS }

// IsBlackBerry reports whether the OS is BlackBerry.
func (u UserAgent) IsBlackBerry() bool { return u.OS.Name == BlackBerry }

// IsGooglebot reports whether the bot is Googlebot.
func (u UserAgent) IsGooglebot() bool { return u.Bot.Name == Googlebot }

// IsTwitterbot reports whether the bot is Twitterbot.
func (u UserAgent) IsTwitterbot() bool { return u.Bot.Name == Twitterbot }

// IsFacebookbot reports whether the bot is facebookexternalhit.
func (u UserAgent) IsFacebookbot() bool { return u.Bot.Name == FacebookExternal }

// IsYandexbot reports whether the bot is YandexBot.
func (u UserAgent) IsYandexbot() bool { return u.Bot.Name == YandexBot }

// Summary returns the same value as String; it is provided for readability at call sites.
func (u UserAgent) Summary() string { return u.String() }
