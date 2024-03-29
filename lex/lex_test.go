package lex

import (
	"bytes"
	"fmt"
	"reflect"
	"testing"
)

// shorter syntax in literals:
type b = []byte
type ts = []Token

func t(tt TokenType, data string) Token {
	return Token{Type: tt, Val: b(data)}
}
func te(tt TokenType, data string, text string, i, j int) Token {
	return Token{
		Type: tt, Val: b(data),
		Err: &Error{text, i, j},
	}
}

var sep = t(TokenSep, ";")

// tab

var lexTab = []struct {
	line   string
	tokens []Token
}{
	{"", ts{}},
	{"   ", ts{t(TokenSpace, "   ")}},
	{"foo", ts{t(TokenData, "foo")}},
	{"foo bar", ts{
		t(TokenData, "foo"), t(TokenSpace, " "),
		t(TokenData, "bar"),
	}},
	{";", ts{sep}},
	{";;", ts{sep, sep}},
	{"; ;", ts{sep, t(TokenSpace, " "), sep}},
	{";;;", ts{sep, sep, sep}},
	{"foo;", ts{t(TokenData, "foo"), sep}},
	{"foo;  bar", ts{
		t(TokenData, "foo"), sep,
		t(TokenSpace, "  "), t(TokenData, "bar"),
	}},
	{`foo; "a;b;c"; "bar"`, ts{
		t(TokenData, "foo"), sep, t(TokenSpace, " "),
		t(TokenData, `"a;b;c"`), sep, t(TokenSpace, " "),
		t(TokenData, `"bar"`),
	}},
	{`"foo`, ts{
		te(TokenError, `"foo`, "unclosed quote", 0, 4),
	}},
	{`"foo` + "\n", ts{
		te(TokenError, `"foo`, "unclosed quote", 0, 4),
	}},
	{`123"foo"`, ts{
		t(TokenData, "123"), t(TokenData, `"foo"`),
	}},
	{`123"foo"456`, ts{
		t(TokenData, "123"), t(TokenData, `"foo"`), t(TokenData, "456"),
	}},
	{`123"foo`, ts{
		t(TokenData, "123"),
		te(TokenError, `"foo`, "unclosed quote", 3, 7),
	}},
	{`123"foo` + "\n", ts{
		t(TokenData, "123"),
		te(TokenError, `"foo`, "unclosed quote", 3, 7),
	}},
	{`123"foo` + "\n" + `456`, ts{
		t(TokenData, "123"),
		te(TokenError, `"foo`, "unclosed quote", 3, 7),
	}},
	{"\x00", ts{
		te(TokenError, "\x00", "binary data", 0, 1),
	}},
	{"  \x00", ts{
		t(TokenSpace, "  "),
		te(TokenError, "\x00", "binary data", 2, 3),
	}},
	{"123\x00", ts{
		te(TokenError, "123\x00", "binary data", 0, 4),
	}},
	{"123\x00456", ts{
		te(TokenError, "123\x00", "binary data", 0, 4),
	}},
	{"123\"\x00", ts{
		t(TokenData, "123"),
		te(TokenError, "\"\x00", "binary data", 3, 5),
	}},
	{"123\"\x00456", ts{
		t(TokenData, "123"),
		te(TokenError, "\"\x00", "binary data", 3, 5),
	}},
}

// pretty-print

func (t Token) String() string {
	if t.Err != nil {
		return fmt.Sprintf("{%d %q %q}", t.Type, t.Val, t.Err)
	}
	return fmt.Sprintf("{%d %q}", t.Type, t.Val)
}

// tests

func TestLex(t *testing.T) {
	for i, tc := range lexTab {
		res := LexTokens(b(tc.line), DefaultConfig)
		if !reflect.DeepEqual(res, tc.tokens) {
			t.Errorf("tc[%d] mismatch\nhave %v\nwant %v",
				i, res, tc.tokens)
		}
	}
}

func FuzzLex(f *testing.F) {
	for _, tc := range lexTab {
		f.Add(tc.line)
	}

	f.Fuzz(func(t *testing.T, data string) {
		end := len(data)
		var buf bytes.Buffer

		for _, tok := range LexTokens(b(data), DefaultConfig) {
			buf.Write(tok.Val)
			if tok.Type == TokenError {
				end = tok.Err.End
			}
		}
		res := buf.String()
		if res != data[:end] {
			t.Errorf("mismatch\nhave `%v`\nwant `%v`", res, data)
		}
	})
}
