package bargle

import (
	"strings"
)

// Creates a positional argument. Positional arguments are parsed based on their relative position
// in the argument stream.
func Positional(metavar string, u Unmarshaler) Arg {
	// Yo... Rust is way better.
	return &positional{u: u, metavar: metavar}
}

type positional struct {
	u       Unmarshaler
	metavar string
}

func (me positional) Metavar() string {
	return me.metavar
}

var _ Metavar = (*positional)(nil)

func (me *positional) ArgInfo() ArgInfo {
	return ArgInfo{
		ArgType:       ArgTypePos,
		MatchingForms: me.u.ArgTypes(),
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
