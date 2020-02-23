package colt

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
		c := Colt{
			Separator: ';',
			Quote:     '"',
			Unquote:   true,
			Selection: tc.sel,
			Command:   cmd,
			Stdout:    o,
			Stderr:    ioutil.Discard,
		}
		c.ProcessLine(b)
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
			c := Colt{
				Separator: ';',
				Quote:     '"',
				Selection: tc.sel,
				Command:   cmd,
				Stdout:    ioutil.Discard,
				Stderr:    ioutil.Discard,
			}
			b := []byte(tc.input)
			c.ProcessLine(b)
		}
	}
}
