package useragent

import (
	"encoding/json"
	"os"
	"testing"
)

type conformanceCase struct {
	Name           string       `json:"name"`
	UA             string       `json:"ua"`
	Browser        string       `json:"browser,omitempty"`
	BrowserEmpty   bool         `json:"browserEmpty,omitempty"`
	BrowserVersion string       `json:"browserVersion,omitempty"`
	OS             string       `json:"os,omitempty"`
	OSEmpty        bool         `json:"osEmpty,omitempty"`
	OSVersion      string       `json:"osVersion,omitempty"`
	Device         *DeviceType  `json:"device,omitempty"`
	Vendor         string       `json:"vendor,omitempty"`
	Model          string       `json:"model,omitempty"`
	CPU            string       `json:"cpu,omitempty"`
	Engine         string       `json:"engine,omitempty"`
	BotDetected    *bool        `json:"botDetected,omitempty"`
	BotName        string       `json:"botName,omitempty"`
	BotVersion     string       `json:"botVersion,omitempty"`
	BotURL         string       `json:"botURL,omitempty"`
	BotCategory    *BotCategory `json:"botCategory,omitempty"`
	App            string       `json:"app,omitempty"`
	AppVersion     string       `json:"appVersion,omitempty"`
	InApp          *bool        `json:"inApp,omitempty"`
	WebView        *bool        `json:"webview,omitempty"`
	Unknown        *bool        `json:"unknown,omitempty"`
}

type conformanceFile struct {
	Cases []conformanceCase `json:"cases"`
}

func TestConformanceCoverage(t *testing.T) {
	data, err := os.ReadFile("testdata/conformance.json")
	if err != nil {
		t.Fatal(err)
	}
	var file conformanceFile
	if err := json.Unmarshal(data, &file); err != nil {
		t.Fatal(err)
	}
	var browsers, systems, devices, bots, apps map[string]bool
	browsers = make(map[string]bool)
	systems = make(map[string]bool)
	devices = make(map[string]bool)
	bots = make(map[string]bool)
	apps = make(map[string]bool)
	for _, c := range file.Cases {
		if c.Browser != "" {
			browsers[c.Browser] = true
		}
		if c.OS != "" {
			systems[c.OS] = true
		}
		if c.Device != nil {
			devices[c.Device.String()] = true
		}
		if c.BotName != "" {
			bots[c.BotName] = true
		}
		if c.App != "" {
			apps[c.App] = true
		}
	}
	checks := []struct {
		name string
		got  int
		min  int
	}{
		{"browser families", len(browsers), 12},
		{"operating systems", len(systems), 7},
		{"device types", len(devices), 6},
		{"bots", len(bots), 5},
		{"apps", len(apps), 2},
	}
	for _, c := range checks {
		if c.got < c.min {
			t.Errorf("%s: got %d, want at least %d", c.name, c.got, c.min)
		}
	}
}

func TestConformanceCorpus(t *testing.T) {
	data, err := os.ReadFile("testdata/conformance.json")
	if err != nil {
		t.Fatal(err)
	}
	var file conformanceFile
	if err := json.Unmarshal(data, &file); err != nil {
		t.Fatal(err)
	}
	if len(file.Cases) < 24 {
		t.Fatalf("conformance corpus has only %d cases", len(file.Cases))
	}
	seen := make(map[string]bool, len(file.Cases))
	for _, c := range file.Cases {
		c := c
		t.Run(c.Name, func(t *testing.T) {
			if c.Name == "" {
				t.Fatal("case has no name")
			}
			if seen[c.Name] {
				t.Fatalf("duplicate case name %q", c.Name)
			}
			seen[c.Name] = true
			if c.UA == "" {
				t.Fatal("case has no user agent")
			}
			got := Parse(c.UA)
			if got.Raw != c.UA {
				t.Errorf("Raw = %q, want %q", got.Raw, c.UA)
			}
			if c.BrowserEmpty {
				if got.Browser.Name != "" {
					t.Errorf("Browser.Name = %q, want empty", got.Browser.Name)
				}
			} else if c.Browser != "" && got.Browser.Name != c.Browser {
				t.Errorf("Browser.Name = %q, want %q", got.Browser.Name, c.Browser)
			}
			if c.BrowserVersion != "" && got.Browser.Version != c.BrowserVersion {
				t.Errorf("Browser.Version = %q, want %q", got.Browser.Version, c.BrowserVersion)
			}
			if c.OSEmpty {
				if got.OS.Name != "" {
					t.Errorf("OS.Name = %q, want empty", got.OS.Name)
				}
			} else if c.OS != "" && got.OS.Name != c.OS {
				t.Errorf("OS.Name = %q, want %q", got.OS.Name, c.OS)
			}
			if c.OSVersion != "" && got.OS.Version != c.OSVersion {
				t.Errorf("OS.Version = %q, want %q", got.OS.Version, c.OSVersion)
			}
			if c.Device != nil && got.Device.Type != *c.Device {
				t.Errorf("Device.Type = %v, want %v", got.Device.Type, *c.Device)
			}
			if c.Vendor != "" && got.Device.Vendor != c.Vendor {
				t.Errorf("Device.Vendor = %q, want %q", got.Device.Vendor, c.Vendor)
			}
			if c.Model != "" && got.Device.Model != c.Model {
				t.Errorf("Device.Model = %q, want %q", got.Device.Model, c.Model)
			}
			if c.CPU != "" && got.CPU.Architecture != c.CPU {
				t.Errorf("CPU.Architecture = %q, want %q", got.CPU.Architecture, c.CPU)
			}
			if c.Engine != "" && got.Engine.Name != c.Engine {
				t.Errorf("Engine.Name = %q, want %q", got.Engine.Name, c.Engine)
			}
			if c.BotDetected != nil && got.Bot.Detected != *c.BotDetected {
				t.Errorf("Bot.Detected = %v, want %v", got.Bot.Detected, *c.BotDetected)
			}
			if c.BotName != "" && got.Bot.Name != c.BotName {
				t.Errorf("Bot.Name = %q, want %q", got.Bot.Name, c.BotName)
			}
			if c.BotVersion != "" && got.Bot.Version != c.BotVersion {
				t.Errorf("Bot.Version = %q, want %q", got.Bot.Version, c.BotVersion)
			}
			if c.BotURL != "" && got.Bot.URL != c.BotURL {
				t.Errorf("Bot.URL = %q, want %q", got.Bot.URL, c.BotURL)
			}
			if c.BotCategory != nil && got.Bot.Category != *c.BotCategory {
				t.Errorf("Bot.Category = %v, want %v", got.Bot.Category, *c.BotCategory)
			}
			if c.App != "" && got.App.Name != c.App {
				t.Errorf("App.Name = %q, want %q", got.App.Name, c.App)
			}
			if c.AppVersion != "" && got.App.Version != c.AppVersion {
				t.Errorf("App.Version = %q, want %q", got.App.Version, c.AppVersion)
			}
			if c.InApp != nil && got.App.InApp != *c.InApp {
				t.Errorf("App.InApp = %v, want %v", got.App.InApp, *c.InApp)
			}
			if c.WebView != nil && got.WebView != *c.WebView {
				t.Errorf("WebView = %v, want %v", got.WebView, *c.WebView)
			}
			if c.Unknown != nil && got.IsUnknown() != *c.Unknown {
				t.Errorf("IsUnknown() = %v, want %v", got.IsUnknown(), *c.Unknown)
			}
		})
	}
}
