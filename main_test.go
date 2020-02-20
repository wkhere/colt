package main

import (
	"bytes"
	"io/ioutil"
	"testing"
)

var (
	cmd       = []string{"perl", "-e", `print "[", uc($ARGV[0]), "]"`}
	cmdWithLF = []string{"perl", "-e", `print "[", uc($ARGV[0]), "]\n"`}
)

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
	{3, " aa; bb; cc; dd ;", " aa; bb; [CC]; dd ;"},
	{3, " aa; bb; cc ; dd ;", " aa; bb; [CC] ; dd ;"},
	{4, " aa; bb;cc ; dd ;", " aa; bb;cc ; [DD] ;"},
	{1, "multi word; bb", "[MULTI WORD]; bb"},
	{1, "multi word ; bb", "[MULTI WORD] ; bb"},
	{1, " multi word; bb", " [MULTI WORD]; bb"},
	{1, " multi word ; bb", " [MULTI WORD] ; bb"},
	{1, " multi word  ; bb", " [MULTI WORD]  ; bb"},
	{1, `"quoted; thing"; bb`, `[QUOTED; THING]; bb`},
	{1, `"quoted; thing" ; bb`, `[QUOTED; THING] ; bb`},
	{1, ` "quoted; thing"; bb`, ` [QUOTED; THING]; bb`},
	{1, ` "quoted; thing" ; bb`, ` [QUOTED; THING] ; bb`},
	{0, "aa;bb", "aa;bb"},
}

func testWithCmd(cmd []string, t *testing.T) {
	t.Helper()
	for i, tc := range tab {
		b := []byte(tc.input)
		o := new(bytes.Buffer)
		p := columnProc{
			separator: ';',
			quote:     '"',
			unquote:   true,
			selection: tc.sel,
			command:   cmd,
			stdout:    o,
			stderr:    ioutil.Discard,
		}
		p.process(b)
		res := o.String()
		if res != tc.output {
			t.Errorf("tc[%d] mismatch\nhave %q\nwant %q", i, res, tc.output)
		}
	}
}

func TestProcess(t *testing.T) {
	testWithCmd(cmd, t)
}

func TestProcessWithLF(t *testing.T) {
	testWithCmd(cmdWithLF, t)
}

func BenchmarkProcess(b *testing.B) {
	for i := 0; i < b.N; i++ {
		for _, tc := range tab {
			p := columnProc{
				separator: ';',
				quote:     '"',
				selection: tc.sel,
				command:   cmd,
				stdout:    ioutil.Discard,
				stderr:    ioutil.Discard,
			}
			b := []byte(tc.input)
			p.process(b)
		}
	}
}
