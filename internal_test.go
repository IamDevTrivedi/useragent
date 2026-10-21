package useragent

import (
	"encoding/json"
	"fmt"
	"strings"
	"testing"
)

func TestCPUArchitectureTaxonomy(t *testing.T) {
	tests := []struct {
		name string
		raw  string
		want string
	}{
		{"arm64", "Mozilla/5.0 (X11; aarch64)", "arm64"},
		{"x86_64", "Mozilla/5.0 (X11; x86_64)", "amd64"},
		{"amd64", "Mozilla/5.0 (X11; amd64)", "amd64"},
		{"win64", "Mozilla/5.0 (Windows NT 10.0; Win64)", "amd64"},
		{"x64", "Mozilla/5.0 (Windows NT 10.0; x64)", "amd64"},
		{"wow64", "Mozilla/5.0 (Windows NT 6.1; WOW64)", "386"},
		{"i686", "Mozilla/5.0 (X11; i686)", "386"},
		{"i386", "Mozilla/5.0 (X11; i386)", "386"},
		{"i86pc", "Mozilla/5.0 (X11; i86pc)", "386"},
		{"armv7l", "Mozilla/5.0 (X11; armv7l)", "arm"},
		{"armv8l", "Mozilla/5.0 (X11; armv8l)", "arm"},
		{"armv7", "Mozilla/5.0 (X11; armv7)", "arm"},
		{"ppc64", "Mozilla/5.0 (X11; ppc64)", "ppc64"},
		{"riscv64", "Mozilla/5.0 (X11; riscv64)", "riscv64"},
		{"mips64", "Mozilla/5.0 (X11; mips64)", "mips"},
		{"mips", "Mozilla/5.0 (X11; mips)", "mips"},
		{"none", "Mozilla/5.0", ""},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Parse(tt.raw).CPU.Architecture; got != tt.want {
				t.Errorf("CPU = %q, want %q", got, tt.want)
			}
		})
	}

	ua := UserAgent{CPU: CPU{Architecture: "wasm"}}
	ts := tokenize("x86_64")
	detectCPU(&ua, &ts)
	if ua.CPU.Architecture != "wasm" {
		t.Errorf("existing CPU was overwritten: %q", ua.CPU.Architecture)
	}
}

func TestInferVendor(t *testing.T) {
	tests := []struct {
		model string
		raw   string
		want  string
	}{
		{"SM-G960F", "", "Samsung"},
		{"GT-I9300", "", "Samsung"},
		{"SGH-I747", "", "Samsung"},
		{"SCH-I535", "", "Samsung"},
		{"LM-G710", "", "LG"},
		{"LG-H870", "", "LG"},
		{"CPH2451", "", "OPPO"},
		{"RMX3630", "", "Realme"},
		{"VOG-L29", "", "Huawei"},
		{"ELE-L29", "", "Huawei"},
		{"ANE-LX1", "", "Huawei"},
		{"JNY-LX1", "", "Huawei"},
		{"MAR-LX1A", "", "Huawei"},
		{"XT2041", "", "Motorola"},
		{"TB-X606F", "", "Lenovo"},
		{"Pixel 8", "", "Google"},
		{"Nexus 5X", "", "Google"},
		{"Redmi Note 12", "", "Xiaomi"},
		{"POCO F5", "", "Xiaomi"},
		{"Mi 11", "", "Xiaomi"},
		{"MI 11", "", "Xiaomi"},
		{"M2101K9G", "", "Xiaomi"},
		{"BRAVIA", "", "Sony"},
		{"Xperia 5", "", "Sony"},
		{"Galaxy S23", "SAMSUNG", "Samsung"},
		{"P30", "HUAWEI", "Huawei"},
		{"Mi 11", "Xiaomi", "Xiaomi"},
		{"Find X6", "OPPO", "OPPO"},
		{"V27", "vivo", "vivo"},
		{"GT Neo", "realme", "Realme"},
		{"3310", "Nokia", "Nokia"},
		{"Xperia", "Sony", "Sony"},
		{"Unknown", "", ""},
	}
	for _, tt := range tests {
		ts := tokenize(tt.raw)
		if got := inferVendor(tt.model, &ts); got != tt.want {
			t.Errorf("inferVendor(%q, %q) = %q, want %q", tt.model, tt.raw, got, tt.want)
		}
	}
}

