package bargle

import (
	"fmt"
)

type UnaryOption[T any] struct {
	optionDefaults
	Unmarshaler UnaryUnmarshaler[T]
	Longs       []string
	Shorts      []rune
	Required    bool
	parsed      bool
}

func (me UnaryOption[T]) Value() T {
	return me.Unmarshaler.Value()
}

func (me *UnaryOption[T]) Init() error {
	return initNilUnmarshalerUsingReflect(&me.Unmarshaler, nil)
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

func (me *UnaryOption[T]) AfterParse(Context) error {
	me.parsed = true
	return nil
}

func (me *UnaryOption[T]) Satisfied() bool {
	return !me.Required || me.parsed
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
	u UnaryUnmarshaler[T]
}

func (me unaryMatchResult[T]) Parse(args Args) error {
	if args.Len() == 0 {
		return missingArgument
	}
	arg := args.Pop()
	err := me.u.UnaryUnmarshal(arg)
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
			return unaryMatchResult[T]{baseMatchResult{_args, me, args.Clone().Pop()}, me.Unmarshaler}
		}
	}
	for _, s := range me.Shorts {
		_args := args.Clone()
		gv := &ShortParser{Short: s, CanUnary: true}
		if gv.Match(_args) {
			return unaryMatchResult[T]{baseMatchResult{_args, me, args.Clone().Pop()}, me.Unmarshaler}
		}
	}
	return noMatch
}

func (me *UnaryOption[T]) Parse(args Args) error {
	return me.Unmarshaler.UnaryUnmarshal(args.Pop())
}
