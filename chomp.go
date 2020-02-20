package main

func chomp(b []byte) []byte {
	for len(b) > 0 {
		l := len(b) - 1
		if b[l] == '\n' || b[l] == '\r' {
			b = b[:l]
		} else {
			break
		}
	}
	return b
}
