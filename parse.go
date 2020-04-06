package messageformat

import (
	"fmt"
	"strconv"
)

// Argument is either named argument or positional argument.
type Argument struct {
	Name  string
	Index int
}

// Node is an semantic item.
type Node interface {
	messageFormatNode()
}

// TextNode is a text segment.
type TextNode struct {
	Value string
}

func (_ TextNode) messageFormatNode() {}

// NoneArgNode is `{Argument}`.
type NoneArgNode struct {
	Arg Argument
}

func (_ NoneArgNode) messageFormatNode() {}

// SelectClause is `Keyword {message}`.
type SelectClause struct {
	Keyword string
	Nodes   []Node
}

// SelectArgNode is `{Argument, select, SelectClause+}`.
type SelectArgNode struct {
	Arg     Argument
	Clauses []SelectClause
}

func (_ SelectArgNode) messageFormatNode() {}

// PluralClause is `(keyword | =ExplicitValue) {message}`.
type PluralClause struct {
	// If Keyword is empty, use ExplicitValue
	Keyword       string
	ExplicitValue int
	Nodes         []Node
}

// PluralArgNode is `{Argument, plural | selectordinal, [offset:number] PluralClause+}`.
type PluralArgNode struct {
	Arg Argument
	// plural or selectordinal
	Kind    string
	Offset  int
	Clauses []PluralClause
}

func (_ PluralArgNode) messageFormatNode() {}

// PoundNode is `#`.
type PoundNode struct{}

func (_ PoundNode) messageFormatNode() {}

// Parse parses the pattern s into message.
func Parse(s string) ([]Node, error) {
	p := parser{lexer: newLexer(s)}
	p.lexer.isInPluralStyle = p.isInPluralStyle
	return p.parse(s)
}

type parser struct {
	lexer      *lexer
	tokens     []Token
	poundStack []bool
}

func (p *parser) pushPoundStack(pound bool) {
	p.poundStack = append(p.poundStack, pound)
}

func (p *parser) popPoundStack() {
	p.poundStack = p.poundStack[0 : len(p.poundStack)-1]
}

func (p *parser) isInPluralStyle() bool {
	if len(p.poundStack) <= 0 {
		return false
	}
	top := p.poundStack[len(p.poundStack)-1]
	return top
}

func (p *parser) parse(s string) ([]Node, error) {
	nodes, err := p.parseMessage(TokenTypeEOF, false)
	if err != nil {
		return nil, err
	}
	return nodes, nil
}

func (p *parser) next() (*Token, error) {
	if len(p.tokens) <= 0 {
		err := p.lexer.Lex()
		if err != nil {
			return nil, err
		}
		p.tokens = p.lexer.Output
		p.lexer.Output = nil
	}
	t := p.tokens[0]
	p.tokens = p.tokens[1:]
	return &t, nil
}

func (p *parser) putBack(t *Token) {
	p.tokens = append([]Token{*t}, p.tokens...)
}

func (p *parser) expect(types ...TokenType) (*Token, error) {
	token, err := p.next()
	if err != nil {
		return nil, err
	}
	for _, t := range types {
		if token.Type == t {
			return token, nil
		}
	}
	return nil, fmt.Errorf("unexpected token: %v", token)
}

func (p *parser) expectWord(words ...string) (*Token, error) {
	word, err := p.expect(TokenTypeWord)
	if err != nil {
		return nil, err
	}
	for _, w := range words {
		if word.Value == w {
			return word, nil
		}
	}
	return nil, fmt.Errorf("unexpected token: %v", word)
}

func (p *parser) parseMessage(endToken TokenType, pound bool) ([]Node, error) {
	p.pushPoundStack(pound)
	defer p.popPoundStack()
	textNodes, err := p.parseMessageText()
	if err != nil {
		return nil, err
	}
	nodes, err := p.parseArgMessageText(endToken)
	if err != nil {
		return nil, err
	}

	var out []Node
	out = append(out, textNodes...)
	out = append(out, nodes...)
	return out, nil
}

func (p *parser) parseMessageText() ([]Node, error) {
	textNode, err := p.parseMessageText0()
	if err != nil {
		return nil, err
	}

	var nodes []Node
	for {
		token, err := p.next()
		if err != nil {
			return nil, err
		}
		if token.Type != TokenTypePound {
			p.putBack(token)
			break
		}
		nodes = append(nodes, PoundNode{})
		textNode, err := p.parseMessageText0()
		if err != nil {
			return nil, err
		}
		nodes = append(nodes, textNode)
	}

	var out []Node
	out = append(out, textNode)
	out = append(out, nodes...)
	return out, nil
}

func (p *parser) parseMessageText0() (Node, error) {
	text, err := p.expect(TokenTypeText)
	if err != nil {
		return nil, err
	}
	return TextNode{Value: text.Value}, nil
}

