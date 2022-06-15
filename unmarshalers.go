package bargle

import (
	"fmt"
	"strconv"

	"golang.org/x/exp/constraints"
)

type UnaryUnmarshaler[T any] interface {
	UnaryUnmarshal(s string, t *T) error
	TargetHelp() string
}

type basicBuiltinUnaryUnmarshalTarget interface {
	~string | ~int16 | ~uint16
}

type builtinUnaryUnmarshalTarget interface {
	basicBuiltinUnaryUnmarshalTarget
}

type BuiltinUnaryUnmarshaler[T builtinUnaryUnmarshalTarget] struct{}

func parseIntType[R constraints.Integer, I constraints.Integer](s string, bits int, f func(string, int, int) (I, error)) (ret R, err error) {
	var i I
	i, err = f(s, 0, bits)
	ret = R(i)
	return
}

func (me BuiltinUnaryUnmarshaler[T]) UnaryUnmarshal(s string, t *T) (err error) {
	switch p := any(t).(type) {
	case *string:
		*p = s
		return nil
	case *int16:
		*p, err = parseIntType[int16](s, 16, strconv.ParseInt)
	case *uint16:
		*p, err = parseIntType[uint16](s, 16, strconv.ParseUint)
	default:
		panic(fmt.Sprintf("builtin unary unmarshaler: unsupported type %T", *t))
	}
	return
}

func (me BuiltinUnaryUnmarshaler[T]) TargetHelp() string {
	var t T
	return fmt.Sprintf("(%T)", t)
}

type Unmarshaler[T any] interface {
	Unmarshal(args Args, t *T) error
}

type Slice[T any] struct {
	Unmarshaler interface {
		UnaryUnmarshaler[T]
		TargetHelper
	}
}

func (s2 Slice[T]) TargetHelp() string {
	return s2.Unmarshaler.TargetHelp() + "..."
}

func NewSlice[T any](u interface {
	UnaryUnmarshaler[T]
	TargetHelper
}) Slice[T] {
	return Slice[T]{Unmarshaler: u}
}

func (sl Slice[T]) UnaryUnmarshal(s string, slice *[]T) error {
	var t T
	err := sl.Unmarshaler.UnaryUnmarshal(s, &t)
	if err != nil {
		return err
	}
	*slice = append(*slice, t)
	return nil
}

type String struct{}

var _ UnaryUnmarshaler[string] = String{}

func (s2 String) TargetHelp() string {
	return "(string)"
}

func (String) Unmarshal(args Args, s *string) error {
	*s = args.Pop()
	return nil
}

func (String) UnaryUnmarshal(arg string, s *string) error {
	*s = arg
	return nil
}

func NewString() *String {
	return &String{}
}
