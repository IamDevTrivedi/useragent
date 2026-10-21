package useragent

import (
	"strconv"
	"strings"
)

// maxVersionComponent bounds parsed components so hostile input cannot
// overflow int. Real user-agent versions never approach this value.
const maxVersionComponent = 1_000_000_000

// Version is a parsed numeric version. Missing components are zero.
type Version struct {
	Major int `json:"major"`
	Minor int `json:"minor"`
	Patch int `json:"patch"`
}

// ParseVersion reads up to three numeric runs from s, ignoring leading
// non-digits. It accepts underscores as separators ("16_3_1") and stops after
// the third component. The raw string is kept by the caller when needed.
func ParseVersion(s string) Version {
	if s == "" {
		return Version{}
	}
	s = strings.ReplaceAll(s, "_", ".")
	var v Version
	comp := 0
	for i := 0; i < len(s) && comp < 3; {
		for i < len(s) && !isDigit(s[i]) {
			i++
		}
		if i >= len(s) {
			break
		}
		n := 0
		for i < len(s) && isDigit(s[i]) {
			d := int(s[i] - '0')
			if n > (maxVersionComponent-d)/10 {
				n = maxVersionComponent
			} else {
				n = n*10 + d
			}
			i++
		}
		switch comp {
		case 0:
			v.Major = n
		case 1:
			v.Minor = n
		case 2:
			v.Patch = n
		}
		comp++
		for i < len(s) && s[i] != '.' {
			i++
		}
	}
	return v
}

func isDigit(c byte) bool { return c >= '0' && c <= '9' }

// String returns "major.minor.patch".
func (v Version) String() string {
	return strconv.Itoa(v.Major) + "." + strconv.Itoa(v.Minor) + "." + strconv.Itoa(v.Patch)
}

// Short returns "major.minor".
func (v Version) Short() string { return strconv.Itoa(v.Major) + "." + strconv.Itoa(v.Minor) }

// IsZero reports whether every component is zero.
func (v Version) IsZero() bool { return v.Major == 0 && v.Minor == 0 && v.Patch == 0 }

// Compare returns -1 if v is older than other, 0 if equal, and +1 if newer.
func (v Version) Compare(other Version) int {
	switch {
	case v.Major != other.Major:
		return cmpInt(v.Major, other.Major)
	case v.Minor != other.Minor:
		return cmpInt(v.Minor, other.Minor)
	default:
		return cmpInt(v.Patch, other.Patch)
	}
}

// AtLeast reports whether v is greater than or equal to major.minor.patch.
func (v Version) AtLeast(major, minor, patch int) bool {
	return v.Compare(Version{Major: major, Minor: minor, Patch: patch}) >= 0
}

func cmpInt(a, b int) int {
	switch {
	case a < b:
		return -1
	case a > b:
		return 1
	default:
		return 0
	}
}
