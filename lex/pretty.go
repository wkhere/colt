package lex

import "fmt"

func (e *Error) Error() string {
	return fmt.Sprintf("%s at [%d:%d]", e.Text, e.Start, e.End)
}
