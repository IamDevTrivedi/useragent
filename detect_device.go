package useragent

import "strings"

func detectDevice(ua *UserAgent, ts *tokenSet) {
	if ua.Device.Type != DeviceUnknown {
		return
	}

	// Desktop operating systems cannot be a TV, wearable, or non-Xbox
	// console. Short-circuiting here keeps the many TV/wearable markers out
	// of the common desktop path.
	switch ua.OS.Name {
	case Windows:
		if ts.has("Xbox") {
			ua.Device = consoleDevice(ts)
		} else {
			ua.Device.Type = DeviceDesktop
		}
		return
	case MacOS, ChromeOS, FreeBSD, OpenBSD, NetBSD, Solaris:
		ua.Device.Type = DeviceDesktop
		return
	}

	if isConsole(ts) {
		ua.Device = consoleDevice(ts)
		return
	}
	if isTV(ts) {
		vendor, model := "", ""
		if ua.OS.Name == Android {
			vendor, model = androidModel(ts)
		}
		if v := tvVendor(ts); v != "" {
			vendor = v
		}
		if model == "" {
			model = tvModel(ts)
		}
		ua.Device = Device{Type: DeviceTV, Vendor: vendor, Model: model}
		return
	}
	if isWearable(ts) {
		ua.Device = wearableDevice(ts)
		return
	}

	switch ua.OS.Name {
	case IOS:
		switch {
		case ts.contains("iPad"):
			ua.Device = Device{Type: DeviceTablet, Vendor: "Apple", Model: "iPad"}
		case ts.contains("iPod"):
			ua.Device = Device{Type: DeviceMobile, Vendor: "Apple", Model: "iPod"}
		case ts.contains("iPhone"):
			ua.Device = Device{Type: DeviceMobile, Vendor: "Apple", Model: "iPhone"}
		default:
			ua.Device = Device{Type: DeviceMobile, Vendor: "Apple"}
		}
		return
	case Android:
		vendor, model := androidModel(ts)
		typ := DeviceMobile
		if isAndroidTablet(ts, model) {
			typ = DeviceTablet
		}
		ua.Device = Device{Type: typ, Vendor: vendor, Model: model}
		return
	case WindowsPhone, KaiOS, BlackBerry, HarmonyOS:
		ua.Device.Type = DeviceMobile
		return
	case Tizen, WebOS:
		ua.Device.Type = DeviceMobile
		return
	case Windows, MacOS, Linux, ChromeOS, FreeBSD, OpenBSD, NetBSD, Solaris:
		ua.Device.Type = DeviceDesktop
		return
	}

	switch {
	case ts.contains("Tablet") || ts.contains("iPad"):
		ua.Device.Type = DeviceTablet
	case ts.contains("Mobile") || isMobileBrowser(ua.Browser.Name):
		ua.Device.Type = DeviceMobile
	case ts.contains("SonyEricsson") || ts.contains("NetFront"):
		ua.Device.Type = DeviceMobile
	}
}

func isConsole(ts *tokenSet) bool {
	return ts.contains("PlayStation") || ts.contains("Xbox") || ts.contains("Nintendo")
}

func consoleDevice(ts *tokenSet) Device {
	switch {
	case ts.contains("PlayStation"):
		model := "PlayStation"
		for _, m := range [...]string{"PlayStation 5", "PlayStation 4", "PlayStation Vita", "PlayStation 3"} {
			if ts.contains(m) {
				model = m
				break
			}
		}
		return Device{Type: DeviceConsole, Vendor: "Sony", Model: model}
	case ts.contains("Xbox"):
		return Device{Type: DeviceConsole, Vendor: "Microsoft", Model: "Xbox"}
	default:
		model := "Nintendo"
		for _, m := range [...]string{"Nintendo Switch", "WiiU", "Wii", "3DS"} {
			if ts.contains(m) {
				model = m
				break
			}
		}
		return Device{Type: DeviceConsole, Vendor: "Nintendo", Model: model}
	}
}

