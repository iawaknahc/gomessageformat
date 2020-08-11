package messageformat

import (
	templateparse "text/template/parse"

	"golang.org/x/text/language"
)

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
