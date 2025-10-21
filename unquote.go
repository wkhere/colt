package colt

import "unicode/utf8"

// use strconv.Unquote instead?
// - might be fine, but also that func unquotes any pair of '"`,
// then, I should stop passing it here, stop making it configurable,
// remove from args, from the Colt struct and also from lex.Config

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
