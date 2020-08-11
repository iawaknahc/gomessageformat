package messageformat

import (
	"fmt"
	"strconv"
	templateparse "text/template/parse"

	"golang.org/x/text/language"
)

type argumentOffset struct {
	Name   string
	Offset int
}

const TemplateRuntimeFuncName = "__messageformat__"

func TemplateRuntimeFunc(typ string, args ...interface{}) interface{} {
	switch typ {
	case "select":
		value := args[0]
		valueString, err := formatValue(value)
		if err != nil {
			panic(fmt.Errorf("messageformat: failed to cast value to string: %w", err))
		}
		keyword := args[1].(string)
		return valueString == keyword
	case "plural":
		tag := args[0].(string)
		offset := args[1].(int)
		value := args[2]
		keyword := args[3].(string)
		explicitVaue := args[4].(int)

		if keyword == "" {
			match, err := matchExplicitValue(value, explicitVaue)
			if err != nil {
				panic(fmt.Errorf("messageformat: failed to match explicit value: %w", err))
			}
			return match
		}

		offsetValue, err := offsetValue(value, offset)
		if err != nil {
			panic(fmt.Errorf("messageformat: failed to compute offset value: %w", err))
		}
		pluralForm, err := Cardinal(language.Make(tag), offsetValue)
		if err != nil {
			panic(fmt.Errorf("messageformat: failed to compute plural form: %w", err))
		}
		return pluralForm == keyword
	case "selectordinal":
		tag := args[0].(string)
		offset := args[1].(int)
		value := args[2]
		keyword := args[3].(string)
		explicitVaue := args[4].(int)

		if keyword == "" {
			match, err := matchExplicitValue(value, explicitVaue)
			if err != nil {
				panic(fmt.Errorf("messageformat: failed to match explicit value: %w", err))
			}
			return match
		}

		offsetValue, err := offsetValue(value, offset)
		if err != nil {
			panic(fmt.Errorf("messageformat: failed to compute offset value: %w", err))
		}
		pluralForm, err := Ordinal(language.Make(tag), offsetValue)
		if err != nil {
			panic(fmt.Errorf("messageformat: failed to compute plural form: %w", err))
		}
		return pluralForm == keyword
	case "pound":
		value := args[0]
		offset := args[1].(int)
		offsetValue, err := offsetValue(value, offset)
		if err != nil {
			panic(fmt.Errorf("messageformat: failed to compute offset value: %w", err))
		}
		offsetValueString, err := formatValue(offsetValue)
		if err != nil {
			panic(fmt.Errorf("messageformat: failed to cast offset value to string: %w", err))
		}
		return offsetValueString
	default:
		panic("messageformat: unexpected argument type: " + typ)
	}
}

func FormatTemplateParseTree(tag language.Tag, pattern string) (tree *templateparse.Tree, err error) {
	nodes, err := Parse(pattern)
	if err != nil {
		return
	}

	parseTree := templateparse.New("tree", nil)
	parseTree.Root = &templateparse.ListNode{
		NodeType: templateparse.NodeList,
	}

	formatter := &templateParseTreeFormatter{
		Tree: parseTree,
		Tag:  tag,
	}

	err = formatter.Format(formatter.Tree.Root, nodes, nil)
	if err != nil {
		return
	}

	tree = formatter.Tree
	return
}

type templateParseTreeFormatter struct {
	Tree *templateparse.Tree
	Tag  language.Tag
}

func (f *templateParseTreeFormatter) Format(root *templateparse.ListNode, nodes []Node, argOffset *argumentOffset) (err error) {
	for _, node := range nodes {
		switch node := node.(type) {
		case TextNode:
			err = f.FormatTextNode(root, node)
		case NoneArgNode:
			err = f.FormatNoneArgNode(root, node)
		case SelectArgNode:
			err = f.FormatSelectArgNode(root, node)
		case PluralArgNode:
			err = f.FormatPluralArgNode(root, node)
		case PoundNode:
			err = f.FormatPoundNode(root, argOffset)
		}
		if err != nil {
			return
		}
	}
	return
}

