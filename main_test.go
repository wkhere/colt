package main

import (
	"bytes"
	"testing"
)

func TestProcess(t *testing.T) {
	var cmd = []string{"perl", "-e", "print uc $ARGV[0]"}

	var tab = []struct {
		sel    int
		input  string
		output string
	}{
		{1, "aa", "AA"},
		{1, "aa;bb", "AA;bb"},
		{1, "aa; bb", "AA; bb"},
		{1, "aa ;bb", "AA ;bb"},
		{1, "aa ; bb", "AA ; bb"},
		{1, " aa; bb", " AA; bb"},
		{2, "aa;bb", "aa;BB"},
		{2, "aa;bb ", "aa;BB "},
		{2, "aa; bb", "aa; BB"},
		{2, "aa ;bb", "aa ;BB"},
		{2, "aa ; bb", "aa ; BB"},
		{3, "aa;bb", "aa;bb"},
		{-1, "aa", "AA"},
		{-1, "aa;bb", "aa;BB"},
		{-2, "aa;bb", "AA;bb"},
		{-3, "aa;bb", "aa;bb"},
		{0, "aa;bb", "aa;bb"},
	}

	for i, tc := range tab {
		var b bytes.Buffer
		p := columnProc{
			separator: ';',
			selection: tc.sel,
			command:   cmd,
			output:    &b,
		}
		p.process(tc.input)
		if res := b.String(); res != tc.output {
			t.Errorf("tc[%d] mismatch\ngot %v\nexp %v", i, res, tc.output)
		}
	}
}
