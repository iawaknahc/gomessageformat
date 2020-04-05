package messageformat

import (
	"golang.org/x/text/feature/plural"
	"golang.org/x/text/language"
)

func formToString(form plural.Form) string {
	switch form {
	case plural.Other:
		return "other"
	case plural.Zero:
		return "zero"
	case plural.One:
		return "one"
	case plural.Two:
		return "two"
	case plural.Few:
		return "few"
	case plural.Many:
		return "many"
	}
	panic("unreachable")
}

func Cardinal(lang language.Tag, number interface{}) (out string, err error) {
	i, v, w, f, t, err := IVWFT(number)
	if err != nil {
		return
	}
	form := plural.Cardinal.MatchPlural(lang, i, v, w, f, t)
	return formToString(form), nil
}

func Ordinal(lang language.Tag, number interface{}) (out string, err error) {
	i, v, w, f, t, err := IVWFT(number)
	if err != nil {
		return
	}
	form := plural.Ordinal.MatchPlural(lang, i, v, w, f, t)
	return formToString(form), nil
}
