package bargle

import (
	"fmt"
)

type Positional[T any] struct {
	posDefaults
	U              UnaryUnmarshaler[T]
	Name           string
	Desc           string
	ok             bool
	AfterParseFunc AfterParseParamFunc
}

func (me Positional[T]) Value() T {
	return me.U.Value()
}

func (me *Positional[T]) Init() error {
	initNilUnmarshalerUsingReflect(&me.U)
	return nil
}

func (me *Positional[T]) Parse(args Args) error {
	return me.U.UnaryUnmarshal(args.Pop())
}

func (me *Positional[T]) Match(args Args) MatchResult {
	if me.ok {
		return noMatch
	}
	mr := unaryMatchResult[T]{
		u: me.U,
	}
	mr.args = args.Clone()
	mr.param = me
	mr.match = args.Pop()
	return mr
}

func (me *Positional[T]) AfterParse(ctx Context) error {
	me.ok = true
	if me.AfterParseFunc != nil {
		return me.AfterParseFunc(ctx)
	}
	return nil
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
