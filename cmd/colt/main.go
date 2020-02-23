package main

import (
	"bufio"
	"fmt"
	"os"

	"github.com/wkhere/colt"
)

func main() {
	c := colt.Colt{Stdout: os.Stdout, Stderr: os.Stderr}

	err := parseArgs(&c, os.Args[1:])
	if err != nil {
		die2(err)
	}

	scanner := bufio.NewScanner(os.Stdin)

	for scanner.Scan() {
		err := c.ProcessLine(scanner.Bytes())
		if _, ok := err.(colt.Warning); ok {
			warn(err)
			continue
		}
		dieIf(err)

		fmt.Fprintln(c.Stdout)
	}

	dieIf(scanner.Err())
}
