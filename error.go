package colt

import "fmt"

func errorf(format string, a ...interface{}) error {
	return fmt.Errorf(format, a...)
}

type Warning struct {
	error
}

func (w Warning) Error() string {
	return fmt.Sprintf("WARN %v\n", w.error)
}
