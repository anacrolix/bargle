package bargle

import (
	"fmt"
	"net/url"

	g "github.com/anacrolix/generics"
)

func UnaryUnmarshalFunc[T any](t *T, f func(string) (T, error)) unaryUnmarshalFunc[T] {
	return unaryUnmarshalFunc[T]{
		t: t,
		f: f,
	}
}

var _ interface {
	Unmarshaler
	UnmarshalerValuer
} = unaryUnmarshalFunc[*url.URL]{}

type unaryUnmarshalFunc[T any] struct {
	t *T
	f func(string) (T, error)
}

func (me unaryUnmarshalFunc[T]) ArgTypes() []string {
	var t T
	return g.Singleton(fmt.Sprintf("%T", t))
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

func (me unaryUnmarshalFunc[T]) Value() any {
	return *me.t
}
