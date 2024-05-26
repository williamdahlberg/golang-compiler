package emitter

import (
	"fmt"
	"os"
	"strings"
)

func NewEmitter(path string) *Emitter {
	e := Emitter{}
	e.FullPath = path
	return &e
}

type Emitter struct {
	FullPath    string
	Header      string
	Code        string
	Indentation int
}

func (e *Emitter) Emit(code string) {
	e.Code += strings.Repeat(" ", e.Indentation) + code
}

func (e *Emitter) EmitLine(code string) {
	e.Code += strings.Repeat(" ", e.Indentation) + code + "\n"
}

func (e *Emitter) HeaderLine(code string) {
	e.Header += code + "\n"
}

func (e *Emitter) WriteFile() {
	fmt.Println("--- WRITING FILE")
	fmt.Println(e.Code)
	err := os.WriteFile(e.FullPath, []byte(e.Code), 0666)
	if err != nil {
		panic("could not write file")
	}
}