func TestParseAndroidModelHelpers(t *testing.T) {
	tests := []struct {
		comment string
		raw     string
		vendor  string
		model   string
		ok      bool
	}{
		{"SAMSUNG SM-A310F", "", "Samsung", "SM-A310F", true},
		{"SM-G960F", "", "Samsung", "SM-G960F", true},
		{"SM-G960F/DS", "", "Samsung", "SM-G960F", true},
		{"Galaxy S23", "SAMSUNG", "Samsung", "Galaxy S23", true},
		{"VOG-L29", "", "Huawei", "VOG-L29", true},
		{"", "", "", "", false},
		{"a", "", "", "", false},
		{"Build/MMB29P", "", "", "", false},
		{"http://example.com", "", "", "", false},
		{"4GB", "", "", "", false},
		{"en-US", "", "", "", false},
		{"1080x2166", "", "", "", false},
		{"13", "", "", "", false},
		{"Xyz", "", "", "", false},
		{strings.Repeat("x", 49), "", "", "", false},
	}
	for _, tt := range tests {
		ts := tokenize(tt.raw)
		vendor, model, ok := parseAndroidModel(tt.comment, &ts)
		if ok != tt.ok || vendor != tt.vendor || model != tt.model {
			t.Errorf("parseAndroidModel(%q, %q) = %q, %q, %v; want %q, %q, %v",
				tt.comment, tt.raw, vendor, model, ok, tt.vendor, tt.model, tt.ok)
		}
	}
}

func TestIsNoiseComment(t *testing.T) {
	noise := []string{
		"", "u", "wv", "mobile", "tablet", "tv", "phone", "linux", "x11",
		"windows", "macintosh", "android", "compatible", "build", "hybrid()",
		"globaldist", "googleplaystore", "watch", "cros", "khtml, like gecko",
		"u; en", "Build/MMB29P", "Version/4.0", "Name:HeyTapBrowser",
		"http://example.com", "GlobalDist", "GooglePlayStore", "13", "13.0",
		"1080x2166", "4GB", "4.5GB", "en-US", "Chrome/120", "Safari/537",
		"AppleWebKit/537", "Gecko/20100101", "Opera/9.8", "Edg/120",
		strings.Repeat("x", 49),
	}
	for _, s := range noise {
		if !isNoiseComment(s) {
			t.Errorf("isNoiseComment(%q) = false, want true", s)
		}
	}
	notNoise := []string{"SM-G960F", "SAMSUNG SM-A310F", "1234", "Pixel 8", "M2101K9G"}
	for _, s := range notNoise {
		if isNoiseComment(s) {
			t.Errorf("isNoiseComment(%q) = true, want false", s)
		}
	}
}

func TestIsMobileBrowserNames(t *testing.T) {
	names := []string{
		MobileSafari, OperaMini, OperaTouch, UCBrowser, QQBrowser, HeyTapBrowser,
		MiuiBrowser, HuaweiBrowser, SamsungBrowser, NetFront, AndroidBrowser,
		Brave, Chrome, Firefox, Safari, CocCoc, YandexBrowser, BaiduBrowser,
		SogouBrowser, QuarkBrowser, Browser360,
	}
	for _, name := range names {
		if !isMobileBrowser(name) {
			t.Errorf("isMobileBrowser(%q) = false", name)
		}
	}
	if isMobileBrowser("") || isMobileBrowser("curl") {
		t.Error("desktop/unknown browser reported as mobile")
	}
}

