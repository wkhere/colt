package main

func unquote(s, q string) string {
	l := len(q)
	if len(s) < 2*l {
		return s
	}
	if s[:l] == q && s[len(s)-l:len(s)] == q {
		return s[l : len(s)-l]
	}
	return s
}
