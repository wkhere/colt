package main

import "testing"

type a = []string

func TestSelectionArg(t *testing.T) {
	tab := []struct {
		opts      []string
		selection int
		err       error
	}{
		{a{"-1"}, -1, nil},
		{a{"+1"}, 1, nil},
		{a{"+BAD"}, 0, usageErr},
	}

	for i, tc := range tab {
		p := new(columnProc)
		err := p.parseArgs(append(tc.opts, "echo"))
		if tc.err != nil {
			if err == nil {
				t.Errorf("tc[%d] have nil, want error", i)
			}
			continue
		}
		if p.selection != tc.selection {
			t.Errorf("tc[%d] mismatch\nhave %v\nwant %v",
				i, p.selection, tc.selection)
		}
	}

}
