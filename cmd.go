package bargle

import (
	"fmt"
)

type Command struct {
	Name string
	Func ContextFunc
	Desc string
}

var _ interface {
	ParamParser
} = Command{}

func (me Command) Parse(ctx Context) error {
	cmd := ctx.Args().Pop()
	defer recoverType(func(err controlError) {
		panic(controlError{fmt.Errorf("%s: %w", cmd, err)})
	})
	if cmd != me.Name {
		return noMatch
	}
	// This doesn't start a new try scope. That's probably wrong.
	me.Func(ctx)
	return nil
}

func (me Command) Help(f HelpFormatter) {
	f.AddCommand(me.Name)
}
