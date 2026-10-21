package uahttp

import (
	"net/http"
	"testing"
)

func BenchmarkParseRequest(b *testing.B) {
	r := &http.Request{Header: make(http.Header)}
	r.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.6099.109 Safari/537.36")
	r.Header.Set("Sec-CH-UA", `"Chromium";v="120", "Google Chrome";v="120"`)
	r.Header.Set("Sec-CH-UA-Platform", `"Windows"`)
	r.Header.Set("Sec-CH-UA-Platform-Version", `"10.0.0"`)
	r.Header.Set("Sec-CH-UA-Mobile", "?0")
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_ = Parse(r)
	}
}

func BenchmarkParseClientHints(b *testing.B) {
	h := make(http.Header)
	h.Set("Sec-CH-UA", `"Chromium";v="120", "Google Chrome";v="120"`)
	h.Set("Sec-CH-UA-Platform", `"Windows"`)
	h.Set("Sec-CH-UA-Platform-Version", `"10.0.0"`)
	h.Set("Sec-CH-UA-Mobile", "?0")
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_ = ParseClientHints(h)
	}
}