func TestOSVersionFallbacks(t *testing.T) {
	var ts tokenSet
	tests := []struct {
		name string
		got  string
		want string
	}{
		{"ios iphone os", iosVersion(&ts, "Mozilla/5.0 (iPhone OS 12_0)"), "12.0"},
		{"ios cpu os", iosVersion(&ts, "Mozilla/5.0 (CPU OS 17_2)"), "17.2"},
		{"ios cpu iphone", iosVersion(&ts, "Mozilla/5.0 (CPU iPhone OS 16_3_1)"), "16.3.1"},
		{"ios slash", iosVersion(&ts, "Darwin iOS/16.3.1"), "16.3.1"},
		{"chromeos", chromeOSVersion(&ts, "X11; CrOS x86_64 14150.74.0)"), "14150.74.0"},
		{"chromeos punctuation", chromeOSVersion(&ts, "X11; CrOS x86_64 14388.62.1);"), "14388.62.1"},
		{"chromeos none", chromeOSVersion(&ts, "X11; CrOS x86_64"), ""},
		{"tizen", tizenVersion(&ts, "Tizen 6.5"), "6.5"},
		{"webos", webOSVersion(&ts, "webOS/6.0"), "6.0"},
		{"webos upper", webOSVersion(&ts, "Web0S/2.0"), "2.0"},
		{"harmony", harmonyVersion(&ts, "HarmonyOS 4.0"), "4.0"},
		{"kaios", kaiVersion(&ts, "KaiOS/2.5"), "2.5"},
		{"blackberry", blackBerryVersion("BlackBerry 9900 Version/7.0.0.187"), "7.0.0.187"},
		{"windows phone os", windowsPhoneVersion(&ts, "Windows Phone OS 7.0"), "7.0"},
		{"windows phone", windowsPhoneVersion(&ts, "Windows Phone 8.1"), "8.1"},
		{"mac", macVersion(&ts, "Mac OS X 10.15.7"), "10.15.7"},
		{"linux", linuxVersion("Linux 6.1.0"), "6.1.0"},
		{"android", androidVersion(&ts, "Android 13"), "13"},
		{"android app", androidVersion(&ts, "Android App"), ""},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.got != tt.want {
				t.Errorf("got %q, want %q", tt.got, tt.want)
			}
		})
	}
}

func TestIOSDetectionHelpers(t *testing.T) {
	tests := []struct {
		name string
		raw  string
		want bool
	}{
		{"iphone", "iPhone", true},
		{"ipad", "iPad", true},
		{"ipod", "iPod", true},
		{"token iphone os", "Mobile/15E148", true},
		{"token ios darwin", "iOS Darwin", true},
		{"token ios cfnetwork", "iOS CFNetwork", true},
		{"token ios slash", "iOS/17.0", true},
		{"token ios alone", "iOS", false},
		{"none", "Mozilla/5.0", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ts := tokenize(tt.raw)
			if tt.name == "token iphone os" {
				ts.add("iPhone OS", "17.0", true)
			}
			if got := isIOS(&ts); got != tt.want {
				t.Errorf("isIOS(%q) = %v, want %v", tt.raw, got, tt.want)
			}
		})
	}
}

func TestTokenSetHelpers(t *testing.T) {
	var ts tokenSet
	ts.add("Mobile Safari", "537", false)
	ts.add("Chrome", "120", false)
	ts.add("Empty", "", false)
	ts.add("", "ignored", false)

	if got := ts.get("mobile safari"); got != "537" {
		t.Errorf("get = %q", got)
	}
	if !ts.has("CHROME") || ts.has("missing") {
		t.Error("has failed")
	}
	if !ts.hasAny("missing", "Chrome") || ts.hasAny("missing", "nada") {
		t.Error("hasAny failed")
	}
	if !ts.hasAll("Chrome", "Mobile Safari") || ts.hasAll("Chrome", "missing") || ts.hasAll() {
		t.Error("hasAll failed")
	}
	if !ts.hasPrefix("Mobile") || ts.hasPrefix("Nope") {
		t.Error("hasPrefix failed")
	}
	if !ts.hasSuffix("Safari") || ts.hasSuffix("Nope") {
		t.Error("hasSuffix failed")
	}
	if !ts.rawHas("537") || !ts.rawHas("Chro") || ts.rawHas("Nope") {
		t.Error("rawHas failed")
	}
	if !hasSuffixFold("Mobile Safari", "safari") || hasSuffixFold("a", "long") {
		t.Error("hasSuffixFold failed")
	}
	if !containsFold("Mobile Safari", "safari") || !containsFold("anything", "") || containsFold("a", "long") {
		t.Error("containsFold failed")
	}
	if !equalFold("Chrome", "chrome") || equalFold("Chrome", "chrom") {
		t.Error("equalFold failed")
	}
	if lowerASCII('A') != 'a' || lowerASCII('z') != 'z' {
		t.Error("lowerASCII failed")
	}

	for i := 0; i < maxTokens+5; i++ {
		ts.add(fmt.Sprintf("k%d", i), "", false)
	}
	if ts.n != maxTokens {
		t.Errorf("token limit = %d, want %d", ts.n, maxTokens)
	}
	for i := 0; i < maxComments+5; i++ {
		ts.addComment(fmt.Sprintf("c%d", i))
	}
	if ts.cn != maxComments {
		t.Errorf("comment limit = %d, want %d", ts.cn, maxComments)
	}
}

