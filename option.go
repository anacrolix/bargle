package bargle

import (
	"github.com/anacrolix/generics"
)

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
