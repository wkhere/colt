package colt

import (
	"fmt"
	"reflect"
	"testing"
)

// shorter syntax in literals:
type b = []byte
type ts = []token

var sep = token{tokenSep, b(";")}

var lexTab = []struct {
	line   string
	tokens []token
}{
	{"", nil},
	{"   ", ts{{tokenSpace, b("   ")}}},
	{"foo", ts{{tokenData, b("foo")}}},
	{"foo bar", ts{
		{tokenData, b("foo")}, {tokenSpace, b(" ")},
		{tokenData, b("bar")}}},
	{";", ts{sep}},
	{"foo;", ts{{tokenData, b("foo")}, sep}},
	{"foo;  bar", ts{
		{tokenData, b("foo")}, sep,
		{tokenSpace, b("  ")}, {tokenData, b("bar")},
	}},
	{`foo; "a;b;c"; "bar"`, ts{
		{tokenData, b("foo")}, sep, {tokenSpace, b(" ")},
		{tokenData, b(`"a;b;c"`)}, sep, {tokenSpace, b(" ")},
		{tokenData, b(`"bar"`)},
	}},
	{`"foo`, ts{
		{tokenError, b(`"foo`)},
	}},
	{`"foo` + "\n", ts{
		{tokenError, b(`"foo`)},
	}},
}

func (t token) String() string {
	return fmt.Sprintf("{%d %q}", t.typ, t.val)
}

var eq = reflect.DeepEqual

func TestLex(t *testing.T) {
	for i, tc := range lexTab {
		res := lexTokens(b(tc.line), ';', '"').flatten()
		if !eq(res, tc.tokens) {
			t.Errorf("tc[%d] mismatch\nhave %v\nwant %v",
				i, res, tc.tokens)
		}
	}
}
