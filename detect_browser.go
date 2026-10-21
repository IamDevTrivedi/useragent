package useragent

import "strings"

func tridentToIE(trident string) string {
	switch {
	case strings.HasPrefix(trident, "7."):
		return "11.0"
	case strings.HasPrefix(trident, "6."):
		return "10.0"
	case strings.HasPrefix(trident, "5."):
		return "9.0"
	case strings.HasPrefix(trident, "4."):
		return "8.0"
	default:
		return trident
	}
}

func detectBrowser(ua *UserAgent, ts *tokenSet) {
	if ua.Browser.Name != "" {
		return
	}
	raw := ua.Raw

	// App-labelled vendor UAs that do not identify themselves as a regular
	// browser should keep their own name rather than falling through to
	// AppleWebKit.
	if ts.has("QXSEC") && ts.contains("QXSEC Android App") {
		version := firstNonEmpty(productVersion(raw, "QxsecApp"), firstVersionAfter(raw, "QXSEC Android App "))
		setBrowser(ua, "QXSEC Android App", version)
		return
	}

	// Embedded runtimes identify themselves as Chromium too, but their own
	// product name is more useful to callers.
	switch {
	case ts.has("QtWebEngine"):
		version := firstNonEmpty(tokenValue(ts, "QtWebEngine"), productVersion(raw, "QtWebEngine"))
		setBrowser(ua, "QtWebEngine", version)
		return
	case ts.has("Electron"):
		version := firstNonEmpty(tokenValue(ts, "Electron"), productVersion(raw, "Electron"))
		setBrowser(ua, "Electron", version)
		return
	case ts.has("CEF"):
		version := firstNonEmpty(tokenValue(ts, "CEF"), productVersion(raw, "CEF"))
		setBrowser(ua, "CEF", version)
		return
	}

	// Regional and vendor browsers take precedence over the Chromium family.
	switch {
	case ts.hasAny("QQBrowser", "MQQBrowser"):
		version := firstNonEmpty(tokenValue(ts, "QQBrowser", "MQQBrowser"), productVersion(raw, "QQBrowser"), productVersion(raw, "MQQBrowser"))
		setBrowser(ua, QQBrowser, version)
		return
	case ts.has("UCBrowser"):
		version := firstNonEmpty(tokenValue(ts, "UCBrowser"), productVersion(raw, "UCBrowser"))
		setBrowser(ua, UCBrowser, version)
		return
	case ts.has("HeyTapBrowser"):
		version := firstNonEmpty(tokenValue(ts, "HeyTapBrowser"), productVersion(raw, "HeyTapBrowser"))
		setBrowser(ua, HeyTapBrowser, version)
		return
	case ts.has("SamsungBrowser"):
		version := firstNonEmpty(tokenValue(ts, "SamsungBrowser"), productVersion(raw, "SamsungBrowser"))
		setBrowser(ua, SamsungBrowser, version)
		return
	case ts.has("MiuiBrowser"):
		version := firstNonEmpty(tokenValue(ts, "MiuiBrowser"), productVersion(raw, "MiuiBrowser"))
		setBrowser(ua, MiuiBrowser, version)
		return
	case ts.has("HuaweiBrowser"):
		version := firstNonEmpty(tokenValue(ts, "HuaweiBrowser"), productVersion(raw, "HuaweiBrowser"))
		setBrowser(ua, HuaweiBrowser, version)
		return
	case ts.has("BaiduBrowser") || ts.has("BIDUBrowser"):
		version := firstNonEmpty(tokenValue(ts, "BaiduBrowser", "BIDUBrowser"), productVersion(raw, "BaiduBrowser"))
		setBrowser(ua, BaiduBrowser, version)
		return
	case ts.has("SogouMobileBrowser") || ts.has("MetaSr"):
		version := firstNonEmpty(tokenValue(ts, "SogouMobileBrowser", "MetaSr"), productVersion(raw, "SogouMobileBrowser"))
		setBrowser(ua, SogouBrowser, version)
		return
	case ts.has("Quark"):
		version := firstNonEmpty(tokenValue(ts, "Quark"), productVersion(raw, "Quark"))
		setBrowser(ua, QuarkBrowser, version)
		return
	case ts.has("QihooBrowser") || ts.has("360SE") || ts.has("360EE"):
		version := firstNonEmpty(tokenValue(ts, "QihooBrowser", "360SE", "360EE"), productVersion(raw, "QihooBrowser"))
		setBrowser(ua, Browser360, version)
		return
	case ts.has("coc_coc_browser"):
		version := firstNonEmpty(tokenValue(ts, "coc_coc_browser"), productVersion(raw, "coc_coc_browser"))
		setBrowser(ua, CocCoc, version)
		return
	case ts.has("YaBrowser"):
		version := firstNonEmpty(tokenValue(ts, "YaBrowser"), productVersion(raw, "YaBrowser"))
		setBrowser(ua, YandexBrowser, version)
		return
	case ts.has("Whale"):
		version := firstNonEmpty(tokenValue(ts, "Whale"), productVersion(raw, "Whale"))
		setBrowser(ua, Whale, version)
		return
	case ts.has("DuckDuckGo"):
		version := firstNonEmpty(tokenValue(ts, "DuckDuckGo"), productVersion(raw, "DuckDuckGo"))
		setBrowser(ua, DuckDuckGo, version)
		return
	case ts.has("ArkWeb"):
		version := firstNonEmpty(tokenValue(ts, "ArkWeb"), productVersion(raw, "ArkWeb"))
		setBrowser(ua, "ArkWeb", version)
		return
	case ts.has("Brave") || ts.has("Brave Chrome"):
		version := firstNonEmpty(tokenValue(ts, "Brave", "Brave Chrome"), tokenValue(ts, "Chrome", "Chromium"), productVersion(raw, "Brave"))
		setBrowser(ua, Brave, version)
		return
	case ts.has("Vivaldi"):
		version := firstNonEmpty(tokenValue(ts, "Vivaldi"), productVersion(raw, "Vivaldi"))
		setBrowser(ua, Vivaldi, version)
		return
	case ts.has("Arc"):
		version := firstNonEmpty(tokenValue(ts, "Arc"), productVersion(raw, "Arc"))
		setBrowser(ua, Arc, version)
		return
	}

	// Opera variants must be checked before Chrome and Safari.
	switch {
	case ts.has("Opera Mini"):
		version := firstNonEmpty(tokenValue(ts, "Opera Mini"), productVersion(raw, "Opera Mini"))
		setBrowser(ua, OperaMini, version)
		return
	case ts.has("Opera Touch") || ts.has("OPT"):
		version := firstNonEmpty(tokenValue(ts, "Opera Touch", "OPT"), productVersion(raw, "Opera Touch"), productVersion(raw, "OPT"))
		setBrowser(ua, OperaTouch, version)
		return
	case ts.has("OPiOS"):
		version := firstNonEmpty(tokenValue(ts, "OPiOS"), productVersion(raw, "OPiOS"))
		setBrowser(ua, Opera, version)
		return
	case ts.has("OPR"):
		version := firstNonEmpty(tokenValue(ts, "OPR"), productVersion(raw, "OPR"))
		setBrowser(ua, Opera, version)
		return
	case ts.has("Opera"):
		version := firstNonEmpty(tokenValue(ts, "Opera"), productVersion(raw, "Opera"), tokenValue(ts, "Version"))
		setBrowser(ua, Opera, version)
		return
	}

	switch {
	case ts.hasAny("EdgiOS", "EdgA", "Edg", "Edge"):
		version := firstNonEmpty(tokenValue(ts, "EdgiOS", "EdgA", "Edg", "Edge"), productVersion(raw, "Edg"), productVersion(raw, "Edge"))
		setBrowser(ua, Edge, version)
		return
	case ts.hasAny("FxiOS", "Firefox", "Focus", "Klar", "Waterfox", "PaleMoon", "Basilisk"):
		name := Firefox
		switch {
		case ts.has("Waterfox"):
			name = "Waterfox"
		case ts.has("PaleMoon"):
			name = "Pale Moon"
		case ts.has("Basilisk"):
			name = "Basilisk"
		}
		version := firstNonEmpty(tokenValue(ts, "FxiOS", "Firefox", "Focus", "Klar", "Waterfox", "PaleMoon", "Basilisk"), productVersion(raw, "Firefox"), productVersion(raw, "FxiOS"))
		setBrowser(ua, name, version)
		return
	case ts.hasAny("CriOS", "Chrome", "Chromium"):
		version := firstNonEmpty(tokenValue(ts, "CriOS", "Chrome", "Chromium"), productVersion(raw, "Chrome"), productVersion(raw, "Chromium"))
		name := Chrome
		if !ts.has("Chrome") && !ts.has("CriOS") && ts.has("Chromium") {
			name = Chromium
		}
		if ts.has("CriOS") {
			name = Chrome
		}
		setBrowser(ua, name, version)
		return
	case ts.has("MSIE") || ts.has("Trident"):
		version := firstNonEmpty(tokenValue(ts, "MSIE"), rawVersion(raw, "rv:"))
		if version == "" {
			version = tridentToIE(tokenValue(ts, "Trident"))
		}
		setBrowser(ua, InternetExplorer, version)
		return
	}

	// Android stock browser before generic Safari: it has Version/4.0 and a
	// Safari token but no Chrome/Chromium token.
	if ua.OS.Name == Android && ts.has("Version") && (ts.has("Safari") || ts.has("Mobile Safari")) &&
		!ts.has("Chrome") && !ts.has("Chromium") && !ts.has("CriOS") {
		version := firstNonEmpty(tokenValue(ts, "Version"), tokenValue(ts, "Safari"))
		setBrowser(ua, AndroidBrowser, version)
		return
	}

	// BlackBerry's WebKit browser identifies itself with Safari and Version,
	// but the product name should remain BlackBerry.
	if ua.OS.Name == BlackBerry && (ts.has("Safari") || ts.has("Version")) {
		version := firstNonEmpty(tokenValue(ts, "Version"), tokenValue(ts, "Safari"))
		setBrowser(ua, BlackBerry, version)
		return
	}

	if ts.has("Safari") || ts.has("Mobile Safari") || ts.has("TV Safari") {
		version := ""
		if ua.Device.Type == DeviceTV {
			version = firstNonEmpty(tokenValue(ts, "Safari", "TV Safari", "Mobile Safari"))
		} else {
			version = firstNonEmpty(tokenValue(ts, "Version"), tokenValue(ts, "Safari", "Mobile Safari", "TV Safari"))
		}
		setBrowser(ua, Safari, version)
		return
	}

	switch {
	case ts.has("CFNetwork"):
		version := firstNonEmpty(tokenValue(ts, "CFNetwork"), productVersion(raw, "CFNetwork"))
		setBrowser(ua, CFNetwork, version)
	case ts.has("Seafile"):
		version := firstNonEmpty(tokenValue(ts, "Seafile"), productVersion(raw, "Seafile"))
		setBrowser(ua, "Seafile", version)
	case ts.has("NetFront"):
		version := firstNonEmpty(tokenValue(ts, "NetFront"), productVersion(raw, "NetFront"))
		setBrowser(ua, NetFront, version)
	case ts.has("Konqueror"):
		version := firstNonEmpty(tokenValue(ts, "Konqueror"), productVersion(raw, "Konqueror"))
		setBrowser(ua, Konqueror, version)
	case ts.has("QtWebEngine"):
		version := firstNonEmpty(tokenValue(ts, "QtWebEngine"), productVersion(raw, "QtWebEngine"))
		setBrowser(ua, "QtWebEngine", version)
	case ts.has("Electron"):
		version := firstNonEmpty(tokenValue(ts, "Electron"), productVersion(raw, "Electron"))
		setBrowser(ua, "Electron", version)
	case ts.has("CEF"):
		version := firstNonEmpty(tokenValue(ts, "CEF"), productVersion(raw, "CEF"))
		setBrowser(ua, "CEF", version)
	case ts.has("Lynx"):
		version := firstNonEmpty(tokenValue(ts, "Lynx"), productVersion(raw, "Lynx"))
		setBrowser(ua, "Lynx", version)
	case ts.has("w3m"):
		version := firstNonEmpty(tokenValue(ts, "w3m"), productVersion(raw, "w3m"))
		setBrowser(ua, "w3m", version)
	case ts.has("Links"):
		version := firstNonEmpty(tokenValue(ts, "Links"), productVersion(raw, "Links"))
		setBrowser(ua, "Links", version)
	case ts.has("AppleWebKit"):
		version := firstNonEmpty(tokenValue(ts, "AppleWebKit"), productVersion(raw, "AppleWebKit"))
		setBrowser(ua, AppleWebKit, version)
	}
}
