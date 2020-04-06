package messageformat

import (
	"testing"

	"golang.org/x/text/language"
)

func TestFormatPositional(t *testing.T) {
	en := language.Make("en")
	test := func(pattern string, expected string, args ...interface{}) {
		actual, err := FormatPositional(en, pattern, args...)
		if err != nil {
			t.Errorf("err: %v\n", err)
		} else {
			if actual != expected {
				t.Errorf("%v: %v != %v\n", pattern, actual, expected)
			}
		}
	}

	// No arguments
	test("Hello", "Hello")

	// 1 argument at the end of the pattern
	test("Hello {0}", "Hello John", "John")
	// 1 argument in the middle of the pattern
	test("Hello {0}, how are you?", "Hello John, how are you?", "John")
	// 1 argument at the beginning of the pattern
	test("{0}, how are you?", "John, how are you?", "John")

	// 2 reordered arguments.
	test("Hello {1}, I am {0}", "Hello John, I am Jane", "Jane", "John")

	// Simple select
	test(`{0, select,
		male {He jumps over the lazy dog}
		female {She jumps over the lazy dog}
		other {They jump over the lazy dog}}`,
		"He jumps over the lazy dog",
		"male")
	test(`{0, select,
		male {He jumps over the lazy dog}
		female {She jumps over the lazy dog}
		other {They jump over the lazy dog}}`,
		"She jumps over the lazy dog",
		"female")
	test(`{0, select,
		male {He jumps over the lazy dog}
		female {She jumps over the lazy dog}
		other {They jump over the lazy dog}}`,
		"They jump over the lazy dog",
		"other")
	test(`{0, select,
		male {He jumps over the lazy dog}
		female {She jumps over the lazy dog}
		other {They jump over the lazy dog}}`,
		"They jump over the lazy dog",
		false)

	// Nested select
	test(`{0, select,
		male {He jumps over {1}}
		female {She jumps over {1}}
		other {They jump over {1}}}`,
		"He jumps over the lazy dog",
		"male",
		"the lazy dog")
}
