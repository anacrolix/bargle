package bargle

import (
	"github.com/anacrolix/generics"
)

type Option[T any] struct {
	value generics.Option[T]
	u     UnaryUnmarshaler[T]
}

func (o Option[T]) UnaryUnmarshal(s string) error {
	err := o.u.UnaryUnmarshal(s)
	if err != nil {
		return err
	}
	o.value = generics.Some(o.u.Value())
	return nil
}

func (o Option[T]) TargetHelp() string {
	return o.u.TargetHelp()
}