func (f *templateParseTreeFormatter) FormatTextNode(root *templateparse.ListNode, node TextNode) (err error) {
	root.Nodes = append(root.Nodes, &templateparse.TextNode{
		NodeType: templateparse.NodeText,
		Text:     []byte(node.Value),
	})
	return
}

func (f *templateParseTreeFormatter) FormatNoneArgNode(root *templateparse.ListNode, node NoneArgNode) (err error) {
	root.Nodes = append(root.Nodes, &templateparse.ActionNode{
		NodeType: templateparse.NodeAction,
		Pipe: &templateparse.PipeNode{
			NodeType: templateparse.NodePipe,
			Cmds: []*templateparse.CommandNode{
				&templateparse.CommandNode{
					NodeType: templateparse.NodeCommand,
					Args: []templateparse.Node{
						&templateparse.FieldNode{
							NodeType: templateparse.NodeField,
							Ident:    []string{node.Arg.Name},
						},
					},
				},
			},
		},
	})
	return
}

func (f *templateParseTreeFormatter) FormatSelectArgNode(root *templateparse.ListNode, node SelectArgNode) (err error) {
	var nonOtherClauses []SelectClause
	var otherClause *SelectClause

	for _, clause := range node.Clauses {
		c := clause
		if clause.Keyword == "other" {
			otherClause = &c
		} else {
			nonOtherClauses = append(nonOtherClauses, c)
		}
	}
	if otherClause == nil {
		err = fmt.Errorf("missing select other clause: %v", node.Arg.Name)
	}

	currRoot := root
	for _, clause := range nonOtherClauses {
		node := &templateparse.IfNode{
			BranchNode: templateparse.BranchNode{
				NodeType: templateparse.NodeIf,
				// This is the if condition
				Pipe: &templateparse.PipeNode{
					NodeType: templateparse.NodePipe,
					Cmds: []*templateparse.CommandNode{
						&templateparse.CommandNode{
							NodeType: templateparse.NodeCommand,
							Args: []templateparse.Node{
								&templateparse.IdentifierNode{
									NodeType: templateparse.NodeIdentifier,
									Ident:    TemplateRuntimeFuncName,
								},
								&templateparse.StringNode{
									NodeType: templateparse.NodeString,
									Quoted:   strconv.Quote("select"),
									Text:     "select",
								},
								&templateparse.FieldNode{
									NodeType: templateparse.NodeField,
									Ident:    []string{node.Arg.Name},
								},
								&templateparse.StringNode{
									NodeType: templateparse.NodeString,
									Quoted:   strconv.Quote(clause.Keyword),
									Text:     clause.Keyword,
								},
							},
						},
					},
				},
				List: &templateparse.ListNode{
					NodeType: templateparse.NodeList,
					Nodes:    []templateparse.Node{},
				},
				// If there is "else if", it contains ONE IfNode.
				// If there is "else", it contains child.
				// Otherwise ElseList itself is nil.
				ElseList: &templateparse.ListNode{
					NodeType: templateparse.NodeList,
					Nodes:    []templateparse.Node{},
				},
			},
		}

		// Recursively format the if body.
		err = f.Format(node.BranchNode.List, clause.Nodes, nil)
		if err != nil {
			return
		}

		// Append the node to current root.
		currRoot.Nodes = append(currRoot.Nodes, node)

		// Adjust current root to the ElseList of the node.
		currRoot = node.BranchNode.ElseList
	}

	// Construct the final else.
	err = f.Format(currRoot, otherClause.Nodes, nil)
	if err != nil {
		return
	}

	return
}

