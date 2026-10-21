package useragent

import "testing"

func TestParseVersion(t *testing.T) {
	tests := []struct {
		in   string
		want Version
	}{
		{"", Version{}},
		{"120.0.0.0", Version{120, 0, 0}},
		{"10_15_7", Version{10, 15, 7}},
		{"16.3.1", Version{16, 3, 1}},
		{"12.11.5-gn", Version{12, 11, 5}},
		{"5.4", Version{5, 4, 0}},
		{"v1.2.3", Version{1, 2, 3}},
		{"7~bl", Version{7, 0, 0}},
		{"2024", Version{2024, 0, 0}},
		{"abc", Version{}},
		{"1..2", Version{1, 2, 0}},
		{"1.2.3.4", Version{1, 2, 3}},
		{"0.0.0", Version{}},
	}
	for _, tt := range tests {
		t.Run(tt.in, func(t *testing.T) {
			if got := ParseVersion(tt.in); got != tt.want {
				t.Errorf("ParseVersion(%q) = %+v, want %+v", tt.in, got, tt.want)
			}
		})
	}
}

func TestVersionMethods(t *testing.T) {
	if got := (Version{1, 2, 3}).String(); got != "1.2.3" {
		t.Errorf("String() = %q", got)
	}
	if got := (Version{1, 2, 3}).Short(); got != "1.2" {
		t.Errorf("Short() = %q", got)
	}
	if !(Version{}).IsZero() {
		t.Error("zero version IsZero() = false")
	}
	if (Version{0, 0, 1}).IsZero() {
		t.Error("patch version IsZero() = true")
	}

	tests := []struct {
		a, b Version
		want int
	}{
		{Version{1, 2, 3}, Version{1, 2, 3}, 0},
		{Version{1, 2, 3}, Version{2, 0, 0}, -1},
		{Version{2, 0, 0}, Version{1, 9, 9}, 1},
		{Version{1, 2, 3}, Version{1, 3, 0}, -1},
		{Version{1, 2, 4}, Version{1, 2, 3}, 1},
	}
	for _, tt := range tests {
		if got := tt.a.Compare(tt.b); got != tt.want {
			t.Errorf("%v.Compare(%v) = %d, want %d", tt.a, tt.b, got, tt.want)
		}
	}
	if !(Version{120, 0, 0}).AtLeast(119, 9, 9) {
		t.Error("AtLeast should be true")
	}
	if (Version{1, 2, 3}).AtLeast(1, 2, 4) {
		t.Error("AtLeast should be false")
	}
}

func TestVersionInParseResult(t *testing.T) {
	ua := Parse("Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.6099.109 Safari/537.36")
	if ua.Browser.VersionNo != (Version{120, 0, 6099}) {
		t.Errorf("Browser.VersionNo = %+v", ua.Browser.VersionNo)
	}
	if ua.OS.VersionNo != (Version{10, 0, 0}) {
		t.Errorf("OS.VersionNo = %+v", ua.OS.VersionNo)
	}
}
