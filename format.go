package messageformat

import (
	"fmt"
	"strconv"
	"strings"

	"golang.org/x/text/language"
)

func FormatPositional(tag language.Tag, pattern string, args ...interface{}) (out string, err error) {
	nodes, err := Parse(pattern)
	if err != nil {
		return
	}

	formatter := &formatter{
		Buf:  &strings.Builder{},
		Tag:  tag,
		Args: sliceArgsToNamedArgs(args),
	}

	err = formatter.Format(nodes, nil)
	if err != nil {
		return
	}

	out = formatter.Buf.String()
	return
}

func sliceArgsToNamedArgs(args []interface{}) (out map[string]interface{}) {
	out = make(map[string]interface{})
	for idx, val := range args {
		name := strconv.Itoa(idx)
		out[name] = val
	}
	return
}

type formatter struct {
	Buf  *strings.Builder
	Tag  language.Tag
	Args map[string]interface{}
}

func (f *formatter) Format(nodes []Node, pluralArgument *Argument) (err error) {
	for _, inode := range nodes {
		switch node := inode.(type) {
		case TextNode:
			err = f.FormatTextNode(node)
		case NoneArgNode:
			err = f.FormatNoneArgNode(node)
		case SelectArgNode:
			err = f.FormatSelectArgNode(node)
		}
		if err != nil {
			return
		}
	}
	return
}

func (f *formatter) ResolveArgument(arg Argument) (name string, value interface{}, err error) {
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

func (f *formatter) FormatValue(argName string, value interface{}) (out string, err error) {
	switch v := value.(type) {
	case int8:
		out = strconv.FormatInt(int64(v), 10)
	case int16:
		out = strconv.FormatInt(int64(v), 10)
	case int32:
		out = strconv.FormatInt(int64(v), 10)
	case int64:
		out = strconv.FormatInt(v, 10)
	case int:
		out = strconv.FormatInt(int64(v), 10)
	case uint8:
		out = strconv.FormatUint(uint64(v), 10)
	case uint16:
		out = strconv.FormatUint(uint64(v), 10)
	case uint32:
		out = strconv.FormatUint(uint64(v), 10)
	case uint64:
		out = strconv.FormatUint(v, 10)
	case uint:
		out = strconv.FormatUint(uint64(v), 10)
	case float32:
		out = strconv.FormatFloat(float64(v), 'f', -1, 32)
	case float64:
		out = strconv.FormatFloat(v, 'f', -1, 64)
	case string:
		out = v
	case bool:
		out = strconv.FormatBool(v)
	default:
		err = fmt.Errorf("unsupported argument type: %v %T", argName, value)
	}
	return
}

func (f *formatter) FormatTextNode(node TextNode) (err error) {
	f.Buf.WriteString(node.Value)
	return
}

func (f *formatter) FormatNoneArgNode(node NoneArgNode) (err error) {
	var argName string
	var argValue interface{}
	argName, argValue, err = f.ResolveArgument(node.Arg)
	if err != nil {
		return
	}

	var stringValue string
	stringValue, err = f.FormatValue(argName, argValue)
	if err != nil {
		return
	}

	f.Buf.WriteString(stringValue)
	return
}

func (f *formatter) FormatSelectArgNode(node SelectArgNode) (err error) {
	var argName string
	var argValue interface{}
	argName, argValue, err = f.ResolveArgument(node.Arg)
	if err != nil {
		return
	}

	var stringValue string
	stringValue, err = f.FormatValue(argName, argValue)
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
