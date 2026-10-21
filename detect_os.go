package useragent

import "strings"

func detectOS(ua *UserAgent, ts *tokenSet) {
	if ua.OS.Name != "" {
		return
	}
	raw := ua.Raw

	switch {
	case isIOS(ts):
		setOS(ua, IOS, iosVersion(ts, raw))
	case isWindowsPhone(ts):
		setOS(ua, WindowsPhone, windowsPhoneVersion(ts, raw))
	case isHarmonyOS(ts):
		setOS(ua, HarmonyOS, harmonyVersion(ts, raw))
	case isKaiOS(ts):
		setOS(ua, KaiOS, kaiVersion(ts, raw))
	case isBlackBerry(ts):
		setOS(ua, BlackBerry, blackBerryVersion(raw))
	case isAndroid(ts):
		setOS(ua, Android, androidVersion(ts, raw))
	case isChromeOS(ts):
		setOS(ua, ChromeOS, chromeOSVersion(ts, raw))
	case isTizen(ts):
		setOS(ua, Tizen, tizenVersion(ts, raw))
	case isWebOS(ts):
		setOS(ua, WebOS, webOSVersion(ts, raw))
	case isMacOS(ts):
		setOS(ua, MacOS, macVersion(ts, raw))
	case isWindows(ts):
		setOS(ua, Windows, windowsVersion(ts, raw))
	case isFreeBSD(ts):
		setOS(ua, FreeBSD, rawVersion(raw, "FreeBSD "))
	case isOpenBSD(ts):
		setOS(ua, OpenBSD, rawVersion(raw, "OpenBSD "))
	case isNetBSD(ts):
		setOS(ua, NetBSD, rawVersion(raw, "NetBSD "))
	case isSolaris(ts):
		setOS(ua, Solaris, rawVersion(raw, "Solaris "))
	case isLinux(ts):
		setOS(ua, Linux, linuxVersion(raw))
	}
}

func isIOS(ts *tokenSet) bool {
	if ts.has("iPhone") || ts.has("iPad") || ts.has("iPod") || ts.hasPrefix("iPod") {
		return true
	}
	if ts.has("iPhone OS") || ts.has("CPU OS") {
		return true
	}
	return ts.has("iOS") && (ts.has("Darwin") || ts.has("CFNetwork") || ts.contains("iOS/"))
}

func isAndroid(ts *tokenSet) bool {
	return ts.has("Android")
}

func isWindowsPhone(ts *tokenSet) bool {
	return ts.has("Windows Phone") || ts.has("Windows Phone OS")
}

func isWindows(ts *tokenSet) bool {
	return ts.has("Windows NT") || ts.has("Windows")
}

func isChromeOS(ts *tokenSet) bool {
	return ts.has("CrOS")
}

func isTizen(ts *tokenSet) bool {
	return ts.has("Tizen")
}

func isWebOS(ts *tokenSet) bool {
	return ts.has("Web0S") || ts.has("webOS") || ts.has("NetCast")
}

func isHarmonyOS(ts *tokenSet) bool {
	return ts.has("HarmonyOS") || ts.has("OpenHarmony")
}

func isKaiOS(ts *tokenSet) bool {
	return ts.has("KAIOS") || ts.has("KaiOS")
}

func isBlackBerry(ts *tokenSet) bool {
	return ts.has("BlackBerry") || ts.has("BB10")
}

func isMacOS(ts *tokenSet) bool {
	return ts.has("Macintosh") || ts.has("Mac OS X")
}

func isFreeBSD(ts *tokenSet) bool { return ts.rawHas("FreeBSD") }

func isOpenBSD(ts *tokenSet) bool { return ts.rawHas("OpenBSD") }

func isNetBSD(ts *tokenSet) bool { return ts.rawHas("NetBSD") }

func isSolaris(ts *tokenSet) bool {
	return ts.has("Solaris") || ts.rawHas("Solaris") || ts.rawHas("SunOS")
}

