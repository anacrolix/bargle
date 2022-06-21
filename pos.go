package bargle

import (
	"fmt"
)

type Positional[T any] struct {
	posDefaults
	Value          UnaryUnmarshaler[T]
	Name           string
	Desc           string
	ok             bool
	AfterParseFunc AfterParseParamFunc
}

func (me *Positional[T]) Init() error {
	return initNilUnmarshalerUsingReflect(&me.Value, nil)
}

func (me *Positional[T]) Parse(args Args) error {
	return me.Value.UnaryUnmarshal(args.Pop())
}

func (me *Positional[T]) Match(args Args) MatchResult {
	if !me.Value.Matching() {
		return noMatch
	}
	mr := unaryMatchResult[T]{
		u: me.Value,
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
		Values:      me.Value.TargetHelp(),
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
//	return doUnaryUnmarshal(ctx.Args().Pop(), &me.Value, me.Value)
//}

func NewPositional[T any](u UnaryUnmarshaler[T]) *Positional[T] {
	return &Positional[T]{Value: u, Name: "arg"}
}
