package useragent

import "testing"

func TestBrowserTaxonomy(t *testing.T) {
	tests := []struct {
		name    string
		raw     string
		browser string
		version string
	}{
		{"qq", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 Chrome/116.0.0.0 Safari/537.36 QQBrowser/13.5.6297.400", QQBrowser, "13.5.6297.400"},
		{"mqq", "Mozilla/5.0 (Linux; Android 13) AppleWebKit/537.36 Chrome/100.0 Mobile Safari/537.36 MQQBrowser/6.2", QQBrowser, "6.2"},
		{"baidu", "Mozilla/5.0 (Windows NT 6.1; WOW64) AppleWebKit/537.36 Chrome/50.0 Safari/537.36 BIDUBrowser/7.6", BaiduBrowser, "7.6"},
		{"sogou", "Mozilla/5.0 (Windows NT 10.0) AppleWebKit/537.36 Chrome/80.0 Safari/537.36 MetaSr/1.0", SogouBrowser, "1.0"},
		{"quark", "Mozilla/5.0 (Linux; U; Android 12) AppleWebKit/537.36 Chrome/99.0 Quark/5.3.0.180 Mobile Safari/537.36", QuarkBrowser, "5.3.0.180"},
		{"360", "Mozilla/5.0 (Windows NT 10.0) AppleWebKit/537.36 Chrome/86.0 Safari/537.36 QihooBrowser/1.2", Browser360, "1.2"},
		{"yandex", "Mozilla/5.0 (Windows NT 10.0) AppleWebKit/537.36 Chrome/110.0 Safari/537.36 YaBrowser/23.1.0.0", YandexBrowser, "23.1.0.0"},
		{"whale", "Mozilla/5.0 (Windows NT 10.0) AppleWebKit/537.36 Chrome/110.0 Safari/537.36 Whale/3.20.0.0", Whale, "3.20.0.0"},
		{"duckduckgo", "Mozilla/5.0 (iPhone; CPU iPhone OS 16_0 like Mac OS X) AppleWebKit/605.1.15 Version/16.0 DuckDuckGo/7.0 Mobile/15E148 Safari/604.1", DuckDuckGo, "7.0"},
		{"arc", "Mozilla/5.0 (Windows NT 10.0) AppleWebKit/537.36 Chrome/120.0 Safari/537.36 Arc/1.0", Arc, "1.0"},
		{"opera", "Opera/9.80 (Windows NT 6.1; WOW64) Presto/2.12.388 Version/12.16", Opera, "9.80"},
		{"opera-opr", "Mozilla/5.0 (Windows NT 10.0) AppleWebKit/537.36 Chrome/105.0 Safari/537.36 OPR/105.0.0.0", Opera, "105.0.0.0"},
		{"opera-opios", "Mozilla/5.0 (iPhone; CPU iPhone OS 9_3 like Mac OS X) AppleWebKit/601.1.46 (KHTML, like Gecko) OPiOS/14.0.0.104835 Mobile/13E233 Safari/9537.53", Opera, "14.0.0.104835"},
		{"waterfox", "Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:109.0) Gecko/20100101 Waterfox/56.0", "Waterfox", "56.0"},
		{"palemoon", "Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:102.0) Gecko/20100101 PaleMoon/28.0", "Pale Moon", "28.0"},
		{"basilisk", "Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:102.0) Gecko/20100101 Basilisk/2023.01.0", "Basilisk", "2023.01.0"},
		{"focus", "Mozilla/5.0 (Android 13; Mobile; rv:120.0) Gecko/120.0 Firefox/120.0 Focus/120.0", Firefox, "120.0"},
		{"klar", "Mozilla/5.0 (Android 13; Mobile; rv:120.0) Gecko/120.0 Firefox/120.0 Klar/120.0", Firefox, "120.0"},
		{"chromium", "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 Chromium/120.0.0.0 Safari/537.36", Chromium, "120.0.0.0"},
		{"trident", "Mozilla/5.0 (Windows NT 6.1; Trident/7.0)", InternetExplorer, "11.0"},
		{"electron", "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 Chrome/120.0.0.0 Electron/28.0.0 Safari/537.36", "Electron", "28.0.0"},
		{"cef", "Mozilla/5.0 (Windows NT 10.0) AppleWebKit/537.36 Chrome/120.0.0.0 CEF/120.0 Safari/537.36", "CEF", "120.0"},
		{"lynx", "Lynx/2.9.0dev.12 libwww-FM/2.14", "Lynx", "2.9.0dev.12"},
		{"w3m", "w3m/0.5.3+git20230121", "w3m", "0.5.3+git20230121"},
		{"links", "Links (2.29; Linux 6.1 x86_64)", "Links", ""},
		{"webkit", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/605.1.15 (KHTML, like Gecko)", AppleWebKit, "605.1.15"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Parse(tt.raw)
			if got.Browser.Name != tt.browser || got.Browser.Version != tt.version {
				t.Errorf("browser = %q %q, want %q %q", got.Browser.Name, got.Browser.Version, tt.browser, tt.version)
			}
		})
	}
}

