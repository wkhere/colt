//go:build !windows

package main

import "fmt"

var usageErr = fmt.Errorf(`colt - copy input transforming chosen column with a given command.

usage: colt [+N|-N] [-sS] [-u] command ...
where:
    -N|+N - column number, starting from 1;
            when negative, counted from end; default -1 (last column)
    -sS   - 1-character column separator; default ';'
    -u    - unquote column content first
    command - a command, possible with args, for transforming the column`,
)