func isWearable(ts *tokenSet) bool {
	return ts.contains("WatchOS") || ts.contains("Watch OS") ||
		ts.contains("Wear OS") || ts.contains("Apple Watch") ||
		ts.contains("SM-R") || ts.contains("Gear S")
}

func wearableDevice(ts *tokenSet) Device {
	switch {
	case ts.contains("Apple Watch") || ts.contains("WatchOS") || ts.contains("Watch OS"):
		return Device{Type: DeviceWearable, Vendor: "Apple", Model: "Apple Watch"}
	case ts.contains("SM-R") || ts.contains("Gear S"):
		return Device{Type: DeviceWearable, Vendor: "Samsung", Model: "Gear"}
	default:
		return Device{Type: DeviceWearable, Vendor: "Google"}
	}
}

func isTV(ts *tokenSet) bool {
	if ts.has("TV") || ts.has("TV Safari") {
		return true
	}
	markers := [...]string{
		"SMART-TV", "SmartTV", "Smart TV", "HbbTV", "GoogleTV", "Google TV",
		"AppleTV", "Web0S", "webOS.TV", "NetCast", "Roku", "FireTV", "Fire TV",
		"BRAVIA", "Android TV", "ATV", "CrKey", "Chromecast", "SmartHub",
		"Maple", "DLNA", "Large Screen", "Samsung TV", "LG Browser", "Viera",
		"Aquos", "MiBOX", "SHIELD",
	}
	for _, m := range markers {
		if ts.contains(m) {
			return true
		}
	}
	return false
}

func tvVendor(ts *tokenSet) string {
	switch {
	case ts.contains("Samsung"), ts.contains("SMART-TV"), ts.contains("Tizen"):
		return "Samsung"
	case ts.contains("LG"), ts.contains("Web0S"), ts.contains("NetCast"):
		return "LG"
	case ts.contains("Sony"), ts.contains("BRAVIA"):
		return "Sony"
	case ts.contains("Roku"):
		return "Roku"
	case ts.contains("FireTV"), ts.contains("Fire TV"):
		return "Amazon"
	case ts.contains("AppleTV"):
		return "Apple"
	case ts.contains("GoogleTV"), ts.contains("Google TV"), ts.contains("CrKey"),
		ts.contains("Chromecast"), ts.contains("Android TV"):
		return "Google"
	default:
		return ""
	}
}

func tvModel(ts *tokenSet) string {
	for _, m := range [...]string{
		"BRAVIA 4K GB ATV3", "BRAVIA", "Chromecast", "SHIELD Android TV",
		"Apple TV", "Roku", "Fire TV", "SmartTV",
	} {
		if ts.contains(m) {
			return m
		}
	}
	return ""
}

func isMobileBrowser(name string) bool {
	switch name {
	case MobileSafari, OperaMini, OperaTouch, UCBrowser, QQBrowser, HeyTapBrowser,
		MiuiBrowser, HuaweiBrowser, SamsungBrowser, NetFront, AndroidBrowser,
		Brave, Chrome, Firefox, Safari, CocCoc, YandexBrowser, BaiduBrowser,
		SogouBrowser, QuarkBrowser, Browser360:
		return true
	default:
		return false
	}
}

func isAndroidTablet(ts *tokenSet, model string) bool {
	markers := [...]string{
		"Tablet", "TABLET", "iPad", "SM-T", "SM-X", "SM-P", "GT-P", "TB-", "Nexus 7",
		"Nexus 9", "Pixel C", "MediaPad", "MatePad", "MI PAD", "Mi Pad",
		"Lenovo TAB", "Transformer", "Xoom",
	}
	for _, m := range markers {
		if ts.contains(m) || containsFold(model, m) {
			return true
		}
	}
	return false
}