func TestOSTaxonomy(t *testing.T) {
	tests := []struct {
		name    string
		raw     string
		os      string
		version string
	}{
		{"ipod", "Mozilla/5.0 (iPod touch; CPU iPhone OS 12_0 like Mac OS X) AppleWebKit/605.1.15 Version/12.0 Mobile/15E148 Safari/604.1", IOS, "12.0"},
		{"ipad", "Mozilla/5.0 (iPad; CPU OS 17_2 like Mac OS X) AppleWebKit/605.1.15 Version/17.2 Mobile/15E148 Safari/604.1", IOS, "17.2"},
		{"ios-cfnetwork", "CFNetwork/1404.0.5 Darwin/22.3.0 (iPhone iOS/16.3.1)", IOS, "16.3.1"},
		{"windows-phone", "Mozilla/5.0 (Windows Phone 10.0; Android 4.2.1; Lumia 950) AppleWebKit/537.36 Edge/15.0", WindowsPhone, "10.0"},
		{"harmony", "Mozilla/5.0 (Phone; OpenHarmony 5.0) AppleWebKit/537.36 Chrome/114.0.0.0 Safari/537.36 ArkWeb/4.1.6.1 Mobile", HarmonyOS, "5.0"},
		{"kaios", "Mozilla/5.0 (Mobile; rv:48.0) Gecko/48.0 Firefox/48.0 KAIOS/2.5", KaiOS, "2.5"},
		{"blackberry", "Mozilla/5.0 (BlackBerry; U; BlackBerry 9900; en-US) AppleWebKit/534.11+ Version/7.0.0.187 Mobile Safari/534.11+", BlackBerry, "7.0.0.187"},
		{"android-comment", "Mozilla/5.0 (7334MB; 1080x2166; Xiaomi M2101K9G; 13) AppleWebKit/537.36 (KHTML, like Gecko) QXSEC Android App 1.435.390 Phone Hybrid() GooglePlayStore QxsecApp/1.435.390", Android, "13"},
		{"chromeos", "Mozilla/5.0 (X11; CrOS armv7l 14388.62.1) AppleWebKit/537.36 Chrome/94.0.4606.114 Safari/537.36", ChromeOS, "14388.62.1"},
		{"tizen", "Mozilla/5.0 (SMART-TV; LINUX; Tizen 6.5) AppleWebKit/537.36 85.0.4183.93/6.5 TV Safari/537.36", Tizen, "6.5"},
		{"webos", "Mozilla/5.0 (Web0S; Linux/SmartTV) AppleWebKit/537.36 Chrome/38.0 Safari/537.36 webOS/6.0", WebOS, "6.0"},
		{"mac", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/605.1.15 Version/17.0 Safari/605.1.15", MacOS, "10.15.7"},
		{"windows", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 Chrome/120.0 Safari/537.36", Windows, "10.0"},
		{"freebsd", "Mozilla/5.0 (X11; FreeBSD 13.2; rv:115.0) Gecko/20100101 Firefox/115.0", FreeBSD, "13.2"},
		{"openbsd", "Mozilla/5.0 (X11; OpenBSD 7.4; rv:115.0) Gecko/20100101 Firefox/115.0", OpenBSD, "7.4"},
		{"netbsd", "Mozilla/5.0 (X11; NetBSD 10.0; rv:115.0) Gecko/20100101 Firefox/115.0", NetBSD, "10.0"},
		{"solaris", "Mozilla/5.0 (X11; Solaris 11.4; rv:52.0) Gecko/20100101 Firefox/52.0", Solaris, "11.4"},
		{"linux", "Mozilla/5.0 (X11; Linux 6.1.0 x86_64; rv:109.0) Gecko/20100101 Firefox/115.0", Linux, "6.1.0"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Parse(tt.raw)
			if got.OS.Name != tt.os || got.OS.Version != tt.version {
				t.Errorf("OS = %q %q, want %q %q", got.OS.Name, got.OS.Version, tt.os, tt.version)
			}
		})
	}
}

func TestDeviceTaxonomy(t *testing.T) {
	tests := []struct {
		name   string
		raw    string
		typ    DeviceType
		vendor string
		model  string
	}{
		{"ps5", "Mozilla/5.0 (PlayStation 5 3.20) AppleWebKit/605.1.15 (KHTML, like Gecko)", DeviceConsole, "Sony", "PlayStation 5"},
		{"ps4", "Mozilla/5.0 (PlayStation 4 3.11) AppleWebKit/537.73 (KHTML, like Gecko)", DeviceConsole, "Sony", "PlayStation 4"},
		{"vita", "Mozilla/5.0 (PlayStation Vita 3.61) AppleWebKit/537.73 (KHTML, like Gecko)", DeviceConsole, "Sony", "PlayStation Vita"},
		{"ps3", "Mozilla/5.0 (PlayStation 3 4.81) AppleWebKit/531.22.8 (KHTML, like Gecko)", DeviceConsole, "Sony", "PlayStation 3"},
		{"xbox", "Mozilla/5.0 (Xbox; Xbox One) AppleWebKit/537.73 (KHTML, like Gecko)", DeviceConsole, "Microsoft", "Xbox"},
		{"switch", "Mozilla/5.0 (Nintendo Switch; WifiWebAuthApplet) AppleWebKit/609.4 (KHTML, like Gecko)", DeviceConsole, "Nintendo", "Nintendo Switch"},
		{"wiiu", "Mozilla/5.0 (Nintendo WiiU) AppleWebKit/536.30 (KHTML, like Gecko)", DeviceConsole, "Nintendo", "WiiU"},
		{"wii", "Mozilla/5.0 (Nintendo Wii; U) AppleWebKit/534.34 (KHTML, like Gecko)", DeviceConsole, "Nintendo", "Wii"},
		{"3ds", "Mozilla/5.0 (Nintendo 3DS; U) AppleWebKit/536.30 (KHTML, like Gecko)", DeviceConsole, "Nintendo", "3DS"},
		{"apple-watch", "Mozilla/5.0 (Apple Watch; CPU WatchOS 9_0 like Mac OS X) AppleWebKit/605.1.15", DeviceWearable, "Apple", "Apple Watch"},
		{"gear", "Mozilla/5.0 (Linux; Android 9; SM-R800) AppleWebKit/537.36 Chrome/70.0 Safari/537.36", DeviceWearable, "Samsung", "Gear"},
		{"wear-os", "Mozilla/5.0 (Linux; Android 9; Wear OS) AppleWebKit/537.36 Chrome/70.0 Safari/537.36", DeviceWearable, "Google", ""},
		{"tv-samsung", "Mozilla/5.0 (SMART-TV; LINUX; Tizen 6.5) AppleWebKit/537.36 TV Safari/537.36", DeviceTV, "Samsung", ""},
		{"tv-lg", "Mozilla/5.0 (Web0S; Linux/SmartTV) AppleWebKit/537.36 Chrome/38.0 Safari/537.36", DeviceTV, "LG", "SmartTV"},
		{"tv-sony", "Mozilla/5.0 (X11; Linux; BRAVIA 4K GB ATV3) AppleWebKit/537.36 Chrome/70.0 Safari/537.36", DeviceTV, "Sony", "BRAVIA 4K GB ATV3"},
		{"tv-roku", "Roku/DVP-9.10 (519.10E04123A)", DeviceTV, "Roku", "Roku"},
		{"tv-amazon", "Mozilla/5.0 (Linux; Android 9; AFTMM) AppleWebKit/537.36 Chrome/70.0 Safari/537.36 Fire TV", DeviceTV, "Amazon", "AFTMM"},
		{"tv-apple", "Mozilla/5.0 (AppleTV) AppleWebKit/537.36 Apple TV", DeviceTV, "Apple", "Apple TV"},
		{"tv-google", "Mozilla/5.0 (X11; Linux x86_64; Chromecast) AppleWebKit/537.36 Chrome/70.0 Safari/537.36", DeviceTV, "Google", "Chromecast"},
		{"tv-generic", "Mozilla/5.0 (X11; Linux x86_64) HbbTV/1.5.1", DeviceTV, "", ""},
		{"android-tablet-smx", "Mozilla/5.0 (Linux; Android 13; SM-X700) AppleWebKit/537.36 Chrome/120.0 Safari/537.36", DeviceTablet, "Samsung", "SM-X700"},
		{"android-tablet-nexus7", "Mozilla/5.0 (Linux; Android 4.4; Nexus 7 Build/KOT49H) AppleWebKit/537.36 Chrome/30.0 Safari/537.36", DeviceTablet, "Google", "Nexus 7"},
		{"android-phone", "Mozilla/5.0 (Linux; Android 13; Pixel 8) AppleWebKit/537.36 Chrome/120.0 Mobile Safari/537.36", DeviceMobile, "Google", "Pixel 8"},
		{"ios-default", "CFNetwork/1404.0.5 Darwin/22.3.0 (iOS/16.3.1)", DeviceMobile, "Apple", ""},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Parse(tt.raw)
			if got.Device.Type != tt.typ || got.Device.Vendor != tt.vendor || got.Device.Model != tt.model {
				t.Errorf("device = %+v, want type=%v vendor=%q model=%q", got.Device, tt.typ, tt.vendor, tt.model)
			}
		})
	}
}
