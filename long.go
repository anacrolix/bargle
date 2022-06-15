package bargle

import (
	"strings"
)

type LongParser struct {
	Long     string
	CanUnary bool
	gotValue bool
}

func (me LongParser) GotValue() bool {
	return me.gotValue
}

func (me *LongParser) Parse(ctx Context) error {
	args := ctx.Args()
	if args.Len() == 0 {
		return noMatch
	}
	next := args.Pop()
	if !strings.HasPrefix(next, "--") {
		return noMatch
	}
	before, after, found := strings.Cut(next[2:], "=")
	me.gotValue = false
	if me.CanUnary && found {
		me.gotValue = true
		args.Push(after)
	}
	if before != me.Long {
		return noMatch
	}
	return nil
}

func (me LongParser) Help(f *ParamHelp) {
	f.Forms = append(f.Forms, "--"+me.Long)
}
