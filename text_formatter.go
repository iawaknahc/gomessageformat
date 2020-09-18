package messageformat

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"golang.org/x/text/language"

	"github.com/iawaknahc/gomessageformat/icu4c"
)

// FormatPositional parses pattern and format to string with a slice of args.
func FormatPositional(tag language.Tag, pattern string, args ...interface{}) (out string, err error) {
	o := make(map[string]interface{})
	for idx, val := range args {
		name := strconv.Itoa(idx)
		o[name] = val
	}
	return FormatNamed(tag, pattern, o)
}

// FormatNamed parses pattern and format to string with a map of args.
func FormatNamed(tag language.Tag, pattern string, args map[string]interface{}) (out string, err error) {
	nodes, err := Parse(pattern)
	if err != nil {
		return
	}

	formatter := &textFormatter{
		Buf:  &strings.Builder{},
		Tag:  tag,
		Args: args,
	}

	err = formatter.Format(nodes, nil)
	if err != nil {
		return
	}

	out = formatter.Buf.String()
	return
}

type argumentMinusOffset struct {
	Name  string
	Value interface{}
}

type textFormatter struct {
	Buf  *strings.Builder
	Tag  language.Tag
	Args map[string]interface{}
}

func (f *textFormatter) Format(nodes []Node, argMinusOffset *argumentMinusOffset) (err error) {
	for _, inode := range nodes {
		switch node := inode.(type) {
		case TextNode:
			err = f.FormatTextNode(node)
		case NoneArgNode:
			err = f.FormatNoneArgNode(node)
		case DateArgNode:
			err = f.FormatDateArgNode(node)
		case TimeArgNode:
			err = f.FormatTimeArgNode(node)
		case DatetimeArgNode:
			err = f.FormatDatetimeArgNode(node)
		case SelectArgNode:
			err = f.FormatSelectArgNode(node)
		case PluralArgNode:
			err = f.FormatPluralArgNode(node)
		case PoundNode:
			err = f.FormatPoundNode(argMinusOffset)
		}
		if err != nil {
			return
		}
	}
	return
}

func (f *textFormatter) ResolveArgument(arg Argument) (name string, value interface{}, err error) {
	name = arg.Name
	if name == "" {
		name = strconv.Itoa(arg.Index)
	}

	value, ok := f.Args[name]
	if !ok {
		err = fmt.Errorf("unknown argument: %v", name)
		return
	}

	return
}

func (f *textFormatter) FormatValue(argName string, value interface{}) (out string, err error) {
	out, err = formatValue(value)
	if err != nil {
		err = fmt.Errorf("unsupported argument type: %v %T", argName, value)
		return
	}
	return
}

func (f *textFormatter) FormatTextNode(node TextNode) (err error) {
	f.Buf.WriteString(node.Value)
	return
}

func (f *textFormatter) FormatNoneArgNode(node NoneArgNode) (err error) {
	argName, argValue, err := f.ResolveArgument(node.Arg)
	if err != nil {
		return
	}

	stringValue, err := f.FormatValue(argName, argValue)
	if err != nil {
		return
	}

	f.Buf.WriteString(stringValue)
	return
}

func (f *textFormatter) FormatDateArgNode(node DateArgNode) (err error) {
	argName, argValue, err := f.ResolveArgument(node.Arg)
	if err != nil {
		return
	}

	var t *time.Time
	switch v := argValue.(type) {
	case time.Time:
		t = &v
	case *time.Time:
		t = v
	}

	if t == nil {
		err = fmt.Errorf("expected %v (%T) to be time.Time", argName, argValue)
		return
	}

	tz := icu4c.TZName("UTC")
	style := styleToStyle(node.Style)
	out, err := icu4c.FormatDatetime(f.Tag, tz, style, icu4c.DateFormatStyleNone, *t)
	if err != nil {
		return
	}

	f.Buf.WriteString(out)
	return
}

func (f *textFormatter) FormatTimeArgNode(node TimeArgNode) (err error) {
	argName, argValue, err := f.ResolveArgument(node.Arg)
	if err != nil {
		return
	}

	var t *time.Time
	switch v := argValue.(type) {
	case time.Time:
		t = &v
	case *time.Time:
		t = v
	}

	if t == nil {
		err = fmt.Errorf("expected %v (%T) to be time.Time", argName, argValue)
		return
	}

	tz := icu4c.TZName("UTC")
	style := styleToStyle(node.Style)
	out, err := icu4c.FormatDatetime(f.Tag, tz, icu4c.DateFormatStyleNone, style, *t)
	if err != nil {
		return
	}

	f.Buf.WriteString(out)
	return
}

