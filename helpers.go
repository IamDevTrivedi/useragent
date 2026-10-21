package useragent

import "strings"

func setBrowser(ua *UserAgent, name, version string) {
	if name == "" {
		return
	}
	ua.Browser.Name = name
	ua.Browser.Version = normalizeVersion(version)
	ua.Browser.VersionNo = ParseVersion(ua.Browser.Version)
}

func setOS(ua *UserAgent, name, version string) {
	if name == "" {
		return
	}
	ua.OS.Name = name
	ua.OS.Version = normalizeVersion(version)
	ua.OS.VersionNo = ParseVersion(ua.OS.Version)
}

func setEngine(ua *UserAgent, name, version string) {
	if name == "" {
		return
	}
	ua.Engine.Name = name
	ua.Engine.Version = normalizeVersion(version)
}

func setBot(ua *UserAgent, name, version, url string, category BotCategory) {
	ua.Bot.Detected = true
	ua.Bot.Name = name
	ua.Bot.Version = normalizeVersion(version)
	ua.Bot.URL = url
	ua.Bot.Category = category
}

func tokenValue(ts *tokenSet, keys ...string) string {
	for _, k := range keys {
		if v := ts.get(k); v != "" {
			return normalizeVersion(v)
		}
	}
	return ""
}

func rawVersion(raw, marker string) string {
	return normalizeVersion(firstVersionAfter(raw, marker))
}

func browserVersion(raw string, products ...string) string {
	for _, p := range products {
		if v := productVersion(raw, p); v != "" {
			return v
		}
	}
	return ""
}

func firstNonEmpty(values ...string) string {
	for _, v := range values {
		if v != "" {
			return v
		}
	}
	return ""
}

func cleanModel(s string) string {
	s = strings.TrimSpace(s)
	if at := indexFold(s, " Build"); at >= 0 {
		s = s[:at]
	}
	if at := indexFold(s, " build/"); at >= 0 {
		s = s[:at]
	}
	s = strings.TrimSpace(s)
	return strings.Trim(s, "()[]{};,")
}

func indexFold(s, sub string) int {
	if len(sub) == 0 {
		return 0
	}
	if len(sub) > len(s) {
		return -1
	}
	first := lowerASCII(sub[0])
	last := len(s) - len(sub)
	for i := 0; i <= last; {
		idx := strings.IndexByte(s[i:last+1], first)
		if first >= 'a' && first <= 'z' {
			if upper := strings.IndexByte(s[i:last+1], first-('a'-'A')); upper >= 0 && (idx < 0 || upper < idx) {
				idx = upper
			}
		}
		if idx < 0 {
			return -1
		}
		pos := i + idx
		if hasPrefixFold(s[pos:], sub) {
			return pos
		}
		i = pos + 1
	}
	return -1
}
