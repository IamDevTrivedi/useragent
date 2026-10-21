package useragent

import (
	"strings"
	"testing"
)

func TestBroadCorpus(t *testing.T) {
	tests := []struct {
		name string
		raw  string
		want want
	}{
		{"safari mac", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_12_6) AppleWebKit/603.3.8 (KHTML, like Gecko) Version/10.1.2 Safari/603.3.8", want{browser: Safari, browserVersion: "10.1.2", os: MacOS, device: DeviceDesktop, checkDevice: true}},
		{"chrome mac", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_12_6) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/60.0.3112.90 Safari/537.36", want{browser: Chrome, os: MacOS, device: DeviceDesktop, checkDevice: true, engine: "Blink"}},
		{"firefox mac", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10.12; rv:54.0) Gecko/20100101 Firefox/54.0", want{browser: Firefox, browserVersion: "54.0", os: MacOS, device: DeviceDesktop, checkDevice: true, engine: "Gecko"}},
		{"opera mac", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_12_6) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/59.0.3071.115 Safari/537.36 OPR/46.0.2597.57", want{browser: Opera, browserVersion: "46.0.2597.57", os: MacOS, device: DeviceDesktop, checkDevice: true, engine: "Blink"}},
		{"vivaldi mac", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_12_6) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/60.0.3112.91 Safari/537.36 Vivaldi/1.92.917.39", want{browser: Vivaldi, browserVersion: "1.92.917.39", os: MacOS, device: DeviceDesktop, checkDevice: true}},
		{"edge chromium mac", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_12_6) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/79.0.3945.130 Safari/537.36 Edg/79.0.309.71", want{browser: Edge, browserVersion: "79.0.309.71", os: MacOS, device: DeviceDesktop, checkDevice: true, engine: "Blink"}},
		{"chrome windows wow64", "Mozilla/5.0 (Windows NT 6.1; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/59.0.3071.115 Safari/537.36", want{browser: Chrome, os: Windows, device: DeviceDesktop, checkDevice: true, cpu: "386"}},
		{"ie8", "Mozilla/4.0 (compatible; MSIE 8.0; Windows NT 6.1; WOW64; Trident/4.0; SLCC2; .NET CLR 2.0.50727; .NET CLR 3.5.30729; .NET CLR 3.0.30729; Media Center PC 6.0; .NET4.0C; .NET4.0E; InfoPath.2; GWX:RED)", want{browser: InternetExplorer, browserVersion: "8.0", os: Windows, device: DeviceDesktop, checkDevice: true, engine: "Trident"}},
		{"ie6", "Mozilla/4.0 (compatible; MSIE 6.0; Windows NT 5.1; SV1; .NET CLR 1.1.4322) NS8/0.9.6", want{browser: InternetExplorer, browserVersion: "6.0", os: Windows, device: DeviceDesktop, checkDevice: true}},
		{"edge legacy", "Mozilla/5.0 (Windows NT 10.0) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/52.0.2743.116 Safari/537.36 Edge/15.15063", want{browser: Edge, browserVersion: "15.15063", os: Windows, device: DeviceDesktop, checkDevice: true, engine: "EdgeHTML"}},
		{"iphone safari", "Mozilla/5.0 (iPhone; CPU iPhone OS 10_3_2 like Mac OS X) AppleWebKit/603.2.4 (KHTML, like Gecko) Version/10.0 Mobile/14F89 Safari/602.1", want{browser: Safari, browserVersion: "10.0", os: IOS, osVersion: "10.3.2", device: DeviceMobile, checkDevice: true, vendor: "Apple", model: "iPhone"}},
		{"iphone chrome", "Mozilla/5.0 (iPhone; CPU iPhone OS 10_3_2 like Mac OS X) AppleWebKit/603.1.30 (KHTML, like Gecko) CriOS/60.0.3112.89 Mobile/14F89 Safari/602.1", want{browser: Chrome, browserVersion: "60.0.3112.89", os: IOS, device: DeviceMobile, checkDevice: true}},
		{"iphone opera", "Mozilla/5.0 (iPhone; CPU iPhone OS 9_3 like Mac OS X) AppleWebKit/601.1.46 (KHTML, like Gecko) OPiOS/14.0.0.104835 Mobile/13E233 Safari/9537.53", want{browser: Opera, browserVersion: "14.0.0.104835", os: IOS, device: DeviceMobile, checkDevice: true}},
		{"iphone firefox", "Mozilla/5.0 (iPhone; CPU iPhone OS 10_3_2 like Mac OS X) AppleWebKit/603.2.4 (KHTML, like Gecko) FxiOS/8.1.1b4948 Mobile/14F89 Safari/603.2.4", want{browser: Firefox, browserVersion: "8.1.1b4948", os: IOS, device: DeviceMobile, checkDevice: true}},
		{"iphone edge", "Mozilla/5.0 (iPhone; CPU iPhone OS 13_3 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/13.0 EdgiOS/44.11.15 Mobile/15E148 Safari/605.1.15", want{browser: Edge, browserVersion: "44.11.15", os: IOS, device: DeviceMobile, checkDevice: true}},
		{"ipad safari", "Mozilla/5.0 (iPad; CPU OS 10_3_2 like Mac OS X) AppleWebKit/603.2.4 (KHTML, like Gecko) Version/10.0 Mobile/14F89 Safari/602.1", want{browser: Safari, os: IOS, device: DeviceTablet, checkDevice: true, vendor: "Apple", model: "iPad"}},
		{"ipad chrome", "Mozilla/5.0 (iPad; CPU OS 10_3_2 like Mac OS X) AppleWebKit/602.1.50 (KHTML, like Gecko) CriOS/58.0.3029.113 Mobile/14F89 Safari/602.1", want{browser: Chrome, os: IOS, device: DeviceTablet, checkDevice: true}},
		{"android tablet firefox", "Mozilla/5.0 (Android 4.4; Tablet; rv:41.0) Gecko/41.0 Firefox/41.0", want{browser: Firefox, os: Android, device: DeviceTablet, checkDevice: true}},
		{"android chrome tablet", "Mozilla/5.0 (Linux; Android 9; Chrome tablet) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/110.0.0.0 Mobile Safari/537.36", want{browser: Chrome, os: Android, device: DeviceTablet, checkDevice: true, model: "Chrome tablet"}},
		{"android galaxy", "Mozilla/5.0 (Linux; Android 4.3; GT-I9300 Build/JSS15J) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/59.0.3071.125 Mobile Safari/537.36", want{browser: Chrome, os: Android, device: DeviceMobile, checkDevice: true, vendor: "Samsung", model: "GT-I9300"}},
		{"android firefox", "Mozilla/5.0 (Android 4.3; Mobile; rv:54.0) Gecko/54.0 Firefox/54.0", want{browser: Firefox, os: Android, device: DeviceMobile, checkDevice: true}},
		{"android opera", "Mozilla/5.0 (Linux; Android 4.3; GT-I9300 Build/JSS15J) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/55.0.2883.91 Mobile Safari/537.36 OPR/42.9.2246.119956", want{browser: Opera, os: Android, device: DeviceMobile, checkDevice: true}},
		{"opera mini", "Opera/9.80 (Android; Opera Mini/28.0.2254/66.318; U; en) Presto/2.12.423 Version/12.16", want{browser: OperaMini, browserVersion: "28.0.2254/66.318", os: Android, device: DeviceMobile, checkDevice: true, engine: "Presto"}},
		{"android browser", "Mozilla/5.0 (Linux; U; Android 4.3; en-us; GT-I9300 Build/JSS15J) AppleWebKit/534.30 (KHTML, like Gecko) Version/4.0 Mobile Safari/534.30", want{browser: AndroidBrowser, browserVersion: "4.0", os: Android, device: DeviceMobile, checkDevice: true}},
		{"android edge", "Mozilla/5.0 (Linux; Android 10; ONEPLUS A6003) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/73.0.3683.0 Mobile Safari/537.36 EdgA/44.11.4.4140", want{browser: Edge, browserVersion: "44.11.4.4140", os: Android, device: DeviceMobile, checkDevice: true}},
		{"samsung internet", "Mozilla/5.0 (Linux; Android 6.0.1; SAMSUNG SM-A310F/A310FXXU2BQB1 Build/MMB29K) AppleWebKit/537.36 (KHTML, like Gecko) SamsungBrowser/5.4 Chrome/51.0.2704.106 Mobile Safari/537.36", want{browser: SamsungBrowser, browserVersion: "5.4", os: Android, device: DeviceMobile, checkDevice: true, vendor: "Samsung", model: "SM-A310F"}},
		{"lg android", "Mozilla/5.0 (Linux; Android 9; LM-Q630) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/86.0.4240.198 Mobile Safari/537.36", want{browser: Chrome, os: Android, device: DeviceMobile, checkDevice: true, vendor: "LG", model: "LM-Q630"}},
		{"miui x11", "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/534.24 (KHTML, like Gecko) Chrome/79.0.3945.147 Safari/534.24 XiaoMi/MiuiBrowser/12.11.5-gn", want{browser: MiuiBrowser, browserVersion: "12.11.5-gn", os: Android, device: DeviceMobile, checkDevice: true}},
		{"miui redmi", "Mozilla/5.0 (Linux; U; Android 11; ru-ru; Redmi Note 10S Build/RP1A.200720.011) AppleWebKit/537.36 (KHTML, like Gecko) Version/4.0 Chrome/89.0.4389.116 Mobile Safari/537.36 XiaoMi/MiuiBrowser/12.13.2-gn", want{browser: MiuiBrowser, browserVersion: "12.13.2-gn", os: Android, osVersion: "11", device: DeviceMobile, checkDevice: true, vendor: "Xiaomi", model: "Redmi Note 10S"}},
		{"huawei browser", "Mozilla/5.0 (Linux; Android 10; MED-LX9N; HMSCore 6.6.0.311) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/92.0.4515.105 HuaweiBrowser/12.1.0.303 Mobile Safari/537.36", want{browser: HuaweiBrowser, browserVersion: "12.1.0.303", os: Android, device: DeviceMobile, checkDevice: true, model: "MED-LX9N"}},
		{"samsung x11", "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) SamsungBrowser/22.0 Chrome/111.0.5563.116 Safari/537.36", want{browser: SamsungBrowser, browserVersion: "22.0", os: Android, device: DeviceMobile, checkDevice: true}},
		{"oneplus", "Mozilla/5.0 (Linux; Android 9; ONEPLUS A6003 Build/PKQ1.180716.001) AppleWebKit/537.36 (KHTML, like Gecko) Version/4.0 Chrome/71.0.3578.99 Mobile Safari/537.36 OPT/1.14.51", want{browser: OperaTouch, browserVersion: "1.14.51", os: Android, device: DeviceMobile, checkDevice: true, vendor: "OnePlus", model: "A6003"}},
		{"moto", "Mozilla/5.0 (Linux; Android 7.0; Moto G (4)) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/84.0.4143.7 Mobile Safari/537.36 Chrome-Lighthouse", want{browser: Chrome, os: Android, device: DeviceMobile, checkDevice: true, vendor: "Motorola"}},
		{"windows phone", "Mozilla/4.0 (compatible; MSIE 7.0; Windows Phone OS 7.0; Trident/3.1; IEMobile/7.0; NOKIA; Lumia 630)", want{browser: InternetExplorer, browserVersion: "7.0", os: WindowsPhone, osVersion: "7.0", device: DeviceMobile, checkDevice: true}},
		{"konqueror freebsd", "Mozilla/5.0 (compatible; Konqueror/4.5; FreeBSD) KHTML/4.5.4 (like Gecko)", want{browser: Konqueror, browserVersion: "4.5", os: FreeBSD, device: DeviceDesktop, checkDevice: true}},
		{"googlebot mobile", "Mozilla/5.0 (Linux; Android 6.0.1; Nexus 5X Build/MMB29P) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/41.0.2272.96 Mobile Safari/537.36 (compatible; Googlebot/2.1; +http://www.google.com/bot.html)", want{browser: Googlebot, browserVersion: "2.1", os: Android, device: DeviceMobile, checkDevice: true, vendor: "Google", model: "Nexus 5X", botDetected: boolPtr(true), botName: Googlebot, botVersion: "2.1", botURL: "http://www.google.com/bot.html"}},
		{"googlebot desktop", "Mozilla/5.0 (compatible; Googlebot/2.1; +http://www.google.com/bot.html)", want{browser: Googlebot, browserVersion: "2.1", botDetected: boolPtr(true), botName: Googlebot, botVersion: "2.1"}},
		{"applebot", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_5) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/13.1.1 Safari/605.1.15 (Applebot/0.1; +http://www.apple.com/go/applebot)", want{browser: Applebot, browserVersion: "0.1", botDetected: boolPtr(true), botName: Applebot, botVersion: "0.1"}},
		{"twitterbot", "Twitterbot/1.0", want{browser: Twitterbot, browserVersion: "1.0", botDetected: boolPtr(true), botName: Twitterbot, botVersion: "1.0", botCategory: catPtr(BotSocial)}},
		{"facebookexternalhit", "facebookexternalhit/1.1", want{browser: FacebookExternal, browserVersion: "1.1", botDetected: boolPtr(true), botName: FacebookExternal, botVersion: "1.1"}},
		{"facebookcatalog", "facebookcatalog/1.0", want{browser: "facebookcatalog", browserVersion: "1.0", botDetected: boolPtr(true), botName: "facebookcatalog", botVersion: "1.0"}},
		{"semrush", "Mozilla/5.0 (compatible; SemrushBot/7~bl; +http://www.semrush.com/bot.html", want{browser: SemrushBot, browserVersion: "7~bl", botDetected: boolPtr(true), botName: SemrushBot, botVersion: "7~bl"}},
		{"yandex", "Mozilla/5.0 (compatible; YandexBot/3.0; +http://yandex.com/bots) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/81.0.4044.268", want{browser: YandexBot, browserVersion: "3.0", botDetected: boolPtr(true), botName: YandexBot}},
		{"discord", "Mozilla/5.0 (compatible; Discordbot/2.0; +https://discordapp.com)", want{browser: Discordbot, browserVersion: "2.0", botDetected: boolPtr(true), botName: Discordbot}},
		{"bing", "Mozilla/5.0 (compatible; bingbot/2.0; +http://www.bing.com/bingbot.htm)", want{browser: Bingbot, browserVersion: "2.0", botDetected: boolPtr(true), botName: Bingbot}},
		{"yahoo ads", "Mozilla/5.0 (compatible; Yahoo Ad monitoring; https://help.yahoo.com/kb/yahoo-ad-monitoring-SLN24857.html) cnv.aws-prod---sieve.hlfs-rest_client/1681346790-0", want{browser: YahooAdMonitor, botDetected: boolPtr(true), botName: YahooAdMonitor, botCategory: catPtr(BotAds)}},
		{"bytespider", "Mozilla/5.0 (compatible; Bytespider; spider-feedback@bytedance.com) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/70.0.0.0 Safari/537.36", want{browser: Bytespider, botDetected: boolPtr(true), botName: Bytespider}},
		{"mediapartners", "Mozilla/5.0 (Linux; Android 4.0.0; Galaxy Nexus Build/IMM76B) AppleWebKit/537.36 (KHTML, like Gecko; Mediapartners-Google) Chrome/104.0.0.0 Mobile Safari/537.36", want{browser: "Mediapartners-Google", os: Android, botDetected: boolPtr(true), botName: "Mediapartners-Google", botCategory: catPtr(BotAds)}},
		{"adsbot", "Mozilla/5.0 (Linux; Android 5.0; SM-G920A) AppleWebKit (KHTML, like Gecko) Chrome Mobile Safari (compatible; AdsBot-Google-Mobile; +http://www.google.com/mobile/adsbot.html)", want{browser: AdsBotGoogle, os: Android, botDetected: boolPtr(true), botName: AdsBotGoogle, botCategory: catPtr(BotAds)}},
		{"brave", "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Brave Chrome/87.0.4280.101 Safari/537.36", want{browser: Brave, browserVersion: "87.0.4280.101", os: Linux, device: DeviceDesktop, checkDevice: true, engine: "Blink"}},
		{"headless", "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) HeadlessChrome/98.0.4758.0 Safari/537.36", want{browser: HeadlessChrome, browserVersion: "98.0.4758.0", os: Linux, device: DeviceDesktop, checkDevice: true, botDetected: boolPtr(true), botName: HeadlessChrome}},
		{"facebook ios app", "Mozilla/5.0 (iPhone; CPU iPhone OS 15_4_1 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Mobile/19E258 [FBAN/FBIOS;FBDV/iPhone8,2;FBMD/iPhone;FBSN/iOS;FBSV/15.4.1;FBSS/3;FBID/phone;FBLC/fr_FR;FBOP/5]", want{os: IOS, device: DeviceMobile, checkDevice: true, app: "Facebook", appVersion: "FBIOS", inApp: boolPtr(true)}},
		{"facebook android app", "Mozilla/5.0 (Linux; Android 13; SM-T220 Build/TP1A.220624.014; wv) AppleWebKit/537.36 (KHTML, like Gecko) Version/4.0 Chrome/109.0.5414.117 Safari/537.36 [FB_IAB/FB4A;FBAV/400.0.0.37.76;]", want{os: Android, device: DeviceTablet, checkDevice: true, app: "Facebook", appVersion: "400.0.0.37.76", inApp: boolPtr(true), webView: boolPtr(true)}},
		{"instagram", "Mozilla/5.0 (iPhone; CPU iPhone OS 16_3 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Mobile/15E148 Instagram 270.0.0.13.83 (iPhone13,2; iOS 16_3; es_ES; es-ES; scale=3.00; 1170x2532; 445843881) NW/1", want{os: IOS, device: DeviceMobile, checkDevice: true, app: "Instagram", appVersion: "270.0.0.13.83", inApp: boolPtr(true), webView: boolPtr(true)}},
		{"tiktok ios", "Mozilla/5.0 (iPhone; CPU iPhone OS 15_5 like Mac OS ) AppleWebKit/605.1.15 (KHTML, like Gecko) Mobile/15E148 musical_ly_28.2.0 JsSdk/2.0 NetType/WIFI Channel/App Store ByteLocale/es Region/PE RevealType/Dialog isDarkMode/0 WKWebView/1 BytedanceWebview/d8a21c6 FalconTag/D6EBBF89-6D75-4BBD-9304-BF199C6B4DB1", want{os: IOS, device: DeviceMobile, checkDevice: true, app: "TikTok", appVersion: "28.2.0", inApp: boolPtr(true)}},
		{"tiktok android", "Mozilla/5.0 (Linux; Android 10; AGS3K-W09 Build/HUAWEIAGS3K-W09; wv) AppleWebKit/537.36 (KHTML, like Gecko) Version/4.0 Chrome/88.0.4324.93 Safari/537.36 trill_2022803040 JsSdk/1.0 NetType/WIFI Channel/huaweiadsglobal_int AppName/musical_ly app_version/28.3.4 ByteLocale/es ByteFullLocale/es Region/PE BytedanceWebview/d8a21c6", want{os: Android, device: DeviceMobile, checkDevice: true, app: "TikTok", appVersion: "28.3.4", inApp: boolPtr(true)}},
		{"qtwebengine", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_11_4) AppleWebKit/537.36 (KHTML, like Gecko) QtWebEngine/5.6.0 Chrome/45.0.2454.101 Safari/537.36", want{browser: "QtWebEngine", browserVersion: "5.6.0", os: MacOS, device: DeviceDesktop, checkDevice: true}},
		{"go http", "Go-http-client/1.1", want{browser: GoHTTPClient, browserVersion: "1.1", botDetected: boolPtr(true), botName: GoHTTPClient}},
		{"wget", "Wget/1.12 (linux-gnu)", want{browser: Wget, browserVersion: "1.12", botDetected: boolPtr(true), botName: Wget}},
		{"seafile", "Seafile/9.0.2 (Linux)", want{browser: "Seafile", browserVersion: "9.0.2", os: Linux, app: "Seafile", appVersion: "9.0.2", inApp: boolPtr(true)}},
		{"blackberry", "Mozilla/5.0 (BlackBerry; U; BlackBerry 9900; en-US) AppleWebKit/534.11+ (KHTML, like Gecko) Version/7.0.0.187 Mobile Safari/534.11+", want{browser: BlackBerry, browserVersion: "7.0.0.187", os: BlackBerry, device: DeviceMobile, checkDevice: true}},
		{"arkweb", "Mozilla/5.0 (Phone; OpenHarmony 5.0) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/114.0.0.0 Safari/537.36  ArkWeb/4.1.6.1 Mobile", want{browser: "ArkWeb", browserVersion: "4.1.6.1", os: HarmonyOS, device: DeviceMobile, checkDevice: true}},
		{"netfront", "SonyEricssonK310iv/R4DA Browser/NetFront/3.3 Profile/MIDP-2.0 Configuration/CLDC-1.1 UP.Link/6.3.1.13.0", want{browser: NetFront, browserVersion: "3.3", device: DeviceMobile, checkDevice: true}},
		{"numeric model 8092", "Mozilla/5.0 (Linux; Android 10; 8092) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/112.0.0.0 Safari/537.36", want{browser: Chrome, os: Android, device: DeviceMobile, checkDevice: true, model: "8092"}},
		{"vivax tablet", "Mozilla/5.0 (Linux; Android 6.0; VIVAX TABLET TPC-101 3G) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/106.0.0.0 Safari/537.36", want{browser: Chrome, os: Android, device: DeviceTablet, checkDevice: true, vendor: "VIVAX", model: "TABLET TPC-101 3G"}},
		{"lenovo tablet", "Mozilla/5.0 (Linux; Android 8.1.0; Lenovo TB-7104F Build/O11019) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/107.0.5304.91 Safari/537.36", want{browser: Chrome, os: Android, device: DeviceTablet, checkDevice: true, vendor: "Lenovo", model: "TB-7104F"}},
		{"samsung tablet", "Mozilla/5.0 (Linux; Android 4.4.4; SM-T560 Build/KTU84P) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/68.0.3440.91 Safari/537.36", want{browser: Chrome, os: Android, device: DeviceTablet, checkDevice: true, vendor: "Samsung", model: "SM-T560"}},
		{"android no model", "Mozilla/5.0 (Linux; Android 10;)", want{os: Android, device: DeviceMobile, checkDevice: true}},
		{"ps4", "Mozilla/5.0 (PlayStation 4 3.11) AppleWebKit/537.73 (KHTML, like Gecko)", want{device: DeviceConsole, checkDevice: true, vendor: "Sony"}},
		{"xbox", "Mozilla/5.0 (Windows NT 10.0; Win64; x64; Xbox; Xbox One) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/70.0.3538.102 Safari/537.36 Edge/18.19041", want{device: DeviceConsole, checkDevice: true, vendor: "Microsoft"}},
		{"nintendo", "Mozilla/5.0 (Nintendo Switch; WifiWebAuthApplet) AppleWebKit/606.4 (KHTML, like Gecko) NF/6.0.1.15.4 NintendoBrowser/5.1.0.20393", want{device: DeviceConsole, checkDevice: true, vendor: "Nintendo"}},
		{"apple watch", "Mozilla/5.0 (iPhone; CPU iPhone OS 10_3 like Mac OS X) AppleWebKit/602.1.50 (KHTML, like Gecko) Mobile/14E277 WatchOS/3.2", want{device: DeviceWearable, checkDevice: true, vendor: "Apple"}},
		{"wear os", "Mozilla/5.0 (Linux; Android 9; Wear OS) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/74.0.3729.157 Mobile Safari/537.36", want{device: DeviceWearable, checkDevice: true}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ua := Parse(tt.raw)
			check(t, ua, tt.want)
		})
	}
}

