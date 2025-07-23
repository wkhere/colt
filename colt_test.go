package colt

import (
	"bytes"
	"io"
	"io/ioutil"
	"strings"
	"testing"
)

var (
	cmdNoLF   = []string{"perl", "-e", `print "[", uc($ARGV[0]), "]"`}
	cmdWithLF = []string{"perl", "-e", `print "[", uc($ARGV[0]), "]\n"`}
)

type softT struct{ io.Writer }

func (t softT) Copy(b []byte) error { t.Write(b); return nil }

func (t softT) Transform(b []byte) error {
	t.Write([]byte{'['})
	t.Write(bytes.ToUpper(b))
	t.Write([]byte{']'})
	return nil
}

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
	te(1, "", "", "invalid column number"),
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

func testWithSoftT(t *testing.T) {
	t.Helper()
	test(t, func(w io.Writer) Transformer {
		return softT{w}
	})
}

func testWithCmd(cmd []string, t *testing.T) {
	t.Helper()
	test(t, func(w io.Writer) Transformer {
		return &CommandT{
			Command: cmd,
			Stdout:  w,
			Stderr:  ioutil.Discard,
		}
	})
}

func test(t *testing.T, trgen func(io.Writer) Transformer) {
	t.Helper()

	for i, tc := range tab {
		b := []byte(tc.input)
		o := new(bytes.Buffer)
		c := Colt{
			Separator: ';',
			Quote:     '"',
			Unquote:   true,
			Selection: tc.sel,
			T:         trgen(o),
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

func TestProcessSoft(t *testing.T) {
	testWithSoftT(t)
}

func TestProcessCmd(t *testing.T) {
	testWithCmd(cmdNoLF, t)
}

func TestProcessCMDWithLF(t *testing.T) {
	testWithCmd(cmdWithLF, t)
}

func BenchmarkProcessSoft(b *testing.B) {
	for i := 0; i < b.N; i++ {
		for _, tc := range tab {
			c := Colt{
				Separator: ';',
				Quote:     '"',
				Selection: tc.sel,
				T:         softT{io.Discard},
			}
			b := []byte(tc.input)
			c.ProcessLine(b)
		}
	}
}
func BenchmarkProcessCmd(b *testing.B) {
	for i := 0; i < b.N; i++ {
		for _, tc := range tab {
			c := Colt{
				Separator: ';',
				Quote:     '"',
				Selection: tc.sel,
				T: &CommandT{
					Command: cmdNoLF,
					Stdout:  ioutil.Discard,
					Stderr:  ioutil.Discard,
				},
			}
			b := []byte(tc.input)
			c.ProcessLine(b)
		}
	}
}