func isLinux(ts *tokenSet) bool {
	return ts.has("Linux") || ts.rawHas("Linux")
}

func iosVersion(ts *tokenSet, raw string) string {
	if v := tokenValue(ts, "iPhone OS", "CPU OS", "iOS"); v != "" {
		return v
	}
	for _, marker := range [...]string{"iPhone OS ", "CPU OS ", "CPU iPhone OS ", "iOS/"} {
		if v := rawVersion(raw, marker); v != "" {
			return v
		}
	}
	return ""
}

func androidVersion(ts *tokenSet, raw string) string {
	if v := ts.get("Android"); v != "" {
		return normalizeVersion(v)
	}
	if v := rawVersion(raw, "Android "); v != "" {
		return v
	}
	// Some Android app UAs put the version only in the device comment,
	// immediately after the model, for example "...; Xiaomi M2101K9G; 13)".
	// Prefer the last numeric-only comment, which is the platform version in
	// these UAs.
	if ts.has("Android") || ts.contains("Android App") {
		best := ""
		for i := 0; i < ts.cn; i++ {
			if isNumericOnly(ts.comments[i]) {
				best = normalizeVersion(ts.comments[i])
			}
		}
		if best != "" {
			return best
		}
	}
	return ""
}

func windowsVersion(ts *tokenSet, raw string) string {
	if v := tokenValue(ts, "Windows NT"); v != "" {
		return v
	}
	return rawVersion(raw, "Windows NT ")
}

func windowsPhoneVersion(ts *tokenSet, raw string) string {
	if v := tokenValue(ts, "Windows Phone OS", "Windows Phone"); v != "" {
		return v
	}
	return firstNonEmpty(rawVersion(raw, "Windows Phone OS "), rawVersion(raw, "Windows Phone "))
}

func macVersion(ts *tokenSet, raw string) string {
	if v := tokenValue(ts, "Mac OS X"); v != "" {
		return v
	}
	return rawVersion(raw, "Mac OS X ")
}

func chromeOSVersion(ts *tokenSet, raw string) string {
	if v := tokenValue(ts, "CrOS"); v != "" {
		return v
	}
	for i := 0; i+len("CrOS ") <= len(raw); i++ {
		if !hasPrefixFold(raw[i:], "CrOS ") {
			continue
		}
		rest := raw[i+len("CrOS "):]
		if at := strings.IndexAny(rest, ");]"); at >= 0 {
			rest = rest[:at]
		}
		fields := strings.Fields(rest)
		for j := len(fields) - 1; j >= 0; j-- {
			if beginsWithDigit(fields[j]) {
				return normalizeVersion(fields[j])
			}
		}
	}
	return ""
}

func tizenVersion(ts *tokenSet, raw string) string {
	if v := tokenValue(ts, "Tizen"); v != "" {
		return v
	}
	return rawVersion(raw, "Tizen ")
}

func webOSVersion(ts *tokenSet, raw string) string {
	if v := tokenValue(ts, "Web0S", "webOS"); v != "" {
		return v
	}
	return firstNonEmpty(rawVersion(raw, "Web0S/"), rawVersion(raw, "webOS/"), rawVersion(raw, "WebOS/"))
}

func harmonyVersion(ts *tokenSet, raw string) string {
	if v := tokenValue(ts, "HarmonyOS", "OpenHarmony"); v != "" {
		return v
	}
	return firstNonEmpty(rawVersion(raw, "HarmonyOS "), rawVersion(raw, "OpenHarmony "))
}

func kaiVersion(ts *tokenSet, raw string) string {
	if v := tokenValue(ts, "KAIOS"); v != "" {
		return v
	}
	return browserVersion(raw, "KAIOS", "KaiOS")
}

func blackBerryVersion(raw string) string {
	if v := productVersion(raw, "Version"); v != "" {
		return v
	}
	return firstNonEmpty(rawVersion(raw, "BlackBerry "), rawVersion(raw, "BB10 "))
}

func linuxVersion(raw string) string {
	return rawVersion(raw, "Linux ")
}