func (f *textFormatter) FormatDatetimeArgNode(node DatetimeArgNode) (err error) {
	argName, argValue, err := f.ResolveArgument(node.Arg)
	if err != nil {
		return
	}

	var t *time.Time
	switch v := argValue.(type) {
	case time.Time:
		t = &v
	case *time.Time:
		t = v
	}

	if t == nil {
		err = fmt.Errorf("expected %v (%T) to be time.Time", argName, argValue)
		return
	}

	tz := icu4c.TZName("UTC")
	style := styleToStyle(node.Style)
	out, err := icu4c.FormatDatetime(f.Tag, tz, style, style, *t)
	if err != nil {
		return
	}

	f.Buf.WriteString(out)
	return
}

func (f *textFormatter) FormatSelectArgNode(node SelectArgNode) (err error) {
	argName, argValue, err := f.ResolveArgument(node.Arg)
	if err != nil {
		return
	}

	stringValue, err := f.FormatValue(argName, argValue)
	if err != nil {
		return
	}

	var otherClause *SelectClause
	var done bool

	for _, clause := range node.Clauses {
		// Remember the other clause for fallback purpose.
		if clause.Keyword == "other" {
			c := clause
			otherClause = &c
		}

		if clause.Keyword == stringValue {
			done = true
			f.Format(clause.Nodes, nil)
			break
		}
	}

	if !done {
		if otherClause == nil {
			err = fmt.Errorf("missing select other clause: %v", argName)
			return
		}
		f.Format(otherClause.Nodes, nil)
	}

	return
}

func (f *textFormatter) FormatPluralArgNode(node PluralArgNode) (err error) {
	argName, argValue, err := f.ResolveArgument(node.Arg)
	if err != nil {
		return
	}

	offsetValue, err := f.OffsetValue(argName, argValue, node.Offset)
	if err != nil {
		return
	}

	argumentMinusOffset := &argumentMinusOffset{
		Name:  argName,
		Value: offsetValue,
	}

	var otherClause *PluralClause
	var done bool

	for _, clause := range node.Clauses {
		if clause.Keyword == "other" {
			c := clause
			otherClause = &c
		}

		if clause.Keyword == "" {
			var match bool
			match, err = f.MatchExplicitValue(argName, argValue, clause.ExplicitValue)
			if err != nil {
				return
			}
			if match {
				done = true
				f.Format(clause.Nodes, argumentMinusOffset)
				break
			}
		}
	}

	if done {
		return
	}

	pluralFunc := Cardinal
	if node.Kind == "selectordinal" {
		pluralFunc = Ordinal
	}
	pluralForm, err := pluralFunc(f.Tag, offsetValue)
	if err != nil {
		return
	}

	for _, clause := range node.Clauses {
		if clause.Keyword == pluralForm {
			done = true
			f.Format(clause.Nodes, argumentMinusOffset)
			break
		}
	}

	if !done {
		if otherClause == nil {
			err = fmt.Errorf("missing plural other clause: %v", argName)
			return
		}
		f.Format(otherClause.Nodes, argumentMinusOffset)
	}

	return
}

func (f *textFormatter) OffsetValue(argName string, value interface{}, offset int) (out interface{}, err error) {
	out, err = offsetValue(value, offset)
	if err != nil {
		err = fmt.Errorf("expected numeric type: %v %T", argName, value)
		return
	}
	return
}

func (f *textFormatter) MatchExplicitValue(argName string, value interface{}, explicitValue int) (match bool, err error) {
	match, err = matchExplicitValue(value, explicitValue)
	if err != nil {
		err = fmt.Errorf("expected numeric type: %v %T", argName, value)
		return
	}
	return
}

func (f *textFormatter) FormatPoundNode(argumentMinusOffset *argumentMinusOffset) (err error) {
	if argumentMinusOffset == nil {
		err = fmt.Errorf("lexer emitted pound token incorrectly")
		return
	}
	out, err := f.FormatValue(argumentMinusOffset.Name, argumentMinusOffset.Value)
	if err != nil {
		return
	}
	f.Buf.WriteString(out)
	return
}
