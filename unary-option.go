package bargle

import (
	"errors"
	"fmt"
)

type UnaryOption[T any] struct {
	optionDefaults
	Value       T
	Unmarshaler UnaryUnmarshaler[T]
	Longs       []string
	Shorts      []rune
}

func NewUnaryOption[T builtinUnaryUnmarshalTarget](target *T) *UnaryOption[T] {
	return &UnaryOption[T]{Unmarshaler: BuiltinUnaryUnmarshaler[T]{}}
}

func (me *UnaryOption[T]) switchForms() (ret []string) {
	for _, l := range me.Longs {
		ret = append(ret, "--"+l)
	}
	for _, s := range me.Shorts {
		ret = append(ret, "-"+string(s))
	}
	return
}

func (me *UnaryOption[T]) Help() ParamHelp {
	return ParamHelp{
		Forms:  me.switchForms(),
		Values: me.Unmarshaler.TargetHelp(),
	}
}

func (me *UnaryOption[T]) AddLong(long string) *UnaryOption[T] {
	me.Longs = append(me.Longs, long)
	return me
}

func (me *UnaryOption[T]) AddShort(short rune) *UnaryOption[T] {
	me.Shorts = append(me.Shorts, short)
	return me
}

func (me *UnaryOption[T]) Match(args Args) MatchResult {
	return me.matchSwitch(args)
}

type unaryMatchResult[T any] struct {
	baseMatchResult
	u      UnaryUnmarshaler[T]
	target *T
}

func (me unaryMatchResult[T]) Parse(ctx Context) error {
	args := ctx.Args()
	if args.Len() == 0 {
		return missingArgument
	}
	arg := ctx.Args().Pop()
	if me.u == nil {
		return errors.New("no unmarshaler set")
	}
	err := me.u.UnaryUnmarshal(arg, me.target)
	if err != nil {
		err = fmt.Errorf("unmarshalling %q: %w", arg, err)
	}
	return err
}

func (me *UnaryOption[T]) matchSwitch(args Args) MatchResult {
	for _, l := range me.Longs {
		_args := args.Clone()
		gv := &LongParser{Long: l, CanUnary: true}
		if gv.Match(_args) {
			return unaryMatchResult[T]{baseMatchResult{_args, me, args.Clone().Pop()}, me.Unmarshaler, &me.Value}
		}
	}
	for _, s := range me.Shorts {
		_args := args.Clone()
		gv := &ShortParser{Short: s, CanUnary: true}
		if gv.Match(_args) {
			return unaryMatchResult[T]{baseMatchResult{_args, me, args.Clone().Pop()}, me.Unmarshaler, &me.Value}
		}
	}
	return noMatch
}
