package messageformat

import (
	"testing"

	"golang.org/x/text/language"
)

func TestCardinal(t *testing.T) {
	test := func(lang string, number interface{}, expected string) {
		actual, err := Cardinal(language.Make(lang), number)
		if err != nil {
			t.Errorf("err: %v\n", err)
		} else {
			if actual != expected {
				t.Errorf("%s %v: %v != %v\n", lang, number, actual, expected)
			}
		}
	}

	test("en", 0, "other")
	test("en", 1, "one")
	test("en", 1.1, "other")
	test("en", -1, "one")
	test("en", 2, "other")
	test("en", 3, "other")
	test("en", 4, "other")

	test("en-US", 0, "other")
	test("en-US", 1, "one")
	test("en-US", 1.1, "other")
	test("en-US", -1, "one")
	test("en-US", 2, "other")
	test("en-US", 3, "other")
	test("en-US", 4, "other")

	test("zh", 0, "other")
	test("zh", 1, "other")
	test("zh", 2, "other")
	test("zh", 3, "other")
	test("zh", 4, "other")

	test("ja", 0, "other")
	test("ja", 1, "other")
	test("ja", 2, "other")
	test("ja", 3, "other")
	test("ja", 4, "other")
}

func TestOrdinal(t *testing.T) {
	test := func(lang string, number interface{}, expected string) {
		actual, err := Ordinal(language.Make(lang), number)
		if err != nil {
			t.Errorf("err: %v\n", err)
		} else {
			if actual != expected {
				t.Errorf("%s %v: %v != %v\n", lang, number, actual, expected)
			}
		}
	}

	test("en", 0, "other")  // 0th
	test("en", 1, "one")    // 1st
	test("en", 2, "two")    // 2nd
	test("en", 3, "few")    // 3rd
	test("en", 4, "other")  // 4th
	test("en", 10, "other") // 10th
	test("en", 11, "other") // 11th
	test("en", 12, "other") // 12th
	test("en", 13, "other") // 13th
	test("en", 14, "other") // 14th
	test("en", 20, "other") // 20th
	test("en", 21, "one")   // 21st
	test("en", 22, "two")   // 22nd
	test("en", 23, "few")   // 23rd
	test("en", 24, "other") // 24th
}
