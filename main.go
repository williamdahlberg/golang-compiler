package main

import (
	emittermodule "golang/compiler/emitter"
	lexmodule "golang/compiler/lexer"
	parsemodule "golang/compiler/parser"
	"os"
)

func main() {
	lexer := lexmodule.NewLexer(readCode())
	emitter := emittermodule.NewEmitter("out.py")
	parser := parsemodule.NewParser(lexer, emitter)
	parser.Program()
	emitter.WriteFile()

}

func readCode() string {
	// read the whole file at once
	b, err := os.ReadFile("code.gl")
	if err != nil {
		panic(err)
	}

	var s = string(b)
	return s
}
