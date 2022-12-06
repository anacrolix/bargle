package bargle

type ParseContext interface {
	Pop() (string, bool)
	Unmarshal(Unmarshaler) bool
	UnmarshalArg(u Unmarshaler, arg string) bool
}

type parseContext struct {
	args []string
	err  error
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
