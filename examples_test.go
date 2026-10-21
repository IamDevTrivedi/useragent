package useragent_test

import (
	"fmt"

	"github.com/IamDevTrivedi/useragent"
)

func ExampleParse() {
	ua := useragent.Parse("Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.6099.109 Safari/537.36")
	fmt.Println(ua.Browser.Name, ua.OS.Name, ua.Device.Type)
	// Output: Chrome Windows desktop
}

func ExampleUserAgent_String() {
	ua := useragent.Parse("Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36")
	fmt.Println(ua.String())
	// Output: Chrome 120.0.0.0 on Windows 10.0 (desktop)
}

func ExampleNew() {
	p := useragent.New(useragent.WithBrowserRules(useragent.BrowserRule{
		Name:        "MyClient",
		Match:       useragent.Match{Token: "MyClient"},
		VersionKeys: []string{"MyClient"},
	}))
	ua := p.Parse("MyClient/2.5")
	fmt.Println(ua.Browser.Name, ua.Browser.Version)
	// Output: MyClient 2.5
}

func ExampleParser() {
	p := useragent.New(useragent.WithOSRules(useragent.OSRule{
		Name:        "AcmeOS",
		Match:       useragent.Match{Token: "AcmeOS"},
		VersionKeys: []string{"AcmeOS"},
	}))
	ua := p.Parse("AcmeOS/3.2")
	fmt.Println(ua.OS.Name, ua.OS.Version)
	// Output: AcmeOS 3.2
}

func ExampleUserAgent_ApplyClientHints() {
	ua := useragent.Parse("Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36")
	ua.ApplyClientHints(useragent.ClientHints{
		FullVersionList: []useragent.Brand{{Name: "Google Chrome", Version: "120.0.6099.109"}},
		Mobile:          true,
		Platform:        "Android",
		PlatformVersion: "13.0.0",
		Model:           "Pixel 7",
	})
	fmt.Println(ua.Browser.Name, ua.Browser.Version, ua.OS.Name, ua.OS.Version, ua.Device.Type, ua.Device.Model)
	// Output: Chrome 120.0.6099.109 Android 13.0.0 mobile Pixel 7
}

func ExampleParseBrandList() {
	for _, brand := range useragent.ParseBrandList(`"Chromium";v="120", "Google Chrome";v="120"`) {
		fmt.Println(brand.Name, brand.Version)
	}
	// Output:
	// Chromium 120
	// Google Chrome 120
}

func ExampleVersion_Compare() {
	v := useragent.ParseVersion("120.0.6099.109")
	fmt.Println(v.AtLeast(120, 0, 0), v.Compare(useragent.Version{Major: 119}))
	// Output: true 1
}