func (p *parser) parseArgMessageText(endToken TokenType) ([]Node, error) {
	lbraceOrEnd, err := p.expect(TokenTypeLBrace, endToken)
	if err != nil {
		return nil, err
	}
	if lbraceOrEnd.Type == endToken {
		return nil, err
	}
	argNode, err := p.parseArg()
	if err != nil {
		return nil, err
	}
	textNodes, err := p.parseMessageText()
	if err != nil {
		return nil, err
	}
	nodes, err := p.parseArgMessageText(endToken)
	if err != nil {
		return nil, err
	}
	out := []Node{argNode}
	out = append(out, textNodes...)
	out = append(out, nodes...)
	return out, nil
}

func (p *parser) parseArg() (Node, error) {
	argNameOrNumber, err := p.expect(TokenTypeWord, TokenTypeNumber)
	if err != nil {
		return nil, err
	}

	var arg Argument
	if argNameOrNumber.Type == TokenTypeWord {
		arg.Name = argNameOrNumber.Value
	} else {
		i, err := strconv.Atoi(argNameOrNumber.Value)
		if err != nil {
			return nil, err
		}
		arg.Index = i
	}

	rbraceOrComma, err := p.expect(TokenTypeRBrace, TokenTypeComma)
	if err != nil {
		return nil, err
	}

	if rbraceOrComma.Type == TokenTypeRBrace {
		return NoneArgNode{Arg: arg}, nil
	}

	argType, err := p.expectWord("plural", "select", "selectordinal")
	if err != nil {
		return nil, err
	}

	_, err = p.expect(TokenTypeComma)
	if err != nil {
		return nil, err
	}

	switch argType.Value {
	case "plural":
		offset, clauses, err := p.parsePluralStyle()
		if err != nil {
			return nil, err
		}
		return PluralArgNode{Arg: arg, Kind: "plural", Offset: offset, Clauses: clauses}, nil
	case "selectordinal":
		offset, clauses, err := p.parsePluralStyle()
		if err != nil {
			return nil, err
		}
		return PluralArgNode{Arg: arg, Kind: "selectordinal", Offset: offset, Clauses: clauses}, nil
	case "select":
		clauses, err := p.parseSelectStyle()
		if err != nil {
			return nil, err
		}
		return SelectArgNode{Arg: arg, Clauses: clauses}, nil
	}

	panic("unreachable")
}

func (p *parser) parsePluralStyle() (offset int, clauses []PluralClause, err error) {
	for {
		var token *Token
		token, err = p.expect(TokenTypeRBrace, TokenTypeWord, TokenTypeEqual)
		if err != nil {
			return
		}

		// Handle offset: before any clauses
		if len(clauses) <= 0 && token.Type == TokenTypeWord && token.Value == "offset" {
			_, err = p.expect(TokenTypeColon)
			if err != nil {
				return
			}
			var number *Token
			number, err = p.expect(TokenTypeNumber)
			if err != nil {
				return
			}
			offset, err = strconv.Atoi(number.Value)
			if err != nil {
				return
			}
			token, err = p.expect(TokenTypeRBrace, TokenTypeWord, TokenTypeEqual)
			if err != nil {
				return
			}
		}

		if token.Type == TokenTypeRBrace {
			if len(clauses) <= 0 {
				err = fmt.Errorf("no plural clauses")
			}
			return
		}

		clause := PluralClause{}

		// Handle explicit value or keyword
		if token.Type == TokenTypeEqual {
			var number *Token
			number, err = p.expect(TokenTypeNumber)
			if err != nil {
				return
			}
			var value int
			value, err = strconv.Atoi(number.Value)
			if err != nil {
				return
			}
			clause.ExplicitValue = value
		} else {
			clause.Keyword = token.Value
		}

		_, err = p.expect(TokenTypeLBrace)
		if err != nil {
			return
		}

		var nodes []Node
		nodes, err = p.parseMessage(TokenTypeRBrace, true)
		if err != nil {
			return
		}
		clause.Nodes = nodes

		clauses = append(clauses, clause)
	}
}

func (p *parser) parseSelectStyle() ([]SelectClause, error) {
	var clauses []SelectClause
	for {
		rbraceOrWord, err := p.expect(TokenTypeRBrace, TokenTypeWord)
		if err != nil {
			return nil, err
		}
		if rbraceOrWord.Type == TokenTypeRBrace {
			if len(clauses) <= 0 {
				return nil, fmt.Errorf("no select clauses")
			}
			return clauses, nil
		}

		_, err = p.expect(TokenTypeLBrace)
		if err != nil {
			return nil, err
		}

		nodes, err := p.parseMessage(TokenTypeRBrace, false)
		if err != nil {
			return nil, err
		}

		clauses = append(clauses, SelectClause{
			Keyword: rbraceOrWord.Value,
			Nodes:   nodes,
		})
	}
}
