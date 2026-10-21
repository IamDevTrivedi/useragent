package useragent

import "strings"

// tokenize performs a single pass over raw. Product tokens are read outside
// parentheses and square brackets. Comment blocks are stored separately so
// device and OS detection can inspect ordered fields without guessing which
// unrelated product token is a model number.
func tokenize(raw string) tokenSet {
	var ts tokenSet
	ts.raw = raw
	ts.url = extractURL(raw)

	for i := 0; i < len(raw); {
		switch raw[i] {
		case ' ', '\t', '\r', '\n', ',', ';':
			i++
		case '(':
			end := findCommentEnd(raw, i+1, ')')
			parseCommentBlock(raw[i+1:end], &ts)
			if end < len(raw) {
				end++
			}
			i = end
		case '[':
			end := findCommentEnd(raw, i+1, ']')
			parseCommentBlock(raw[i+1:end], &ts)
			if end < len(raw) {
				end++
			}
			i = end
		case ')', ']':
			i++
		default:
			key, value, next := scanProduct(raw, i)
			if key != "" {
				ts.add(key, value, false)
				if isArchWord(value) {
					ts.add(value, "", false)
				}
				// Vendor products are often nested one level, as in
				// "XiaoMi/MiuiBrowser/12.11.5-gn" or
				// "Browser/NetFront/3.3". Expose the inner product too so
				// callers do not need a full raw scan for it.
				if at := strings.IndexByte(value, '/'); at > 0 {
					ts.add(value[:at], value[at+1:], false)
				}
			}
			if next <= i {
				next = i + 1
			}
			i = next
		}
	}
	return ts
}

func findCommentEnd(raw string, start int, closing byte) int {
	depth := 1
	open := byte('(')
	if closing == ']' {
		open = '['
	}
	for i := start; i < len(raw); i++ {
		switch raw[i] {
		case open:
			depth++
		case closing:
			depth--
			if depth == 0 {
				return i
			}
		}
	}
	return len(raw)
}

// scanProduct reads either a bare word or a product such as "Mobile Safari/537"
// or "Chrome/120.0". Multi-word product names are recognized only for a small
// allowlist; otherwise the last word before '/' is used. This prevents a
// preceding app label from becoming part of the browser name.
func scanProduct(raw string, start int) (key, value string, next int) {
	firstEnd := start
	for firstEnd < len(raw) && !isBoundary(raw[firstEnd]) {
		firstEnd++
	}
	if firstEnd == start {
		return "", "", start + 1
	}

	slash := -1
	if at := strings.IndexByte(raw[start:firstEnd], '/'); at >= 0 {
		slash = start + at
	} else if firstEnd < len(raw) && raw[firstEnd] == ' ' && isMultiWordProductStart(raw[start:firstEnd]) {
		// Product names such as "Mobile Safari/537" and "Opera Mini/28"
		// span exactly one space. Do not scan across arbitrary spaces: a
		// bare token like "splash" must remain its own token instead of
		// being folded into the following "Version/10.0".
		j := firstEnd + 1
		for j < len(raw) && !isBoundary(raw[j]) {
			if raw[j] == '/' {
				slash = j
				break
			}
			j++
		}
		if slash >= 0 {
			candidate := strings.TrimSpace(raw[start:slash])
			if !isMultiWordProduct(candidate) {
				slash = -1
			}
		}
	}

	if slash < 0 {
		return raw[start:firstEnd], "", firstEnd
	}

	key = strings.TrimSpace(raw[start:slash])
	if key == "" {
		return "", "", slash + 1
	}

	end := slash + 1
	vstart := end
	for end < len(raw) && !isBoundary(raw[end]) {
		end++
	}
	return key, raw[vstart:end], end
}

func isArchWord(word string) bool {
	switch word {
	case "amd64", "arm64", "aarch64", "Win64", "WOW64", "x64", "x86_64",
		"i686", "i386", "i86pc", "armv7", "armv7l", "armv8l", "ppc64",
		"riscv64", "mips", "mips64":
		return true
	default:
		return false
	}
}

