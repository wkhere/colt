package colt

import (
	"fmt"

	"github.com/wkhere/colt/lex"
	"github.com/wkhere/colt/parse"
)

type Colt struct {
	Separator, Quote rune
	Selection        int
	Unquote          bool
	T                Transformer
}

type Warning struct {
	error
}

func (c *Colt) ProcessLine(line []byte) error {

	lexConf := lex.DefaultConfig
	lexConf.Sep = c.Separator
	lexConf.Quote = c.Quote

	// todo: return tokens as Pos, no slice allocation
	// then pass input to parse and extract chunks
	ts := lex.LexTokens(line, lexConf)
	cols, err := parse.GroupTokens(ts)
	if err != nil {
		return Warning{err}
	}

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
		return -1, fmt.Errorf(
			"invalid column number %d, have %d column(s)",
			col, ncols)
	}
}
