package bargle

import (
	"fmt"
)

type Positional[T any] struct {
	posDefaults
	Value T
	U     UnaryUnmarshaler[T]
	Name  string
	Desc  string
	ok    bool
}

func (me *Positional[T]) Match(args Args) MatchResult {
	mr := unaryMatchResult[T]{
		u:      me.U,
		target: &me.Value,
	}
	mr.args = args
	mr.param = me
	return mr
}

func (me *Positional[T]) Satisfied() bool {
	return me.ok
}

func (me *Positional[T]) Help() ParamHelp {
	return ParamHelp{
		Forms:       []string{fmt.Sprintf("<%v>", me.Name)},
		Description: me.Desc,
	}
}

//func (me *Positional[T]) Parse(ctx Context) error {
//	if ctx.Args().Len() == 0 {
//		return missingArgument
//	}
//	if !ctx.MatchPos() {
//		return noMatch
//	}
//	return doUnaryUnmarshal(ctx.Args().Pop(), &me.Value, me.U)
//}

func NewPositional[T any](u UnaryUnmarshaler[T]) *Positional[T] {
	return &Positional[T]{U: u, Name: "arg"}
}
