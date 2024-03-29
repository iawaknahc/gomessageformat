package messageformat

import (
	"fmt"
	"testing"
	"time"

	"golang.org/x/text/language"
)

func TestFormatNamed(t *testing.T) {
	en := language.Make("en")
	test := func(pattern string, expected string, args map[string]interface{}) {
		actual, err := FormatNamed(en, pattern, args)
		if err != nil {
			t.Errorf("err: %v\n", err)
		} else {
			if actual != expected {
				t.Errorf("%v: %v != %v\n", pattern, actual, expected)
			}
		}
	}

	// No arguments
	test("Hello", "Hello", nil)

	// 1 argument at the end of the pattern
	test("Hello {NAME}", "Hello John", map[string]interface{}{
		"NAME": "John",
	})
	// 1 argument in the middle of the pattern
	test("Hello {NAME}, how are you?", "Hello John, how are you?", map[string]interface{}{
		"NAME": "John",
	})
	// 1 argument at the beginning of the pattern
	test("{NAME}, how are you?", "John, how are you?", map[string]interface{}{
		"NAME": "John",
	})

	// 2 reordered arguments.
	test("Hello {YOU}, I am {ME}", "Hello John, I am Jane", map[string]interface{}{
		"YOU": "John",
		"ME":  "Jane",
	})

	// date arguments.
	test("{T, date, short}", "11/10/09", map[string]interface{}{
		"T": time.Date(2009, time.November, 10, 23, 0, 0, 0, time.UTC),
	})
	test("{T, date, medium}", "Nov 10, 2009", map[string]interface{}{
		"T": time.Date(2009, time.November, 10, 23, 0, 0, 0, time.UTC),
	})
	test("{T, date, long}", "November 10, 2009", map[string]interface{}{
		"T": time.Date(2009, time.November, 10, 23, 0, 0, 0, time.UTC),
	})
	test("{T, date, full}", "Tuesday, November 10, 2009", map[string]interface{}{
		"T": time.Date(2009, time.November, 10, 23, 0, 0, 0, time.UTC),
	})
	test("{T, time, short}", "11:00 PM", map[string]interface{}{
		"T": time.Date(2009, time.November, 10, 23, 0, 0, 0, time.UTC),
	})
	test("{T, time, medium}", "11:00:00 PM", map[string]interface{}{
		"T": time.Date(2009, time.November, 10, 23, 0, 0, 0, time.UTC),
	})
	test("{T, time, long}", "11:00:00 PM UTC", map[string]interface{}{
		"T": time.Date(2009, time.November, 10, 23, 0, 0, 0, time.UTC),
	})
	test("{T, time, full}", "11:00:00 PM Coordinated Universal Time", map[string]interface{}{
		"T": time.Date(2009, time.November, 10, 23, 0, 0, 0, time.UTC),
	})
	test("{T, datetime, short}", "11/10/09, 11:00 PM", map[string]interface{}{
		"T": time.Date(2009, time.November, 10, 23, 0, 0, 0, time.UTC),
	})
	test("{T, datetime, medium}", "Nov 10, 2009, 11:00:00 PM", map[string]interface{}{
		"T": time.Date(2009, time.November, 10, 23, 0, 0, 0, time.UTC),
	})
	test("{T, datetime, long}", "November 10, 2009 at 11:00:00 PM UTC", map[string]interface{}{
		"T": time.Date(2009, time.November, 10, 23, 0, 0, 0, time.UTC),
	})
	test("{T, datetime, full}", "Tuesday, November 10, 2009 at 11:00:00 PM Coordinated Universal Time", map[string]interface{}{
		"T": time.Date(2009, time.November, 10, 23, 0, 0, 0, time.UTC),
	})

	// Simple select
	test(`{GENDER, select,
				male {He jumps over the lazy dog}
				female {She jumps over the lazy dog}
				other {They jump over the lazy dog}}`,
		"He jumps over the lazy dog",
		map[string]interface{}{
			"GENDER": "male",
		})
	test(`{GENDER, select,
				male {He jumps over the lazy dog}
				female {She jumps over the lazy dog}
				other {They jump over the lazy dog}}`,
		"She jumps over the lazy dog",
		map[string]interface{}{
			"GENDER": "female",
		})
	test(`{GENDER, select,
				male {He jumps over the lazy dog}
				female {She jumps over the lazy dog}
				other {They jump over the lazy dog}}`,
		"They jump over the lazy dog",
		map[string]interface{}{
			"GENDER": "other",
		})
	test(`{GENDER, select,
				male {He jumps over the lazy dog}
				female {She jumps over the lazy dog}
				other {They jump over the lazy dog}}`,
		"They jump over the lazy dog",
		map[string]interface{}{
			"GENDER": false,
		})

	// Nested select
	test(`{GENDER, select,
				male {He jumps over {OBJECT}}
				female {She jumps over {OBJECT}}
				other {They jump over {OBJECT}}}`,
		"He jumps over the lazy dog",
		map[string]interface{}{
			"GENDER": "male",
			"OBJECT": "the lazy dog",
		})

	// Simple plural
	test(`{COUNT, plural, one{# cat} other{# cats}}`, "1 cat", map[string]interface{}{
		"COUNT": 1,
	})
	test(`{COUNT, plural, one{# cat} other{# cats}}`, "2 cats", map[string]interface{}{
		"COUNT": 2,
	})

	// plural with explicit value
	test(`{COUNT, plural, =0{no cats} one{# cat} other{# cats}}`, "no cats", map[string]interface{}{
		"COUNT": 0,
	})

	// plural with offset
	test(`{COUNT, plural, offset:1
				=1{Kitty and no other cats}
				one{Kitty and 1 other cat}
				other{Kitty and # other cats}}`,
		"Kitty and no other cats",
		map[string]interface{}{
			"COUNT": 1,
		})
	test(`{COUNT, plural, offset:1
				=1{Kitty and no other cats}
				one{Kitty and 1 other cat}
				other{Kitty and # other cats}}`,
		"Kitty and 1 other cat",
		map[string]interface{}{
			"COUNT": 2,
		})
	test(`{COUNT, plural, offset:1
				=1{Kitty and no other cats}
				one{Kitty and 1 other cat}
				other{Kitty and # other cats}}`,
		"Kitty and 2 other cats",
		map[string]interface{}{
			"COUNT": 3,
		})
	test(`{COUNT, plural, offset:1
				=1{Kitty and no other cats}
				one{Kitty and # other cat}
				other{Kitty and # other cats}}`,
		"Kitty and -1 other cat",
		map[string]interface{}{
			"COUNT": 0,
		})
	test(`{COUNT, plural, offset:1
				=0{No Kitty}
				=1{Kitty and no other cats}
				one{Kitty and # other cat}
				other{Kitty and # other cats}}`,
		"No Kitty",
		map[string]interface{}{
			"COUNT": 0,
		})

	// real world example
	pattern := `{GENDER, select,
		   female {{
		       COUNT, plural, offset:1
		       =0 {{HOST} does not give a party.}
		       =1 {{HOST} invites {GUEST} to her party.}
		       =2 {{HOST} invites {GUEST} and one other person to her party.}
		       other {{HOST} invites {GUEST} and # other people to her party.}}}
		   male {{
		       COUNT, plural, offset:1
		       =0 {{HOST} does not give a party.}
		       =1 {{HOST} invites {GUEST} to his party.}
		       =2 {{HOST} invites {GUEST} and one other person to his party.}
		       other {{HOST} invites {GUEST} and # other people to his party.}}}
		   other {{
		       COUNT, plural, offset:1
		       =0 {{HOST} does not give a party.}
		       =1 {{HOST} invites {GUEST} to their party.}
		       =2 {{HOST} invites {GUEST} and one other person to their party.}
		       other {{HOST} invites {GUEST} and # other people to their party.}}}}`

	gender := "female"
	host := "Jane"
	guest := "John"
	test(pattern, "Jane does not give a party.", map[string]interface{}{
		"GENDER": gender,
		"COUNT":  0,
		"HOST":   host,
		"GUEST":  guest,
	})
	test(pattern, "Jane invites John to her party.", map[string]interface{}{
		"GENDER": gender,
		"COUNT":  1,
		"HOST":   host,
		"GUEST":  guest,
	})
	test(pattern, "Jane invites John and one other person to her party.", map[string]interface{}{
		"GENDER": gender,
		"COUNT":  2,
		"HOST":   host,
		"GUEST":  guest,
	})
	test(pattern, "Jane invites John and 2 other people to her party.", map[string]interface{}{
		"GENDER": gender,
		"COUNT":  3,
		"HOST":   host,
		"GUEST":  guest,
	})

	gender = "male"
	host = "John"
	guest = "Jane"
	test(pattern, "John does not give a party.", map[string]interface{}{
		"GENDER": gender,
		"COUNT":  0,
		"HOST":   host,
		"GUEST":  guest,
	})
	test(pattern, "John invites Jane to his party.", map[string]interface{}{
		"GENDER": gender,
		"COUNT":  1,
		"HOST":   host,
		"GUEST":  guest,
	})
	test(pattern, "John invites Jane and one other person to his party.", map[string]interface{}{
		"GENDER": gender,
		"COUNT":  2,
		"HOST":   host,
		"GUEST":  guest,
	})
	test(pattern, "John invites Jane and 2 other people to his party.", map[string]interface{}{
		"GENDER": gender,
		"COUNT":  3,
		"HOST":   host,
		"GUEST":  guest,
	})

	gender = "unspecified"
	host = "Sam"
	guest = "Alex"
	test(pattern, "Sam does not give a party.", map[string]interface{}{
		"GENDER": gender,
		"COUNT":  0,
		"HOST":   host,
		"GUEST":  guest,
	})
	test(pattern, "Sam invites Alex to their party.", map[string]interface{}{
		"GENDER": gender,
		"COUNT":  1,
		"HOST":   host,
		"GUEST":  guest,
	})
	test(pattern, "Sam invites Alex and one other person to their party.", map[string]interface{}{
		"GENDER": gender,
		"COUNT":  2,
		"HOST":   host,
		"GUEST":  guest,
	})
	test(pattern, "Sam invites Alex and 2 other people to their party.", map[string]interface{}{
		"GENDER": gender,
		"COUNT":  3,
		"HOST":   host,
		"GUEST":  guest,
	})
}

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

	// Simple plural
	test(`{0, plural, one{# cat} other{# cats}}`, "1 cat", 1)
	test(`{0, plural, one{# cat} other{# cats}}`, "2 cats", 2)

	// plural with explicit value
	test(`{0, plural, =0{no cats} one{# cat} other{# cats}}`, "no cats", 0)

	// plural with offset
	test(`{0, plural, offset:1
		=1{Kitty and no other cats}
		one{Kitty and 1 other cat}
		other{Kitty and # other cats}}`,
		"Kitty and no other cats",
		1)
	test(`{0, plural, offset:1
		=1{Kitty and no other cats}
		one{Kitty and 1 other cat}
		other{Kitty and # other cats}}`,
		"Kitty and 1 other cat",
		2)
	test(`{0, plural, offset:1
		=1{Kitty and no other cats}
		one{Kitty and 1 other cat}
		other{Kitty and # other cats}}`,
		"Kitty and 2 other cats",
		3)
	test(`{0, plural, offset:1
		=1{Kitty and no other cats}
		one{Kitty and # other cat}
		other{Kitty and # other cats}}`,
		"Kitty and -1 other cat",
		0)
	test(`{0, plural, offset:1
		=0{No Kitty}
		=1{Kitty and no other cats}
		one{Kitty and # other cat}
		other{Kitty and # other cats}}`,
		"No Kitty",
		0)

	// real world example
	pattern := `{0, select,
  female {{
      1, plural, offset:1
      =0 {{2} does not give a party.}
      =1 {{2} invites {3} to her party.}
      =2 {{2} invites {3} and one other person to her party.}
      other {{2} invites {3} and # other people to her party.}}}
  male {{
      1, plural, offset:1
      =0 {{2} does not give a party.}
      =1 {{2} invites {3} to his party.}
      =2 {{2} invites {3} and one other person to his party.}
      other {{2} invites {3} and # other people to his party.}}}
  other {{
      1, plural, offset:1
      =0 {{2} does not give a party.}
      =1 {{2} invites {3} to their party.}
      =2 {{2} invites {3} and one other person to their party.}
      other {{2} invites {3} and # other people to their party.}}}}`

	gender := "female"
	host := "Jane"
	guest := "John"
	test(pattern, "Jane does not give a party.", gender, 0, host, guest)
	test(pattern, "Jane invites John to her party.", gender, 1, host, guest)
	test(pattern, "Jane invites John and one other person to her party.", gender, 2, host, guest)
	test(pattern, "Jane invites John and 2 other people to her party.", gender, 3, host, guest)

	gender = "male"
	host = "John"
	guest = "Jane"
	test(pattern, "John does not give a party.", gender, 0, host, guest)
	test(pattern, "John invites Jane to his party.", gender, 1, host, guest)
	test(pattern, "John invites Jane and one other person to his party.", gender, 2, host, guest)
	test(pattern, "John invites Jane and 2 other people to his party.", gender, 3, host, guest)

	gender = "unspecified"
	host = "Sam"
	guest = "Alex"
	test(pattern, "Sam does not give a party.", gender, 0, host, guest)
	test(pattern, "Sam invites Alex to their party.", gender, 1, host, guest)
	test(pattern, "Sam invites Alex and one other person to their party.", gender, 2, host, guest)
	test(pattern, "Sam invites Alex and 2 other people to their party.", gender, 3, host, guest)
}

