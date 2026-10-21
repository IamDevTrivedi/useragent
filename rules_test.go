package useragent

import (
	"sync"
	"testing"
)

func TestCustomBrowserRule(t *testing.T) {
	p := New(WithBrowserRules(BrowserRule{
		Name:        "MyClient",
		Match:       Match{Token: "MyClient"},
		VersionKeys: []string{"MyClient"},
		Product:     "MyClient",
	}))
	ua := p.Parse("MyClient/2.5 (X11; Linux x86_64)")
	if ua.Browser.Name != "MyClient" || ua.Browser.Version != "2.5" {
		t.Errorf("browser = %q %q", ua.Browser.Name, ua.Browser.Version)
	}
	if ua.Browser.VersionNo != (Version{2, 5, 0}) {
		t.Errorf("VersionNo = %+v", ua.Browser.VersionNo)
	}
	if ua.OS.Name != Linux || ua.Device.Type != DeviceDesktop {
		t.Errorf("os/device = %q/%v", ua.OS.Name, ua.Device.Type)
	}
}

func TestCustomRulePriority(t *testing.T) {
	p := New(WithBrowserRules(
		BrowserRule{Name: "Low", Match: Match{Token: "MyClient"}, Priority: 1},
		BrowserRule{Name: "High", Match: Match{Token: "MyClient"}, Priority: 10},
		BrowserRule{Name: "Middle", Match: Match{Token: "MyClient"}, Priority: 5},
	))
	ua := p.Parse("MyClient/1.0")
	if ua.Browser.Name != "High" {
		t.Errorf("Browser.Name = %q, want High", ua.Browser.Name)
	}
}

func TestCustomMatchForms(t *testing.T) {
	raw := "X-PrefixThing/1.2 Other/2"
	tests := []struct {
		name  string
		match Match
		want  bool
	}{
		{"token", Match{Token: "X-PrefixThing"}, true},
		{"prefix", Match{Prefix: "X-Prefix"}, true},
		{"suffix", Match{Suffix: "Thing"}, true},
		{"contains", Match{Contains: "Prefix"}, true},
		{"contains value", Match{Contains: "1.2"}, true},
		{"all", Match{All: []string{"X-PrefixThing", "Other"}}, true},
		{"all missing", Match{All: []string{"X-PrefixThing", "Nope"}}, false},
		{"any", Match{Any: []string{"Nope", "Other"}}, true},
		{"any missing", Match{Any: []string{"Nope", "Nada"}}, false},
		{"none", Match{Token: "X-PrefixThing", None: []string{"Other"}}, false},
		{"none absent", Match{Token: "X-PrefixThing", None: []string{"Nope"}}, true},
		{"none only absent", Match{None: []string{"Nope"}}, true},
		{"none only present", Match{None: []string{"X-PrefixThing"}}, false},
		{"zero", Match{}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := New(WithBrowserRules(BrowserRule{Name: "Matched", Match: tt.match}))
			ua := p.Parse(raw)
			got := ua.Browser.Name == "Matched"
			if got != tt.want {
				t.Errorf("match = %v, want %v (browser %q)", got, tt.want, ua.Browser.Name)
			}
		})
	}
}

func TestCustomOSRuleOverridesBuiltIn(t *testing.T) {
	p := New(WithOSRules(OSRule{
		Name:        "Plan9",
		Match:       Match{Token: "Windows NT"},
		VersionKeys: []string{"Windows NT"},
	}))
	ua := p.Parse("Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 Chrome/120.0.0.0 Safari/537.36")
	if ua.OS.Name != "Plan9" || ua.OS.Version != "10.0" {
		t.Errorf("OS = %q %q", ua.OS.Name, ua.OS.Version)
	}
}

func TestCustomBotRule(t *testing.T) {
	p := New(WithBotRules(BotRule{
		Name:        "MyBot",
		Match:       Match{Token: "MyBot"},
		VersionKeys: []string{"MyBot"},
		Category:    BotOther,
	}))
	ua := p.Parse("MyBot/1.2 (+https://example.com/bot)")
	if !ua.Bot.Detected || ua.Bot.Name != "MyBot" || ua.Bot.Version != "1.2" {
		t.Errorf("bot = %+v", ua.Bot)
	}
	if ua.Bot.URL != "https://example.com/bot" {
		t.Errorf("Bot.URL = %q", ua.Bot.URL)
	}
	if ua.Bot.Category != BotOther {
		t.Errorf("Bot.Category = %v", ua.Bot.Category)
	}
	if ua.Browser.Name != "MyBot" {
		t.Errorf("Browser.Name = %q, want MyBot", ua.Browser.Name)
	}
}