func TestTokenizerEdges(t *testing.T) {
	ts := tokenize(`Foo/1 (Bar; Baz/2) [Qux/3] "unterminated (Nested/4)`)
	if !ts.has("Foo") || !ts.has("Bar") || !ts.has("Baz") || !ts.has("Qux") || !ts.has("Nested") {
		t.Errorf("tokenize missed fields: %+v", ts.list[:ts.n])
	}

	if key, value, next := scanProduct("", 0); key != "" || value != "" || next != 1 {
		t.Errorf("scanProduct empty = %q %q %d", key, value, next)
	}
	if key, value, _ := scanProduct("/1", 0); key != "" || value != "" {
		t.Errorf("scanProduct empty key = %q %q", key, value)
	}
	if key, value, _ := scanProduct("Foo/", 0); key != "Foo" || value != "" {
		t.Errorf("scanProduct empty value = %q %q", key, value)
	}
	if key, value, _ := scanProduct("Mobile Safari/537", 0); key != "Mobile Safari" || value != "537" {
		t.Errorf("scanProduct multiword = %q %q", key, value)
	}
	if key, value, _ := scanProduct("Foo Bar/1", 0); key != "Foo" || value != "" {
		t.Errorf("scanProduct non-multiword = %q %q", key, value)
	}

	for _, product := range []string{"Mobile Safari", "TV Safari", "Brave Chrome", "Opera Mini",
		"Opera Touch", "Internet Explorer", "Microsoft Edge", "Google Chrome", "Samsung Internet"} {
		if !isMultiWordProduct(product) {
			t.Errorf("isMultiWordProduct(%q) = false", product)
		}
	}
	if isMultiWordProduct("Not A Product") {
		t.Error("unknown multi-word product accepted")
	}

	urlTests := []struct {
		raw  string
		want string
	}{
		{"see https://example.com/path, next", "https://example.com/path"},
		{"http://example.com/a. more", "http://example.com/a"},
		{"HTTP://EXAMPLE.COM", "HTTP://EXAMPLE.COM"},
		{"https://example.com/a!?;:", "https://example.com/a"},
		{"no url here", ""},
	}
	for _, tt := range urlTests {
		if got := extractURL(tt.raw); got != tt.want {
			t.Errorf("extractURL(%q) = %q, want %q", tt.raw, got, tt.want)
		}
	}
	for _, c := range []byte{'.', ',', ';', ':', '!', '?'} {
		if !isTrailingURLPunct(c) {
			t.Errorf("isTrailingURLPunct(%q) = false", c)
		}
	}
	if isTrailingURLPunct('/') || isTrailingURLPunct('a') {
		t.Error("non-punctuation reported as trailing punctuation")
	}

	if got := firstVersionAfter("plain text", "Chrome/"); got != "" {
		t.Errorf("firstVersionAfter = %q", got)
	}
	if got := firstVersionAfter("Chrome/", "Chrome/"); got != "" {
		t.Errorf("firstVersionAfter empty = %q", got)
	}
	if !isNumericOnly("12.3") || isNumericOnly("12a") || isNumericOnly("") {
		t.Error("isNumericOnly failed")
	}
}

