package bargle

import (
	"fmt"
)

type Parser interface {
	Parse(ctx Context) error
}

type Matcher interface {
	// Should this allow a nil return for no match?
	Match(args Args) MatchResult
}

func SetUnaryDefault(p Param, s string) {
	err := p.Parse(NewArgs([]string{s}))
	if err != nil {
		panic(fmt.Errorf("setting unary default: %w", err))
	}
}
