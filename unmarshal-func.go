package bargle

type unmarshalFunc func(ctx UnmarshalContext) error

func (me unmarshalFunc) Unmarshal(ctx UnmarshalContext) error {
	return me(ctx)
}
