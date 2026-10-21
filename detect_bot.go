package useragent

type botDef struct {
	name        string
	keys        []string
	contains    []string
	versionKeys []string
	category    BotCategory
}

var botDefs = [...]botDef{
	// Social preview and link unfurlers.
	{name: FacebookExternal, keys: []string{"facebookexternalhit"}, category: BotSocial},
	{name: "Facebot", keys: []string{"Facebot"}, category: BotSocial},
	{name: "facebookcatalog", keys: []string{"facebookcatalog"}, category: BotSocial},
	{name: Twitterbot, keys: []string{"Twitterbot"}, category: BotSocial},
	{name: LinkedInBot, keys: []string{"LinkedInBot"}, category: BotSocial},
	{name: Slackbot, keys: []string{"Slackbot", "Slackbot-LinkExpanding"}, category: BotSocial},
	{name: Discordbot, keys: []string{"Discordbot"}, category: BotSocial},
	{name: TelegramBot, keys: []string{"TelegramBot"}, category: BotSocial},
	{name: "Pinterestbot", keys: []string{"Pinterestbot"}, category: BotSocial},
	{name: "redditbot", keys: []string{"redditbot"}, category: BotSocial},

	// Search and advertising crawlers.
	{name: Googlebot, keys: []string{"Googlebot", "Googlebot-Image", "Googlebot-News", "Googlebot-Video"}, category: BotSearch},
	{name: GoogleOther, keys: []string{"GoogleOther"}, category: BotSearch},
	{name: "Google-InspectionTool", keys: []string{"Google-InspectionTool"}, category: BotSearch},
	{name: Bingbot, keys: []string{"bingbot", "BingPreview"}, category: BotSearch},
	{name: YandexBot, keys: []string{"YandexBot"}, category: BotSearch},
	{name: Baiduspider, keys: []string{"Baiduspider", "Baiduspider-image", "Baiduspider-video"}, category: BotSearch},
	{name: DuckDuckBot, keys: []string{"DuckDuckBot"}, category: BotSearch},
	{name: Applebot, keys: []string{"Applebot"}, category: BotSearch},
	{name: "GoogleProber", keys: []string{"GoogleProber"}, category: BotSearch},
	{name: "GoogleProducer", keys: []string{"GoogleProducer"}, category: BotSearch},
	{name: AdsBotGoogle, keys: []string{"AdsBot-Google", "AdsBot-Google-Mobile", "AdsBot-Google-Mobile-Apps"}, category: BotAds},
	{name: "Mediapartners-Google", keys: []string{"Mediapartners-Google"}, category: BotAds},
	{name: YahooAdMonitor, keys: []string{"Yahoo Ad monitoring"}, category: BotAds},

	// AI crawlers.
	{name: GPTBot, keys: []string{"GPTBot"}, category: BotAICrawler},
	{name: ChatGPTUser, keys: []string{"ChatGPT-User"}, category: BotAICrawler},
	{name: OAI, keys: []string{"OAI-SearchBot"}, category: BotAICrawler},
	{name: ClaudeBot, keys: []string{"ClaudeBot"}, category: BotAICrawler},
	{name: ClaudeWeb, keys: []string{"Claude-Web"}, category: BotAICrawler},
	{name: PerplexityBot, keys: []string{"PerplexityBot"}, category: BotAICrawler},
	{name: CCBot, keys: []string{"CCBot"}, category: BotAICrawler},
	{name: Amazonbot, keys: []string{"Amazonbot"}, category: BotAICrawler},
	{name: Bytespider, keys: []string{"Bytespider"}, category: BotAICrawler},

	// SEO and audit crawlers.
	{name: SemrushBot, keys: []string{"SemrushBot"}, category: BotSEOTool},
	{name: AhrefsBot, keys: []string{"AhrefsBot"}, category: BotSEOTool},
	{name: MJ12bot, keys: []string{"MJ12bot"}, category: BotSEOTool},
	{name: DotBot, keys: []string{"DotBot"}, category: BotSEOTool},

	// Uptime and monitoring.
	{name: UptimeRobot, keys: []string{"UptimeRobot"}, category: BotMonitoring},
	{name: Pingdom, contains: []string{"Pingdom"}, category: BotMonitoring},
	{name: StatusCake, keys: []string{"StatusCake"}, category: BotMonitoring},

	// Security scanners.
	{name: Nuclei, keys: []string{"Nuclei"}, category: BotSecurityScanner},
	{name: Censys, keys: []string{"Censys"}, category: BotSecurityScanner},
	{name: Shodan, keys: []string{"Shodan"}, category: BotSecurityScanner},
	{name: ZGrab, keys: []string{"ZGrab"}, category: BotSecurityScanner},
	{name: Masscan, keys: []string{"masscan"}, category: BotSecurityScanner},
	{name: Nmap, keys: []string{"Nmap"}, category: BotSecurityScanner},
	{name: SQLmap, keys: []string{"sqlmap"}, category: BotSecurityScanner},
	{name: Nikto, keys: []string{"Nikto"}, category: BotSecurityScanner},

	// Scrapers and rendering engines.
	{name: Scrapy, keys: []string{"Scrapy"}, category: BotScraper},
	{name: Splash, keys: []string{"splash", "Splash"}, versionKeys: []string{"Version"}, category: BotScraper},
	{name: HTTrack, keys: []string{"HTTrack"}, category: BotScraper},
	{name: "BUbiNG", contains: []string{"BUbiNG"}, category: BotScraper},
	{name: "Google", contains: []string{"Google (+http"}, category: BotOther},

	// Headless browsers and automation frameworks.
	{name: HeadlessChrome, keys: []string{"HeadlessChrome"}, category: BotHeadless},
	{name: PhantomJS, keys: []string{"PhantomJS"}, category: BotHeadless},
	{name: Selenium, keys: []string{"Selenium"}, category: BotHeadless},
	{name: Playwright, keys: []string{"Playwright"}, category: BotHeadless},
	{name: Puppeteer, keys: []string{"Puppeteer"}, category: BotHeadless},

	// HTTP clients and libraries.
	{name: Curl, keys: []string{"curl"}, category: BotHTTPClient},
	{name: Wget, keys: []string{"Wget", "wget"}, category: BotHTTPClient},
	{name: GoHTTPClient, keys: []string{"Go-http-client"}, category: BotHTTPClient},
	{name: PythonRequests, keys: []string{"python-requests"}, category: BotHTTPClient},
	{name: OkHTTP, keys: []string{"okhttp"}, category: BotHTTPClient},
	{name: JavaHTTP, keys: []string{"Java"}, category: BotHTTPClient},
	{name: Axios, keys: []string{"axios"}, category: BotHTTPClient},
	{name: NodeFetch, keys: []string{"node-fetch"}, category: BotHTTPClient},
	{name: Undici, keys: []string{"undici"}, category: BotHTTPClient},
	{name: Postman, keys: []string{"PostmanRuntime"}, category: BotHTTPClient},
	{name: Insomnia, keys: []string{"Insomnia"}, category: BotHTTPClient},
	{name: HTTPie, keys: []string{"HTTPie"}, category: BotHTTPClient},
	{name: "libwww-perl", keys: []string{"libwww-perl"}, category: BotHTTPClient},
	{name: "Apache-HttpClient", keys: []string{"Apache-HttpClient"}, category: BotHTTPClient},
	{name: "GuzzleHttp", keys: []string{"GuzzleHttp"}, category: BotHTTPClient},
	{name: "RestSharp", keys: []string{"RestSharp"}, category: BotHTTPClient},
	{name: "WinHTTP", keys: []string{"WinHTTP"}, category: BotHTTPClient},
}

