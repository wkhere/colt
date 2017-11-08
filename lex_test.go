package main

import (
	"fmt"
	"reflect"
	"testing"
)

// aliases to allow shorter syntax in literals:
type bs = []byte
type ts = []token

var sep = token{tokenSep, bs(";")}

var lexTab = []struct {
	line   []byte
	tokens []token
}{
	{bs(""), nil},
	{bs("   "), ts{{tokenSpace, bs("   ")}}},
	{bs("foo"), ts{{tokenData, bs("foo")}}},
	{bs("foo bar"), ts{
		{tokenData, bs("foo")}, {tokenSpace, bs(" ")},
		{tokenData, bs("bar")}}},
	{bs(";"), ts{sep}},
	{bs("foo;"), ts{{tokenData, bs("foo")}, sep}},
	{bs("foo;  bar"), ts{
		{tokenData, bs("foo")}, sep,
		{tokenSpace, bs("  ")}, {tokenData, bs("bar")},
	}},
}

func (t token) String() string {
	return fmt.Sprintf("{%d %q}", t.typ, string(t.val))
}

func (ts tokenStream) gather() (res []token) {
	for tok := range ts {
		res = append(res, tok)
	}
	return
}

var eq = reflect.DeepEqual

func TestLex(t *testing.T) {
	for i, tc := range lexTab {
		if res := lexTokens(tc.line, ';').gather(); !eq(res, tc.tokens) {
			t.Errorf("tc[%d] mismatch\ngot %v\nexp %v",
				i, res, tc.tokens)
		}
	}
}
