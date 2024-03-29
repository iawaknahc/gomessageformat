package messageformat

import (
	"fmt"
	htmltemplate "html/template"
	"strings"
	"testing"
	"time"

	"golang.org/x/text/language"
)

func TestFormatTemplateParseTree(t *testing.T) {
	en := language.Make("en")

	test := func(pattern string, expected string, args map[string]interface{}) {
		tree, err := FormatTemplateParseTree(en, pattern)
		if err != nil {
			t.Errorf("failed to format html template: %v\n", err)
		} else {
			template := htmltemplate.New("main")
			template.Funcs(htmltemplate.FuncMap{
				TemplateRuntimeFuncName: TemplateRuntimeFunc,
			})
			template, err := template.AddParseTree("main", tree)
			if err != nil {
				t.Errorf("failed to add parse tree: %v\n", err)
			} else {
				var buf strings.Builder
				err = template.Execute(&buf, args)
				if err != nil {
					t.Errorf("failed to execute: %v\n", err)
				} else {
					actual := buf.String()
					if actual != expected {
						t.Errorf("%v: %v != %v\n", pattern, actual, expected)
					}
				}
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

	// Select with only other clause
	test(`{GENDER, select,
				other {They jump over the lazy dog}}`,
		"They jump over the lazy dog",
		map[string]interface{}{
			"GENDER": 0,
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

	// HTML
	test(`Hello <b>{NAME}</b>`, `Hello <b>John</b>`, map[string]interface{}{
		"NAME": "John",
	})

	test(`<a href="{UNSAFE}">Click me</a>`, `<a href="#ZgotmplZ">Click me</a>`, map[string]interface{}{
		"UNSAFE": "javascript:alert()",
	})
}

func TestTemplateUnknownArgument(t *testing.T) {
	en := language.Make("en")
	test := func(pattern string, expected string, args map[string]interface{}) {
		tree, err := FormatTemplateParseTree(en, pattern)
		if err != nil {
			t.Errorf("failed to format html template: %v\n", err)
		} else {
			template := htmltemplate.New("main")
			template.Funcs(htmltemplate.FuncMap{
				TemplateRuntimeFuncName: TemplateRuntimeFunc,
			})
			template, err := template.AddParseTree("main", tree)
			if err != nil {
				t.Errorf("failed to add parse tree: %v\n", err)
			} else {
				var buf strings.Builder
				err = template.Execute(&buf, args)
				if err != nil {
					t.Errorf("failed to execute: %v\n", err)
				} else {
					actual := buf.String()
					if actual != expected {
						t.Errorf("%v: %v != %v\n", pattern, actual, expected)
					}
				}
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

func TestIsEmptyParseTree(t *testing.T) {
	tree, _ := FormatTemplateParseTree(language.Make("en"), "nonempty")
	if IsEmptyParseTree(tree) {
		t.Errorf("nonempty")
	}

	tree, _ = FormatTemplateParseTree(language.Make("en"), "")
	if !IsEmptyParseTree(tree) {
		t.Errorf("empty")
	}
}

func ExampleFormatTemplateParseTree() {
	tree, err := FormatTemplateParseTree(language.English,
		`Hello there! Check <a href="{LINK}">this</a> out!`)
	if err != nil {
		panic(err)
	}

	template := htmltemplate.New("main")
	template.Funcs(htmltemplate.FuncMap{
		TemplateRuntimeFuncName: TemplateRuntimeFunc,
		"makemap": func(pairs ...interface{}) map[string]interface{} {
			out := make(map[string]interface{})
			for i := 0; i < len(pairs); i += 2 {
				key := pairs[i].(string)
				value := pairs[i+1]
				out[key] = value
			}
			return out
		},
	})

	_, err = template.Parse(`<!DOCTYPE html><html><head><title>Hi</title></head><body><p>{{ template "greeting" (makemap "LINK" .Link) }}</p></body></html>`)
	if err != nil {
		panic(err)
	}

	_, err = template.AddParseTree("greeting", tree)
	if err != nil {
		panic(err)
	}
	var buf strings.Builder
	err = template.Execute(&buf, map[string]interface{}{
		"Link": "https://www.example.com",
	})
	if err != nil {
		panic(err)
	}
	fmt.Printf("%s\n", buf.String())
	// Output: <!DOCTYPE html><html><head><title>Hi</title></head><body><p>Hello there! Check <a href="https://www.example.com">this</a> out!</p></body></html>
}
