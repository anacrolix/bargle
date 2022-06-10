package bargle

import (
	"errors"
	"log"
	"os"
)

type Main struct {
	OnError  func(err error)
	deferred []func()
}

func (me *Main) Defer(f func()) {
	me.deferred = append(me.deferred, f)
}

func (me *Main) Run(f ContextFunc) {
	ctx := NewContext(os.Args[1:])
	err := ctx.Run(f)
	if err == nil {
		for _, f := range ctx.actions {
			err = f()
			if err != nil {
				break
			}
		}
	}
	if err != nil {
		if me.OnError != nil {
			me.OnError(err)
		} else {
			log.Printf("error running main: %v", err)
		}
	}
	for i := range me.deferred {
		me.deferred[len(me.deferred)-1-i]()
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
