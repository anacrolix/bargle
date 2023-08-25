package bargle

import (
	g "github.com/anacrolix/generics"
)

type Unmarshaler interface {
	Unmarshal(ctx UnmarshalContext) error
	ArgTypes() []string
}

func String(s *string) Unmarshaler {
	return stringUnmarshaler{s}
}

type stringUnmarshaler struct {
	s *string
}

func (me stringUnmarshaler) ArgTypes() []string {
	return g.Singleton("string")
}

func (me stringUnmarshaler) Unmarshal(ctx UnmarshalContext) (err error) {
	*me.s, err = ctx.Pop()
	return err
}