func (f *templateParseTreeFormatter) FormatPluralArgNode(root *templateparse.ListNode, node PluralArgNode) (err error) {
	var nonOtherClauses []PluralClause
	var otherClause *PluralClause

	for _, clause := range node.Clauses {
		c := clause
		if clause.Keyword == "other" {
			otherClause = &c
		} else {
			nonOtherClauses = append(nonOtherClauses, c)
		}
	}
	if otherClause == nil {
		err = fmt.Errorf("missing select other clause: %v", node.Arg.Name)
	}

	argOffset := &argumentOffset{
		Name:   node.Arg.Name,
		Offset: node.Offset,
	}

	currRoot := root
	for _, clause := range nonOtherClauses {
		node := &templateparse.IfNode{
			BranchNode: templateparse.BranchNode{
				NodeType: templateparse.NodeIf,
				// This is the if condition
				Pipe: &templateparse.PipeNode{
					NodeType: templateparse.NodePipe,
					Cmds: []*templateparse.CommandNode{
						&templateparse.CommandNode{
							NodeType: templateparse.NodeCommand,
							Args: []templateparse.Node{
								&templateparse.IdentifierNode{
									NodeType: templateparse.NodeIdentifier,
									Ident:    TemplateRuntimeFuncName,
								},
								&templateparse.StringNode{
									NodeType: templateparse.NodeString,
									Quoted:   strconv.Quote(node.Kind),
									Text:     node.Kind,
								},
								&templateparse.StringNode{
									NodeType: templateparse.NodeString,
									Quoted:   strconv.Quote(f.Tag.String()),
									Text:     f.Tag.String(),
								},
								makeNumberNode(node.Offset),
								&templateparse.FieldNode{
									NodeType: templateparse.NodeField,
									Ident:    []string{node.Arg.Name},
								},
								&templateparse.StringNode{
									NodeType: templateparse.NodeString,
									Quoted:   strconv.Quote(clause.Keyword),
									Text:     clause.Keyword,
								},
								makeNumberNode(clause.ExplicitValue),
							},
						},
					},
				},
				List: &templateparse.ListNode{
					NodeType: templateparse.NodeList,
					Nodes:    []templateparse.Node{},
				},
				// If there is "else if", it contains ONE IfNode.
				// If there is "else", it contains child.
				// Otherwise ElseList itself is nil.
				ElseList: &templateparse.ListNode{
					NodeType: templateparse.NodeList,
					Nodes:    []templateparse.Node{},
				},
			},
		}

		// Recursively format the if body.
		err = f.Format(node.BranchNode.List, clause.Nodes, argOffset)
		if err != nil {
			return
		}

		// Append the node to current root.
		currRoot.Nodes = append(currRoot.Nodes, node)

		// Adjust current root to the ElseList of the node.
		currRoot = node.BranchNode.ElseList
	}

	// Construct the final else.
	err = f.Format(currRoot, otherClause.Nodes, argOffset)
	if err != nil {
		return
	}

	return
}

func (f *templateParseTreeFormatter) FormatPoundNode(root *templateparse.ListNode, argOffset *argumentOffset) (err error) {
	root.Nodes = append(root.Nodes, &templateparse.ActionNode{
		NodeType: templateparse.NodeAction,
		Pipe: &templateparse.PipeNode{
			NodeType: templateparse.NodePipe,
			Cmds: []*templateparse.CommandNode{
				&templateparse.CommandNode{
					NodeType: templateparse.NodeCommand,
					Args: []templateparse.Node{
						&templateparse.IdentifierNode{
							NodeType: templateparse.NodeIdentifier,
							Ident:    TemplateRuntimeFuncName,
						},
						&templateparse.StringNode{
							NodeType: templateparse.NodeString,
							Quoted:   strconv.Quote("pound"),
							Text:     "pound",
						},
						&templateparse.FieldNode{
							NodeType: templateparse.NodeField,
							Ident:    []string{argOffset.Name},
						},
						makeNumberNode(argOffset.Offset),
					},
				},
			},
		},
	})
	return
}

func makeNumberNode(offset int) *templateparse.NumberNode {
	node := &templateparse.NumberNode{
		NodeType: templateparse.NodeNumber,
	}

	text, err := formatValue(offset)
	if err != nil {
		panic(fmt.Errorf("messageformat: failed to make number node: %w", err))
	}

	node.IsInt = true
	node.Int64 = int64(offset)

	if offset >= 0 {
		node.IsUint = true
		node.Uint64 = uint64(offset)
	}

	node.IsFloat = true
	node.Float64 = float64(offset)

	node.Text = text

	return node
}
