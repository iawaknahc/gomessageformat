package messageformat

import (
	"strings"
	"testing"
)

// Unicode space separators that icu4c may emit in place of an ordinary space.
//
// These are written as numeric literals on purpose: as character literals they
// are indistinguishable from an ordinary space when reading the source.
const (
	noBreakSpace       = 0x00A0
	thinSpace          = 0x2009
	narrowNoBreakSpace = 0x202F
)

// normalizeICUSpaces replaces those separators with an ordinary U+0020 SPACE.
//
// icu4c delegates date and time formatting to CLDR, and CLDR 42 (shipped in
// icu4c 72) changed the separator before AM/PM in en from U+0020 SPACE to
// U+202F NARROW NO-BREAK SPACE. Which separator is used is icu4c's business
// rather than this package's, so the tests compare normalized strings and stay
// green across icu4c versions.
func normalizeICUSpaces(s string) string {
	return strings.Map(func(r rune) rune {
		switch r {
		case noBreakSpace, thinSpace, narrowNoBreakSpace:
			return ' '
		default:
			return r
		}
	}, s)
}

// normalizeICUSpacesEqual reports whether actual and expected are equal after
// normalizing Unicode space separators.
func normalizeICUSpacesEqual(actual, expected string) bool {
	return normalizeICUSpaces(actual) == normalizeICUSpaces(expected)
}

func TestNormalizeICUSpaces(t *testing.T) {
	// The icu4c >= 72 rendering must compare equal to the plain-space expectation.
	if !normalizeICUSpacesEqual("11:00\u202fPM", "11:00 PM") {
		t.Errorf("expected U+202F to normalize to a space")
	}
	if !normalizeICUSpacesEqual("a\u00a0b", "a b") {
		t.Errorf("expected U+00A0 to normalize to a space")
	}
	// Genuine differences must still be reported.
	if normalizeICUSpacesEqual("11:00\u202fPM", "11:00\u202fAM") {
		t.Errorf("must not ignore a real difference")
	}
	if normalizeICUSpacesEqual("11:00 PM", "11:00PM") {
		t.Errorf("must not ignore a missing separator")
	}
}
