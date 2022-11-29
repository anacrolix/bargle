package args

func Positional(u Unmarshaler) positional {
	return positional{u}
}

type positional struct {
	u Unmarshaler
}

func (me positional) Parse(ctx ParseContext) bool {
	return ctx.Unmarshal(me.u)
}
