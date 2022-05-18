package main

import (
	"fmt"
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

func warn(err error) {
	fmt.Fprintf(os.Stderr, "--WARN %v\n", err)
}
