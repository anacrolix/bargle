package bargle

type intUnmarshaler[T interface {
	int | int32
}, U interface {
	int64 | uint64
}] struct {
	t    *T
	f    func(string, int, int) (U, error)
	bits int
}

func (me intUnmarshaler[T, U]) Unmarshal(ctx UnmarshalContext) error {
	arg, err := ctx.Pop()
	if err != nil {
		return err
	}
	u, err := me.f(arg, 0, me.bits)
	if err != nil {
		return err
	}
	*me.t = T(u)
	return nil
}
