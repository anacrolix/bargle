package bargle

import (
	g "github.com/anacrolix/generics"
)

func AppendSlice[T any](s *[]T, uc func(*T) Unmarshaler) Unmarshaler {
	var t T
	u := uc(&t)
	return unmarshalFunc{
		func(ctx UnmarshalContext) error {
			err := u.Unmarshal(ctx)
			if err == nil {
				*s = append(*s, t)
			}
			return err
		},
		g.SliceMap(u.ArgTypes(), func(f string) string { return f + "..." }),
	}
}
