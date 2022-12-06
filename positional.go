package bargle

import "github.com/anacrolix/generics"

func Positional(u Unmarshaler) positional {
	return positional{u}
}

type positional struct {
	u Unmarshaler
}

func (me positional) ArgInfo() ArgInfo {
	return ArgInfo{
		ArgType:       ArgTypePos,
		MatchingForms: generics.Singleton("<value>"),
	}
}

func (me positional) Parse(ctx ParseContext) bool {
	return ctx.Unmarshal(me.u)
}
