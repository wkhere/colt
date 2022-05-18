package colt

import "fmt"

func errorf(format string, a ...interface{}) error {
	return fmt.Errorf(format, a...)
}

type Warning struct {
	error
}
