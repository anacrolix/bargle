package bargle

import (
	"fmt"
	g "github.com/anacrolix/generics"
)

func OptionUnmarshaler[V any](o *g.Option[V], vu func(*V) Unmarshaler) Unmarshaler {
	u := vu(&o.Value)
	var v V
	return unmarshalFunc{
		func(ctx UnmarshalContext) error {
			err := u.Unmarshal(ctx)
			if err == nil {
				o.Ok = true
			}
			return err
		},
		g.Singleton(fmt.Sprintf("%T", v)),
	}
}

func BuiltinOptionUnmarshaler[V BuiltinUnmarshalerType](o *g.Option[V]) Unmarshaler {
	return OptionUnmarshaler(o, func(v *V) Unmarshaler {
		return BuiltinUnmarshaler(v)
	})
}
