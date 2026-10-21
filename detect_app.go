package useragent

func detectApp(ua *UserAgent, ts *tokenSet) {
	raw := ua.Raw
	if ua.App.Name == "" {
		switch {
		case ts.has("QxsecApp"):
			ua.App = App{Name: "QxsecApp", Version: tokenValue(ts, "QxsecApp"), InApp: true}
		case ts.has("QXSEC") && ts.contains("QXSEC Android App"):
			ua.App = App{Name: "QXSEC", InApp: true}
		case ts.has("Seafile"):
			ua.App = App{Name: "Seafile", Version: tokenValue(ts, "Seafile"), InApp: true}
		case ts.has("FBAN") || ts.has("FB_IAB"):
			ua.App = App{Name: "Facebook", Version: tokenValue(ts, "FBAV", "FBAN"), InApp: true}
		case ts.has("Instagram"):
			version := firstNonEmpty(tokenValue(ts, "Instagram"), firstVersionAfter(raw, "Instagram "))
			ua.App = App{Name: "Instagram", Version: version, InApp: true}
		case ts.hasPrefix("musical_ly") || ts.hasPrefix("musically") ||
			ts.hasPrefix("trill") || ts.has("TikTok"):
			version := firstNonEmpty(
				appVersionToken(ts, raw, "app_version", "musical_ly", "trill"),
				firstVersionAfter(raw, "musical_ly_"),
				firstVersionAfter(raw, "trill_"),
			)
			ua.App = App{Name: "TikTok", Version: version, InApp: true}
		case ts.has("MicroMessenger"):
			ua.App = App{Name: "WeChat", Version: tokenValue(ts, "MicroMessenger"), InApp: true}
		case ts.has("QQ"):
			ua.App = App{Name: "QQ", Version: tokenValue(ts, "QQ"), InApp: true}
		case ts.has("Line"):
			ua.App = App{Name: "LINE", Version: tokenValue(ts, "Line"), InApp: true}
		case ts.has("Pinterest"):
			ua.App = App{Name: "Pinterest", Version: tokenValue(ts, "Pinterest"), InApp: true}
		case ts.has("Snapchat"):
			ua.App = App{Name: "Snapchat", Version: tokenValue(ts, "Snapchat"), InApp: true}
		case ts.has("CFNetwork"):
			ua.App = App{Name: "CFNetwork", Version: tokenValue(ts, "CFNetwork"), InApp: true}
		case ts.has("Electron"):
			ua.App = App{Name: "Electron", Version: tokenValue(ts, "Electron"), InApp: true}
		case ts.has("QtWebEngine"):
			ua.App = App{Name: "QtWebEngine", Version: tokenValue(ts, "QtWebEngine"), InApp: true}
		case ts.has("CEF"):
			ua.App = App{Name: "CEF", Version: tokenValue(ts, "CEF"), InApp: true}
		case ts.has("GSA"):
			ua.App = App{Name: "Google", Version: tokenValue(ts, "GSA"), InApp: true}
		case ts.has("WebAppManager"):
			ua.App = App{Name: "WebAppManager", InApp: true}
		}
	}

	if ua.App.InApp || ts.has("wv") || ts.has("Hybrid()") ||
		ts.has("BytedanceWebview") || ts.has("FB_IAB") ||
		ts.has("Instagram") || ts.has("MicroMessenger") {
		ua.WebView = true
	}
}

func appVersionToken(ts *tokenSet, raw string, products ...string) string {
	for _, p := range products {
		if v := tokenValue(ts, p); v != "" {
			return v
		}
		if v := productVersion(raw, p); v != "" {
			return v
		}
	}
	return ""
}