func TestBroadCorpusBrowserNotArbitrary(t *testing.T) {
	// A generic "best match" fallback can return a comment or prefix token.
	// This parser must leave the browser empty unless it can name a real client.
	for _, raw := range []string{
		"() { ignored; }; echo Content-Type: text/plain ; echo ; echo \"bash_cve_2014_6271_rce Output : $((70+91))\"",
		"${jndi:ldap://log4shell-generic-8ZnJfq2XFL3GWyaLyOpT${lower:ten}.w.nessus.org/nessus}",
		"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7)",
		"cloud-run-proxy/source",
	} {
		ua := Parse(raw)
		if ua.Browser.Name != "" {
			t.Errorf("%q: Browser.Name = %q, want empty", raw, ua.Browser.Name)
		}
	}
}

func TestUnknownBrowserHelpers(t *testing.T) {
	ua := Parse("")
	if !ua.IsUnknown() {
		t.Error("empty UA should be unknown")
	}
	if ua.IsChrome() || ua.IsFirefox() || ua.IsSafari() || ua.IsEdge() ||
		ua.IsOpera() || ua.IsAndroid() || ua.IsWindows() || ua.IsMacOS() ||
		ua.IsIOS() || ua.IsLinux() || ua.IsChromeOS() || ua.IsBlackBerry() ||
		ua.IsGooglebot() || ua.IsTwitterbot() || ua.IsFacebookbot() || ua.IsYandexbot() {
		t.Error("zero UA matched a helper")
	}
	if !strings.Contains(ua.String(), "unknown") {
		t.Errorf("String() = %q", ua.String())
	}
}
