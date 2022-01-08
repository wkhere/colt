package lex

import (
	"fmt"
	"reflect"
	"testing"
)

// shorter syntax in literals:
type b = []byte
type ts = []Token

var sep = Token{TokenSep, b(";")}

var lexTab = []struct {
	line   string
	tokens []Token
}{
	{"", nil},
	{"   ", ts{{TokenSpace, b("   ")}}},
	{"foo", ts{{TokenData, b("foo")}}},
	{"foo bar", ts{
		{TokenData, b("foo")}, {TokenSpace, b(" ")},
		{TokenData, b("bar")}}},
	{";", ts{sep}},
	{"foo;", ts{{TokenData, b("foo")}, sep}},
	{"foo;  bar", ts{
		{TokenData, b("foo")}, sep,
		{TokenSpace, b("  ")}, {TokenData, b("bar")},
	}},
	{`foo; "a;b;c"; "bar"`, ts{
		{TokenData, b("foo")}, sep, {TokenSpace, b(" ")},
		{TokenData, b(`"a;b;c"`)}, sep, {TokenSpace, b(" ")},
		{TokenData, b(`"bar"`)},
	}},
	{`"foo`, ts{
		{TokenError, b(`"foo`)},
	}},
	{`"foo` + "\n", ts{
		{TokenError, b(`"foo`)},
	}},
}

func (t Token) String() string {
	return fmt.Sprintf("{%d %q}", t.Type, t.Val)
}

func (ts TokenStream) flatten() (res []Token) {
	for tok := range ts {
		res = append(res, tok)
	}
	return
}

var eq = reflect.DeepEqual

func TestLex(t *testing.T) {
	for i, tc := range lexTab {
		res := LexTokens(b(tc.line), ';', '"').flatten()
		if !eq(res, tc.tokens) {
			t.Errorf("tc[%d] mismatch\nhave %v\nwant %v",
				i, res, tc.tokens)
		}
	}
}
