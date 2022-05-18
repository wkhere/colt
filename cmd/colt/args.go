package main

import (
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
		if len(arg) < 1 {
			continue
		}
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
