package main

import "bytes"

func normalizeColumn(col []token) (res []token) {
	res = make([]token, 0, len(col))
	var i, j int

	for i = 0; i < len(col) && col[i].typ != tokenData; i++ {
		res = append(res, col[i])
	}
	if i == len(col) {
		return res
	}
	for j = len(col) - 1; j > i && col[j].typ != tokenData; j-- {
	}
	// now i is at the first data token and j at the last

	if i == j {
		res = append(res, col[i])
	} else {
		var b bytes.Buffer
		for _, token := range col[i : j+1] {
			b.Write(token.val)
		}
		res = append(res, token{typ: tokenData, val: b.Bytes()})
	}

	for j++; j < len(col); j++ {
		res = append(res, col[j])
	}

	return res
}
