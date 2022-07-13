package bargle

import (
	"fmt"
)

// A Command represents a collection of parameters and behaviour that are parsed together. Some parameters might
// themselves be Commands that are parsed recursively.
type Command struct {
	// Parameters that can be parsed as matched.
	Options []Param
	// Parameters that are parsed in order.
	Positionals []Param
	// A function executed after this command is parsed. Parsing is not yet complete, so any actions should be deferred
	// from this callback.
	AfterParseFunc AfterParseParamFunc
	// Action taken if no subcommand is invoked. If there are subcommands, and none is chosen and there is no
	// DefaultAction, parsing fails.
	DefaultAction func() error
	// A human description of what this command does.
	Desc string
}

func (me Command) Init() error {
	for _, p := range me.AllParams() {
		err := func() error {
			defer func() {
				r := recover()
				if r != nil {
					panic(fmt.Sprintf("initing %v: %v", p, r))
				}
			}()
			err := p.Init()
			if err != nil {
				return err
			}
			return p.Subcommand().Value.Init()
		}()
		if err != nil {
			return err
		}
	}
	return nil
}

func (me Command) AllParams() []Param {
	return append(me.Options, me.Positionals...)
}

func (me Command) HasSubcommands() bool {
	for _, p := range me.AllParams() {
		if p.Subcommand().Ok {
			return true
		}
	}
	return false
}

func (cmd Command) Help() (hf commandHelp) {
	for _, p := range cmd.Options {
		hf.AddOption(p.Help())
	}
	for _, p := range cmd.Positionals {
		subCmd := p.Subcommand()
		if subCmd.Ok {
			cmdHelp := p.Help()
			cmdHelp.Subcommand = subCmd.Value.Help()
			hf.AddCommand(cmdHelp)
		} else {
			hf.AddPositional(p.Help())
		}
	}
	hf.Desc = cmd.Desc
	return
}
