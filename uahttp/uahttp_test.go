package uahttp_test

import (
	"net/http"
	"reflect"
	"testing"

	ua "github.com/IamDevTrivedi/useragent"
	"github.com/IamDevTrivedi/useragent/uahttp"
)

func request(userAgent string) *http.Request {
	r := &http.Request{Header: http.Header{}}
	if userAgent != "" {
		r.Header.Set("User-Agent", userAgent)
	}
	return r
}

func TestParse(t *testing.T) {
	r := request("Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36")
	got := uahttp.Parse(r)
	if got.Browser.Name != ua.Chrome || got.OS.Name != ua.Windows || got.Device.Type != ua.DeviceDesktop {
		t.Errorf("unexpected parse result: %+v", got)
	}
}

func TestParseNilRequest(t *testing.T) {
	if got := uahttp.Parse(nil); got != (ua.UserAgent{}) {
		t.Errorf("Parse(nil) = %+v, want zero", got)
	}
}

func TestParseClientHints(t *testing.T) {
	h := http.Header{}
	h.Set("Sec-CH-UA", `"Chromium";v="120", "Google Chrome";v="120", "Not_A Brand";v="8"`)
	h.Set("Sec-CH-UA-Full-Version-List", `"Chromium";v="120.0.6099.109", "Google Chrome";v="120.0.6099.109", "Not_A Brand";v="8.0.0.0"`)
	h.Set("Sec-CH-UA-Mobile", "?1")
	h.Set("Sec-CH-UA-Platform", `"Android"`)
	h.Set("Sec-CH-UA-Platform-Version", `"13.0.0"`)
	h.Set("Sec-CH-UA-Arch", `"arm64"`)
	h.Set("Sec-CH-UA-Bitness", `"64"`)
	h.Set("Sec-CH-UA-Model", `"Pixel 7"`)
	h.Set("Sec-CH-UA-WoW64", "?0")
	got := uahttp.ParseClientHints(h)
	if len(got.Brands) != 3 || len(got.FullVersionList) != 3 {
		t.Errorf("brand lists = %+v / %+v", got.Brands, got.FullVersionList)
	}
	if !got.Mobile || got.Platform != "Android" || got.PlatformVersion != "13.0.0" ||
		got.Architecture != "arm64" || got.Bitness != "64" ||
		got.Model != "Pixel 7" || got.Wow64 {
		t.Errorf("hints = %+v", got)
	}
}

func TestParseClientHintsNil(t *testing.T) {
	if got := uahttp.ParseClientHints(nil); !reflect.DeepEqual(got, ua.ClientHints{}) {
		t.Errorf("ParseClientHints(nil) = %+v", got)
	}
}

func TestParseFrozenAndroidWithHints(t *testing.T) {
	r := request("Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36")
	r.Header.Set("Sec-CH-UA-Full-Version-List", `"Chromium";v="120.0.6099.109", "Google Chrome";v="120.0.6099.109", "Not_A Brand";v="8.0.0.0"`)
	r.Header.Set("Sec-CH-UA-Mobile", "?1")
	r.Header.Set("Sec-CH-UA-Platform", `"Android"`)
	r.Header.Set("Sec-CH-UA-Platform-Version", `"13.0.0"`)
	r.Header.Set("Sec-CH-UA-Arch", `"arm64"`)
	r.Header.Set("Sec-CH-UA-Model", `"Pixel 7"`)
	got := uahttp.Parse(r)
	if got.Browser.Name != ua.Chrome || got.Browser.Version != "120.0.6099.109" {
		t.Errorf("browser = %+v", got.Browser)
	}
	if got.OS.Name != ua.Android || got.OS.Version != "13.0.0" {
		t.Errorf("os = %+v", got.OS)
	}
	if got.Device.Type != ua.DeviceMobile || got.Device.Model != "Pixel 7" {
		t.Errorf("device = %+v", got.Device)
	}
	if got.CPU.Architecture != "arm64" {
		t.Errorf("cpu = %+v", got.CPU)
	}
}

func TestParseSafariKeepsBrowser(t *testing.T) {
	r := request("Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/17.0 Safari/605.1.15")
	r.Header.Set("Sec-CH-UA-Full-Version-List", `"Google Chrome";v="120.0.6099.109"`)
	r.Header.Set("Sec-CH-UA-Platform", `"macOS"`)
	got := uahttp.Parse(r)
	if got.Browser.Name != ua.Safari || got.Browser.Version != "17.0" {
		t.Errorf("Safari browser was replaced: %+v", got.Browser)
	}
}

func TestParseMalformedHints(t *testing.T) {
	r := request("Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 Chrome/120.0.0.0 Safari/537.36")
	r.Header.Set("Sec-CH-UA", `"unterminated`)
	r.Header.Set("Sec-CH-UA-Platform", `"Fuchsia"`)
	r.Header.Set("Sec-CH-UA-Platform-Version", `"1.0"`)
	got := uahttp.Parse(r)
	if got.Browser.Name != ua.Chrome || got.OS.Name != ua.Windows {
		t.Errorf("malformed hints changed the result: %+v", got)
	}
}

func TestParseWithCustomParser(t *testing.T) {
	p := ua.New(ua.WithBrowserRules(ua.BrowserRule{
		Name:        "MyClient",
		Match:       ua.Match{Token: "MyClient"},
		VersionKeys: []string{"MyClient"},
	}))
	r := request("MyClient/3.4")
	got := uahttp.ParseWith(p, r)
	if got.Browser.Name != "MyClient" || got.Browser.Version != "3.4" {
		t.Errorf("custom parser result: %+v", got.Browser)
	}
	chrome := request("Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 Chrome/120.0.0.0 Safari/537.36")
	if got := uahttp.ParseWith(nil, chrome); got.Browser.Name != ua.Chrome {
		t.Errorf("nil parser should use DefaultParser: %+v", got.Browser)
	}
}

func TestParseEmptyRequest(t *testing.T) {
	got := uahttp.Parse(request(""))
	if got.Raw != "" || !got.IsUnknown() {
		t.Errorf("empty request result: %+v", got)
	}
}
