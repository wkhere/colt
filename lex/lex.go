package lex

import (
	"unicode"
	"unicode/utf8"
)

// lexer interface

type TokenType uint

const (
	TokenError TokenType = iota
	TokenSpace
	TokenSep
	TokenData
)

type Token struct {
	Type TokenType
	Val  []byte
}

type TokenStream <-chan Token

func LexTokens(input []byte, sep, quote rune) TokenStream {
	l := &lexer{
		sep:    sep,
		quote:  quote,
		input:  input,
		tokens: make(chan Token),
	}
	go l.run()
	return l.tokens
}

func (ts TokenStream) Flatten() (res []Token) {
	for tok := range ts {
		res = append(res, tok)
	}
	return
}

func (ts TokenStream) Group() (res [][]Token) {
	res = make([][]Token, 1, 3)
	groupIdx := 0
	for tok := range ts {
		res[groupIdx] = append(res[groupIdx], tok)
		if tok.Type == TokenSep {
			res = append(res, nil)
			groupIdx++
		}
	}
	return
}

// engine

type lexer struct {
	sep, quote rune
	input      []byte
	start, pos int
	lastw      int
	tokens     chan Token
}

type stateFn func(*lexer) stateFn

func (l *lexer) run() {
	for st := lexStart; st != nil; {
		st = st(l)
	}
	close(l.tokens)
}

func (l *lexer) emit(t TokenType) {
	l.tokens <- Token{Type: t, Val: l.input[l.start:l.pos]}
	l.start = l.pos
}

// input-consuming primitives

const (
	cEOF rune = -1
	cLF       = '\n'
)

func (l *lexer) readc() (c rune) {
	c, l.lastw = utf8.DecodeRune(l.input[l.pos:])
	if l.lastw == 0 {
		return cEOF
	}
	l.pos += l.lastw
	return c
}

// backup can be used only once after each readc.
func (l *lexer) backup() {
	l.pos -= l.lastw
}

func (l *lexer) unbackup() {
	l.pos += l.lastw
}

// func (l *lexer) peek() rune {
// 	c := l.readc()
// 	l.backup()
// 	return c
// }

// input-consuming helpers

func (l *lexer) acceptAny(pred func(rune) bool) {
	for pred(l.readc()) {
	}
	l.backup()
}

func (l *lexer) skipUntil(pred func(rune) bool) {
	for {
		if c := l.readc(); c == cEOF || pred(c) {
			break
		}
	}
	l.backup()
}

// state functions

func lexStart(l *lexer) stateFn {
	switch c := l.readc(); {
	case c == cEOF:
		return nil
	case c == l.sep:
		l.emit(TokenSep)
		return lexStart
	case c == l.quote:
		return lexQuoted
	case unicode.IsSpace(c):
		return lexSpace
	default:
		l.backup()
		return lexData
	}
}

func lexSpace(l *lexer) stateFn {
	l.acceptAny(unicode.IsSpace)
	l.emit(TokenSpace)
	return lexStart
}

func lexData(l *lexer) stateFn {
	l.skipUntil(func(c rune) bool {
		return c == l.sep || unicode.IsSpace(c)
	})
	l.emit(TokenData)
	return lexStart
}

func lexQuoted(l *lexer) stateFn {
	var closed bool
	l.skipUntil(func(c rune) bool {
		switch c {
		case l.quote:
			closed = true
			return true
		case cLF:
			return true
		default:
			return false
		}
	})
	if !closed {
		l.emit(TokenError)
		return nil
	}
	l.unbackup()
	l.emit(TokenData)
	return lexStart
}
