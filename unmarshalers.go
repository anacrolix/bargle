package bargle

import (
	"strconv"
)

type UnaryUnmarshaler[T any] interface {
	Unmarshal(s string, t *T) error
	Help(ph *ParamHelp)
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

func NewString() *String {
	return &String{}
}
