package main

import (
	"fmt"
	"os"
	"regexp"
	"strconv"
)

var selectionR = regexp.MustCompile(`^[-+]?\d+$`)

func (p *columnProc) parseArgs(args []string) {
	p.selection = -1
	p.separator = ';'

loop:
	for i, arg := range args {
		switch arg[0] {
		case '-', '+':
			switch {
			case selectionR.MatchString(arg):
				p.selection, _ = strconv.Atoi(arg)

			case arg[:2] == "-d":
				if len(arg[2:]) != 1 {
					usage()
				}
				p.separator = rune(arg[2])

			default:
				usage()
			}
		default:
			p.command = args[i:]
			break loop
		}
	}
	if len(p.command) == 0 {
		usage()
	}
}

func usage() {
	fmt.Fprint(os.Stderr, `colt - copy input transforming chosen column with a given command.

usage: colt [+N|-N] [-dC] command ...
where:
    N - integer denoting column number, starting from 1;
        when negative, counted from end; default -1 (last column)
    C - 1-character column delimiter
    command - a command, possible with args, for transforming the column
`)
	os.Exit(2)
}
