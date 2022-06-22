package bargle

import (
	"bytes"
	"errors"
	"fmt"

	"github.com/anacrolix/generics"
)

type unhandledErr struct {
	arg string
}

func (me unhandledErr) Error() string {
	return fmt.Sprintf("unhandled argument: %q", me.arg)
}

type parseFailure struct{}

func (p parseFailure) Error() string {
	return "parse failure"
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

type tried struct{}

type parseError struct {
	inner error
	arg   generics.Option[string]
	param Param
}

func (me parseError) Unwrap() error {
	return me.inner
}

func (me parseError) Error() string {
	if me.arg.Ok {
		return fmt.Sprintf("parsing %v from %q: %v", me.param, me.arg.Value, me.inner)
	} else {
		return fmt.Sprintf("parsing %v: %v", me.param, me.inner)
	}
}

var missingArgument = errors.New("missing argument")

type exitCodeErrorWrapper struct {
	error
	exitCode int
}

func (me exitCodeErrorWrapper) ExitCode() int {
	return me.exitCode
}

func withExitCode(exitCode int, err error) error {
	return exitCodeErrorWrapper{
		error:    err,
		exitCode: exitCode,
	}
}

func (me exitCodeErrorWrapper) Unwrap() error {
	return me.error
}

type paramError struct {
	msg string
}

func (me paramError) Error() string {
	return me.msg
}

type unsatisfiedParam struct {
	p Param
}

func (me unsatisfiedParam) Error() string {
	var buf bytes.Buffer
	hw := HelpWriter{w: &buf}.Indented()
	me.p.Help().Write(hw)
	return fmt.Sprintf("unsatisfied param:\n%s", buf.Bytes())
}
