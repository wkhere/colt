package colt

import (
	"bytes"

	"github.com/wkhere/colt/lex"
)

func normalizeColumn(col []lex.Token) (res []lex.Token) {
	res = make([]lex.Token, 0, len(col))
	var i, j int

	for i = 0; i < len(col) && col[i].Type != lex.TokenData; i++ {
		res = append(res, col[i])
	}
	if i == len(col) {
		return res
	}
	for j = len(col) - 1; j > i && col[j].Type != lex.TokenData; j-- {
	}
	// now i is at the first data token and j at the last

	if i == j {
		res = append(res, col[i])
	} else {
		var b bytes.Buffer
		for _, token := range col[i : j+1] {
			b.Write(token.Val)
		}
		res = append(res, lex.Token{Type: lex.TokenData, Val: b.Bytes()})
	}

	for j++; j < len(col); j++ {
		res = append(res, col[j])
	}

	return res
}
