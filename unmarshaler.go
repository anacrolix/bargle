package args

type Unmarshaler interface {
	Unmarshal(ctx UnmarshalContext) error
}

func String(s *string) stringUnmarshaler {
	return stringUnmarshaler{s}
}

type stringUnmarshaler struct {
	s *string
}

func (me stringUnmarshaler) Unmarshal(ctx UnmarshalContext) (err error) {
	*me.s, err = ctx.Pop()
	return err
}

func UnmarshalFunc[T any](t *T, f func(string) (T, error)) unmarshalFunc[T] {
	return unmarshalFunc[T]{
		t: t,
		f: f,
	}
}

type unmarshalFunc[T any] struct {
	t *T
	f func(string) (T, error)
}

func (me unmarshalFunc[T]) Unmarshal(ctx UnmarshalContext) (err error) {
	arg, err := ctx.Pop()
	if err != nil {
		return
	}
	t, err := me.f(arg)
	if err != nil {
		return
	}
	*me.t = t
	return
}
