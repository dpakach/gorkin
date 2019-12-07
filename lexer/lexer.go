package lexer

import (
	"strings"
	"github.com/dpakach/gorkin/token"
)

type Lexer struct {
	input string
	position int
	readPosition int
	ch byte
}

func New(input string) * Lexer {
	l := &Lexer{input: input}
	l.readChar()
	return l
}

func (l *Lexer) readChar() {
	if l.readPosition >= len(l.input) {
		l.ch = 0
	} else {
		l.ch = l.input[l.readPosition]
	}
	l.position = l.readPosition
	l.readPosition += 1
}

func newToken(tokenType token.TokenType, ch byte) token.Token {
	return token.Token{Type: tokenType, Literal: string(ch)}
}

func isLetter(ch byte) bool {
	return 'a' <= ch && ch <= 'z' || 'A' <= ch && ch <= 'Z' || ch == '_'
}

func isValidBodyChar(ch byte) bool {
	return isLetter(ch) || ch == ' ' || ch == '_' || ch == '-'
}

func isDigit(ch byte) bool {
	return '0' <= ch && ch <= '9'
}

func isWhitespace(ch byte) bool {
	if ch == ' ' || ch == '\t' || ch == '\n' || ch == '\r' {
		return true
	}
	return false
}

func (l *Lexer) skipWhitespace() {
	for isWhitespace(l.ch) {
		l.readChar()
	}
}

func (l *Lexer) skipWhitespacesTillLineBreak() {
	for l.ch == ' ' || l.ch == '\t' || l.ch == '\r' {
		l.readChar()
	}
}

func (l *Lexer) readNumber() string {
	position := l.position
	for isDigit(l.ch) {
		l.readChar()
	}
	return l.input[position:l.position]
}

func (l *Lexer) peekChar() byte {
	if l.readPosition > len(l.input) {
		return 0
	} else {
		return l.input[l.readPosition]
	}
}

func (l *Lexer) readString() string {
	position := l.position + 1
	for {
		l.readChar()
		if l.ch == '"' || l.ch == 0 {
			break
		}
	}

	return l.input[position:l.position]
}

func (l *Lexer) readExampleValue() string {
	position := l.position + 1
	for {
		l.readChar()
		if l.ch == '>' || l.ch == 0 {
			break
		}
	}

	return l.input[position:l.position]
}

func (l *Lexer) readWord() string {
	position := l.position
	for {
		l.readChar()
		if isLetter(l.ch) || isDigit(l.ch) || l.ch == 0 {
			continue
		}
		break
	}
	return l.input[position:l.position]
}

func (l *Lexer) readBody() string {
	position := l.position

	for isLetter(l.ch) || l.ch == ' '{
		l.readChar()
	}
	return l.input[position:l.position]
}

func (l *Lexer) readTableData() string {
	position := l.position + 1
	for {
		l.readChar()
		if l.peekChar() == '|' || l.ch == 0 {
			break
		}
	}

	return l.input[position:l.position]
}

func (l *Lexer) NextToken() token.Token {
	var tok token.Token
	l.skipWhitespace()
	switch l.ch {
		case':':
			tok = newToken(token.COLON, l.ch)
			l.readChar()
		case '"':
			tok.Type = token.STRING
			tok.Literal = l.readString()
			l.readChar()
		case '@':
			l.readChar()
			tok.Type = token.TAG
			word := l.readWord()
			tok.Literal = word
		case '<':
			word := l.readExampleValue()
			tok.Type = token.TABLE_DATA
			tok.Literal = word
			l.readChar()
		case '|':
			l.skipWhitespacesTillLineBreak()
			if l.peekChar() == '\n' {
				tok.Type = token.TABLE_LINE_BREAK
				tok.Literal = token.TABLE_LINE_BREAK
				l.readChar()
			} else {
				tok.Type = token.TABLE_DATA
				tok.Literal = strings.TrimSpace(l.readTableData())
				l.readChar()
			}
		default:
			if isDigit(l.ch) {
				tok.Literal = l.readNumber()
				tok.Type = token.NUMBER
			} else {
				word := l.readWord()
				if key, ok := token.GherkinKeyword[word]; ok {
					tok.Literal = word
					tok.Type = key
				} else {
					body := l.readBody()
					tok.Literal = strings.TrimSpace(word + body)
					tok.Type = token.STEP_BODY
				}
			}
	}
	return tok
}
