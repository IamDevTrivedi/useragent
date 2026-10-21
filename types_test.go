package useragent

import (
	"encoding/json"
	"reflect"
	"strings"
	"testing"
)

func TestDeviceTypeString(t *testing.T) {
	tests := []struct {
		typ  DeviceType
		want string
	}{
		{DeviceUnknown, "unknown"},
		{DeviceDesktop, "desktop"},
		{DeviceMobile, "mobile"},
		{DeviceTablet, "tablet"},
		{DeviceTV, "tv"},
		{DeviceConsole, "console"},
		{DeviceWearable, "wearable"},
		{DeviceEmbedded, "embedded"},
		{DeviceType(200), "unknown"},
	}
	for _, tt := range tests {
		if got := tt.typ.String(); got != tt.want {
			t.Errorf("DeviceType(%d).String() = %q, want %q", tt.typ, got, tt.want)
		}
	}
}

func TestBotCategoryString(t *testing.T) {
	tests := []struct {
		cat  BotCategory
		want string
	}{
		{BotUnknown, "unknown"},
		{BotSearch, "search"},
		{BotSocial, "social"},
		{BotAICrawler, "ai-crawler"},
		{BotSEOTool, "seo-tool"},
		{BotMonitoring, "monitoring"},
		{BotSecurityScanner, "security-scanner"},
		{BotScraper, "scraper"},
		{BotHTTPClient, "http-client"},
		{BotHeadless, "headless"},
		{BotAds, "ads"},
		{BotOther, "other"},
		{BotCategory(200), "unknown"},
	}
	for _, tt := range tests {
		if got := tt.cat.String(); got != tt.want {
			t.Errorf("BotCategory(%d).String() = %q, want %q", tt.cat, got, tt.want)
		}
	}
}

func TestEnumJSON(t *testing.T) {
	data, err := json.Marshal(DeviceMobile)
	if err != nil {
		t.Fatal(err)
	}
	if string(data) != `"mobile"` {
		t.Errorf("DeviceType JSON = %s", data)
	}
	var dt DeviceType
	if err := json.Unmarshal(data, &dt); err != nil {
		t.Fatal(err)
	}
	if dt != DeviceMobile {
		t.Errorf("DeviceType unmarshal = %v", dt)
	}
	if err := json.Unmarshal([]byte(`3`), &dt); err != nil {
		t.Fatal(err)
	}
	if dt != DeviceTablet {
		t.Errorf("numeric DeviceType = %v", dt)
	}
	if err := json.Unmarshal([]byte(`"martian"`), &dt); err == nil {
		t.Error("unknown device type should fail")
	}

	data, err = json.Marshal(BotSearch)
	if err != nil {
		t.Fatal(err)
	}
	if string(data) != `"search"` {
		t.Errorf("BotCategory JSON = %s", data)
	}
	var bc BotCategory
	if err := json.Unmarshal(data, &bc); err != nil {
		t.Fatal(err)
	}
	if bc != BotSearch {
		t.Errorf("BotCategory unmarshal = %v", bc)
	}
	if err := json.Unmarshal([]byte(`"search-bot"`), &bc); err == nil {
		t.Error("unknown bot category should fail")
	}
}

func TestUserAgentJSONRoundTrip(t *testing.T) {
	ua := Parse("Mozilla/5.0 (Linux; Android 6.0.1; Nexus 5X Build/MMB29P) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/118.0.5993.117 Mobile Safari/537.36 (compatible; Googlebot/2.1; +http://www.google.com/bot.html)")
	data, err := json.Marshal(ua)
	if err != nil {
		t.Fatal(err)
	}
	var got UserAgent
	if err := json.Unmarshal(data, &got); err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(ua, got) {
		t.Errorf("round trip mismatch:\n got %+v\nwant %+v\njson %s", got, ua, data)
	}
	if !strings.Contains(string(data), `"type":"mobile"`) {
		t.Errorf("device type not encoded as a string: %s", data)
	}
	if !strings.Contains(string(data), `"category":"search"`) {
		t.Errorf("bot category not encoded as a string: %s", data)
	}
}

func TestUserAgentString(t *testing.T) {
	tests := []struct {
		raw  string
		want []string
	}{
		{
			"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36",
			[]string{"Chrome 120.0.0.0", "on Windows 10.0", "(desktop)"},
		},
		{
			"",
			[]string{"unknown"},
		},
		{
			"Mozilla/5.0 (Linux; Android 6.0.1; Nexus 5X Build/MMB29P) AppleWebKit/537.36 Chrome/118.0.5993.117 Mobile Safari/537.36 (compatible; Googlebot/2.1; +http://www.google.com/bot.html)",
			[]string{"bot Googlebot 2.1", "on Android 6.0.1", "(mobile: Nexus 5X)"},
		},
	}
	for _, tt := range tests {
		got := Parse(tt.raw).String()
		for _, part := range tt.want {
			if !strings.Contains(got, part) {
				t.Errorf("String() = %q, missing %q", got, part)
			}
		}
	}
	if got := (UserAgent{}).Summary(); got != "unknown" {
		t.Errorf("Summary() = %q", got)
	}
}

