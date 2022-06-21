package bargle

import (
	"fmt"

	"github.com/anacrolix/generics"
)

type Option[T any] struct {
	value generics.Option[T]
	u     UnaryUnmarshaler[T]
}

func (o *Option[T]) UnaryUnmarshal(s string) error {
	err := initNilUnmarshalerUsingReflect(&o.u)
	if err != nil {
		return fmt.Errorf("initing inner value unmarshaler: %w", err)
	}
	err = o.u.UnaryUnmarshal(s)
	if err != nil {
		return err
	}
	o.value.Set(o.u.Value())
	return nil
}

func (o Option[T]) TargetHelp() string {
	return o.u.TargetHelp()
}

func (me Option[T]) Value() generics.Option[T] {
	return me.value
}
