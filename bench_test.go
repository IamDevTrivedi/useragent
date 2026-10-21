package useragent

import (
	"strings"
	"testing"
)

var benchmarkUA UserAgent

const (
	benchmarkDesktop   = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.6099.109 Safari/537.36"
	benchmarkMobile    = "Mozilla/5.0 (Linux; Android 13; Pixel 7 Pro Build/TQ3A.230805.001; wv) AppleWebKit/537.36 (KHTML, like Gecko) Version/4.0 Chrome/120.0.6099.144 Mobile Safari/537.36"
	benchmarkBot       = "Mozilla/5.0 (compatible; Googlebot/2.1; +http://www.google.com/bot.html)"
	benchmarkMalformed = "() { ignored; }; echo Content-Type: text/plain"
)

var benchmarkLongInput = strings.Repeat(benchmarkDesktop+" ", 100)

func BenchmarkParseDesktop(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		benchmarkUA = DefaultParser.Parse(benchmarkDesktop)
	}
}

func BenchmarkParseMobile(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		benchmarkUA = DefaultParser.Parse(benchmarkMobile)
	}
}

func BenchmarkParseBot(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		benchmarkUA = DefaultParser.Parse(benchmarkBot)
	}
}

func BenchmarkParseParallel(b *testing.B) {
	b.ReportAllocs()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			benchmarkUA = DefaultParser.Parse(benchmarkDesktop)
		}
	})
}

func BenchmarkParseWithCustomRules(b *testing.B) {
	p := New(
		WithBrowserRules(BrowserRule{Name: "MyClient", Match: Match{Token: "MyClient"}, Priority: 10}),
		WithBotRules(BotRule{Name: "MyBot", Match: Match{Token: "MyBot"}, Category: BotOther, Priority: 10}),
	)
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		benchmarkUA = p.Parse(benchmarkDesktop)
	}
}

func BenchmarkParseMalformed(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		benchmarkUA = DefaultParser.Parse(benchmarkMalformed)
	}
}

func BenchmarkParseLongInput(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		benchmarkUA = DefaultParser.Parse(benchmarkLongInput)
	}
}

func BenchmarkTokenizer(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_ = tokenize(benchmarkDesktop)
	}
}

func BenchmarkVersionParse(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_ = ParseVersion("120.0.6099.109")
	}
}

func BenchmarkParseBrandList(b *testing.B) {
	const hints = `"Chromium";v="120", "Google Chrome";v="120", "Not_A Brand";v="8"`
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_ = ParseBrandList(hints)
	}
}

func BenchmarkApplyClientHints(b *testing.B) {
	base := Parse(benchmarkDesktop)
	h := ClientHints{
		FullVersionList: []Brand{{Name: "Google Chrome", Version: "120.0.6099.109"}},
		Platform:        "Windows",
		PlatformVersion: "10.0.0",
		Architecture:    "x86_64",
	}
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		ua := base
		ua.ApplyClientHints(h)
		benchmarkUA = ua
	}
}
