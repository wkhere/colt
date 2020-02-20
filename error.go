package main

import (
	"fmt"
	"io"
	"os"
)

func die2(err error) {
	fmt.Fprintln(os.Stderr, err)
	os.Exit(2)
}

func die(err error) {
	fmt.Fprintln(os.Stderr, "colt:", err)
	os.Exit(1)
}

func dieIf(err error) {
	if err != nil {
		die(err)
	}
}

func fwarn(w io.Writer, err error) {
	fmt.Fprintf(w, "WARN %v\n", err)
}

func errorf(format string, a ...interface{}) error {
	return fmt.Errorf(format, a...)
}
