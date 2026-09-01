package bargle

// Returns an Unmarshaler from a function, and the argument types it consumes for help output.
func UnmarshalFunc(f func(ctx UnmarshalContext) error, argTypes ...string) Unmarshaler {
	return unmarshalFunc{f, argTypes}
}

// Returns u with the argument types it reports for help output replaced. Useful when the Go type
// name an Unmarshaler derives isn't meaningful to a user, such as a duration or a byte count
// parsed from a string.
func WithArgTypes(u Unmarshaler, argTypes ...string) Unmarshaler {
	return unmarshalFunc{u.Unmarshal, argTypes}
}

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
