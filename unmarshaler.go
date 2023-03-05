package bargle

type Unmarshaler interface {
	Unmarshal(ctx UnmarshalContext) error
}

func String(s *string) Unmarshaler {
	return stringUnmarshaler{s}
}

type stringUnmarshaler struct {
	s *string
}

func (me stringUnmarshaler) Unmarshal(ctx UnmarshalContext) (err error) {
	*me.s, err = ctx.Pop()
	return err
}
