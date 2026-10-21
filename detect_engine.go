package useragent

func detectEngine(ua *UserAgent, ts *tokenSet) {
	if ua.Engine.Name != "" {
		return
	}
	raw := ua.Raw

	chromeVersion := firstNonEmpty(
		tokenValue(ts, "Chrome", "Chromium", "CriOS", "HeadlessChrome"),
		productVersion(raw, "Chrome"),
		productVersion(raw, "Chromium"),
	)

	if ua.Browser.Name == Edge && ts.has("Edge") && !ts.has("Edg") {
		setEngine(ua, "EdgeHTML", firstNonEmpty(tokenValue(ts, "Edge"), productVersion(raw, "Edge")))
		return
	}

	if ts.has("Presto") {
		setEngine(ua, "Presto", firstNonEmpty(tokenValue(ts, "Presto"), tokenValue(ts, "Opera")))
		return
	}

	if isBlinkBrowser(ua.Browser.Name) || ts.has("HeadlessChrome") {
		setEngine(ua, "Blink", chromeVersion)
		return
	}

	if ts.has("Trident") || ts.has("MSIE") {
		setEngine(ua, "Trident", firstNonEmpty(tokenValue(ts, "Trident"), tokenValue(ts, "MSIE")))
		return
	}

	if ts.hasAny("Firefox", "FxiOS", "Waterfox", "PaleMoon", "Basilisk") || ts.has("Gecko") {
		version := firstNonEmpty(
			tokenValue(ts, "Firefox", "FxiOS", "Waterfox", "PaleMoon", "Basilisk"),
			tokenValue(ts, "Gecko"),
			productVersion(raw, "Firefox"),
		)
		setEngine(ua, "Gecko", version)
		return
	}

	if ts.has("AppleWebKit") {
		setEngine(ua, "WebKit", firstNonEmpty(tokenValue(ts, "AppleWebKit"), productVersion(raw, "AppleWebKit")))
		return
	}
}

func isBlinkBrowser(name string) bool {
	switch name {
	case Chrome, Chromium, Edge, Opera, OperaTouch, OperaMini, Brave, Vivaldi,
		Arc, CocCoc, YandexBrowser, Whale, DuckDuckGo, QQBrowser, UCBrowser,
		HeyTapBrowser, SamsungBrowser, MiuiBrowser, HuaweiBrowser, BaiduBrowser,
		SogouBrowser, QuarkBrowser, Browser360, "QtWebEngine", "Electron", "CEF":
		return true
	default:
		return false
	}
}