func TestAppTaxonomy(t *testing.T) {
	tests := []struct {
		name    string
		raw     string
		app     string
		version string
		inApp   bool
		webview bool
	}{
		{"qxsec-app", "Mozilla/5.0 (Linux; Android 13; 22021211RC) AppleWebKit/537.36 Chrome/121.0 Mobile Safari/537.36 QxsecApp/1.435.390", "QxsecApp", "1.435.390", true, true},
		{"qxsec-raw", "Mozilla/5.0 (Linux; Android 13; 22021211RC) AppleWebKit/537.36 QXSEC Android App 1.0 Phone Hybrid()", "QXSEC", "", true, true},
		{"facebook", "Mozilla/5.0 (iPhone; CPU iPhone OS 16_0 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Mobile/15E148 [FBAN/FBIOS;FBAV/123.0.0.1;FBBV/1]", "Facebook", "123.0.0.1", true, true},
		{"facebook-iab", "Mozilla/5.0 (Linux; Android 13) AppleWebKit/537.36 [FB_IAB/FB4A;FBAV/124.0.0.1;]", "Facebook", "124.0.0.1", true, true},
		{"instagram", "Mozilla/5.0 (iPhone; CPU iPhone OS 16_0 like Mac OS X) AppleWebKit/605.1.15 Mobile/15E148 Instagram 123.0", "Instagram", "123.0", true, true},
		{"tiktok", "Mozilla/5.0 (Linux; Android 13) AppleWebKit/537.36 musical_ly_28.3.4", "TikTok", "28.3.4", true, true},
		{"wechat", "Mozilla/5.0 (Linux; Android 13) AppleWebKit/537.36 MicroMessenger/8.0.40", "WeChat", "8.0.40", true, true},
		{"qq", "Mozilla/5.0 (Linux; Android 13) AppleWebKit/537.36 QQ/9.1.0", "QQ", "9.1.0", true, true},
		{"line", "Mozilla/5.0 (iPhone; CPU iPhone OS 16_0 like Mac OS X) AppleWebKit/605.1.15 Line/13.0", "LINE", "13.0", true, true},
		{"pinterest", "Mozilla/5.0 (iPhone; CPU iPhone OS 16_0 like Mac OS X) AppleWebKit/605.1.15 Pinterest/12.0", "Pinterest", "12.0", true, true},
		{"snapchat", "Mozilla/5.0 (iPhone; CPU iPhone OS 16_0 like Mac OS X) AppleWebKit/605.1.15 Snapchat/12.0", "Snapchat", "12.0", true, true},
		{"electron", "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 Chrome/120.0 Electron/28.0 Safari/537.36", "Electron", "28.0", true, true},
		{"qtwebengine", "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 QtWebEngine/5.6.0 Chrome/45.0 Safari/537.36", "QtWebEngine", "5.6.0", true, true},
		{"cef", "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 CEF/120.0 Chrome/120.0 Safari/537.36", "CEF", "120.0", true, true},
		{"gsa", "Mozilla/5.0 (iPhone; CPU iPhone OS 16_0 like Mac OS X) AppleWebKit/605.1.15 GSA/123.0", "Google", "123.0", true, true},
		{"webappmanager", "Mozilla/5.0 (Web0S; Linux/SmartTV) AppleWebKit/537.36 WebAppManager", "WebAppManager", "", true, true},
		{"bytedance", "Mozilla/5.0 (Linux; Android 13) AppleWebKit/537.36 BytedanceWebview", "", "", false, true},
		{"hybrid", "Mozilla/5.0 (Linux; Android 13) AppleWebKit/537.36 Hybrid()", "", "", false, true},
		{"none", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 Chrome/120.0 Safari/537.36", "", "", false, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Parse(tt.raw)
			if got.App.Name != tt.app || got.App.Version != tt.version ||
				got.App.InApp != tt.inApp || got.WebView != tt.webview {
				t.Errorf("app = %+v webview=%v, want %q %q inApp=%v webview=%v",
					got.App, got.WebView, tt.app, tt.version, tt.inApp, tt.webview)
			}
		})
	}
}

