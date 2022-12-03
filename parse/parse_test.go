package parse

import (
	"reflect"
	"testing"

	. "github.com/wkhere/colt/lex"
)

// shorter syntax

type b = []byte
type tg = [][]Token
type ts = []Token

func t(tt TokenType, data string) Token {
	return Token{Type: tt, Val: b(data)}
}

var (
	sep  = t(TokenSep, ";")
	d1   = t(TokenData, "a")
	d2   = t(TokenData, "bb")
	d11  = t(TokenData, "aa")
	d12  = t(TokenData, "abb")
	d111 = t(TokenData, "aaa")
	d121 = t(TokenData, "abba")
	d1s1 = t(TokenData, "a a")
	d1s2 = t(TokenData, "a bb")
	spc  = t(TokenSpace, " ")
)

// tabs

var grouptab = []struct {
	tokens ts
	group  tg
}{
	{ts{}, tg{}},
	{ts{d1}, tg{ts{d1}}},
	{ts{d1, d2}, tg{ts{d1, d2}}},
	{ts{sep}, tg{ts{sep}}}, // strange but it is how it is
	{ts{d1, sep}, tg{ts{d1, sep}}},
	{ts{d1, sep, d2}, tg{ts{d1, sep}, ts{d2}}},
	{ts{d1, d2, sep, d1}, tg{ts{d1, d2, sep}, ts{d1}}},
	{ts{d1, sep, d1, d2}, tg{ts{d1, sep}, ts{d1, d2}}},
	{ts{d1, d2, sep, d1, d2}, tg{ts{d1, d2, sep}, ts{d1, d2}}},
	{ts{d1, sep, d2, sep}, tg{ts{d1, sep}, ts{d2, sep}}},
	{ts{d1, sep, d2, sep, d1}, tg{ts{d1, sep}, ts{d2, sep}, ts{d1}}},
	{ts{d1, sep, sep}, tg{ts{d1, sep}, ts{sep}}}, // also strange
	{ts{spc}, tg{ts{spc}}},
	{ts{spc, sep}, tg{ts{spc, sep}}},
	{ts{sep, spc}, tg{ts{sep}, ts{spc}}},
	{ts{spc, sep, spc}, tg{ts{spc, sep}, ts{spc}}},
	{ts{d1, sep, spc, sep, d2}, tg{ts{d1, sep}, ts{spc, sep}, ts{d2}}},
	{ts{sep, sep}, tg{ts{sep}, ts{sep}}},
	{ts{sep, sep, spc}, tg{ts{sep}, ts{sep}, ts{spc}}},
	{ts{sep, sep, d1}, tg{ts{sep}, ts{sep}, ts{d1}}},
	{ts{d1, sep, sep}, tg{ts{d1, sep}, ts{sep}}},
	{ts{d1, sep, sep, spc}, tg{ts{d1, sep}, ts{sep}, ts{spc}}},
}

var normtab = []struct {
	col, norm ts
}{
	{ts{}, ts{}},
	{ts{spc}, ts{spc}},
	{ts{sep}, ts{sep}},
	{ts{sep, spc}, ts{sep, spc}}, // does not occur IRL as sep will be last
	{ts{spc, sep}, ts{spc, sep}},
	{ts{d1, spc}, ts{d1, spc}},
	{ts{d1, sep}, ts{d1, sep}},
	{ts{d1, d1}, ts{d11}},
	{ts{d1, d2}, ts{d12}},
	{ts{d1, d1, d1}, ts{d111}},
	{ts{d1, d2, d1}, ts{d121}},
	{ts{spc, d1, d2}, ts{spc, d12}},
	{ts{d1, d2, spc}, ts{d12, spc}},
	{ts{spc, d1, d2, spc}, ts{spc, d12, spc}},
	{ts{spc, spc, d1, d2, spc}, ts{spc, spc, d12, spc}},
	{ts{spc, d1, d2, spc, spc}, ts{spc, d12, spc, spc}},
	{ts{spc, spc, d1, d2, spc, spc}, ts{spc, spc, d12, spc, spc}},
	{ts{d1, spc, d1}, ts{d1s1}}, // important:
	// ^^ space surrounded by data should be incorporated into data
	{ts{d1, spc, d2}, ts{d1s2}},                               // as above
	{ts{spc, d1, spc, d2}, ts{spc, d1s2}},                     // as above
	{ts{d1, spc, d2, spc}, ts{d1s2, spc}},                     // as above
	{ts{spc, d1, spc, d2, spc}, ts{spc, d1s2, spc}},           // as above
	{ts{spc, spc, d1, spc, d2, spc}, ts{spc, spc, d1s2, spc}}, // as above
}

// tests

func TestParseGroup(t *testing.T) {
	for i, tc := range grouptab {
		res, err := GroupTokens(tc.tokens)
		if err != nil {
			panic(err)
		}
		if !reflect.DeepEqual(res, tc.group) {
			t.Errorf("tc#%d mismatch\nhave %v\nwant %v", i, res, tc.group)
		}
	}
}

func TestParseNormalize(t *testing.T) {
	for i, tc := range normtab {
		res := NormalizeColumn(tc.col)
		if !reflect.DeepEqual(res, tc.norm) {
			t.Errorf("tc#%d mismatch\nhave %v\nwant %v", i, res, tc.norm)
		}
	}
}
