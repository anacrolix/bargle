package bargle

import "github.com/anacrolix/generics"

func Keyword(arg string) keyword {
	return keyword{arg}
}

type keyword struct {
	arg string
}

func (me keyword) ArgInfo() ArgInfo {
	return ArgInfo{
		ArgType:       ArgTypePos,
		MatchingForms: generics.Singleton(me.arg),
	}
}

func (me keyword) Parse(ctx ParseContext) bool {
	arg, ok := ctx.Pop()
	return ok && arg == me.arg
}
