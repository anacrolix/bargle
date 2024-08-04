package bargle

import (
	"fmt"
	"strconv"

	g "github.com/anacrolix/generics"
	"golang.org/x/exp/constraints"
)

type floatUnmarshaler[T interface {
	constraints.Float
}] struct {
	t    *T
	bits int
}

func (me floatUnmarshaler[T]) ArgTypes() []string {
	var t T
	return g.Singleton(fmt.Sprintf("%T", t))
}

func (me floatUnmarshaler[T]) Unmarshal(ctx UnmarshalContext) error {
	arg, err := ctx.Pop()
	if err != nil {
		return err
	}
	t, err := unmarshalFloat[T](arg, me.bits)
	if err != nil {
		return err
	}
	*me.t = t
	return nil
}

func unmarshalFloat[
	T constraints.Float,
](
	arg string,
	bitSize int,
) (
	T, error,
) {
	u, intErr := strconv.ParseFloat(arg, bitSize)
	return T(u), intErr
}
