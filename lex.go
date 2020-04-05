package messageformat

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"strconv"
	"unicode"
)

var ErrUnterminatedQuotedString = errors.New("unterminated quoted string")
var ErrLeadingZeroNumber = errors.New("number must not have leading zero")

type TokenType int

const (
	TokenTypeEOF TokenType = iota
	TokenTypeText
	TokenTypeWord
	TokenTypeNumber
	TokenTypeLBrace
	TokenTypeRBrace
	TokenTypeComma
	TokenTypeEqual
	TokenTypePound
	TokenTypeColon
)

type Token struct {
	Type  TokenType
	Value string
}

func (t Token) String() string {
	switch t.Type {
	case TokenTypeEOF:
		return "<EOF>"
	case TokenTypeText:
		return strconv.Quote(t.Value)
	case TokenTypeWord:
		return t.Value
	case TokenTypeNumber:
		return t.Value
	case TokenTypeLBrace:
		return "{"
	case TokenTypeRBrace:
		return "}"
	case TokenTypeComma:
		return ","
	case TokenTypeEqual:
		return "="
	case TokenTypePound:
		return "#"
	case TokenTypeColon:
		return ":"
	default:
		panic("unreachable")
	}
}

// TODO(lex): optional quoting
type lexer struct {
	input *bytes.Buffer
	// arg tells whether the next lex call is LexText or LexArg.
	arg             bool
	isInPluralStyle func() bool
	Output          []Token
}

func newLexer(s string) *lexer {
	return &lexer{
		input: bytes.NewBufferString(s),
	}
}

func (l *lexer) Lex() error {
	if l.arg {
		return l.LexArg()
	}
	return l.LexText(&bytes.Buffer{})
}

func (l *lexer) LexText(buf *bytes.Buffer) error {
	for {
		ch, err := l.input.ReadByte()
		if errors.Is(err, io.EOF) {
			l.outText(buf.String())
			l.out(TokenTypeEOF)
			l.arg = false
			return nil
		} else if err != nil {
			return err
		}
		switch ch {
		case '\'':
			return l.lexQuotedText(buf, &bytes.Buffer{})
		case '{':
			l.outText(buf.String())
			l.out(TokenTypeLBrace)
			l.arg = true
			return nil
		case '}':
			l.outText(buf.String())
			l.out(TokenTypeRBrace)
			l.arg = true
			return nil
		case '#':
			if l.isInPluralStyle != nil && l.isInPluralStyle() {
				l.outText(buf.String())
				l.out(TokenTypePound)
				l.arg = false
				return nil
			} else {
				buf.WriteByte(ch)
			}
		default:
			buf.WriteByte(ch)
		}
	}
}

func (l *lexer) lexQuotedText(textBuf *bytes.Buffer, quoteBuf *bytes.Buffer) error {
	for {
		ch, err := l.input.ReadByte()
		if errors.Is(err, io.EOF) {
			return ErrUnterminatedQuotedString
		} else if err != nil {
			return err
		}
		switch ch {
		case '\'':
			// ''
			if quoteBuf.Len() == 0 {
				textBuf.WriteByte('\'')
				return l.LexText(textBuf)
			}
			// look ahead one byte to see if it is '
			nextCh, err := l.input.ReadByte()
			if errors.Is(err, io.EOF) {
				textBuf.Write(quoteBuf.Bytes())
				l.outText(textBuf.String())
				l.out(TokenTypeEOF)
				return nil
			} else if err != nil {
				return err
			}
			switch nextCh {
			case '\'':
				textBuf.WriteByte('\'')
			default:
				l.input.UnreadByte()
				textBuf.Write(quoteBuf.Bytes())
				return l.LexText(textBuf)
			}
		default:
			quoteBuf.WriteByte(ch)
		}
	}
}

func (l *lexer) LexArg() error {
	for {
		ch, err := l.input.ReadByte()
		if errors.Is(err, io.EOF) {
			l.out(TokenTypeEOF)
			return nil
		} else if err != nil {
			return err
		}

		// Skip whitespace
		if unicode.IsSpace(rune(ch)) {
			continue
		}

		switch ch {
		case '{':
			l.out(TokenTypeLBrace)
			l.arg = false
			return nil
		case '}':
			l.out(TokenTypeRBrace)
			l.arg = false
			return nil
		case ',':
			l.out(TokenTypeComma)
			return nil
		case '=':
			l.out(TokenTypeEqual)
			return nil
		case ':':
			l.out(TokenTypeColon)
			return nil
		}

		if ch >= '0' && ch <= '9' {
			return l.lexNumber(bytes.NewBuffer([]byte{ch}))
		}

		if (ch >= 'a' && ch <= 'z') || (ch >= 'A' && ch <= 'Z') || ch == '_' {
			return l.lexWord(bytes.NewBuffer([]byte{ch}))
		}

		return fmt.Errorf("unexpected character: %v", strconv.QuoteRune(rune(ch)))
	}
}

func (l *lexer) lexNumber(buf *bytes.Buffer) error {
	for {
		ch, err := l.input.ReadByte()
		if errors.Is(err, io.EOF) {
			l.outNumber(buf.String())
			l.out(TokenTypeEOF)
			return nil
		} else if err != nil {
			return err
		}

		if ch >= '0' && ch <= '9' {
			if buf.String() == "0" {
				return ErrLeadingZeroNumber
			}
			buf.WriteByte(ch)
		} else {
			l.input.UnreadByte()
			l.outNumber(buf.String())
			return nil
		}
	}
}

func (l *lexer) lexWord(buf *bytes.Buffer) error {
	for {
		ch, err := l.input.ReadByte()
		if errors.Is(err, io.EOF) {
			l.outWord(buf.String())
			l.out(TokenTypeEOF)
			return nil
		} else if err != nil {
			return err
		}

		if (ch >= 'a' && ch <= 'z') || (ch >= 'A' && ch <= 'Z') || (ch >= '0' && ch <= '9') || ch == '_' {
			buf.WriteByte(ch)
		} else {
			l.input.UnreadByte()
			l.outWord(buf.String())
			return nil
		}
	}
}

func (l *lexer) outText(s string) {
	l.Output = append(l.Output, Token{Type: TokenTypeText, Value: s})
}

func (l *lexer) outNumber(s string) {
	l.Output = append(l.Output, Token{Type: TokenTypeNumber, Value: s})
}

func (l *lexer) outWord(s string) {
	l.Output = append(l.Output, Token{Type: TokenTypeWord, Value: s})
}

func (l *lexer) out(t TokenType) {
	l.Output = append(l.Output, Token{Type: t})
}
