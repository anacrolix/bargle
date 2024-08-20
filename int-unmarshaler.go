package bargle

import (
	"errors"
	"fmt"
	"math"
	"strconv"

	g "github.com/anacrolix/generics"
	"golang.org/x/exp/constraints"
)

type intUnmarshaler[T interface {
	constraints.Integer
}, U interface {
	int64 | uint64
}] struct {
	t    *T
	f    func(string, int, int) (U, error)
	bits int
}

func (me intUnmarshaler[T, U]) ArgTypes() []string {
	var t T
	return g.Singleton(fmt.Sprintf("%T", t))
}

func (me intUnmarshaler[T, U]) Unmarshal(ctx UnmarshalContext) error {
	arg, err := ctx.Pop()
	if err != nil {
		return err
	}
	t, err := unmarshalInt[T, U](arg, me.f, me.bits)
	if err != nil {
		return err
	}
	*me.t = t
	return nil
}

func unmarshalInt[
	T constraints.Integer,
	U interface{ int64 | uint64 },
](
	arg string,
	intParser func(string, int, int) (U, error),
	intBits int,
) (
	T, error,
) {
parseInt:
	u, intErr := intParser(arg, 0, intBits)
	if intErr == nil {
		return T(u), nil
	}
	newArg, floatErr := removeExponent(arg)
	if floatErr != nil {
		return 0, errors.Join(intErr, floatErr)
	}
	if newArg != arg {
		arg = newArg
		goto parseInt
	}
	return 0, intErr
}

// Parses as a float and removes exponent, ensuring that the value is representable as an integer.
func removeExponent(arg string) (ret string, err error) {
	f64, err := strconv.ParseFloat(arg, 64)
	if err != nil {
		return
	}
	// Nudge the float, and see if the value moves by more than 1 integer. If so we are out of
	// precision and the integer we parsed may not be representable by a float.
	i, f := math.Modf(math.Nextafter(f64, 0))
	if f == 0 && math.Abs(f64-i) > 1 {
		err = errors.New("insufficient precision to represent as a float")
		return
	}
	ret = strconv.FormatFloat(f64, 'f', -1, 64)
	return
}

func (me intUnmarshaler[T, U]) Value() any {
	return *me.t
}

var _ UnmarshalerValuer = intUnmarshaler[int, int64]{}
