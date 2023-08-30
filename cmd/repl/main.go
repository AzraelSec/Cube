package main

import (
	"bufio"
	"fmt"
	"io"
	"os"

	"github.com/AzraelSec/cube/pkg/evaluator"
	"github.com/AzraelSec/cube/pkg/lexer"
	"github.com/AzraelSec/cube/pkg/object"
	"github.com/AzraelSec/cube/pkg/parser"
)

const prompt = ">>"

func main() {
	start(os.Stdin, os.Stdout)
}

func start(in io.Reader, out io.Writer) {
	scanner := bufio.NewScanner(os.Stdin)
	env := object.NewEnvironment()

	for {
		fmt.Print(prompt)
		sn := scanner.Scan()
		if !sn {
			return
		}

		line := scanner.Text()
		l := lexer.New(line)
		p := parser.New(l)

		prog := p.ParseProgram()
		if len(p.Errors()) != 0 {
			printParserErrors(out, p.Errors())
			continue
		}

		evaluated := evaluator.Eval(prog, env)
		if evaluated != nil {
			io.WriteString(out, evaluated.Inspect())
			io.WriteString(out, "\n")
		}
	}
}

func printParserErrors(out io.Writer, errors []string) {
	for _, msg := range errors {
		io.WriteString(out, "\t"+msg+"\n")
	}
}
