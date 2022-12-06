package bargle

import "errors"

var (
	ErrNoArgs            = errors.New("no arguments remaining")
	ErrExpectedArguments = errors.New("expected arguments")
)
