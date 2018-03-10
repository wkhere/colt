package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"os"
	"os/exec"
)

func dieIf(err error) {
	if err != nil {
		fmt.Fprintln(os.Stderr, "colt:", err)
		os.Exit(1)
	}
}

type columnProc struct {
	separator      rune
	selection      int
	command        []string
	stdout, stderr io.Writer
}

func main() {
	p := columnProc{stdout: os.Stdout, stderr: os.Stderr}

	p.parseArgs(os.Args[1:])

	scanner := bufio.NewScanner(os.Stdin)

	for scanner.Scan() {
		p.process(scanner.Text())
		fmt.Fprintln(p.stdout)
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
		p.warn(err)
		fmt.Fprint(p.stdout, line)
		return
	}

	for i, col := range cols {
		for _, token := range col {
			if token.typ == tokenData && i == selectedIdx {
				p.processData(token.val)
				continue
			}
			fmt.Fprint(p.stdout, token.val)
		}
	}
}

func (p *columnProc) processData(s string) {
	var b bytes.Buffer
	cmd := exec.Command(p.command[0], append(p.command[1:], s)...)
	cmd.Env = append(os.Environ(), "COLOR=1")
	cmd.Stdout = &b
	cmd.Stderr = p.stderr
	err := cmd.Run()
	dieIf(err)
	p.stdout.Write(chomp(b.Bytes()))
}

func (p *columnProc) warn(err error) {
	fmt.Fprintf(p.stderr, "WARN %v\n", err)
}

func chomp(b []byte) []byte {
outer:
	for len(b) > 0 {
		switch l := len(b) - 1; b[l] {
		case '\n', '\r':
			b = b[:l]
		default:
			break outer
		}
	}
	return b
}
