package main

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
	val string
}

type tokenStream <-chan token

func lexTokens(input string, sep rune) tokenStream {
	l := &lexer{
		sep:    sep,
		input:  input,
		tokens: make(chan token),
	}
	go l.run()
	return l.tokens
}

func (ts tokenStream) gather() (res []token) {
	for tok := range ts {
		res = append(res, tok)
	}
	return
}

// engine

type lexer struct {
	sep        rune
	input      string
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

const cEOF rune = -1

func (l *lexer) readc() (c rune) {
	c, l.lastw = utf8.DecodeRuneInString(l.input[l.pos:])
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
