
[![Build Status](https://github.com/wkhere/colt/actions/workflows/go.yml/badge.svg?branch=master)](https://github.com/wkhere/colt/actions/workflows/go.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/wkhere/colt)](https://goreportcard.com/report/github.com/wkhere/colt)
[![Go Reference](https://pkg.go.dev/badge/github.com/wkhere/colt.svg)](https://pkg.go.dev/github.com/wkhere/colt)

## colt

![](https://small.shootingsportsmedia.com/52544.jpg "BANG!")

`colt` is COLumn Transformer.

`colt` can be used as a command-line tool or as a library.
Below is the description of the tool. For the library,
see [the reference](https://pkg.go.dev/github.com/wkhere/colt).

`colt +N command` transforms the Nth column of input with given command,
leaving other columns as they are.

`colt -N command` does the same for the Nth column counting from the end.

`colt command` transforms the last column (`-1` is the default).

`colt -s':' command` specifies the column separator
(the default is `;`, it must be 1 character).

`colt -u command` unquotes the content of a column before transformation.

`command` should accept input data as an argument and print the output.
