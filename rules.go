package useragent

// Match selects a token set. All non-zero fields must match. A zero Match does
// not match anything, which prevents accidental catch-all custom rules.
type Match struct {
	// Token matches an exact token key, case-insensitively.
	Token string
	// Prefix matches when any token key starts with Prefix.
	Prefix string
	// Suffix matches when any token key ends with Suffix.
	Suffix string
	// Contains matches when any token key or value contains Contains.
	Contains string
	// All requires every listed token key to be present.
	All []string
	// Any requires at least one listed token key to be present.
	Any []string
	// None requires none of the listed token keys to be present.
	None []string
}

func (m Match) matches(ts *tokenSet) bool {
	set := false
	if m.Token != "" {
		set = true
		if !ts.has(m.Token) {
			return false
		}
	}
	if m.Prefix != "" {
		set = true
		if !ts.hasPrefix(m.Prefix) {
			return false
		}
	}
	if m.Suffix != "" {
		set = true
		if !ts.hasSuffix(m.Suffix) {
			return false
		}
	}
	if m.Contains != "" {
		set = true
		if !ts.rawHas(m.Contains) {
			return false
		}
	}
	if len(m.All) != 0 {
		set = true
		if !ts.hasAll(m.All...) {
			return false
		}
	}
	if len(m.Any) != 0 {
		set = true
		found := false
		for _, k := range m.Any {
			if ts.has(k) {
				found = true
				break
			}
		}
		if !found {
			return false
		}
	}
	if len(m.None) != 0 {
		set = true
		for _, k := range m.None {
			if ts.has(k) {
				return false
			}
		}
	}
	return set
}

// BrowserRule adds or overrides browser detection. VersionKeys are token keys
// tried in order. Product extracts a version from raw "Product/version" syntax
// after VersionKeys are exhausted.
type BrowserRule struct {
	Name        string
	Match       Match
	VersionKeys []string
	Product     string
	Priority    int
}

// OSRule adds or overrides operating-system detection.
type OSRule struct {
	Name        string
	Match       Match
	VersionKeys []string
	Product     string
	Priority    int
}

// BotRule adds or overrides bot detection.
type BotRule struct {
	Name        string
	Match       Match
	VersionKeys []string
	Product     string
	Category    BotCategory
	URL         string
	Priority    int
}

// EngineRule adds or overrides engine detection.
type EngineRule struct {
	Name        string
	VersionKeys []string
	Product     string
	Match       Match
	Priority    int
}

// DeviceRule adds or overrides device detection.
type DeviceRule struct {
	Type     DeviceType
	Vendor   string
	Model    string
	Match    Match
	Priority int
}

// Option configures a Parser. Options only affect construction; a Parser is
// immutable and safe for concurrent use after New returns.
type Option func(*Parser)

// WithBrowserRules appends custom browser rules. Higher Priority rules are
// evaluated first; the first matching rule wins.
func WithBrowserRules(rules ...BrowserRule) Option {
	return func(p *Parser) { p.browserRules = append(p.browserRules, rules...) }
}

// WithOSRules appends custom operating-system rules.
func WithOSRules(rules ...OSRule) Option {
	return func(p *Parser) { p.osRules = append(p.osRules, rules...) }
}

// WithBotRules appends custom bot rules.
func WithBotRules(rules ...BotRule) Option {
	return func(p *Parser) { p.botRules = append(p.botRules, rules...) }
}

// WithEngineRules appends custom engine rules.
func WithEngineRules(rules ...EngineRule) Option {
	return func(p *Parser) { p.engineRules = append(p.engineRules, rules...) }
}

// WithDeviceRules appends custom device rules.
func WithDeviceRules(rules ...DeviceRule) Option {
	return func(p *Parser) { p.deviceRules = append(p.deviceRules, rules...) }
}
