package colt

import (
	"bytes"
	"io"
	"os"
	"os/exec"

	"github.com/wkhere/colt/lex"
	"github.com/wkhere/colt/parse"
)

type Colt struct {
	Separator, Quote rune
	Selection        int
	Unquote          bool
	Command          []string
	Stdout, Stderr   io.Writer
}

func (c *Colt) ProcessLine(line []byte) error {

	cols := parse.GroupTokens(lex.LexTokens(line, c.Separator, c.Quote))

	selectedIdx, err := setupIdx(c.Selection, len(cols))
	if err != nil {
		c.Stdout.Write(line)
		return Warning{err}
	}

	for i, col := range cols {
		col = parse.NormalizeColumn(col)
		for _, token := range col {
			if token.Type == lex.TokenData && i == selectedIdx {
				err := c.ProcessData(token.Val)
				if err != nil {
					return err
				}
				continue
			}
			c.Stdout.Write(token.Val)
		}
	}
	return nil
}

func (c *Colt) ProcessData(d []byte) error {
	if c.Unquote {
		d = unquote(d, c.Quote)
	}
	var b bytes.Buffer
	cmd := exec.Command(c.Command[0], append(c.Command[1:], string(d))...)
	cmd.Env = append(os.Environ(), "COLOR=1")
	cmd.Stdout = &b
	cmd.Stderr = c.Stderr
	err := cmd.Run()
	if err != nil {
		return err
	}
	c.Stdout.Write(chomp(b.Bytes()))
	return nil
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
