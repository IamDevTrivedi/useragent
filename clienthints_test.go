package useragent

import (
	"reflect"
	"testing"
)

func TestParseBrandList(t *testing.T) {
	tests := []struct {
		name string
		in   string
		want []Brand
	}{
		{"empty", "", nil},
		{"standard", `"Chromium";v="120", "Google Chrome";v="120", "Not_A Brand";v="8"`, []Brand{
			{"Chromium", "120"}, {"Google Chrome", "120"}, {"Not_A Brand", "8"},
		}},
		{"comma in brand", `"Not, A Brand";v="8", "Google Chrome";v="120"`, []Brand{
			{"Not, A Brand", "8"}, {"Google Chrome", "120"},
		}},
		{"missing version", `"Chromium", "Google Chrome";v="120"`, []Brand{
			{"Chromium", ""}, {"Google Chrome", "120"},
		}},
		{"malformed entries", `garbage, "A";v="1", "B";v="2", "unterminated`, []Brand{
			{"A", "1"}, {"B", "2"},
		}},
		{"spaces", `  "Chromium" ; v="120"  `, []Brand{{"Chromium", "120"}}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ParseBrandList(tt.in); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ParseBrandList(%q) = %+v, want %+v", tt.in, got, tt.want)
			}
		})
	}
}

func TestApplyClientHintsFullVersion(t *testing.T) {
	ua := Parse("Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36")
	ua.ApplyClientHints(ClientHints{
		Brands: []Brand{
			{Name: "Chromium", Version: "120"},
			{Name: "Google Chrome", Version: "120"},
			{Name: "Not_A Brand", Version: "8"},
		},
		FullVersionList: []Brand{
			{Name: "Chromium", Version: "120.0.6099.109"},
			{Name: "Google Chrome", Version: "120.0.6099.109"},
			{Name: "Not_A Brand", Version: "8.0.0.0"},
		},
		Platform:        "Windows",
		PlatformVersion: "10.0.0",
		Architecture:    "x86_64",
	})
	if ua.Browser.Name != Chrome || ua.Browser.Version != "120.0.6099.109" {
		t.Errorf("browser = %q %q", ua.Browser.Name, ua.Browser.Version)
	}
	if ua.Browser.VersionNo != (Version{120, 0, 6099}) {
		t.Errorf("VersionNo = %+v", ua.Browser.VersionNo)
	}
	if ua.OS.Name != Windows || ua.OS.Version != "10.0.0" {
		t.Errorf("OS = %q %q", ua.OS.Name, ua.OS.Version)
	}
	if ua.CPU.Architecture != "amd64" {
		t.Errorf("CPU = %q", ua.CPU.Architecture)
	}
}

func TestApplyClientHintsAndroidFrozenUA(t *testing.T) {
	ua := Parse("Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36")
	ua.ApplyClientHints(ClientHints{
		FullVersionList: []Brand{{Name: "Google Chrome", Version: "120.0.6099.109"}},
		Mobile:          true,
		Platform:        "Android",
		PlatformVersion: "13.0.0",
		Architecture:    "arm64",
		Model:           "Pixel 7",
	})
	if ua.Browser.Name != Chrome || ua.Browser.Version != "120.0.6099.109" {
		t.Errorf("browser = %q %q", ua.Browser.Name, ua.Browser.Version)
	}
	if ua.OS.Name != Android || ua.OS.Version != "13.0.0" {
		t.Errorf("OS = %q %q", ua.OS.Name, ua.OS.Version)
	}
	if ua.Device.Type != DeviceMobile || ua.Device.Model != "Pixel 7" {
		t.Errorf("device = %+v", ua.Device)
	}
	if ua.CPU.Architecture != "arm64" {
		t.Errorf("CPU = %q", ua.CPU.Architecture)
	}
}

func TestApplyClientHintsKeepsSafari(t *testing.T) {
	ua := Parse("Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/17.0 Safari/605.1.15")
	ua.ApplyClientHints(ClientHints{
		FullVersionList: []Brand{{Name: "Google Chrome", Version: "120.0.6099.109"}},
		Platform:        "macOS",
		PlatformVersion: "14.0.0",
	})
	if ua.Browser.Name != Safari || ua.Browser.Version != "17.0" {
		t.Errorf("Safari was replaced: %+v", ua.Browser)
	}
	if ua.OS.Name != MacOS || ua.OS.Version != "14.0.0" {
		t.Errorf("OS = %q %q", ua.OS.Name, ua.OS.Version)
	}
}

