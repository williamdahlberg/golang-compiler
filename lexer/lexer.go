package lexer

import (
	"fmt"
	"slices"
	"unicode"
)

func Lex(s string) (tokens []Token) {
	l := Lexer{s, 1, -1}
	l.nextChar()
	for l.curChar != 0 {
		tokens = append(tokens, l.GetToken())
	}
	return
}

func NewLexer(s string) Lexer {
	l := Lexer{s, 1, -1}
	l.nextChar()
	return l
}

// LEXER
type Lexer struct {
	source  string
	curChar byte
	curPos  int
}

func (l *Lexer) nextChar() {
	l.curPos += 1
	if l.curPos >= len(l.source) {
		l.curChar = 0
	} else {
		l.curChar = l.source[l.curPos]
	}
}

func (l *Lexer) peek() byte {
	if l.curPos+1 >= len(l.source) {
		return 0
	}
	return l.source[l.curPos+1]
}

func (l *Lexer) GetToken() (token Token) {
	l.skipWhitespace()
	l.skipComment()
	switch {
	case l.curChar == 0:
		token = Token{"EOF", string(l.curChar)}
	case l.curChar == '\n':
		token = Token{"NEWLINE", string(l.curChar)}
	case l.curChar == '+':
		token = Token{"PLUS", string(l.curChar)}
	case l.curChar == '-':
		token = Token{"MINUS", string(l.curChar)}
	case l.curChar == '*':
		token = Token{"ASTERISK", string(l.curChar)}
	case l.curChar == '/':
		token = Token{"SLASH", string(l.curChar)}

	case l.curChar == '=':
		switch l.peek() {
		case '=':
			token = Token{"EQEQ", str(l.curChar, l.peek())}
			l.nextChar()
		default:
			token = Token{"EQ", str(l.curChar)}
		}

	case l.curChar == '>':
		switch l.peek() {
		case '=':
			token = Token{"GTEQ", str(l.curChar, l.peek())}
			l.nextChar()
		default:
			token = Token{"GT", str(l.curChar)}
		}

	case l.curChar == '<':
		switch l.peek() {
		case '=':
			token = Token{"LTEQ", str(l.curChar, l.peek())}
			l.nextChar()
		default:
			token = Token{"LT", str(l.curChar)}
		}

	case l.curChar == '!':
		switch l.peek() {
		case '=':
			token = Token{"NOTEQ", str(l.curChar, l.peek())}
			l.nextChar()
		default:
			panic(fmt.Sprintf("Expected !=, got !%s", str(l.peek())))
		}

	case l.curChar == '"':
		l.nextChar()
		startPos := l.curPos

		disallowedChars := []byte{'\r', '\n', '\t', '\\', '%'}

		for ; l.curChar != '"'; l.nextChar() {
			if slices.Contains(disallowedChars, l.curChar) {
				panic(fmt.Sprintf("Illegal character in string: %s", str(l.curChar)))
			}
		}
		tokenText := l.source[startPos:l.curPos]
		token = Token{"STRING", tokenText}

	case unicode.IsDigit(rune(l.curChar)):
		startPos := l.curPos
		for ; unicode.IsDigit(rune(l.peek())); l.nextChar() {
		}
		if l.peek() == '.' {
			l.nextChar()
			if !unicode.IsDigit(rune(l.peek())) {
				panic("Illegal character in number")
			}
			for ; unicode.IsDigit(rune(l.peek())); l.nextChar() {
			}
		}
		tokenText := l.source[startPos : l.curPos+1]
		token = Token{"NUMBER", tokenText}

	case unicode.IsLetter(rune(l.curChar)):
		startPos := l.curPos
		for ; unicode.IsLetter(rune(l.peek())) || unicode.IsDigit(rune(l.peek())); l.nextChar() {
		}
		tokenText := l.source[startPos : l.curPos+1]
		tokenType := KindType(tokenText)
		if tokenType == "KEYWORD" {
			token = Token{tokenText, tokenText}
		} else {
			token = Token{"IDENT", tokenText}
		}

	default:
		panic(fmt.Sprintf("Lexing error on symbol %s", string(l.curChar)))
	}
	l.nextChar()
	return
}

func (l *Lexer) skipWhitespace() {
	whitespaces := []byte{' ', '\t', '\r'}
	for ; slices.Contains(whitespaces, l.curChar); l.nextChar() {
	}
}

func (l *Lexer) skipComment() {
	if l.curChar == '#' {
		for ; l.curChar != '\n'; l.nextChar() {
		}
	}
}

// TOKEN
type Token struct {
	Kind string
	Text string
}

func (t Token) String() string {
	return fmt.Sprintf("Kind: %s, txt: '%s'", t.Kind, t.Text)
}

// HLEP
func KindType(Kind string) string {
	basics := []string{
		"NEWLINE", "NUMBER",
		"IDENT", "STRING"}

	keywords := []string{
		"LABEL", "GOTO",
		"PRINT", "INPUT",
		"LET", "IF",
		"THEN", "ENDIF",
		"WHILE", "REPEAT",
		"ENDWHILE",
	}

	operators := []string{
		"EQ", "PLUS",
		"MINUS", "ASTERISK",
		"SLASH", "EQEQ",
		"NOTEQ", "LT",
		"LTEQ", "GT",
		"GTEQ",
	}

	if Kind == "EOF" {
		return "EOF"
	}

	if slices.Contains(basics, Kind) {
		return "BASIC"
	}

	if slices.Contains(keywords, Kind) {
		return "KEYWORD"
	}

	if slices.Contains(operators, Kind) {
		return "OPERATOR"
	}

	return "OTHER"
}

func str(bytes ...byte) (s string) {
	for _, b := range bytes {
		s += string(b)
	}
	return
}
