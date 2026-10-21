// Command ua parses user-agent strings from arguments, a file, or stdin.
package main

import (
	"bufio"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
	"sort"
	"strings"
	"text/tabwriter"

	ua "github.com/IamDevTrivedi/useragent"
)

func main() {
	if err := run(os.Args[1:], os.Stdin, os.Stdout, os.Stderr); err != nil {
		fmt.Fprintln(os.Stderr, "ua:", err)
		os.Exit(1)
	}
}

type config struct {
	json    bool
	pretty  bool
	summary bool
	file    string
	logMode string
	args    []string
}

func run(args []string, stdin io.Reader, stdout, stderr io.Writer) error {
	var cfg config
	fs := flag.NewFlagSet("ua", flag.ContinueOnError)
	fs.SetOutput(stderr)
	fs.BoolVar(&cfg.json, "json", false, "write one JSON object per line")
	fs.BoolVar(&cfg.pretty, "pretty", false, "indent JSON output")
	fs.BoolVar(&cfg.summary, "summary", false, "print aggregate counts")
	fs.StringVar(&cfg.file, "file", "", "read user-agent strings or log lines from file")
	fs.StringVar(&cfg.logMode, "log", "", "extract the user agent from combined or common log lines")
	if err := fs.Parse(args); err != nil {
		return err
	}
	cfg.args = fs.Args()
	if cfg.json && cfg.summary {
		return errors.New("-json and -summary cannot be combined")
	}
	switch cfg.logMode {
	case "", "combined", "common":
	default:
		return fmt.Errorf("unknown -log format %q (want combined or common)", cfg.logMode)
	}

	inputs, err := collectInputs(cfg, stdin)
	if err != nil {
		return err
	}
	parsed := make([]ua.UserAgent, 0, len(inputs))
	for _, in := range inputs {
		raw := in
		if cfg.logMode != "" {
			raw = extractUA(in, cfg.logMode)
		}
		parsed = append(parsed, ua.Parse(raw))
	}

	if cfg.summary {
		return printSummary(stdout, summarize(parsed))
	}
	if cfg.json {
		enc := json.NewEncoder(stdout)
		if cfg.pretty {
			enc.SetIndent("", "  ")
		}
		for _, v := range parsed {
			if err := enc.Encode(v); err != nil {
				return err
			}
		}
		return nil
	}
	return printTable(stdout, parsed)
}

func collectInputs(cfg config, stdin io.Reader) ([]string, error) {
	var inputs []string
	if len(cfg.args) != 0 {
		inputs = append(inputs, cfg.args...)
		return inputs, nil
	}
	r := stdin
	if cfg.file != "" {
		f, err := os.Open(cfg.file)
		if err != nil {
			return nil, err
		}
		defer func() { _ = f.Close() }()
		r = f
	}
	sc := bufio.NewScanner(r)
	sc.Buffer(make([]byte, 64*1024), 4*1024*1024)
	for sc.Scan() {
		line := strings.TrimSpace(sc.Text())
		if line == "" {
			continue
		}
		inputs = append(inputs, line)
	}
	if err := sc.Err(); err != nil {
		return nil, err
	}
	return inputs, nil
}

// extractUA returns the user-agent field from an Apache-style log line.
// Combined logs put the user agent in the last quoted field. Common logs do
// not contain a user agent, so the whole line is returned unchanged.
func extractUA(line, mode string) string {
	if mode != "combined" {
		return line
	}
	start := -1
	end := -1
	for i := len(line) - 1; i >= 0; i-- {
		if line[i] != '"' {
			continue
		}
		if end < 0 {
			end = i
			continue
		}
		start = i
		break
	}
	if start >= 0 && end > start {
		return line[start+1 : end]
	}
	return strings.TrimSpace(line)
}

func printTable(w io.Writer, parsed []ua.UserAgent) error {
	tw := tabwriter.NewWriter(w, 0, 4, 2, ' ', 0)
	if _, err := fmt.Fprintln(tw, "BROWSER\tVERSION\tOS\tDEVICE\tBOT\tAPP\tRAW"); err != nil {
		return err
	}
	for _, v := range parsed {
		bot := ""
		if v.Bot.Detected {
			bot = v.Bot.Name
		}
		app := ""
		if v.App.Name != "" {
			app = v.App.Name
		}
		if _, err := fmt.Fprintf(tw, "%s\t%s\t%s\t%s\t%s\t%s\t%s\n",
			v.Browser.Name, v.Browser.Version, v.OS.Name, v.Device.Type, bot, app, v.Raw); err != nil {
			return err
		}
	}
	return tw.Flush()
}

// Summary contains aggregate parser output.
type Summary struct {
	Total          int            `json:"total"`
	Browsers       map[string]int `json:"browsers"`
	Operating      map[string]int `json:"operating_systems"`
	Devices        map[string]int `json:"devices"`
	Bots           map[string]int `json:"bots"`
	BotCategories  map[string]int `json:"bot_categories"`
	UnknownBrowser int            `json:"unknown_browser"`
}

func summarize(parsed []ua.UserAgent) Summary {
	s := Summary{
		Browsers:      map[string]int{},
		Operating:     map[string]int{},
		Devices:       map[string]int{},
		Bots:          map[string]int{},
		BotCategories: map[string]int{},
	}
	for _, v := range parsed {
		s.Total++
		if v.Browser.Name == "" {
			s.UnknownBrowser++
		} else {
			s.Browsers[v.Browser.Name]++
		}
		if v.OS.Name != "" {
			s.Operating[v.OS.Name]++
		}
		s.Devices[v.Device.Type.String()]++
		if v.Bot.Detected {
			s.Bots[v.Bot.Name]++
			s.BotCategories[v.Bot.Category.String()]++
		}
	}
	return s
}

func printSummary(w io.Writer, s Summary) error {
	if _, err := fmt.Fprintf(w, "total: %d\n", s.Total); err != nil {
		return err
	}
	if err := printCounts(w, "browsers", s.Browsers); err != nil {
		return err
	}
	if err := printCounts(w, "operating systems", s.Operating); err != nil {
		return err
	}
	if err := printCounts(w, "devices", s.Devices); err != nil {
		return err
	}
	if err := printCounts(w, "bots", s.Bots); err != nil {
		return err
	}
	if err := printCounts(w, "bot categories", s.BotCategories); err != nil {
		return err
	}
	if s.UnknownBrowser != 0 {
		_, err := fmt.Fprintf(w, "unknown browser: %d\n", s.UnknownBrowser)
		return err
	}
	return nil
}

func printCounts(w io.Writer, label string, counts map[string]int) error {
	keys := make([]string, 0, len(counts))
	for k := range counts {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	if _, err := fmt.Fprintf(w, "%s:\n", label); err != nil {
		return err
	}
	for _, k := range keys {
		if _, err := fmt.Fprintf(w, "  %s: %d\n", k, counts[k]); err != nil {
			return err
		}
	}
	return nil
}
