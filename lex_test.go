package messageformat

import (
	"bytes"
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
}

func TestLexArg(t *testing.T) {
	lexArg(t,
		"{ arg, plural, offset:1 =0 {} =1 {} one{} other{} }",
		"{", "arg", ",", "plural", ",", "offset", ":", "1", "=", "0", "{", "}", "=", "1", "{", "}", "one", "{", "}", "other", "{", "}", "}")
}
