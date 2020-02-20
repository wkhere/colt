package main

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

var selectionR = regexp.MustCompile(`^[-+]\d+$`)

func (p *columnProc) parseArgs(args []string) error {
	p.selection = -1
	p.separator = ';'
	p.quote = '"'

loop:
	for i, arg := range args {
		switch arg[0] {
		case '-', '+':
			switch {
			case selectionR.MatchString(arg):
				p.selection, _ = strconv.Atoi(arg)

			case strings.HasPrefix(arg, "-s"):
				if len(arg[2:]) != 1 {
					return usageErr
				}
				p.separator = rune(arg[2])

			case arg == "-u":
				p.unquote = true

			default:
				return usageErr
			}
		default:
			p.command = args[i:]
			break loop
		}
	}
	if len(p.command) == 0 {
		return usageErr
	}

	return nil
}

var usageErr = fmt.Errorf(`colt - copy input transforming chosen column with a given command.

usage: colt [+N|-N] [-sS] [-u] command ...
where:
    -N|+N - column number, starting from 1;
            when negative, counted from end; default -1 (last column)
    -sS   - 1-character column separator; default ';'
    -u    - unquote column content first
    command - a command, possible with args, for transforming the column`,
)
