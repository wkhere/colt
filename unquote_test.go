package main

import (
	"testing"
)

func TestUnquote(t *testing.T) {
	var tab = []struct {
		input string
		want  string
	}{
		{``, ``},
		{`""`, ``},
		{`"`, `"`},
		{`foo`, `foo`},
		{`"foo"`, `foo`},
		{`"with spaces inside"`, `with spaces inside`},
		{`"unfinished`, `"unfinished`},
		{`unstarted"`, `unstarted"`},
		{`unstarted2"..`, `unstarted2"..`},
		{`bro"ken`, `bro"ken`},
		{`."head"`, `."head"`},
		{`"tail".`, `"tail".`},
		{`."headtail".`, `."headtail".`},
	}

	for i, tc := range tab {
		b := []byte(tc.input)
		res := unquote(b, '"')
		if string(res) != tc.want {
			t.Errorf("tc[%d] mismatch\nhave %v\nwant %v", i, res, tc.want)
		}
	}
}
