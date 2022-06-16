package bargle

import (
	"fmt"
)

type Positional[T any] struct {
	Value T
	U     UnaryUnmarshaler[T]
	Name  string
	Desc  string
}

func (me *Positional[T]) Help(f HelpFormatter) {
	f.AddPositional(ParamHelp{
		Forms:       []string{fmt.Sprintf("<%v>", me.Name)},
		Description: me.Desc,
	})
}

func (me *Positional[T]) Parse(ctx Context) error {
	if ctx.Args().Len() == 0 {
		return missingArgument
	}
	if !ctx.MatchPos() {
		return noMatch
	}
	return doUnaryUnmarshal(ctx.Args().Pop(), &me.Value, me.U)
}

func NewPositional[T any](u UnaryUnmarshaler[T]) *Positional[T] {
	return &Positional[T]{U: u, Name: "arg"}
}
