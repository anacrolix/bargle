package bargle

import (
	"errors"
)

func recoverType[T any](with func(err T)) {
	r := recover()
	if r == nil {
		return
	}
	t, ok := r.(T)
	if !ok {
		panic(r)
	}
	with(t)
}

func recoverErrorAs[T error](with func(t T)) {
	recoverType(func(err error) {
		var t T
		if errors.As(err, &t) {
			with(t)
		} else {
			panic(err)
		}
	})
}

func runDeferred(deferred []func()) {
	for i := range deferred {
		deferred[len(deferred)-1-i]()
	}
}
