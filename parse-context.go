package bargle

type ParseContext interface {
	NumArgs() int
	Pop() (string, bool)
	Unmarshal(Unmarshaler) bool
	UnmarshalArg(u Unmarshaler, arg string) bool
	// Whether "--" has been given.
	PositionalOnly() bool
	PeekArgs() []string
}

type parseContext struct {
	args []string
	err  error
	// "--" or something having the effect of not checking for switches has been set.
	posOnly bool
}

func (me *parseContext) PeekArgs() []string {
	return me.args
}

func (me *parseContext) PositionalOnly() bool {
	return me.posOnly
}

func (me *parseContext) NumArgs() int {
	return len(me.args)
}

func (me *parseContext) UnmarshalArg(u Unmarshaler, arg string) bool {
	if me.err != nil {
		return false
	}
	ctx := unmarshalContext{args: []string{arg}, explicitValue: true}
	me.err = u.Unmarshal(&ctx)
	return me.err == nil
}

func (me *parseContext) Pop() (arg string, ok bool) {
	if len(me.args) == 0 {
		return
	}
	arg = me.args[0]
	me.args = me.args[1:]
	ok = true
	return
}

func (me *parseContext) SetError(err error) {
	me.err = err
}

func (me *parseContext) Unmarshal(u Unmarshaler) bool {
	if me.err != nil {
		return false
	}
	ctx := unmarshalContext{
		args: me.args,
	}
	me.err = u.Unmarshal(&ctx)
	me.args = ctx.args
	return me.err == nil
}
