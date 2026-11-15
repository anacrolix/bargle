package bargle

import (
	"strings"
)

// Creates a positional argument. Positional arguments are parsed based on their relative position
// in the argument stream. It matches at most once: see Positionals for one that repeats.
func Positional(metavar string, u Unmarshaler) *positional {
	// Yo... Rust is way better.
	return &positional{u: u, metavar: metavar}
}

// Creates a positional argument that matches as many times as it's offered. Pass an Unmarshaler
// that accumulates, such as AppendSlice.
func Positionals(metavar string, u Unmarshaler) *positional {
	return &positional{u: u, metavar: metavar, repeat: true}
}

type positional struct {
	u       Unmarshaler
	metavar string
	repeat  bool
	// Maybe this should be done with a ParseCount wrapper type?
	parseCount int
}

func (me positional) Metavar() string {
	return me.metavar
}

var _ interface {
	Arg
	Metavar
	ParseCounter
} = (*positional)(nil)

func (me *positional) ArgInfo() ArgInfo {
	return ArgInfo{
		ArgType:       ArgTypePos,
		MatchingForms: me.u.ArgTypes(),
	}
}

// The number of times the argument has matched.
func (me *positional) ParseCount() int {
	return me.parseCount
}

func (me *positional) Parse(ctx ParseContext) bool {
	if me.parseCount != 0 && !me.repeat {
		return false
	}
	if ctx.NumArgs() < 1 {
		return false
	}
	// I'm not sure where to put this. It could go in the Parser arg parsing wrappers, in
	// positional.Parse, or maybe in ParseContext or UnmarshalContext to filter out individual args.
	if !ctx.PositionalOnly() && strings.HasPrefix(ctx.PeekArgs()[0], "-") {
		return false
	}
	parsed := ctx.Unmarshal(me.u)
	if parsed {
		me.parseCount++
	}
	return parsed
}

//func (me *Positional[T]) Parse(ctx Context) error {
//	if ctx.Args().Len() == 0 {
//		return missingArgument
//	}
//	if !ctx.MatchPos() {
//		return noMatch
//	}
//	return doUnaryUnmarshal(ctx.Args().Pop(), &me.Value, me.Value)
//}

func NewPositional(u Unmarshaler) *positional {
	return &positional{u: u, metavar: "arg"}
}
