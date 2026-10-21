package uahttp

import (
	"net/http"
	"testing"
)

func FuzzParseRequest(f *testing.F) {
	f.Add(
		"Mozilla/5.0 (Windows NT 10.0; Win64; x64) Chrome/120.0.0.0 Safari/537.36",
		`"Chromium";v="120", "Google Chrome";v="120"`,
		`"Chromium";v="120.0.6099.109", "Google Chrome";v="120.0.6099.109"`,
		`"Windows"`,
		`"10.0.0"`,
		`"x86"`,
		`"64"`,
		`""`,
		"?0",
		"?0",
	)
	f.Add("", "", "", "", "", "", "", "", "", "")
	f.Fuzz(func(t *testing.T, rawUA, brands, fullVersionList, platform, platformVersion, arch, bitness, model, mobile, wow64 string) {
		r := &http.Request{Header: make(http.Header)}
		r.Header.Set("User-Agent", rawUA)
		r.Header.Set("Sec-CH-UA", brands)
		r.Header.Set("Sec-CH-UA-Full-Version-List", fullVersionList)
		r.Header.Set("Sec-CH-UA-Platform", platform)
		r.Header.Set("Sec-CH-UA-Platform-Version", platformVersion)
		r.Header.Set("Sec-CH-UA-Arch", arch)
		r.Header.Set("Sec-CH-UA-Bitness", bitness)
		r.Header.Set("Sec-CH-UA-Model", model)
		r.Header.Set("Sec-CH-UA-Mobile", mobile)
		r.Header.Set("Sec-CH-UA-WoW64", wow64)

		got := Parse(r)
		if got.Raw != rawUA {
			t.Fatalf("Raw = %q, want %q", got.Raw, rawUA)
		}
		_ = got.String()
		_ = got.IsUnknown()
	})
}
