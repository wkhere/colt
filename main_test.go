package main

import (
	"bytes"
	"io/ioutil"
	"testing"
)

var cmd = []string{"perl", "-e", `print "[", uc($ARGV[0]), "]"`}

var tab = []struct {
	sel    int
	input  string
	output string
}{
	{1, "", ""},
	{1, ";", ";"},
	{1, ";;", ";;"},
	{1, " ", " "},
	{1, "  ", "  "},
	{1, "aa", "[AA]"},
	{1, "aa  ", "[AA]  "},
	{1, "  aa", "  [AA]"},
	{1, "  aa  ", "  [AA]  "},
	{1, "aa;bb", "[AA];bb"},
	{1, "aa; bb", "[AA]; bb"},
	{1, "aa ;bb", "[AA] ;bb"},
	{1, "aa ; bb", "[AA] ; bb"},
	{1, " aa; bb", " [AA]; bb"},
	{2, "aa;bb", "aa;[BB]"},
	{2, "aa;bb ", "aa;[BB] "},
	{2, "aa; bb", "aa; [BB]"},
	{2, "aa ;bb", "aa ;[BB]"},
	{2, "aa ; bb", "aa ; [BB]"},
	{3, "aa;bb", "aa;bb"},
	{-1, "aa", "[AA]"},
	{-1, "aa;bb", "aa;[BB]"},
	{-2, "aa;bb", "[AA];bb"},
	{-3, "aa;bb", "aa;bb"},
	{1, " aa; bb;cc ; dd ;", " [AA]; bb;cc ; dd ;"},
	{2, " aa; bb;cc ; dd ;", " aa; [BB];cc ; dd ;"},
	{3, " aa; bb;cc ; dd ;", " aa; bb;[CC] ; dd ;"},
	{4, " aa; bb;cc ; dd ;", " aa; bb;cc ; [DD] ;"},
	{1, "multi word; bb", "[MULTI WORD]; bb"},
	{1, "multi word ; bb", "[MULTI WORD] ; bb"},
	{1, " multi word; bb", " [MULTI WORD]; bb"},
	{1, " multi word ; bb", " [MULTI WORD] ; bb"},
	{1, " multi word  ; bb", " [MULTI WORD]  ; bb"},
	{0, "aa;bb", "aa;bb"},
}

func TestProcess(t *testing.T) {
	for i, tc := range tab {
		var b bytes.Buffer
		p := columnProc{
			separator: ';',
			selection: tc.sel,
			command:   cmd,
			stdout:    &b,
			stderr:    ioutil.Discard,
		}
		p.process(tc.input)
		if res := b.String(); res != tc.output {
			t.Errorf("tc[%d] mismatch\nhave %v\nwant %v", i, res, tc.output)
		}
	}
}

func BenchmarkProcess(b *testing.B) {
	for i := 0; i < b.N; i++ {
		for _, tc := range tab {
			p := columnProc{
				separator: ';',
				selection: tc.sel,
				command:   cmd,
				stdout:    ioutil.Discard,
				stderr:    ioutil.Discard,
			}
			p.process(tc.input)
		}
	}
}
