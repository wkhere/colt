package parse

import (
	"bytes"

	"github.com/wkhere/colt/lex"
)

// GroupTokens returns slice of groups of tokens, each group
// being a tokens slice, divided by separator (lex.TokenSep).
// Separator, if exists, is a last token in such group.
func GroupTokens(tt []lex.Token) (res [][]lex.Token) {
	res = make([][]lex.Token, 1, 3)
	groupIdx := 0
	for _, tok := range tt {
		res[groupIdx] = append(res[groupIdx], tok)
		if tok.Type == lex.TokenSep {
			res = append(res, nil)
			groupIdx++
		}
	}
	if last := len(res) - 1; res[last] == nil {
		res = res[:last]
	}
	return
}

func NormalizeColumn(col []lex.Token) (res []lex.Token) {
	res = make([]lex.Token, 0, len(col))
	var i, j int

	for i = 0; i < len(col) && col[i].Type != lex.TokenData; i++ {
		res = append(res, col[i])
	}
	if i == len(col) {
		return res
		// todo: return error or bool indicating no data?
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
