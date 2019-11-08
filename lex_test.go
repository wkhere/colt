package main

import (
	"fmt"
	"reflect"
	"testing"
)

// shorter syntax in literals:
type ts = []token

var sep = token{tokenSep, ";"}

var lexTab = []struct {
	line   string
	tokens []token
}{
	{"", nil},
	{"   ", ts{{tokenSpace, "   "}}},
	{"foo", ts{{tokenData, "foo"}}},
	{"foo bar", ts{
		{tokenData, "foo"}, {tokenSpace, " "},
		{tokenData, "bar"}}},
	{";", ts{sep}},
	{"foo;", ts{{tokenData, "foo"}, sep}},
	{"foo;  bar", ts{
		{tokenData, "foo"}, sep,
		{tokenSpace, "  "}, {tokenData, "bar"},
	}},
	{`foo; "a;b;c"; "bar"`, ts{
		{tokenData, "foo"}, sep, {tokenSpace, " "},
		{tokenData, `"a;b;c"`}, sep, {tokenSpace, " "},
		{tokenData, `"bar"`},
	}},
}

func (t token) String() string {
	return fmt.Sprintf("{%d %q}", t.typ, t.val)
}

var eq = reflect.DeepEqual

func TestLex(t *testing.T) {
	for i, tc := range lexTab {
		if res := lexTokens(tc.line, ';', '"').flatten(); !eq(res, tc.tokens) {
			t.Errorf("tc[%d] mismatch\nhave %v\nwant %v",
				i, res, tc.tokens)
		}
	}
}
