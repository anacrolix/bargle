package bargle

import g "github.com/anacrolix/generics"

func OptionUnmarshaler[V any](o *g.Option[V], vu func(*V) Unmarshaler) Unmarshaler {
	u := vu(&o.Value)
	return unmarshalFunc(func(ctx UnmarshalContext) error {
		err := u.Unmarshal(ctx)
		if err == nil {
			o.Ok = true
		}
		return err
	})
}

func BuiltinOptionUnmarshaler[V BuiltinUnmarshalerType](o *g.Option[V]) Unmarshaler {
	return OptionUnmarshaler(o, func(v *V) Unmarshaler {
		return BuiltinUnmarshaler(v)
	})
}
