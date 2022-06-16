package bargle

import (
	"fmt"
)

type Subcommand struct {
	Commands map[string]ContextFunc
}

func (me Subcommand) Help(f HelpFormatter) {
	for name := range me.Commands {
		f.AddCommand(name, "use bargle.Command")
	}
}

func (me Subcommand) Parse(ctx Context) error {
	cmd := ctx.Args().Pop()
	defer recoverType(func(err controlError) {
		panic(controlError{fmt.Errorf("%s: %w", cmd, err)})
	})
	f, ok := me.Commands[cmd]
	if !ok {
		return noMatch
	}
	f(ctx)
	return nil
}
