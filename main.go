package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"os/exec"
)

func dieIf(err error) {
	if err != nil {
		fmt.Fprintln(os.Stderr, "cols:", err)
		os.Exit(1)
	}
}

type columnProc struct {
	sep     rune
	field   int
	command []string
	output  *os.File
}

func main() {
	p := columnProc{';', 1, nil, os.Stdout}

	{
		fieldFlag := flag.Int("f", 1, "field to extract")

		flag.Parse()

		p.field = *fieldFlag
		p.command = flag.Args()
	}

	scanner := bufio.NewScanner(os.Stdin)

	for scanner.Scan() {
		p.processFields(scanner.Text())
		fmt.Fprintln(p.output)
	}

	dieIf(scanner.Err())
}

func (p *columnProc) processFields(line string) {
	i := 1
	for token := range lexTokens(line, p.sep) {
		switch token.typ {
		case tokenData:
			if i == p.field {
				p.processData(token.val)
				continue
			}
		case tokenSep:
			i++
		case tokenSpace:
		}
		p.output.WriteString(token.val)
	}
}

func (p *columnProc) processData(s string) {
	cmd := exec.Command(p.command[0], append(p.command[1:], s)...)
	cmd.Stdout = p.output
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	dieIf(err)
}
