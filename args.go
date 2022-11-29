package args

import "os"

type Arg interface {
	Parse(ctx ParseContext) bool
}

type Input struct {
	args []string
}

func NewParser() *Parser {
	return &Parser{
		args: os.Args[1:],
	}
}