func isMultiWordProductStart(word string) bool {
	switch {
	case equalFold(word, "Mobile"), equalFold(word, "TV"),
		equalFold(word, "Brave"), equalFold(word, "Opera"),
		equalFold(word, "Internet"), equalFold(word, "Microsoft"),
		equalFold(word, "Samsung"), equalFold(word, "Google"):
		return true
	default:
		return false
	}
}

func isMultiWordProduct(key string) bool {
	switch {
	case equalFold(key, "Mobile Safari"),
		equalFold(key, "TV Safari"),
		equalFold(key, "Brave Chrome"),
		equalFold(key, "Opera Mini"),
		equalFold(key, "Opera Touch"),
		equalFold(key, "Internet Explorer"),
		equalFold(key, "Microsoft Edge"),
		equalFold(key, "Google Chrome"),
		equalFold(key, "Samsung Internet"):
		return true
	default:
		return false
	}
}

func parseCommentBlock(block string, ts *tokenSet) {
	start := 0
	for i := 0; i <= len(block); i++ {
		if i != len(block) && block[i] != ';' {
			continue
		}
		field := strings.TrimSpace(block[start:i])
		start = i + 1
		if field == "" {
			continue
		}
		ts.addComment(field)
		if isURLField(field) {
			continue
		}
		ts.add(field, "", true)
		parseCommentField(field, ts)
	}
}

func isURLField(field string) bool {
	return hasPrefixFold(field, "http://") ||
		hasPrefixFold(field, "https://") ||
		hasPrefixFold(field, "+http://") ||
		hasPrefixFold(field, "+https://")
}

func parseCommentField(field string, ts *tokenSet) {
	// Comment fields can contain multi-word products, for example
	// "Opera Mini/28.0.2254/66.318". Parse the leading product before the
	// generic space-delimited product scan below.
	if key, value, _ := scanProduct(field, 0); key != "" && strings.IndexByte(key, ' ') > 0 {
		ts.add(key, value, true)
		if at := strings.IndexByte(value, '/'); at > 0 {
			ts.add(value[:at], value[at+1:], true)
		}
	}

	switch {
	case hasPrefixFold(field, "Android "):
		ts.add("Android", normalizeVersion(field[len("Android "):]), true)
	case hasPrefixFold(field, "Windows NT "):
		ts.add("Windows NT", normalizeVersion(field[len("Windows NT "):]), true)
	case hasPrefixFold(field, "Windows Phone OS "):
		ts.add("Windows Phone OS", normalizeVersion(field[len("Windows Phone OS "):]), true)
	case hasPrefixFold(field, "Windows Phone "):
		ts.add("Windows Phone", normalizeVersion(field[len("Windows Phone "):]), true)
	case hasPrefixFold(field, "CrOS "):
		parseChromeOSField(field, ts)
	case hasPrefixFold(field, "Tizen "):
		ts.add("Tizen", normalizeVersion(field[len("Tizen "):]), true)
	case hasPrefixFold(field, "OpenHarmony "):
		ts.add("OpenHarmony", normalizeVersion(field[len("OpenHarmony "):]), true)
	case hasPrefixFold(field, "HarmonyOS "):
		ts.add("HarmonyOS", normalizeVersion(field[len("HarmonyOS "):]), true)
	case hasPrefixFold(field, "CPU iPhone OS "):
		ts.add("iPhone OS", normalizeVersion(markerValue(field, "CPU iPhone OS ")), true)
		ts.add("iPhone", "", true)
	case hasPrefixFold(field, "CPU OS "):
		ts.add("CPU OS", normalizeVersion(markerValue(field, "CPU OS ")), true)
	case hasPrefixFold(field, "iPhone OS "):
		ts.add("iPhone OS", normalizeVersion(markerValue(field, "iPhone OS ")), true)
		ts.add("iPhone", "", true)
	case hasPrefixFold(field, "Intel Mac OS X "):
		ts.add("Mac OS X", normalizeVersion(markerValue(field, "Intel Mac OS X ")), true)
		ts.add("Macintosh", "", true)
	case hasPrefixFold(field, "Mac OS X "):
		ts.add("Mac OS X", normalizeVersion(markerValue(field, "Mac OS X ")), true)
	case hasPrefixFold(field, "MSIE "):
		ts.add("MSIE", normalizeVersion(markerValue(field, "MSIE ")), true)
	}

	// A few comment fields contain one or more product tokens, for example
	// "iPhone iOS/16.3.1", "FBDV/iPhone8,2", or "Build/MMB29P".
	for i := 0; i < len(field); {
		word, next := nextSpaceField(field, i)
		if word != "" {
			if at := strings.IndexByte(word, '/'); at > 0 && at < len(word)-1 && !isURLField(word) {
				ts.add(word[:at], word[at+1:], true)
				if isArchWord(word[at+1:]) {
					ts.add(word[at+1:], "", true)
				}
			}
			if isArchWord(word) && !ts.has(word) {
				ts.add(word, "", true)
			}
		}
		if next <= i {
			break
		}
		i = next
	}
}

