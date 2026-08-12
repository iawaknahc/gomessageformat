package messageformat

import (
	"bytes"
	"errors"
	"reflect"
	"testing"
)

func lexText(t *testing.T, input string, tokens ...string) {
	l := newLexer(input)
	var actual []string
	var err error
loop:
	for {
		err = l.LexText(&bytes.Buffer{})
		if err != nil {
			break
		}
		for _, token := range l.Output {
			if token.Type == TokenTypeEOF {
				break loop
			}
			actual = append(actual, token.String())
		}
		l.Output = nil
	}

	if err != nil {
		t.Errorf("err: %v\n", err)
	}

	if !reflect.DeepEqual(actual, tokens) {
		t.Errorf("expected: %v\n", tokens)
		t.Errorf("actual: %v\n", actual)
	}
}

func lexArg(t *testing.T, input string, tokens ...string) {
	l := newLexer(input)
	var actual []string
	var err error
loop:
	for {
		err = l.LexArg()
		if err != nil {
			break
		}
		for _, token := range l.Output {
			if token.Type == TokenTypeEOF {
				break loop
			}
			actual = append(actual, token.String())
		}
		l.Output = nil
	}

	if err != nil {
		t.Errorf("err: %v\n", err)
	}

	if !reflect.DeepEqual(actual, tokens) {
		t.Errorf("expected: %v\n", tokens)
		t.Errorf("actual: %v\n", actual)
	}
}

func TestLexText(t *testing.T) {
	lexText(t,
		"",
		`""`)
	lexText(t,
		"a",
		`"a"`)
	lexText(t,
		"Hello {",
		`"Hello "`, "{", `""`)
	lexText(t,
		"''",
		`"'"`)
	lexText(t,
		"'a'",
		`"a"`)
	lexText(t,
		"a'a'",
		`"aa"`)
	lexText(t,
		"'a'a",
		`"aa"`)
	// '' is a literal apostrophe outside quoted text.
	lexText(t,
		"it''s",
		`"it's"`)
	// '' is a literal apostrophe inside quoted text, and the quoted text continues.
	lexText(t,
		"This '{isn''t}' obvious",
		`"This {isn't} obvious"`)
	// Two consecutive quoted texts, separated by a literal apostrophe.
	lexText(t,
		"'{a}''{b}'",
		`"{a}'{b}"`)
	// Quoting the syntax characters.
	lexText(t,
		"'{'",
		`"{"`)
	lexText(t,
		"'}'",
		`"}"`)
	lexText(t,
		"'#'",
		`"#"`)
	// Apostrophe-doubling is resolved before quoting, so this is two
	// literal apostrophes, not a quoted text containing one.
	lexText(t,
		"'''' ",
		`"'' "`)
}

func TestLexTextError(t *testing.T) {
	lexTextError := func(input string, expected error) {
		l := newLexer(input)
		var err error
		for {
			err = l.LexText(&bytes.Buffer{})
			if err != nil {
				break
			}
			var eof bool
			for _, token := range l.Output {
				if token.Type == TokenTypeEOF {
					eof = true
				}
			}
			l.Output = nil
			if eof {
				break
			}
		}
		if !errors.Is(err, expected) {
			t.Errorf("%v: expected %v, actual %v\n", input, expected, err)
		}
	}

	lexTextError("'", ErrUnterminatedQuotedString)
	lexTextError("a'", ErrUnterminatedQuotedString)
	lexTextError("'''", ErrUnterminatedQuotedString)
}

func TestLexArg(t *testing.T) {
	lexArg(t,
		"{ arg, plural, offset:1 =0 {} =1 {} one{} other{} }",
		"{", "arg", ",", "plural", ",", "offset", ":", "1", "=", "0", "{", "}", "=", "1", "{", "}", "one", "{", "}", "other", "{", "}", "}")
}
