package bargle

type UnmarshalContext interface {
	Pop() (string, error)
	// Value was specifically linked to this context, such as --like=this.
	HaveExplicitValue() bool
}

type unmarshalContext struct {
	args          []string
	explicitValue bool
}

func (me *unmarshalContext) HaveExplicitValue() bool {
	return me.explicitValue
}

func (me *unmarshalContext) Pop() (arg string, err error) {
	if len(me.args) == 0 {
		err = ErrNoArgs
		return
	}
	arg = me.args[0]
	me.args = me.args[1:]
	return
}
