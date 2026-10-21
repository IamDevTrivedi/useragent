package main

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"

	ua "github.com/IamDevTrivedi/useragent"
)

func TestExtractUA(t *testing.T) {
	tests := []struct {
		name string
		line string
		mode string
		want string
	}{
		{
			"combined",
			`127.0.0.1 - - [10/Oct/2000:13:55:36 -0700] "GET / HTTP/1.0" 200 2326 "http://example.com/" "Mozilla/5.0 (X11; Linux x86_64) Chrome/120.0.0.0 Safari/537.36"`,
			"combined",
			"Mozilla/5.0 (X11; Linux x86_64) Chrome/120.0.0.0 Safari/537.36",
		},
		{"common", `127.0.0.1 - - [10/Oct/2000:13:55:36 -0700] "GET / HTTP/1.0" 200 2326`, "common", `127.0.0.1 - - [10/Oct/2000:13:55:36 -0700] "GET / HTTP/1.0" 200 2326`},
		{"no quotes", "plain user agent", "combined", "plain user agent"},
		{"empty quotes", `a "" b`, "combined", ""},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := extractUA(tt.line, tt.mode); got != tt.want {
				t.Errorf("extractUA() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestSummarize(t *testing.T) {
	parsed := []ua.UserAgent{
		ua.Parse("Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 Chrome/120.0.0.0 Safari/537.36"),
		ua.Parse("Mozilla/5.0 (compatible; Googlebot/2.1; +http://www.google.com/bot.html)"),
		ua.Parse(""),
	}
	got := summarize(parsed)
	if got.Total != 3 {
		t.Errorf("Total = %d", got.Total)
	}
	if got.Browsers["Chrome"] != 1 || got.Browsers["Googlebot"] != 1 {
		t.Errorf("Browsers = %+v", got.Browsers)
	}
	if got.UnknownBrowser != 1 {
		t.Errorf("UnknownBrowser = %d", got.UnknownBrowser)
	}
	if got.Operating["Windows"] != 1 {
		t.Errorf("Operating = %+v", got.Operating)
	}
	if got.Devices["desktop"] != 1 || got.Devices["unknown"] != 2 {
		t.Errorf("Devices = %+v", got.Devices)
	}
	if got.Bots["Googlebot"] != 1 || got.BotCategories["search"] != 1 {
		t.Errorf("bots = %+v / %+v", got.Bots, got.BotCategories)
	}
}

func TestRunJSON(t *testing.T) {
	var out, errOut bytes.Buffer
	raw := "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 Chrome/120.0.0.0 Safari/537.36"
	if err := run([]string{"-json", raw}, strings.NewReader(""), &out, &errOut); err != nil {
		t.Fatal(err)
	}
	s := out.String()
	for _, want := range []string{`"raw":"Mozilla/5.0`, `"browser"`, `"name":"Chrome"`, `"type":"desktop"`} {
		if !strings.Contains(s, want) {
			t.Errorf("JSON output missing %q:\n%s", want, s)
		}
	}
	if strings.Count(strings.TrimSpace(s), "\n") != 0 {
		t.Errorf("expected one JSON object, got:\n%s", s)
	}
}

func TestRunPrettyJSON(t *testing.T) {
	var out, errOut bytes.Buffer
	if err := run([]string{"-json", "-pretty", "Twitterbot/1.0"}, strings.NewReader(""), &out, &errOut); err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(out.String(), "\n  \"raw\"") {
		t.Errorf("pretty JSON missing indentation:\n%s", out.String())
	}
}

func TestRunSummary(t *testing.T) {
	var out, errOut bytes.Buffer
	raw := "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 Chrome/120.0.0.0 Safari/537.36"
	if err := run([]string{"-summary", raw}, strings.NewReader(""), &out, &errOut); err != nil {
		t.Fatal(err)
	}
	s := out.String()
	for _, want := range []string{"total: 1", "browsers:", "Chrome: 1", "devices:", "desktop: 1", "bots:", "unknown browser: 0"} {
		if want == "unknown browser: 0" {
			if strings.Contains(s, want) {
				t.Errorf("summary should omit zero unknown browsers:\n%s", s)
			}
			continue
		}
		if !strings.Contains(s, want) {
			t.Errorf("summary missing %q:\n%s", want, s)
		}
	}
}

func TestRunTable(t *testing.T) {
	var out, errOut bytes.Buffer
	if err := run([]string{"Twitterbot/1.0"}, strings.NewReader(""), &out, &errOut); err != nil {
		t.Fatal(err)
	}
	s := out.String()
	for _, want := range []string{"BROWSER", "VERSION", "OS", "DEVICE", "BOT", "Twitterbot"} {
		if !strings.Contains(s, want) {
			t.Errorf("table missing %q:\n%s", want, s)
		}
	}
}

func TestRunStdin(t *testing.T) {
	var out, errOut bytes.Buffer
	in := strings.NewReader("\nTwitterbot/1.0\n\nGooglebot/2.1\n")
	if err := run(nil, in, &out, &errOut); err != nil {
		t.Fatal(err)
	}
	if got := strings.Count(out.String(), "\n"); got != 3 { // header + two rows
		t.Errorf("line count = %d:\n%s", got, out.String())
	}
}

func TestRunFile(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "agents.txt")
	if err := os.WriteFile(path, []byte("Twitterbot/1.0\nGooglebot/2.1\n"), 0o600); err != nil {
		t.Fatal(err)
	}
	var out, errOut bytes.Buffer
	if err := run([]string{"-file", path}, strings.NewReader(""), &out, &errOut); err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(out.String(), "Twitterbot") || !strings.Contains(out.String(), "Googlebot") {
		t.Errorf("file output:\n%s", out.String())
	}
}

func TestRunCombinedLog(t *testing.T) {
	var out, errOut bytes.Buffer
	line := `127.0.0.1 - - [10/Oct/2000:13:55:36 -0700] "GET / HTTP/1.0" 200 2326 "http://example.com/" "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 Chrome/120.0.0.0 Safari/537.36"`
	if err := run([]string{"-log", "combined", line}, strings.NewReader(""), &out, &errOut); err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(out.String(), "Chrome") || strings.Contains(out.String(), "127.0.0.1") {
		t.Errorf("combined log output:\n%s", out.String())
	}
}

func TestRunErrors(t *testing.T) {
	var out, errOut bytes.Buffer
	if err := run([]string{"-json", "-summary", "x"}, strings.NewReader(""), &out, &errOut); err == nil {
		t.Error("expected -json/-summary conflict")
	}
	if err := run([]string{"-log", "nope", "x"}, strings.NewReader(""), &out, &errOut); err == nil {
		t.Error("expected unknown log format error")
	}
	if err := run([]string{"-file", filepath.Join(t.TempDir(), "missing")}, strings.NewReader(""), &out, &errOut); err == nil {
		t.Error("expected missing file error")
	}
	if err := run([]string{"-unknown-flag"}, strings.NewReader(""), &out, &errOut); err == nil {
		t.Error("expected flag parse error")
	}
}
