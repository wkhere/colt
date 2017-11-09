package main

import (
	"bufio"
	"fmt"
	"io"
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
	separator rune
	selection int
	command   []string
	output    io.Writer
}

func main() {
	p := columnProc{output: os.Stdout}

	p.parseArgs(os.Args[1:])

	scanner := bufio.NewScanner(os.Stdin)

	for scanner.Scan() {
		p.process(scanner.Text())
		fmt.Fprintln(p.output)
	}

	dieIf(scanner.Err())
}

func setupIdx(col, ncols int) (int, error) {
	switch {
	case col < 0 && col >= -ncols:
		return ncols + col, nil
	case col > 0 && col <= ncols:
		return col - 1, nil
	default:
		return -1, fmt.Errorf(
			"invalid column selector #%d for %d columns",
			col, ncols)
	}
}

func (p *columnProc) process(line string) {

	cols := lexTokens(line, p.separator).group()

	selectedIdx, err := setupIdx(p.selection, len(cols))
	if err != nil {
		warn(err)
		fmt.Fprint(p.output, line)
		return
	}

	for i, col := range cols {
		for _, token := range col {
			if token.typ == tokenData && i == selectedIdx {
				p.processData(token.val)
				continue
			}
			fmt.Fprint(p.output, token.val)
		}
	}
}

func (p *columnProc) processData(s string) {
	cmd := exec.Command(p.command[0], append(p.command[1:], s)...)
	cmd.Stdout = p.output
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	dieIf(err)
}

func warn(err error) {
	fmt.Fprintf(os.Stderr, "WARN %v\n", err)
}