// nextSpaceField returns the next space- or tab-delimited word and the index
// just after it. It avoids the slice allocated by strings.Fields in the hot
// tokenizer path.
func nextSpaceField(s string, i int) (string, int) {
	for i < len(s) && (s[i] == ' ' || s[i] == '\t') {
		i++
	}
	start := i
	for i < len(s) && s[i] != ' ' && s[i] != '\t' {
		i++
	}
	return s[start:i], i
}

func parseChromeOSField(field string, ts *tokenSet) {
	rest := strings.TrimSpace(field[len("CrOS "):])
	version := ""
	for end := len(rest); end > 0; {
		start := end
		for start > 0 && rest[start-1] != ' ' && rest[start-1] != '\t' {
			start--
		}
		part := rest[start:end]
		if beginsWithDigit(part) {
			version = normalizeVersion(part)
			break
		}
		end = start
		for end > 0 && (rest[end-1] == ' ' || rest[end-1] == '\t') {
			end--
		}
	}
	ts.add("CrOS", version, true)
	for i := 0; i < len(rest); {
		p, next := nextSpaceField(rest, i)
		switch p {
		case "x86_64", "i686", "i386", "armv7l", "armv8l", "aarch64":
			ts.add("CrOS "+p, "", true)
		}
		if next <= i {
			break
		}
		i = next
	}
}

func markerValue(field, marker string) string {
	rest := strings.TrimSpace(field[len(marker):])
	if at := strings.IndexByte(rest, ' '); at >= 0 {
		return rest[:at]
	}
	return rest
}

func beginsWithDigit(s string) bool {
	return len(s) != 0 && isDigit(s[0])
}

func extractURL(raw string) string {
	start := -1
	for off := 0; off+len("http") <= len(raw); {
		at := indexFold(raw[off:], "http")
		if at < 0 {
			break
		}
		i := off + at
		if hasPrefixFold(raw[i:], "https://") || hasPrefixFold(raw[i:], "http://") {
			start = i
			break
		}
		off = i + len("http")
	}
	if start < 0 {
		return ""
	}
	end := start
	for end < len(raw) && !isURLTerminator(raw[end]) {
		end++
	}
	for end > start && isTrailingURLPunct(raw[end-1]) {
		end--
	}
	return raw[start:end]
}

func isURLTerminator(c byte) bool {
	switch c {
	case ' ', '\t', '\r', '\n', ')', ']', ';', ',', '"', '\'', '<', '>':
		return true
	default:
		return false
	}
}

func isTrailingURLPunct(c byte) bool {
	switch c {
	case '.', ',', ';', ':', '!', '?':
		return true
	default:
		return false
	}
}
