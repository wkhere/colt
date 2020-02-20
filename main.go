package main

import (
	"bufio"
	"bytes"
	"io"
	"os"
	"os/exec"
)

type columnProc struct {
	separator, quote rune
	selection        int
	unquote          bool
	command          []string
	stdout, stderr   io.Writer
}

func main() {
	p := columnProc{stdout: os.Stdout, stderr: os.Stderr}

	err := p.parseArgs(os.Args[1:])
	if err != nil {
		die2(err)
	}

	scanner := bufio.NewScanner(os.Stdin)

	for scanner.Scan() {
		p.process(scanner.Bytes())
		io.WriteString(p.stdout, "\n")
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
		return -1, errorf(
			"invalid column selector #%d for %d columns",
			col, ncols)
	}
}

func (p *columnProc) process(line []byte) {

	cols := lexTokens(line, p.separator, p.quote).group()

	selectedIdx, err := setupIdx(p.selection, len(cols))
	if err != nil {
		fwarn(p.stderr, err)
		p.stdout.Write(line)
		return
	}

	for i, col := range cols {
		col = normalizeColumn(col)
		for _, token := range col {
			if token.typ == tokenData && i == selectedIdx {
				p.processData(token.val)
				continue
			}
			p.stdout.Write(token.val)
		}
	}
}

func (p *columnProc) processData(d []byte) {
	if p.unquote {
		d = unquote(d, p.quote)
	}
	var b bytes.Buffer
	cmd := exec.Command(p.command[0], append(p.command[1:], string(d))...)
	cmd.Env = append(os.Environ(), "COLOR=1")
	cmd.Stdout = &b
	cmd.Stderr = p.stderr
	err := cmd.Run()
	dieIf(err)
	p.stdout.Write(chomp(b.Bytes()))
}
