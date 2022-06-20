package bargle

import (
	"github.com/anacrolix/generics"
)

type Option[T any] struct {
	u UnaryUnmarshaler[*T]
}

func (o Option[T]) UnaryUnmarshal(s string, t *generics.Option[T]) error {
	err := o.u.UnaryUnmarshal(s, &t.Value)
	if err != nil {
		return err
	}
	t.Ok = true
	return nil
}

func (o Option[T]) TargetHelp() string {
	return o.u.TargetHelp()
}

//func (o Option[T]) Unmarshal(args Args, t *generics.Option[T]) error {
//	err := o.u.UnaryUnmarshal(args.Pop(), &t.Value)
//	if err != nil {
//		return err
//	}
//	t.Ok = true
//	return nil
//}

func NewOption[T any](u UnaryUnmarshaler[*T]) *Option[T] {
	return &Option[T]{u: u}
}