func TestCustomBotURLOverride(t *testing.T) {
	p := New(WithBotRules(BotRule{
		Name:     "MyBot",
		Match:    Match{Token: "MyBot"},
		Category: BotOther,
		URL:      "https://override.example/",
	}))
	ua := p.Parse("MyBot/1.0 (+https://example.com/original)")
	if ua.Bot.URL != "https://override.example/" {
		t.Errorf("Bot.URL = %q", ua.Bot.URL)
	}
}

func TestCustomEngineRule(t *testing.T) {
	p := New(WithEngineRules(EngineRule{
		Name:        "EngineX",
		Match:       Match{Token: "EngineX"},
		VersionKeys: []string{"EngineX"},
	}))
	ua := p.Parse("Mozilla/5.0 (X11; Linux x86_64) Chrome/120.0.0.0 EngineX/9.9")
	if ua.Engine.Name != "EngineX" || ua.Engine.Version != "9.9" {
		t.Errorf("engine = %q %q", ua.Engine.Name, ua.Engine.Version)
	}
}

func TestCustomDeviceRule(t *testing.T) {
	p := New(WithDeviceRules(DeviceRule{
		Type:   DeviceWearable,
		Vendor: "Acme",
		Model:  "X1",
		Match:  Match{Token: "MyPhone"},
	}))
	ua := p.Parse("Mozilla/5.0 (Windows NT 10.0; Win64; x64) MyPhone/1.0")
	if ua.Device.Type != DeviceWearable || ua.Device.Vendor != "Acme" || ua.Device.Model != "X1" {
		t.Errorf("device = %+v", ua.Device)
	}
}

func TestCustomRulesOverrideBuiltIns(t *testing.T) {
	raw := "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36"
	p := New(
		WithBrowserRules(BrowserRule{Name: "MyChrome", Match: Match{Token: "Chrome"}}),
		WithOSRules(OSRule{Name: "MyOS", Match: Match{Token: "Windows NT"}}),
		WithEngineRules(EngineRule{Name: "MyEngine", Match: Match{Token: "Chrome"}}),
		WithDeviceRules(DeviceRule{Type: DeviceMobile, Match: Match{Token: "Chrome"}}),
	)
	ua := p.Parse(raw)
	if ua.Browser.Name != "MyChrome" || ua.OS.Name != "MyOS" ||
		ua.Engine.Name != "MyEngine" || ua.Device.Type != DeviceMobile {
		t.Errorf("custom overrides failed: %+v", ua)
	}
}

func TestParserZeroValueAndDefault(t *testing.T) {
	var p Parser
	ua := p.Parse("Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 Chrome/120.0.0.0 Safari/537.36")
	if ua.Browser.Name != Chrome {
		t.Errorf("zero parser browser = %q", ua.Browser.Name)
	}
	if got, want := Parse("Twitterbot/1.0"), DefaultParser.Parse("Twitterbot/1.0"); got != want {
		t.Errorf("package Parse differs from DefaultParser.Parse: %+v vs %+v", got, want)
	}
	p2 := New(nil)
	if p2.Parse("Twitterbot/1.0").Bot.Name != Twitterbot {
		t.Error("nil option changed parser behavior")
	}
}

func TestParserConcurrentUse(t *testing.T) {
	p := New(WithBrowserRules(BrowserRule{Name: "MyClient", Match: Match{Token: "MyClient"}}))
	raw := "MyClient/1.0 (X11; Linux x86_64)"
	var wg sync.WaitGroup
	for i := 0; i < 32; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			ua := p.Parse(raw)
			if ua.Browser.Name != "MyClient" {
				t.Errorf("Browser.Name = %q", ua.Browser.Name)
			}
		}()
	}
	wg.Wait()
}
