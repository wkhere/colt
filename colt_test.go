package colt

import (
	"bytes"
	"io/ioutil"
	"strings"
	"testing"
)

var (
	cmd       = []string{"perl", "-e", `print "[", uc($ARGV[0]), "]"`}
	cmdWithLF = []string{"perl", "-e", `print "[", uc($ARGV[0]), "]\n"`}
)

type testcase struct {
	sel    int
	input  string
	output string
	err    string
}

func t(sel int, i, o string) testcase {
	return testcase{sel: sel, input: i, output: o}
}
func te(sel int, i, o string, e string) testcase {
	return testcase{sel: sel, input: i, output: o, err: e}
}

var tab = []testcase{
	t(1, "", ""),
	t(1, ";", ";"),
	t(1, ";;", ";;"),
	t(1, " ", " "),
	t(1, "  ", "  "),
	t(1, "aa", "[AA]"),
	t(1, "aa  ", "[AA]  "),
	t(1, "  aa", "  [AA]"),
	t(1, "  aa  ", "  [AA]  "),
	t(1, "aa;bb", "[AA];bb"),
	t(1, "aa; bb", "[AA]; bb"),
	t(1, "aa ;bb", "[AA] ;bb"),
	t(1, "aa ; bb", "[AA] ; bb"),
	t(1, " aa; bb", " [AA]; bb"),
	t(2, "aa;bb", "aa;[BB]"),
	t(2, "aa;bb ", "aa;[BB] "),
	t(2, "aa; bb", "aa; [BB]"),
	t(2, "aa ;bb", "aa ;[BB]"),
	t(2, "aa ; bb", "aa ; [BB]"),
	te(3, "aa;bb", "aa;bb", "invalid column number"),
	t(-1, "aa", "[AA]"),
	t(-1, "aa;bb", "aa;[BB]"),
	t(-2, "aa;bb", "[AA];bb"),
	te(-3, "aa;bb", "aa;bb", "invalid column number"),
	t(1, " aa; bb;cc ; dd ;", " [AA]; bb;cc ; dd ;"),
	t(2, " aa; bb;cc ; dd ;", " aa; [BB];cc ; dd ;"),
	t(3, " aa; bb;cc ; dd ;", " aa; bb;[CC] ; dd ;"),
	t(3, " aa; bb; cc; dd ;", " aa; bb; [CC]; dd ;"),
	t(3, " aa; bb; cc ; dd ;", " aa; bb; [CC] ; dd ;"),
	t(4, " aa; bb;cc ; dd ;", " aa; bb;cc ; [DD] ;"),
	t(1, "multi word; bb", "[MULTI WORD]; bb"),
	t(1, "multi word ; bb", "[MULTI WORD] ; bb"),
	t(1, " multi word; bb", " [MULTI WORD]; bb"),
	t(1, " multi word ; bb", " [MULTI WORD] ; bb"),
	t(1, " multi word  ; bb", " [MULTI WORD]  ; bb"),
	t(1, `"quoted; thing"; bb`, `[QUOTED; THING]; bb`),
	t(1, `"quoted; thing" ; bb`, `[QUOTED; THING] ; bb`),
	t(1, ` "quoted; thing"; bb`, ` [QUOTED; THING]; bb`),
	t(1, ` "quoted; thing" ; bb`, ` [QUOTED; THING] ; bb`),
	te(1, `"unclosed quote`, "", "unclosed quote"),
	te(0, "aa;bb", "aa;bb", "invalid column number"),
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
		err := c.ProcessLine(b)
		res := o.String()

		switch {
		case tc.err == "" && err != nil:
			t.Errorf("tc#%d unexpected error: %v", i, err)
		case tc.err != "" && err == nil:
			t.Errorf(
				"tc#%d no error; want error with substring: %v",
				i, tc.err)
		case tc.err != "" && err != nil:
			if !strings.Contains(err.Error(), tc.err) {
				t.Errorf(
					"tc#%d error mismatch\nhave: %v\nwant substring: %v",
					i, err, tc.err)
			}
		default:
			if res != tc.output {
				t.Errorf("tc#%d mismatch\nhave %q\nwant %q",
					i, res, tc.output)
			}
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
