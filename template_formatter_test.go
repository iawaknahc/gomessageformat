package messageformat

import (
	htmltemplate "html/template"
	"strings"
	"testing"

	"golang.org/x/text/language"
)

func TestFormatTemplateParseTree(t *testing.T) {
	en := language.Make("en")

	test := func(pattern string, expected string, args map[string]interface{}) {
		tree, err := FormatTemplateParseTree(en, pattern)
		if err != nil {
			t.Errorf("failed to format html template: %v\n", err)
		} else {
			template, err := htmltemplate.New("main").AddParseTree("main", tree)
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
}
