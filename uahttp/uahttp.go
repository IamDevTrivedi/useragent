// Package uahttp integrates the useragent parser with net/http requests and
// User-Agent Client Hints. It has no dependencies outside the standard library
// and the parent module.
package uahttp

import (
	"net/http"
	"strconv"
	"strings"

	ua "github.com/IamDevTrivedi/useragent"
)

// Parse parses the User-Agent header and client hints from r using
// ua.DefaultParser. A nil request returns the zero UserAgent.
func Parse(r *http.Request) ua.UserAgent {
	return ParseWith(ua.DefaultParser, r)
}

// ParseWith is Parse with a custom parser.
func ParseWith(p *ua.Parser, r *http.Request) ua.UserAgent {
	if r == nil {
		return ua.UserAgent{}
	}
	if p == nil {
		p = ua.DefaultParser
	}
	out := p.Parse(r.UserAgent())
	out.ApplyClientHints(ParseClientHints(r.Header))
	return out
}

// ParseClientHints extracts Sec-CH-UA-* headers. Missing or malformed values
// are ignored; this function never panics.
func ParseClientHints(h http.Header) ua.ClientHints {
	if h == nil {
		return ua.ClientHints{}
	}
	return ua.ClientHints{
		Brands:          ua.ParseBrandList(h.Get("Sec-CH-UA")),
		FullVersionList: ua.ParseBrandList(h.Get("Sec-CH-UA-Full-Version-List")),
		Mobile:          parseBoolHint(h.Get("Sec-CH-UA-Mobile")),
		Platform:        unquoteHint(h.Get("Sec-CH-UA-Platform")),
		PlatformVersion: unquoteHint(h.Get("Sec-CH-UA-Platform-Version")),
		Architecture:    unquoteHint(h.Get("Sec-CH-UA-Arch")),
		Bitness:         unquoteHint(h.Get("Sec-CH-UA-Bitness")),
		Model:           unquoteHint(h.Get("Sec-CH-UA-Model")),
		Wow64:           parseBoolHint(h.Get("Sec-CH-UA-WoW64")),
	}
}

func parseBoolHint(s string) bool {
	switch strings.TrimSpace(strings.ToLower(s)) {
	case "?1", "1", "true":
		return true
	default:
		return false
	}
}

func unquoteHint(s string) string {
	s = strings.TrimSpace(s)
	if len(s) >= 2 && s[0] == '"' && s[len(s)-1] == '"' {
		if v, err := strconv.Unquote(s); err == nil {
			return v
		}
		return s[1 : len(s)-1]
	}
	return s
}