func TestClientHintBrandSelection(t *testing.T) {
	tests := []struct {
		name    string
		brands  []Brand
		want    string
		version string
	}{
		{"chrome", []Brand{{"Not_A Brand", "8"}, {"Chromium", "120"}, {"Google Chrome", "121"}}, Chrome, "121"},
		{"edge", []Brand{{"Microsoft Edge", "121"}}, Edge, "121"},
		{"opera", []Brand{{"Opera", "106"}}, Opera, "106"},
		{"brave", []Brand{{"Brave", "121"}}, Brave, "121"},
		{"vivaldi", []Brand{{"Vivaldi", "6.5"}}, Vivaldi, "6.5"},
		{"samsung", []Brand{{"Samsung Internet", "23"}}, SamsungBrowser, "23"},
		{"samsung-alt", []Brand{{"SamsungBrowser", "23"}}, SamsungBrowser, "23"},
		{"chromium", []Brand{{"Chromium", "120"}}, Chromium, "120"},
		{"webview", []Brand{{"Android WebView", "120"}}, Chrome, "120"},
		{"edge-fallback", []Brand{{"Edge Beta", "121"}}, Edge, "121"},
		{"opera-fallback", []Brand{{"Opera GX", "106"}}, Opera, "106"},
		{"chrome-fallback", []Brand{{"Chrome Mobile", "121"}}, Chrome, "121"},
		{"chromium-fallback", []Brand{{"Some Chromium Build", "121"}}, Chromium, "121"},
		{"none", []Brand{{"Not_A Brand", "8"}}, "", ""},
		{"empty", nil, "", ""},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			name, version := chooseClientHintBrand(tt.brands)
			if name != tt.want || version != tt.version {
				t.Errorf("choose = %q %q, want %q %q", name, version, tt.want, tt.version)
			}
		})
	}
}

func TestRuleProductFallbacks(t *testing.T) {
	p := New(
		WithBrowserRules(BrowserRule{Name: "MyBrowser", Match: Match{Token: "MyApp"}, Product: "MyApp"}),
		WithOSRules(OSRule{Name: "MyOS", Match: Match{Token: "MyApp"}, Product: "MyOS"}),
		WithBotRules(BotRule{Name: FacebookExternal, Match: Match{Token: "MyApp"}, Product: "MyBot", Category: BotOther}),
		WithEngineRules(EngineRule{Name: "MyEngine", Match: Match{Token: "MyApp"}, Product: "MyEngine"}),
	)
	ua := p.Parse("MyApp/1.0 MyOS/2.0 MyBot/3.0 MyEngine/4.0")
	if ua.Browser.Name != "MyBrowser" || ua.Browser.Version != "1.0" {
		t.Errorf("browser = %+v", ua.Browser)
	}
	if ua.OS.Name != "MyOS" || ua.OS.Version != "2.0" {
		t.Errorf("os = %+v", ua.OS)
	}
	if !ua.Bot.Detected || ua.Bot.Name != FacebookExternal || ua.Bot.Version != "3.0" {
		t.Errorf("bot = %+v", ua.Bot)
	}
	if ua.Engine.Name != "MyEngine" || ua.Engine.Version != "4.0" {
		t.Errorf("engine = %+v", ua.Engine)
	}
}

func TestEnumUnmarshalAll(t *testing.T) {
	devices := map[string]DeviceType{
		"unknown": DeviceUnknown, "desktop": DeviceDesktop, "mobile": DeviceMobile,
		"tablet": DeviceTablet, "tv": DeviceTV, "console": DeviceConsole,
		"wearable": DeviceWearable, "embedded": DeviceEmbedded,
	}
	for name, want := range devices {
		var got DeviceType
		if err := json.Unmarshal([]byte(`"`+name+`"`), &got); err != nil {
			t.Errorf("device %q: %v", name, err)
		} else if got != want {
			t.Errorf("device %q = %v, want %v", name, got, want)
		}
	}
	if err := json.Unmarshal([]byte(`true`), new(DeviceType)); err == nil {
		t.Error("invalid device JSON accepted")
	}

	categories := map[string]BotCategory{
		"unknown": BotUnknown, "search": BotSearch, "social": BotSocial,
		"ai-crawler": BotAICrawler, "seo-tool": BotSEOTool, "monitoring": BotMonitoring,
		"security-scanner": BotSecurityScanner, "scraper": BotScraper,
		"http-client": BotHTTPClient, "headless": BotHeadless, "ads": BotAds,
		"other": BotOther,
	}
	for name, want := range categories {
		var got BotCategory
		if err := json.Unmarshal([]byte(`"`+name+`"`), &got); err != nil {
			t.Errorf("category %q: %v", name, err)
		} else if got != want {
			t.Errorf("category %q = %v, want %v", name, got, want)
		}
	}
	if err := json.Unmarshal([]byte(`true`), new(BotCategory)); err == nil {
		t.Error("invalid bot category JSON accepted")
	}
}
