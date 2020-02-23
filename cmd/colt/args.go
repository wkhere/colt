package main

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/wkhere/colt"
)

var selectionR = regexp.MustCompile(`^[-+]\d+$`)

func parseArgs(c *colt.Colt, args []string) error {
	c.Selection = -1
	c.Separator = ';'
	c.Quote = '"'

loop:
	for i, arg := range args {
		switch arg[0] {
		case '-', '+':
			switch {
			case selectionR.MatchString(arg):
				c.Selection, _ = strconv.Atoi(arg)

			case strings.HasPrefix(arg, "-s"):
				if len(arg[2:]) != 1 {
					return usageErr
				}
				c.Separator = rune(arg[2])

			case arg == "-u":
				c.Unquote = true

			default:
				return usageErr
			}
		default:
			c.Command = args[i:]
			break loop
		}
	}
	if len(c.Command) == 0 {
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
