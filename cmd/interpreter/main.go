package main

import (
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/AzraelSec/cube/pkg/evaluator"
	"github.com/AzraelSec/cube/pkg/lexer"
	"github.com/AzraelSec/cube/pkg/object"
	"github.com/AzraelSec/cube/pkg/parser"
)

func main() {
	if len(os.Args) < 2 {
		help(os.Args[0])
		return
	}

	if !strings.HasSuffix(os.Args[1], ".cb") {
		fmt.Printf("wrong file suffix in %s", os.Args[1])
		return
	}

	file, err := os.Open(os.Args[1])
	if err != nil {
		fmt.Printf("impossible to open the file %s: %v", os.Args[1], err)
		return
	}

	env := object.NewEnvironment()
	content, err := io.ReadAll(file)
	if err != nil {
		fmt.Println("impossible to read file content")
		return
	}

	scontent := string(content)

	l := lexer.New(scontent)
	p := parser.New(l)

	prog := p.ParseProgram()
	if len(p.Errors()) != 0 {
		fmt.Println("Errors found:")
		printParserErrors(os.Stderr, p.Errors())
		return
	}

	if evaluated := evaluator.Eval(prog, env); evaluated != nil && evaluated.Type() == object.ERROR_OBJ {
		fmt.Println(evaluated.Inspect())
		return
	}
}

func help(exec string) {
	fmt.Printf("usage: %s [file.cb]", exec)
}

func printParserErrors(out io.Writer, errors []string) {
	for idx, msg := range errors {
		io.WriteString(out, fmt.Sprintf("\t%d: %s\n", idx, msg))
	}
}
