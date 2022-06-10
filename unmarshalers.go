package bargle

import (
	"fmt"
	"strconv"
)

type UnaryUnmarshaler[T any] interface {
	UnaryUnmarshal(s string, t *T) error
	Help(ph *ParamHelp)
}

type basicBuiltinUnaryUnmarshalTarget interface {
	string | int16 | uint16
}

type builtinUnaryUnmarshalTarget interface {
	//basicBuiltinUnaryUnmarshalTarget | []basicBuiltinUnaryUnmarshalTarget
}

type BuiltinUnaryUnmarshaler[T builtinUnaryUnmarshalTarget] struct{}

func (me BuiltinUnaryUnmarshaler[T]) UnaryUnmarshal(s string, t *T) error {
	switch p := any(t).(type) {
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

func (me BuiltinUnaryUnmarshaler[T]) Help(ph *ParamHelp) {
	var t T
	ph.Values = fmt.Sprintf("(%T)", t)
}

type Unmarshaler[T any] interface {
	Unmarshal(args Args, t *T) error
}

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

var _ UnaryUnmarshaler[string] = String{}

func (s2 String) Help(ph *ParamHelp) {
	//TODO implement me
	panic("implement me")
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