func TestUserAgentHelpers(t *testing.T) {
	ua := Parse("Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36")
	if !ua.IsChrome() || !ua.IsWindows() || !ua.IsDesktop() {
		t.Errorf("helpers failed for %+v", ua)
	}
	if ua.IsFirefox() || ua.IsSafari() || ua.IsEdge() || ua.IsOpera() ||
		ua.IsAndroid() || ua.IsMacOS() || ua.IsIOS() || ua.IsLinux() ||
		ua.IsChromeOS() || ua.IsBlackBerry() || ua.IsBot() {
		t.Errorf("unexpected helper match for %+v", ua)
	}
	if ua.IsUnknown() {
		t.Error("Chrome UA should not be unknown")
	}
	if !(Device{Type: DeviceMobile}).IsMobile() {
		t.Error("Device.IsMobile() = false")
	}

	bot := Parse("Mozilla/5.0 (compatible; Googlebot/2.1; +http://www.google.com/bot.html)")
	if !bot.IsGooglebot() || !bot.IsBot() {
		t.Errorf("Googlebot helpers failed for %+v", bot)
	}
	if !Parse("Twitterbot/1.0").IsTwitterbot() {
		t.Error("IsTwitterbot() = false")
	}
	if !Parse("facebookexternalhit/1.1").IsFacebookbot() {
		t.Error("IsFacebookbot() = false")
	}
	if !Parse("Mozilla/5.0 (compatible; YandexBot/3.0; +http://yandex.com/bots)").IsYandexbot() {
		t.Error("IsYandexbot() = false")
	}

	if !Parse("Mozilla/5.0 (iPhone; CPU iPhone OS 16_0 like Mac OS X) AppleWebKit/605.1.15 Version/16.0 Mobile/15E148 Safari/604.1").IsIOS() {
		t.Error("IsIOS() = false")
	}
	if !Parse("Mozilla/5.0 (Linux; Android 13; Pixel 7) AppleWebKit/537.36 Chrome/120.0.0.0 Mobile Safari/537.36").IsAndroid() {
		t.Error("IsAndroid() = false")
	}
	if !Parse("Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/605.1.15 Version/17.0 Safari/605.1.15").IsMacOS() {
		t.Error("IsMacOS() = false")
	}
	if !Parse("Mozilla/5.0 (X11; CrOS x86_64 14541.0.0) AppleWebKit/537.36 Chrome/120.0.0.0 Safari/537.36").IsChromeOS() {
		t.Error("IsChromeOS() = false")
	}
	if !Parse("Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 Chrome/120.0.0.0 Safari/537.36").IsLinux() {
		t.Error("IsLinux() = false")
	}
	if !Parse("Mozilla/5.0 (BlackBerry; U; BlackBerry 9900; en-US) AppleWebKit/534.11+ Version/7.0.0.187 Mobile Safari/534.11+").IsBlackBerry() {
		t.Error("IsBlackBerry() = false")
	}
	if !Parse("Mozilla/5.0 (iPad; CPU OS 16_0 like Mac OS X) AppleWebKit/605.1.15 Version/16.0 Mobile/15E148 Safari/604.1").IsTablet() {
		t.Error("IsTablet() = false")
	}
	if !Parse("Mozilla/5.0 (SMART-TV; Linux; Tizen 6.5) AppleWebKit/537.36 TV Safari/537.36").IsTV() {
		t.Error("IsTV() = false")
	}
	if !Parse("Mozilla/5.0 (PlayStation 4 3.11) AppleWebKit/537.73 (KHTML, like Gecko)").IsConsole() {
		t.Error("IsConsole() = false")
	}
	if !Parse("Mozilla/5.0 (Linux; Android 10; SM-G960F; wv) AppleWebKit/537.36 Version/4.0 Chrome/80.0.0.0 Mobile Safari/537.36").IsWebView() {
		t.Error("IsWebView() = false")
	}
}

func TestZeroUserAgent(t *testing.T) {
	var ua UserAgent
	if !ua.IsUnknown() {
		t.Error("zero UserAgent should be unknown")
	}
	if ua.String() != "unknown" {
		t.Errorf("zero String() = %q", ua.String())
	}
	if ua.Browser.Is("") {
		t.Error("empty browser should not match empty name")
	}
	if (OS{}).Is("") {
		t.Error("empty OS should not match empty name")
	}
	if (Browser{Name: Chrome}).Is(Chrome) == false {
		t.Error("Browser.Is failed")
	}
	if (OS{Name: Android}).Is(Android) == false {
		t.Error("OS.Is failed")
	}
}