func androidModel(ts *tokenSet) (vendor, model string) {
	if ts.cn == 0 {
		return "", ""
	}
	start := -1
	for i := 0; i < ts.cn; i++ {
		c := ts.comments[i]
		if hasPrefixFold(c, "Android ") || equalFold(c, "Android") {
			start = i
			break
		}
	}
	if start < 0 && (ts.has("Android") || ts.contains("Android App")) {
		start = 0
	}
	if start < 0 {
		return "", ""
	}

	for i := start; i < ts.cn; i++ {
		c := ts.comments[i]
		if i == start && (hasPrefixFold(c, "Android ") || equalFold(c, "Android")) {
			continue
		}
		v, m, ok := parseAndroidModel(c, ts)
		if ok {
			return v, m
		}
	}
	return "", ""
}

func parseAndroidModel(c string, ts *tokenSet) (string, string, bool) {
	c = cleanModel(c)
	if c == "" || isNoiseComment(c) {
		return "", "", false
	}
	if at := strings.IndexByte(c, '/'); at > 0 {
		c = strings.TrimSpace(c[:at])
	}
	if c == "" || isNoiseComment(c) {
		return "", "", false
	}
	if len(c) < 2 {
		return "", "", false
	}

	vendor, model := splitVendorModel(c)
	if vendor == "" {
		vendor = inferVendor(c, ts)
	}
	if model == "" {
		return "", "", false
	}
	if !hasDigit(model) && vendor == "" && len(model) < 4 {
		return "", "", false
	}
	if isNoiseComment(model) {
		return "", "", false
	}
	return vendor, strings.TrimSpace(model), true
}

func splitVendorModel(c string) (vendor, model string) {
	type prefix struct {
		prefix string
		vendor string
		strip  bool
	}
	prefixes := [...]prefix{
		{"SAMSUNG ", "Samsung", true},
		{"Samsung ", "Samsung", true},
		{"HUAWEI ", "Huawei", true},
		{"Huawei ", "Huawei", true},
		{"HONOR ", "Honor", true},
		{"Honor ", "Honor", true},
		{"Xiaomi ", "Xiaomi", true},
		{"XIAOMI ", "Xiaomi", true},
		{"OPPO ", "OPPO", true},
		{"Oppo ", "OPPO", true},
		{"vivo ", "vivo", true},
		{"Vivo ", "vivo", true},
		{"OnePlus ", "OnePlus", true},
		{"ONEPLUS ", "OnePlus", true},
		{"Sony ", "Sony", true},
		{"SonyEricsson", "Sony", true},
		{"Realme ", "Realme", true},
		{"realme ", "Realme", true},
		{"Lenovo ", "Lenovo", true},
		{"ASUS ", "Asus", true},
		{"Asus ", "Asus", true},
		{"VIVAX ", "VIVAX", true},
		{"Vivax ", "VIVAX", true},
		{"Pixel ", "Google", false},
		{"Nexus ", "Google", false},
		{"Redmi ", "Xiaomi", false},
		{"POCO ", "Xiaomi", false},
		{"Poco ", "Xiaomi", false},
		{"Mi ", "Xiaomi", false},
		{"MI ", "Xiaomi", false},
		{"Moto ", "Motorola", false},
		{"Nokia ", "Nokia", false},
		{"LG ", "LG", false},
		{"HTC ", "HTC", false},
		{"ZTE ", "ZTE", false},
		{"TCL ", "TCL", false},
	}
	for _, p := range prefixes {
		if hasPrefixFold(c, p.prefix) {
			if p.strip {
				return p.vendor, strings.TrimSpace(c[len(p.prefix):])
			}
			return p.vendor, c
		}
	}
	return "", c
}