type botKeyEntry struct {
	key   string
	order int
}

var botKeyBuckets [256][]botKeyEntry
var botContainsDefs []int

func init() {
	for i := range botDefs {
		d := &botDefs[i]
		for _, k := range d.keys {
			if k != "" {
				c := lowerASCII(k[0])
				botKeyBuckets[c] = append(botKeyBuckets[c], botKeyEntry{key: k, order: i})
			}
		}
		if len(d.contains) != 0 {
			botContainsDefs = append(botContainsDefs, i)
		}
	}
}

func botContains(ts *tokenSet, marker string) bool {
	switch marker {
	case "Pingdom", "BUbiNG":
		return ts.has(marker) || ts.hasPrefix(marker)
	case "Google (+http":
		return ts.has("Google") && ts.contains(marker)
	default:
		return ts.contains(marker)
	}
}

func detectBot(ua *UserAgent, ts *tokenSet) {
	if ua.Bot.Detected {
		return
	}

	best := -1
	for i := 0; i < ts.n; i++ {
		key := ts.list[i].key
		if key == "" {
			continue
		}
		for _, e := range botKeyBuckets[lowerASCII(key[0])] {
			if (best < 0 || e.order < best) && equalFold(key, e.key) {
				best = e.order
			}
		}
	}
	for _, idx := range botContainsDefs {
		if best >= 0 && idx > best {
			continue
		}
		for _, marker := range botDefs[idx].contains {
			if botContains(ts, marker) {
				if best < 0 || idx < best {
					best = idx
				}
				break
			}
		}
	}
	if best < 0 {
		return
	}

	d := &botDefs[best]
	version := tokenValue(ts, d.versionKeys...)
	if version == "" && len(d.versionKeys) == 0 {
		version = tokenValue(ts, d.keys...)
	}
	setBot(ua, d.name, version, ts.url, d.category)
}
