package colt

import "testing"

func TestChomp(t *testing.T) {
	var tab = []struct{ input, want string }{
		{"", ""},
		{"\n", ""},
		{"\n\n", ""},
		{"\r\n", ""},
		{"\r\n\r\n", ""},
		{"aa", "aa"},
		{"aa\n", "aa"},
		{"aa\r\n", "aa"},
		{"aa \n", "aa "},
		{"aa \r\n", "aa "},
	}

	for i, tc := range tab {
		res := string(chomp([]byte(tc.input)))
		if res != tc.want {
			t.Errorf("tc#%d mismatch\nhave: `%s`\nwant: `%s`",
				i, res, tc.want)
		}
	}
}
