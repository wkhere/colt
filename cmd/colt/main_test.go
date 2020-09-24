package main

import (
	"bytes"
	"io"
	"io/ioutil"
	"os"
)

func ExampleMain() {

	os.Args = os.Args[:1]
	os.Args = append(os.Args,
		`perl`, `-e`, `print "[", uc($ARGV[0]), "]"`)

	b := new(bytes.Buffer)
	b.WriteString("a;b;c\n")
	b.WriteString("a;b;c d\n")
	b.WriteString(`a;b;"c;d"` + "\n")

	feed(&os.Stdin, b)
	main()

	// Output:
	// a;b;[C]
	// a;b;[C D]
	// a;b;["C;D"]
}

func feed(fp **os.File, b io.Reader) {
	pr, pw, err := os.Pipe()
	if err != nil {
		panic(err)
	}
	*fp = pr
	ioutil.ReadAll(io.TeeReader(b, pw))
	pw.Close()
}
