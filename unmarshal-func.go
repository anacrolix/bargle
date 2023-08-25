package bargle

type unmarshalFunc struct {
	f        func(ctx UnmarshalContext) error
	argTypes []string
}

func (me unmarshalFunc) ArgTypes() []string {
	return me.argTypes
}

func (me unmarshalFunc) Unmarshal(ctx UnmarshalContext) error {
	return me.f(ctx)
}
