package args

import "strconv"

type boolUnmarshaler struct {
	b *bool
}

// TODO: Bools/flags are special, this should probably not take an arg, and an extra interface exist
// for inline values only.
func (b boolUnmarshaler) Unmarshal(ctx UnmarshalContext) (err error) {
	if !ctx.HaveExplicitValue() {
		*b.b = true
		return
	}
	arg, err := ctx.Pop()
	if err != nil {
		return err
	}
	*b.b, err = strconv.ParseBool(arg)
	return
}
