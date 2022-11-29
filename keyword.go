package args

func Keyword(arg string) keyword {
	return keyword{arg}
}

type keyword struct {
	arg string
}

func (me keyword) Parse(ctx ParseContext) bool {
	arg, ok := ctx.Pop()
	return ok && arg == me.arg
}
