package lex

import (
	"unicode"
	"unicode/utf8"
)

// lexer api

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
	Err  *Error
}

type Error struct {
	Text       string
	Start, End int
}

type Config struct {
	Sep, Quote rune
	BufSize    int
}

var DefaultConfig = Config{
	Sep:     ';',
	Quote:   '"',
	BufSize: 20,
}

func LexTokens(input []byte, conf Config) []Token {
	l := &lexer{
		sep:    conf.Sep,
		quote:  conf.Quote,
		input:  input,
		tokens: make([]Token, 0, conf.BufSize),
	}
	l.run()
	return l.tokens
}

// engine

type lexer struct {
	sep, quote rune
	input      []byte
	start, pos int
	lastw      int
	tokens     []Token
}

type stateFn func(*lexer) stateFn

func (l *lexer) run() {
	for st := lexStart; st != nil; {
		st = st(l)
	}
}

func (l *lexer) emit(t TokenType) {
	l.tokens = append(l.tokens, Token{Type: t, Val: l.input[l.start:l.pos]})
	l.start = l.pos
}

func (l *lexer) emitError(text string) {
	l.tokens = append(l.tokens, Token{
		Type: TokenError, Val: l.input[l.start:l.pos],
		Err: &Error{text, l.start, l.pos},
	})
	l.start = l.pos
}

// input-consuming primitives

const (
	cEOF rune = -1
	cBin      = 0
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
	var bin bool
	l.skipUntil(func(c rune) bool {
		switch {
		case c == cBin:
			bin = true
			return true
		case c == l.sep || c == l.quote || unicode.IsSpace(c):
			return true
		default:
			return false
		}
	})
	if bin {
		l.unbackup()
		l.emitError("binary data")
		return nil
	}
	l.emit(TokenData)
	return lexStart
}

func lexQuoted(l *lexer) stateFn {
	var bin, closed bool
	l.skipUntil(func(c rune) bool {
		switch c {
		case cBin:
			bin = true
			return true
		case l.quote:
			closed = true
			return true
		case cLF:
			return true
		default:
			return false
		}
	})
	if bin {
		l.unbackup()
		l.emitError("binary data")
		return nil
	}
	if !closed {
		l.emitError("unclosed quote")
		return nil
	}
	l.unbackup()
	l.emit(TokenData)
	return lexStart
}
