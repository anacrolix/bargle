package bargle

import (
	"errors"
	"fmt"
)

var noMatch = errors.New("no match")

type unhandledErr struct {
	arg string
}

func (me unhandledErr) Error() string {
	return fmt.Sprintf("unhandled argument: %q", me.arg)
}

type ExitCoder interface {
	ExitCode() int
}

type userError error

type controlError struct {
	error
}

func (me controlError) Unwrap() error {
	return me.error
}

func (controlError) ExitCode() int {
	return 2
}

type success struct{}

func (success) Error() string {
	return "success"
}
