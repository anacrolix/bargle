package bargle

import (
	"github.com/anacrolix/generics"
)

func NewOption[T any](target *generics.Option[T], u UnaryUnmarshaler) Option[T] {
	ret := Option[T]{
		TargetOk:    &target.Ok,
		TargetValue: &target.Value,
	}
	initNilUnmarshalerUsingReflect(&ret.ValueUnmarshaler, &target.Value)
	return ret
}

type Option[T any] struct {
	TargetOk         *bool
	TargetValue      *T
	ValueUnmarshaler UnaryUnmarshaler
}

func (o Option[T]) UnaryUnmarshal(s string) error {
	err := o.ValueUnmarshaler.UnaryUnmarshal(s)
	if err != nil {
		return err
	}
	*o.TargetOk = true
	return nil
}

func (o Option[T]) TargetHelp() string {
	return o.ValueUnmarshaler.TargetHelp()
}

func (me Option[T]) Matching() bool {
	return !*me.TargetOk
}

func (me Option[T]) Value() generics.Option[T] {
	return generics.Option[T]{
		Ok:    *me.TargetOk,
		Value: *me.TargetValue,
	}
}
