package useragent

import "strings"

const (
	maxTokens   = 64
	maxComments = 32
)

type token struct {
	key     string
	value   string
	comment bool
}

type tokenSet struct {
	list     [maxTokens]token
	n        int
	comments [maxComments]string
	cn       int
	url      string
	raw      string
	lowerRaw string
	bloom    uint64
}

func (ts *tokenSet) add(key, value string, comment bool) {
	if key == "" || ts.n >= maxTokens {
		return
	}
	ts.list[ts.n] = token{key: key, value: value, comment: comment}
	ts.n++
	ts.bloom |= keyMask(key)
}

// keyMask is a cheap membership filter for the small token set. It combines
// the first byte, last byte, and length, which are enough to reject nearly all
// misses without hashing the whole key on every lookup.
func keyMask(key string) uint64 {
	if key == "" {
		return 0
	}
	return 1<<(lowerASCII(key[0])&63) |
		1<<((lowerASCII(key[len(key)-1])+17)&63) |
		1<<((uint64(len(key))+29)&63)
}

func (ts *tokenSet) addComment(c string) {
	if c == "" || ts.cn >= maxComments {
		return
	}
	ts.comments[ts.cn] = c
	ts.cn++
}

func (ts *tokenSet) get(key string) string {
	mask := keyMask(key)
	if ts.bloom&mask != mask {
		return ""
	}
	for i := 0; i < ts.n; i++ {
		if equalFold(ts.list[i].key, key) {
			return ts.list[i].value
		}
	}
	return ""
}

func (ts *tokenSet) has(key string) bool {
	mask := keyMask(key)
	if ts.bloom&mask != mask {
		return false
	}
	for i := 0; i < ts.n; i++ {
		if equalFold(ts.list[i].key, key) {
			return true
		}
	}
	return false
}

func (ts *tokenSet) hasAny(keys ...string) bool {
	for _, k := range keys {
		if ts.has(k) {
			return true
		}
	}
	return false
}

func (ts *tokenSet) hasAll(keys ...string) bool {
	for _, k := range keys {
		if !ts.has(k) {
			return false
		}
	}
	return len(keys) != 0
}

func (ts *tokenSet) hasPrefix(prefix string) bool {
	for i := 0; i < ts.n; i++ {
		if hasPrefixFold(ts.list[i].key, prefix) {
			return true
		}
	}
	return false
}

func (ts *tokenSet) hasSuffix(suffix string) bool {
	for i := 0; i < ts.n; i++ {
		if hasSuffixFold(ts.list[i].key, suffix) {
			return true
		}
	}
	return false
}

func (ts *tokenSet) rawHas(sub string) bool {
	for i := 0; i < ts.n; i++ {
		if containsFold(ts.list[i].key, sub) || containsFold(ts.list[i].value, sub) {
			return true
		}
	}
	return false
}

// contains reports whether the lowercased raw user agent contains sub. The
// lowercase copy is built lazily once, on the first detector that needs it.
func (ts *tokenSet) contains(sub string) bool {
	if ts.lowerRaw == "" && ts.raw != "" {
		ts.lowerRaw = strings.ToLower(ts.raw)
	}
	return containsLower(ts.lowerRaw, sub)
}

func containsLower(s, sub string) bool {
	if len(sub) == 0 {
		return true
	}
	if len(sub) > len(s) {
		return false
	}
	first := lowerASCII(sub[0])
	last := len(s) - len(sub)
	for i := 0; i <= last; {
		idx := strings.IndexByte(s[i:last+1], first)
		if idx < 0 {
			return false
		}
		pos := i + idx
		if hasPrefixFold(s[pos:], sub) {
			return true
		}
		i = pos + 1
	}
	return false
}

func equalFold(a, b string) bool {
	if len(a) != len(b) {
		return false
	}
	for i := 0; i < len(a); i++ {
		if lowerASCII(a[i]) != lowerASCII(b[i]) {
			return false
		}
	}
	return true
}

func hasPrefixFold(s, prefix string) bool {
	if len(s) < len(prefix) {
		return false
	}
	for i := 0; i < len(prefix); i++ {
		if lowerASCII(s[i]) != lowerASCII(prefix[i]) {
			return false
		}
	}
	return true
}

func hasSuffixFold(s, suffix string) bool {
	if len(s) < len(suffix) {
		return false
	}
	off := len(s) - len(suffix)
	for i := 0; i < len(suffix); i++ {
		if lowerASCII(s[off+i]) != lowerASCII(suffix[i]) {
			return false
		}
	}
	return true
}

func containsFold(s, sub string) bool { return indexFold(s, sub) >= 0 }

func lowerASCII(c byte) byte {
	if c >= 'A' && c <= 'Z' {
		return c + ('a' - 'A')
	}
	return c
}

func isBoundary(c byte) bool {
	switch c {
	case ' ', '\t', '\r', '\n', '(', ')', '[', ']', ';', ',':
		return true
	default:
		return false
	}
}

func trimSpace(s string) string { return strings.TrimSpace(s) }

func normalizeVersion(s string) string {
	return strings.ReplaceAll(trimSpace(s), "_", ".")
}

// productVersion extracts a product version from raw without allocating if the
// product is present. It looks for product followed by '/' and reads until a
// boundary. This is used for tokens embedded inside larger products, such as
// XiaoMi/MiuiBrowser/12.11.5-gn.
func productVersion(raw, product string) string {
	for off := 0; off+len(product) <= len(raw); {
		at := indexFold(raw[off:], product)
		if at < 0 {
			return ""
		}
		i := off + at
		j := i + len(product)
		for j < len(raw) && (raw[j] == ' ' || raw[j] == '\t') {
			j++
		}
		if j < len(raw) && raw[j] == '/' {
			j++
			start := j
			for j < len(raw) && !isBoundary(raw[j]) {
				j++
			}
			if j > start {
				return normalizeVersion(raw[start:j])
			}
		}
		off = i + 1
	}
	return ""
}

// firstVersionAfter finds the first numeric version-like token after marker.
func firstVersionAfter(raw, marker string) string {
	for off := 0; off+len(marker) <= len(raw); {
		at := indexFold(raw[off:], marker)
		if at < 0 {
			return ""
		}
		i := off + at
		j := i + len(marker)
		for j < len(raw) && (raw[j] == ' ' || raw[j] == '\t' || raw[j] == '/') {
			j++
		}
		start := j
		for j < len(raw) && (isDigit(raw[j]) || raw[j] == '.' || raw[j] == '_') {
			j++
		}
		if j > start {
			return normalizeVersion(raw[start:j])
		}
		off = i + 1
	}
	return ""
}

func isNumericOnly(s string) bool {
	if s == "" {
		return false
	}
	for i := 0; i < len(s); i++ {
		if !isDigit(s[i]) && s[i] != '.' {
			return false
		}
	}
	return true
}
