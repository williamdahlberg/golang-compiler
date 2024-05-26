package parser

import (
	"fmt"
	"golang/compiler/emitter"
	"golang/compiler/lexer"
	"slices"
)

func NewParser(lexer lexer.Lexer, emitter *emitter.Emitter) Parser {
	p := Parser{}
	p.lexer = lexer
	p.emitter = emitter
	p.nextToken()
	p.nextToken()
	return p
}

type Parser struct {
	lexer          lexer.Lexer
	emitter        *emitter.Emitter
	curToken       lexer.Token
	peekToken      lexer.Token
	symbols        []string
	labelsDeclared []string
	labelsGotoed   []string
}

func (p *Parser) nextToken() {
	p.curToken = p.peekToken
	p.peekToken = p.lexer.GetToken()
}

func (p *Parser) match(kind string) {
	if p.curToken.Kind != kind {
		p.abort(fmt.Sprintf("Expected %s, got %s", kind, p.curToken.Kind))
	}
	p.nextToken()
}

func (p *Parser) abort(message string) {
	panic(message)
}

func (p *Parser) Program() {
	fmt.Println("PROGRAM")

	for p.curToken.Kind == "NEWLINE" {
		p.nextToken()
	}

	for p.curToken.Kind != "EOF" {
		p.statement()
	}

	for _, label := range p.labelsGotoed {
		if !slices.Contains(p.labelsDeclared, label) {
			p.abort("Tried to GOTO non-existing label: " + label)
		}
	}
}

func (p *Parser) statement() {
	switch p.curToken.Kind {
	case "PRINT":
		p.nextToken()
		if p.curToken.Kind == "STRING" {
			p.emitter.EmitLine(fmt.Sprintf("print(\"%s\")", p.curToken.Text))
			p.nextToken()
		} else {
			p.emitter.Emit("print(")
			p.expression()
			p.emitter.EmitLine(")")
		}

	case "IF":
		p.nextToken()
		p.emitter.Emit("if ")
		p.comparison()

		p.match("THEN")
		p.nl()
		p.emitter.EmitLine(":")

		p.emitter.Indentation += 2

		for p.curToken.Kind != "ENDIF" {
			p.statement()
		}

		p.match("ENDIF")
		p.emitter.Indentation -= 2

	case "WHILE":
		p.nextToken()
		p.emitter.Emit("while ")
		p.comparison()

		p.match("REPEAT")
		p.nl()
		p.emitter.EmitLine(":")
		p.emitter.Indentation += 2

		for p.curToken.Kind != "ENDWHILE" {
			p.statement()
		}

		p.match("ENDWHILE")
		p.emitter.Indentation -= 2

	// case "LABEL":
	// 	fmt.Println("STATEMENT-LABEL")
	// 	p.nextToken()
	// 	if slices.Contains(p.labelsDeclared, p.curToken.Text) {
	// 		p.abort("Label exists: " + p.curToken.Text)
	// 	}
	// 	p.labelsDeclared = append(p.labelsDeclared, p.curToken.Text)

	// 	p.match("IDENT")

	// case "GOTO":
	// 	fmt.Println("STATEMENT-GOTO")
	// 	p.nextToken()
	// 	p.labelsGotoed = append(p.labelsGotoed, p.curToken.Text)
	// 	p.match("IDENT")

	case "LET":
		p.nextToken()
		if !slices.Contains(p.symbols, p.curToken.Text) {
			p.symbols = append(p.symbols, p.curToken.Text)
		}
		p.emitter.Emit(p.curToken.Text + " = ")
		p.match("IDENT")
		p.match("EQ")
		p.expression()
		p.emitter.EmitLine("")

	case "INPUT":
		p.nextToken()
		if !slices.Contains(p.symbols, p.curToken.Text) {
			p.symbols = append(p.symbols, p.curToken.Text)
		}
		p.emitter.EmitLine(p.curToken.Text + " = int(input())")
		p.match("IDENT")

	default:
		p.abort(fmt.Sprintf("Invalid statemnt at %s (%s)", p.curToken.Text, p.curToken.Kind))
	}

	p.nl()

}

func (p *Parser) nl() {
	if p.curToken.Kind == "EOF" {
		return
	}
	p.match("NEWLINE")
	for p.curToken.Kind == "NEWLINE" {
		p.nextToken()
	}
}

func (p *Parser) expression() {
	p.term()
	for p.curToken.Kind == "PLUS" || p.curToken.Kind == "MINUS" {
		p.emitter.Emit(p.curToken.Text)
		p.nextToken()
		p.term()
	}
}

func (p *Parser) term() {
	p.unary()
	for p.curToken.Kind == "ASTERISK" || p.curToken.Kind == "SLASH" {
		p.emitter.Emit(p.curToken.Text)
		p.nextToken()
		p.unary()
	}
}

func (p *Parser) unary() {
	if p.curToken.Kind == "PLUS" || p.curToken.Kind == "MINUS" {
		p.emitter.Emit(p.curToken.Text)
		p.nextToken()
	}
	p.primary()
}

func (p *Parser) primary() {
	switch p.curToken.Kind {
	case "NUMBER":
		p.emitter.Emit(p.curToken.Text)
		p.nextToken()
	case "IDENT":
		if !slices.Contains(p.symbols, p.curToken.Text) {
			p.abort("Reference variable before assigning: " + p.curToken.Text)
		}
		p.emitter.Emit(p.curToken.Text)
		p.nextToken()
	default:
		p.abort("Unexpected token at " + p.curToken.Text)
	}
}

func (p *Parser) comparison() {
	p.expression()
	if p.isComparisonOperator() {
		p.emitter.Emit(p.curToken.Text)
		p.nextToken()
		p.expression()
	} else {
		p.abort(fmt.Sprintf("Expected comparison operator at: %s ", p.curToken.Text))
	}
	for p.isComparisonOperator() {
		p.emitter.Emit(p.curToken.Text)
		p.nextToken()
		p.expression()
	}
}

func (p *Parser) isComparisonOperator() bool {
	return (p.curToken.Kind == "GT" ||
		p.curToken.Kind == "GTEQ" ||
		p.curToken.Kind == "LT" ||
		p.curToken.Kind == "LTEQ" ||
		p.curToken.Kind == "EQEQ" ||
		p.curToken.Kind == "NOTEQ")
}
