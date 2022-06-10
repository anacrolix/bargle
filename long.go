package bargle

import (
	"strings"
)

type LongParser struct {
	Long     string
	CanUnary bool
}

func (me LongParser) Parse(ctx Context) error {
	args := ctx.Args()
	next := args.Pop()
	if !strings.HasPrefix(next, "--") {
		return noMatch
	}
	before, after, found := strings.Cut(next[2:], "=")
	if me.CanUnary && found {
		args.Push(after)
	}
	if before != me.Long {
		return noMatch
	}
	return nil
}
