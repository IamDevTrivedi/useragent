package useragent

import "strings"

// Brand is one entry from Sec-CH-UA or Sec-CH-UA-Full-Version-List.
type Brand struct {
	Name    string `json:"name"`
	Version string `json:"version"`
}

// ClientHints contains the subset of User-Agent Client Hints used by the
// parser. Empty fields are ignored. Callers that only have an http.Request can
// use the uahttp subpackage to populate this value.
type ClientHints struct {
	Brands          []Brand `json:"brands,omitempty"`
	FullVersionList []Brand `json:"full_version_list,omitempty"`
	Mobile          bool    `json:"mobile,omitempty"`
	Platform        string  `json:"platform,omitempty"`
	PlatformVersion string  `json:"platform_version,omitempty"`
	Architecture    string  `json:"architecture,omitempty"`
	Bitness         string  `json:"bitness,omitempty"`
	Model           string  `json:"model,omitempty"`
	Wow64           bool    `json:"wow64,omitempty"`
}

// ParseBrandList parses a Sec-CH-UA style brand list:
//
//	"Chromium";v="120", "Google Chrome";v="120", "Not_A Brand";v="8"
//
// Malformed entries are skipped. The returned slice may be nil.
func ParseBrandList(s string) []Brand {
	var out []Brand
	for len(s) != 0 {
		item := s
		if at := brandListComma(s); at >= 0 {
			item = s[:at]
			s = s[at+1:]
		} else {
			s = ""
		}
		item = strings.TrimSpace(item)
		if item == "" {
			continue
		}
		start := strings.IndexByte(item, '"')
		if start < 0 {
			continue
		}
		end := strings.IndexByte(item[start+1:], '"')
		if end < 0 {
			continue
		}
		brand := item[start+1 : start+1+end]
		rest := item[start+1+end+1:]
		version := ""
		if vi := indexFold(rest, `v="`); vi >= 0 {
			vrest := rest[vi+3:]
			if ve := strings.IndexByte(vrest, '"'); ve >= 0 {
				version = vrest[:ve]
			}
		}
		if brand != "" {
			out = append(out, Brand{Name: brand, Version: version})
		}
	}
	return out
}

func brandListComma(s string) int {
	quoted := false
	for i := 0; i < len(s); i++ {
		switch s[i] {
		case '"':
			quoted = !quoted
		case ',':
			if !quoted {
				return i
			}
		}
	}
	return -1
}

type brandName struct {
	brand string
	name  string
}

var clientHintBrands = [...]brandName{
	{"Google Chrome", Chrome},
	{"Microsoft Edge", Edge},
	{"Opera", Opera},
	{"Brave", Brave},
	{"Vivaldi", Vivaldi},
	{"Samsung Internet", SamsungBrowser},
	{"SamsungBrowser", SamsungBrowser},
	{"Chromium", Chromium},
	{"Android WebView", Chrome},
}

func chooseClientHintBrand(list []Brand) (name, version string) {
	for _, want := range clientHintBrands {
		for _, b := range list {
			if equalFold(b.Name, want.brand) {
				return want.name, b.Version
			}
		}
	}
	for _, b := range list {
		switch {
		case containsFold(b.Name, "Edge"):
			return Edge, b.Version
		case containsFold(b.Name, "Opera"):
			return Opera, b.Version
		case containsFold(b.Name, "Chrome"):
			return Chrome, b.Version
		case containsFold(b.Name, "Chromium"):
			return Chromium, b.Version
		}
	}
	return "", ""
}

// ApplyClientHints merges User-Agent Client Hints into ua. It is safe to call
// on a nil receiver. Client hints are treated as higher-quality version data
// than a frozen UA string, but an already-detected Safari/Firefox client is not
// replaced by a generic Chromium brand.
func (u *UserAgent) ApplyClientHints(h ClientHints) {
	if u == nil {
		return
	}

	list := h.FullVersionList
	full := true
	if len(list) == 0 {
		list = h.Brands
		full = false
	}
	brandName, brandVersion := chooseClientHintBrand(list)
	if brandName != "" && shouldApplyHintBrowser(u.Browser.Name, brandName) {
		u.Browser.Name = brandName
	}
	if brandVersion != "" && u.Browser.Name == brandName {
		if full || u.Browser.Version == "" {
			u.Browser.Version = normalizeVersion(brandVersion)
			u.Browser.VersionNo = ParseVersion(u.Browser.Version)
		}
	}

	if h.Platform != "" {
		p := normalizeHintPlatform(h.Platform)
		if p != "" {
			if u.OS.Name == "" || u.OS.Name == Linux || u.OS.Name == p {
				u.OS.Name = p
				if h.PlatformVersion != "" {
					u.OS.Version = normalizeVersion(h.PlatformVersion)
					u.OS.VersionNo = ParseVersion(u.OS.Version)
				}
			}
		}
	}

	if h.Mobile {
		if u.Device.Type == DeviceUnknown || u.Device.Type == DeviceDesktop {
			u.Device.Type = DeviceMobile
		}
	}
	if h.Model != "" {
		u.Device.Model = h.Model
		if u.Device.Type == DeviceUnknown || u.Device.Type == DeviceDesktop {
			u.Device.Type = DeviceMobile
		}
	}
	if h.Architecture != "" {
		// Client hints are more authoritative than the frozen UA string.
		u.CPU.Architecture = normalizeHintArch(h.Architecture)
	}
}

func shouldApplyHintBrowser(current, hinted string) bool {
	if current == "" || current == AppleWebKit || current == Chromium {
		return true
	}
	if current == Chrome {
		return hinted == Chrome || hinted == Edge || hinted == Opera ||
			hinted == Brave || hinted == Vivaldi || hinted == SamsungBrowser
	}
	return isBlinkBrowser(current) && hinted != Chromium
}

func normalizeHintPlatform(s string) string {
	switch {
	case equalFold(s, "Windows"):
		return Windows
	case equalFold(s, "macOS"), equalFold(s, "Mac OS X"):
		return MacOS
	case equalFold(s, "Android"):
		return Android
	case equalFold(s, "iOS"):
		return IOS
	case equalFold(s, "Linux"):
		return Linux
	case equalFold(s, "Chrome OS"), equalFold(s, "ChromeOS"):
		return ChromeOS
	default:
		return ""
	}
}

func normalizeHintArch(s string) string {
	switch {
	case equalFold(s, "x86"), equalFold(s, "i386"), equalFold(s, "i686"):
		return "386"
	case equalFold(s, "x86_64"), equalFold(s, "amd64"), equalFold(s, "x64"):
		return "amd64"
	case equalFold(s, "arm"):
		return "arm"
	case equalFold(s, "arm64"), equalFold(s, "aarch64"):
		return "arm64"
	default:
		return s
	}
}
