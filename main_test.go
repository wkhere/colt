package main

import (
	"bytes"
	"testing"
)

func TestProcessFields(t *testing.T) {
	var cmd = []string{"echo", "-n", "AA"}

	var tab = []struct {
		field  int
		input  string
		output string
	}{
		{1, "aa", "AA aa"},
		{1, "aa;bb", "AA aa;bb"},
		{2, "aa;bb", "aa;AA bb"},
		{3, "aa;bb", "aa;bb"},
	}

	for i, tc := range tab {
		var b bytes.Buffer
		p := columnProc{';', tc.field, cmd, &b}
		p.processFields(tc.input)
		if res := b.String(); res != tc.output {
			t.Errorf("tc[%d] mismatch\ngot %v\nexp %v", i, res, tc.output)
		}
	}
}
