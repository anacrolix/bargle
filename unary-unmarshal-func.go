package bargle

func UnaryUnmarshalFunc[T any](t *T, f func(string) (T, error)) unaryUnmarshalFunc[T] {
	return unaryUnmarshalFunc[T]{
		t: t,
		f: f,
	}
}

type unaryUnmarshalFunc[T any] struct {
	t *T
	f func(string) (T, error)
}

func (me unaryUnmarshalFunc[T]) Unmarshal(ctx UnmarshalContext) (err error) {
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
