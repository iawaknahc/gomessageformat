package messageformat

import (
	"fmt"
	"strconv"
	"strings"

	"golang.org/x/text/language"
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

func (f *textFormatter) OffsetValue(argName string, value interface{}, offset int) (offsetValue interface{}, err error) {
	switch v := value.(type) {
	case int8:
		offsetValue = int8(int64(v) - int64(offset))
	case int16:
		offsetValue = int16(int64(v) - int64(offset))
	case int32:
		offsetValue = int32(int64(v) - int64(offset))
	case int64:
		offsetValue = int64(int64(v) - int64(offset))
	case int:
		offsetValue = int(int64(v) - int64(offset))
	case uint8:
		offsetValue = uint8(int64(v) - int64(offset))
	case uint16:
		offsetValue = uint16(int64(v) - int64(offset))
	case uint32:
		offsetValue = uint32(int64(v) - int64(offset))
	case uint64:
		offsetValue = uint64(int64(v) - int64(offset))
	case uint:
		offsetValue = uint(int64(v) - int64(offset))
	case float32:
		offsetValue = float32(float32(v) - float32(offset))
	case float64:
		offsetValue = float64(float64(v) - float64(offset))
	case string:
		var f64 float64
		f64, err = strconv.ParseFloat(v, 64)
		if err != nil {
			return
		}
		offsetValue = strconv.FormatFloat(f64-float64(offset), 'f', -1, 64)
	default:
		err = fmt.Errorf("expected numeric type: %v %T", argName, value)
	}
	return
}

func (f *textFormatter) MatchExplicitValue(argName string, value interface{}, explicitValue int) (match bool, err error) {
	switch v := value.(type) {
	case int8:
		match = int64(v) == int64(explicitValue)
	case int16:
		match = int64(v) == int64(explicitValue)
	case int32:
		match = int64(v) == int64(explicitValue)
	case int64:
		match = int64(v) == int64(explicitValue)
	case int:
		match = v == explicitValue
	case uint8:
		match = int64(v) == int64(explicitValue)
	case uint16:
		match = int64(v) == int64(explicitValue)
	case uint32:
		match = int64(v) == int64(explicitValue)
	case uint64:
		match = int64(v) == int64(explicitValue)
	case uint:
		match = int64(v) == int64(explicitValue)
	case float32:
		match = float32(v) == float32(explicitValue)
	case float64:
		match = float64(v) == float64(explicitValue)
	case string:
		var f64 float64
		f64, err = strconv.ParseFloat(v, 64)
		if err != nil {
			return
		}
		match = f64 == float64(explicitValue)
	default:
		err = fmt.Errorf("expected numeric type: %v %T", argName, value)
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
