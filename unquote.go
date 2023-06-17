package colt

import "unicode/utf8"

// use strconv.Unquote instead?
// - might be fine, but also that func unquotes any pair of '"`,
// then, we should actually fix it (stop thinking about making it
// configurable, remove from Colt struct and possibly from lex.Config)

func unquote(b []byte, q rune) []byte {
	l := utf8.RuneLen(q)
	if len(b) < 2*l {
		return b
	}
	eqq := func(p []byte) bool {
		c, n := utf8.DecodeRune(p)
		return c == q && n == len(p)
	}
	if eqq(b[:l]) && eqq(b[len(b)-l:]) {
		return b[l : len(b)-l]
	}
	return b
}
