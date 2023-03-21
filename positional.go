package bargle

import (
	"fmt"
	"strings"

	g "github.com/anacrolix/generics"
)

func Positional(u Unmarshaler) positional {
	return positional{u}
}

type positional struct {
	u Unmarshaler
}

func (me positional) ArgInfo() ArgInfo {
	return ArgInfo{
		ArgType:       ArgTypePos,
		MatchingForms: g.Singleton(fmt.Sprintf("%T", me.u)),
	}
}

func (me positional) Parse(ctx ParseContext) bool {
	if ctx.NumArgs() < 1 {
		return false
	}
	// I'm not sure where to put this. It could go in the Parser arg parsing wrappers, in
	// positional.Parse, or maybe in ParseContext or UnmarshalContext to filter out individual args.
	if !ctx.PositionalOnly() && strings.HasPrefix(ctx.PeekArgs()[0], "-") {
		return false
	}
	return ctx.Unmarshal(me.u)
}