func TestTextUnknownArgument(t *testing.T) {
	en := language.Make("en")
	test := func(pattern string, expected string, args map[string]interface{}) {
		actual, err := FormatNamed(en, pattern, args)
		if err != nil {
			t.Errorf("err: %v\n", err)
		} else {
			if actual != expected {
				t.Errorf("%v: %v != %v\n", pattern, actual, expected)
			}
		}
	}

	test("Hello {NAME}", "Hello ", nil)
	test("Hello {T, date, short} Hello", "Hello  Hello", nil)
	test("Hello {T, time, short} Hello", "Hello  Hello", nil)
	test("Hello {T, datetime, short} Hello", "Hello  Hello", nil)
	test("Hello {GENDER, select, male {he} female {she} other {they}}", "Hello they", nil)
	test("Hello {COUNT, plural, one {# cat} other {# cats}}", "Hello 0 cats", nil)
}

func ExampleFormatPositional() {
	numFiles := 1
	out, err := FormatPositional(
		language.English,
		`{0, plural,
			=0 {There are no files on disk.}
			=1 {There is only 1 file on disk.}
			other {There are # files on disk.}
		}`,
		numFiles,
	)
	if err != nil {
		panic(err)
	}
	fmt.Printf("%s\n", out)
	// Output: There is only 1 file on disk.
}
