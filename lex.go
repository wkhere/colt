package colt

import (
	"unicode"
	"unicode/utf8"
)

// lexer interface

type tokenType uint

const (
	tokenError tokenType = iota
	tokenSpace
	tokenSep
	tokenData
)

type token struct {
	typ tokenType
	val []byte
}

type tokenStream <-chan token

func lexTokens(input []byte, sep, quote rune) tokenStream {
	l := &lexer{
		sep:    sep,
		quote:  quote,
		input:  input,
		tokens: make(chan token),
	}
	go l.run()
	return l.tokens
}

func (ts tokenStream) flatten() (res []token) {
	for tok := range ts {
		res = append(res, tok)
	}
	return
}

func (ts tokenStream) group() (res [][]token) {
	res = make([][]token, 1, 3)
	groupIdx := 0
	for tok := range ts {
		res[groupIdx] = append(res[groupIdx], tok)
		if tok.typ == tokenSep {
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
	tokens     chan token
}

type stateFn func(*lexer) stateFn

func (l *lexer) run() {
	for st := lexStart; st != nil; {
		st = st(l)
	}
	close(l.tokens)
}

func (l *lexer) emit(t tokenType) {
	l.tokens <- token{typ: t, val: l.input[l.start:l.pos]}
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
		l.emit(tokenSep)
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
	l.emit(tokenSpace)
	return lexStart
}

func lexData(l *lexer) stateFn {
	l.skipUntil(func(c rune) bool {
		return c == l.sep || unicode.IsSpace(c)
	})
	l.emit(tokenData)
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
		l.emit(tokenError)
		return nil
	}
	l.unbackup()
	l.emit(tokenData)
	return lexStart
}