func inferVendor(model string, ts *tokenSet) string {
	switch {
	case hasPrefixFold(model, "SM-"), hasPrefixFold(model, "GT-"),
		hasPrefixFold(model, "SGH-"), hasPrefixFold(model, "SCH-"):
		return "Samsung"
	case hasPrefixFold(model, "LM-"), hasPrefixFold(model, "LG-"):
		return "LG"
	case hasPrefixFold(model, "CPH"):
		return "OPPO"
	case hasPrefixFold(model, "RMX"):
		return "Realme"
	case hasPrefixFold(model, "VOG-"), hasPrefixFold(model, "ELE-"),
		hasPrefixFold(model, "ANE-"), hasPrefixFold(model, "JNY-"),
		hasPrefixFold(model, "MAR-"):
		return "Huawei"
	case hasPrefixFold(model, "XT"):
		return "Motorola"
	case hasPrefixFold(model, "TB-"):
		return "Lenovo"
	case hasPrefixFold(model, "Pixel"), hasPrefixFold(model, "Nexus"):
		return "Google"
	case hasPrefixFold(model, "Redmi"), hasPrefixFold(model, "POCO"),
		hasPrefixFold(model, "Mi "), hasPrefixFold(model, "MI "),
		hasPrefixFold(model, "M2"):
		return "Xiaomi"
	case hasPrefixFold(model, "BRAVIA"), hasPrefixFold(model, "Xperia"):
		return "Sony"
	}
	switch {
	case ts.contains("SAMSUNG"):
		return "Samsung"
	case ts.contains("HUAWEI"):
		return "Huawei"
	case ts.contains("Xiaomi"):
		return "Xiaomi"
	case ts.contains("OPPO"):
		return "OPPO"
	case ts.contains("vivo"):
		return "vivo"
	case ts.contains("realme"):
		return "Realme"
	case ts.contains("Nokia"):
		return "Nokia"
	case ts.contains("Sony"):
		return "Sony"
	}
	return ""
}

var noiseComments = [...]string{
	"", "u", "wv", "mobile", "tablet", "tv", "phone", "linux", "x11",
	"windows", "macintosh", "android", "compatible", "build", "hybrid()",
	"globaldist", "googleplaystore", "watch", "cros", "khtml, like gecko",
}

func isNoiseComment(c string) bool {
	if len(c) > 48 {
		return true
	}
	for _, word := range noiseComments {
		if equalFold(c, word) {
			return true
		}
	}
	if len(c) >= 2 && lowerASCII(c[0]) == 'u' && c[1] == ';' {
		return true
	}
	if hasPrefixFold(c, "build/") || hasPrefixFold(c, "version/") ||
		hasPrefixFold(c, "name:") {
		return true
	}
	if containsFold(c, "http") || containsFold(c, "globaldist") ||
		containsFold(c, "googleplaystore") {
		return true
	}
	if isNumericLike(c) {
		// A bare short number is an OS version; 3-8 digit numbers are often
		// device model codes (for example 8092).
		return len(c) < 3 || strings.IndexByte(c, '.') >= 0
	}
	if isResolution(c) || isMemory(c) || isLanguage(c) {
		return true
	}
	switch {
	case containsFold(c, "chrome/"), containsFold(c, "safari/"),
		containsFold(c, "applewebkit"), containsFold(c, "gecko/"),
		containsFold(c, "opera/"), containsFold(c, "edg/"):
		return true
	}
	return false
}

func isNumericLike(s string) bool {
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

func isResolution(s string) bool {
	seen := false
	for i := 0; i < len(s); i++ {
		switch {
		case isDigit(s[i]):
			seen = true
		case s[i] == 'x' || s[i] == 'X':
			if !seen {
				return false
			}
		default:
			return false
		}
	}
	return seen
}

func isMemory(s string) bool {
	if len(s) < 3 {
		return false
	}
	suffix := s[len(s)-2:]
	if !equalFold(suffix, "mb") && !equalFold(suffix, "gb") {
		return false
	}
	prefix := s[:len(s)-2]
	for i := 0; i < len(prefix); i++ {
		if !isDigit(prefix[i]) && prefix[i] != '.' {
			return false
		}
	}
	return prefix != ""
}

func isLanguage(s string) bool {
	if len(s) < 2 || len(s) > 6 {
		return false
	}
	dash := false
	for i := 0; i < len(s); i++ {
		c := s[i]
		switch {
		case c == '-':
			dash = true
		case c >= 'a' && c <= 'z', c >= 'A' && c <= 'Z':
		default:
			return false
		}
	}
	return dash && !hasDigit(s)
}

func hasDigit(s string) bool {
	for i := 0; i < len(s); i++ {
		if isDigit(s[i]) {
			return true
		}
	}
	return false
}
