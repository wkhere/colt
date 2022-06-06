package colt

import (
	"github.com/wkhere/colt/lex"
	"github.com/wkhere/colt/parse"
)

type Colt struct {
	Separator, Quote rune
	Selection        int
	Unquote          bool
	T                Transformer
}

func (c *Colt) ProcessLine(line []byte) error {

	tt, err := lex.LexTokens(line, c.Separator, c.Quote).Gather()
	if err != nil {
		return Warning{err}
	}

	cols := parse.GroupTokens(tt)

	selectedIdx, err := setupIdx(c.Selection, len(cols))
	if err != nil {
		c.T.Copy(line)
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
			c.T.Copy(token.Val)
		}
	}
	return nil
}

func (c *Colt) ProcessData(d []byte) error {
	if c.Unquote {
		d = unquote(d, c.Quote)
	}
	return c.T.Transform(d)
}

func setupIdx(col, ncols int) (int, error) {
	switch {
	case col < 0 && col >= -ncols:
		return ncols + col, nil
	case col > 0 && col <= ncols:
		return col - 1, nil
	default:
		return -1, errorf(
			"invalid column number %d, have %d column(s)",
			col, ncols)
	}
}
