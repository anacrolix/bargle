package bargle

import (
	"errors"
	"log"
	"os"
)

type Main struct {
	OnError  func(err error)
	deferred []func()
	Command
	NoDefaultHelpSubcommand bool
}

func (me *Main) Defer(f func()) {
	me.deferred = append(me.deferred, f)
}

func (me *Main) Run() {
	ctx := NewContext(os.Args[1:])
	if me.Command.HasSubcommands() && !me.NoDefaultHelpSubcommand {
		HelpCommand{}.AddToCommand(&me.Command)
	}
	err := ctx.Run(me.Command)
	if err != nil {
		if me.OnError != nil {
			me.OnError(err)
		} else {
			log.Printf("error running main: %v", err)
		}
	}
	var (
		exitCoder ExitCoder
		exitCode  int
	)
	if errors.As(err, &exitCoder) {
		exitCode = exitCoder.ExitCode()
	} else if err != nil {
		exitCode = 1
	}
	os.Exit(exitCode)
}
