package messageformat

import (
	"fmt"
	"strconv"
	templateparse "text/template/parse"

	"golang.org/x/text/language"
)

const TemplateRuntimeFuncName = "__messageformat__"

func TemplateRuntimeFunc(typ string, args ...interface{}) bool {
	switch typ {
	case "select":
		value := args[0]
		valueString, err := formatValue(value)
		if err != nil {
			panic(fmt.Errorf("messageformat: failed to cast value to string: %w", err))
		}
		keyword := args[1].(string)
		return valueString == keyword
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

func (f *templateParseTreeFormatter) Format(root *templateparse.ListNode, nodes []Node, argMinusOffset *argumentMinusOffset) (err error) {
	for _, node := range nodes {
		switch node := node.(type) {
		case TextNode:
			err = f.FormatTextNode(root, node)
		case NoneArgNode:
			err = f.FormatNoneArgNode(root, node)
		case SelectArgNode:
			err = f.FormatSelectArgNode(root, node)
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
