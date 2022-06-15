package bargle

import (
	"fmt"
)

type UnaryOption[T any] struct {
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

func (me *UnaryOption[T]) Help(f HelpFormatter) {
	ph := ParamHelp{
		Forms: me.switchForms(),
	}
	ph.Values = me.Unmarshaler.TargetHelp()
	f.AddOption(ph)
}

func (me *UnaryOption[T]) AddLong(long string) *UnaryOption[T] {
	me.Longs = append(me.Longs, long)
	return me
}

func (me *UnaryOption[T]) AddShort(short rune) *UnaryOption[T] {
	me.Shorts = append(me.Shorts, short)
	return me
}

func (me *UnaryOption[T]) Parse(ctx Context) error {
	if !me.matchSwitch(ctx) {
		return noMatch
	}
	arg := ctx.Args().Pop()
	if me.Unmarshaler == nil {
		return fmt.Errorf("unary option %s has no unmarshaler", me.switchForms())
	}
	err := me.Unmarshaler.UnaryUnmarshal(arg, &me.Value)
	if err != nil {
		err = fmt.Errorf("unmarshalling %q: %w", arg, err)
	}
	return err
}

func (me UnaryOption[T]) matchSwitch(ctx Context) bool {
	for _, l := range me.Longs {
		if ctx.Match(&LongParser{Long: l, CanUnary: true}) {
			return true
		}
	}
	for _, s := range me.Shorts {
		if ctx.Match(&ShortParser{Short: s, CanUnary: true}) {
			return true
		}
	}
	return false
}
