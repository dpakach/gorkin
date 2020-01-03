package lexer

import (
	"github.com/dpakach/gorkin/token"
	"io/ioutil"
	"strings"
)

// Lexer is the Lexer object for reading through the Gherkin input
type Lexer struct {
	input         string
	position      int
	readPosition  int
	ch            byte
	currentLineNo int
	FilePath      string
}

// New Creates a new Lexer object for given input
func New(input string) *Lexer {
	l := &Lexer{input: input, currentLineNo: 1}
	l.readChar()
	return l
}

// NewFromFile Creates a new Lexer object for given feature file
func NewFromFile(path string) *Lexer {
	dat, err := ioutil.ReadFile(path)
	if err != nil {
		panic(err)
	}

	l := &Lexer{input: string(dat), FilePath: path}
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
	l.readPosition++
}

func newToken(tokenType token.Type, ch byte) token.Token {
	return token.Token{Type: tokenType, Literal: string(ch)}
}

func isLetter(ch byte) bool {
	return 'a' <= ch && ch <= 'z' || 'A' <= ch && ch <= 'Z' || ch == '_' || ch == '-' || ch == '.' || ch == '+'
}

func isValidBodyChar(ch byte) bool {
	return isLetter(ch) || ch == ' ' || ch == '_' || ch == '-'
}

func isDigit(ch byte) bool {
	return '0' <= ch && ch <= '9'
}

func isWhitespace(ch byte) bool {
	if ch == ' ' || ch == '\t' || ch == '\r' {
		return true
	}
	return false
}

func (l *Lexer) skipWhitespace() {
	for isWhitespace(l.ch) {
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
	}
	return l.input[l.readPosition]
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

func (l *Lexer) readPyString() string {
	position := l.position + 1
	for {
		l.readChar()
		if l.ch == '\n' {
			l.currentLineNo++
		}
		if (l.ch == '"' && l.peekChar() == '"') || l.ch == 0 {
			if l.peekChar() == '"' {
				break
			}
		}
	}

	res := l.input[position:l.position]
	return res
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
		if isLetter(l.ch) || isDigit(l.ch) {
			continue
		}
		break
	}
	return l.input[position:l.position]
}

func (l *Lexer) readTillLineBreak() string {
	position := l.position
	for l.ch != '\n' {
		if l.ch == 0 {
			break
		}
		l.readChar()
	}
	return l.input[position:l.position]
}

func (l *Lexer) readBody() string {
	position := l.position

	for isValidBodyChar(l.ch) {
		l.readChar()
	}
	return l.input[position:l.position]
}

func (l *Lexer) readTableData() string {
	position := l.position
	for {
		l.readChar()
		if l.peekChar() == '|' || l.ch == 0 {
			break
		}
	}

	return l.input[position:l.position]
}

// NextToken returns the next token in the lexer
func (l *Lexer) NextToken() token.Token {
	var tok token.Token
	l.skipWhitespace()
	switch l.ch {
	case 0:
		tok.Literal = token.EOF.String()
		tok.Type = token.EOF
		l.readChar()
	case '#':
		l.readChar()
		l.skipWhitespace()
		body := l.readTillLineBreak()
		tok.Type = token.COMMENT
		tok.Literal = strings.TrimSpace(body)
	case ':':
		tok = newToken(token.COLON, l.ch)
		l.readChar()
	case '"':
		if l.peekChar() != '"' {
			tok.Type = token.STRING
			tok.Literal = l.readString()
			l.readChar()
		} else {
			l.readChar()
			if l.peekChar() != '"' {
				tok.Type = token.STRING
				tok.Literal = ""
				l.readChar()
			} else {
				tok.LineNumber = l.currentLineNo
				l.readChar()
				tok.Type = token.PYSTRING
				tok.Literal = strings.TrimSpace(l.readPyString())
				l.readChar()
				l.readChar()
				l.readChar()
			}
		}
	case '@':
		l.readChar()
		tok.Type = token.TAG
		word := l.readWord()
		tok.Literal = word
	case '\n':
		l.readChar()
		tok.Type = token.NEWLINE
		tok.Literal = token.NEWLINE.String()
	case '<':
		word := l.readExampleValue()
		tok.Type = token.TABLEDATA
		tok.Literal = word
		l.readChar()
	case '|':
		l.readChar()
		l.skipWhitespace()
		if l.ch == '\n' {
			tok.Type = token.NEWLINE
			tok.Literal = token.NEWLINE.String()
			l.readChar()
		} else if l.ch == 0 {
			tok.Type = token.EOF
			tok.Literal = token.EOF.String()
			l.readChar()
		} else if l.ch == '|' {
			tok.Type = token.TABLEDATA
			tok.Literal = ""
		} else {
			tok.Type = token.TABLEDATA
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
				tok.Type = token.STEPBODY
			}
		}
	}
	if tok.LineNumber == 0 {
		tok.LineNumber = l.currentLineNo
	}
	if tok.Type == token.NEWLINE {
		l.currentLineNo++
	}
	return tok
}
