package bargle

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/anacrolix/generics"
)

type LongParser struct {
	Long     string
	CanUnary bool
}

func (me LongParser) Parse(ctx Context) error {
	args := ctx.Args()
	next := args.Pop()
	if !strings.HasPrefix(next, "--") {
		return noMatch
	}
	before, after, found := strings.Cut(next[2:], "=")
	if me.CanUnary && found {
		args.Push(after)
	}
	if before != me.Long {
		return noMatch
	}
	return nil
}

type Subcommand struct {
	Commands map[string]ContextFunc
}

func (me Subcommand) Parse(ctx Context) error {
	cmd := ctx.Args().Pop()
	defer recoverType(func(err controlError) {
		panic(controlError{fmt.Errorf("%s: %w", cmd, err)})
	})
	f, ok := me.Commands[cmd]
	if !ok {
		return noMatch
	}
	f(ctx)
	return nil
}

type Choice[T any] struct {
	Choices map[string]T
}

func (me Choice[T]) Unmarshal(choice string, t *T) error {
	var ok bool
	*t, ok = me.Choices[choice]
	if !ok {
		return noMatch
	}
	return nil
}

func (me Choice[T]) Add(name string, value T) {
	generics.MakeMapIfNil(&me.Choices)
	me.Choices[name] = value
}

func (me Choice[T]) Get(key string) T {
	return me.Choices[key]
}

type targetElem interface {
	string | int16
}

type target[T targetElem] interface {
	*T | *[]T
}

func unmarshalTarget[E targetElem, T target[E]](target T, s string) error {
	switch interface{}(target).(type) {
	case *string:
	case *[]string:
	}
	return nil
}

func (me *Positional[T]) Parse(ctx Context) error {
	return me.u.Unmarshal(ctx.Args(), &me.Value)
}

type unaryOption[T any] struct {
	Value       T
	Unmarshaler UnaryUnmarshaler[T]
	Longs       []string
	Shorts      []rune
}

func UnaryOption[T any](target UnaryUnmarshaler[T]) *unaryOption[T] {
	return &unaryOption[T]{Unmarshaler: target}
}

func (me *unaryOption[T]) AddLong(long string) *unaryOption[T] {
	me.Longs = append(me.Longs, long)
	return me
}

func (me *unaryOption[T]) AddShort(short rune) *unaryOption[T] {
	me.Shorts = append(me.Shorts, short)
	return me
}

func (me *unaryOption[T]) Parse(ctx Context) error {
	if !me.matchSwitch(ctx) {
		return noMatch
	}
	return me.Unmarshaler.Unmarshal(ctx.Args().Pop(), &me.Value)
}

func (me unaryOption[T]) matchSwitch(ctx Context) bool {
	for _, l := range me.Longs {
		if ctx.Match(LongParser{Long: l, CanUnary: true}) {
			return true
		}
	}
	// TODO: Short parsing
	return false
}

type UnaryUnmarshaler[T any] interface {
	Unmarshal(s string, t *T) error
}

type BuiltinUnaryUnmarshaler[T interface {
	string | int16
}] struct {
	Value *T
}

func (me BuiltinUnaryUnmarshaler[T]) Unmarshal(s string) error {
	switch p := any(me.Value).(type) {
	case *string:
		*p = s
	case *int16:
		i64, err := strconv.ParseInt(s, 0, 16)
		if err != nil {
			return err
		}
		*p = int16(i64)
	}
	panic("unreachable")
}

type Unmarshaler[T any] interface {
	Unmarshal(args Args, t *T) error
}

type Valuer[T any] interface {
	Value() *T
}

type UnmarshalValuer[T any] interface {
	Unmarshaler[T]
	Valuer[T]
}

type Positional[T any] struct {
	Value T
	u     Unmarshaler[T]
}

//func (me *Positional[T, V]) Value() *V {
//	return me.value.Value()
//}

type Slice[T any] struct {
	U Unmarshaler[T]
}

func NewSlice[T any](u Unmarshaler[T]) Slice[T] {
	return Slice[T]{U: u}
}

func (s Slice[T]) Unmarshal(args Args, slice *[]T) error {
	var t T
	err := s.U.Unmarshal(args, &t)
	if err != nil {
		return err
	}
	*slice = append(*slice, t)
	return nil
}

type String struct{}

func (String) Unmarshal(args Args, s *string) error {
	*s = args.Pop()
	return nil
}

func NewPositional[T any](u Unmarshaler[T]) *Positional[T] {
	return &Positional[T]{u: u}
}

func NewString() *String {
	return &String{}
}

func NewChoice[T any](choices map[string]T) *Choice[T] {
	return &Choice[T]{Choices: choices}
}

type Choices[T any] map[string]T

type Option[T any] struct {
	u Unmarshaler[T]
}

func (o Option[T]) Unmarshal(args Args, t *generics.Option[T]) error {
	err := o.u.Unmarshal(args, &t.Value)
	if err != nil {
		return err
	}
	t.Ok = true
	return nil
}

func NewOption[T any](u Unmarshaler[T]) *Option[T] {
	return &Option[T]{u: u}
}