func TestApplyClientHintsPlatformConflict(t *testing.T) {
	ua := Parse("Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 Chrome/120.0.0.0 Safari/537.36")
	ua.ApplyClientHints(ClientHints{Platform: "Android", PlatformVersion: "13.0.0"})
	if ua.OS.Name != Windows || ua.OS.Version != "10.0" {
		t.Errorf("OS = %q %q", ua.OS.Name, ua.OS.Version)
	}
}

func TestApplyClientHintsMajorOnly(t *testing.T) {
	ua := Parse("")
	ua.ApplyClientHints(ClientHints{
		Brands: []Brand{{Name: "Google Chrome", Version: "120"}},
	})
	if ua.Browser.Name != Chrome || ua.Browser.Version != "120" {
		t.Errorf("browser = %q %q", ua.Browser.Name, ua.Browser.Version)
	}
}

func TestApplyClientHintsArchitecture(t *testing.T) {
	tests := []struct {
		in   string
		want string
	}{
		{"x86", "386"},
		{"i686", "386"},
		{"x86_64", "amd64"},
		{"amd64", "amd64"},
		{"arm", "arm"},
		{"arm64", "arm64"},
		{"aarch64", "arm64"},
		{"mips", "mips"},
	}
	for _, tt := range tests {
		ua := Parse("Twitterbot/1.0")
		ua.ApplyClientHints(ClientHints{Architecture: tt.in})
		if ua.CPU.Architecture != tt.want {
			t.Errorf("Architecture %q => %q, want %q", tt.in, ua.CPU.Architecture, tt.want)
		}
	}
}

func TestApplyClientHintsPlatforms(t *testing.T) {
	tests := []struct {
		in   string
		want string
	}{
		{"Windows", Windows},
		{"macOS", MacOS},
		{"Mac OS X", MacOS},
		{"Android", Android},
		{"iOS", IOS},
		{"Linux", Linux},
		{"Chrome OS", ChromeOS},
		{"ChromeOS", ChromeOS},
		{"Fuchsia", ""},
	}
	for _, tt := range tests {
		ua := Parse("")
		ua.ApplyClientHints(ClientHints{Platform: tt.in, PlatformVersion: "1.2.3"})
		if ua.OS.Name != tt.want {
			t.Errorf("Platform %q => %q, want %q", tt.in, ua.OS.Name, tt.want)
		}
		if tt.want != "" && ua.OS.Version != "1.2.3" {
			t.Errorf("Platform %q version = %q", tt.in, ua.OS.Version)
		}
	}
}

func TestApplyClientHintsEmptyAndNil(t *testing.T) {
	var ua *UserAgent
	ua.ApplyClientHints(ClientHints{Platform: "Android"}) // must not panic

	before := Parse("Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 Chrome/120.0.0.0 Safari/537.36")
	after := before
	after.ApplyClientHints(ClientHints{})
	if !reflect.DeepEqual(before, after) {
		t.Errorf("empty hints changed UA:\n%+v\n%+v", before, after)
	}
}

func TestApplyClientHintsSpecificBrowsers(t *testing.T) {
	tests := []struct {
		brand string
		want  string
	}{
		{"Microsoft Edge", Edge},
		{"Opera", Opera},
		{"Brave", Brave},
		{"Vivaldi", Vivaldi},
		{"Samsung Internet", SamsungBrowser},
		{"Chromium", Chrome},
		{"Android WebView", Chrome},
	}
	for _, tt := range tests {
		ua := Parse("Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 Chrome/120.0.0.0 Safari/537.36")
		ua.ApplyClientHints(ClientHints{Brands: []Brand{{Name: tt.brand, Version: "120"}}})
		if ua.Browser.Name != tt.want {
			t.Errorf("brand %q => %q, want %q", tt.brand, ua.Browser.Name, tt.want)
		}
	}
}
