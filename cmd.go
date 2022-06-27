package bargle

import (
	"fmt"
)

type Command struct {
	Options        []Param
	Positionals    []Param
	AfterParseFunc AfterParseParamFunc
	// Action taken if no subcommand is invoked.
	DefaultAction func() error
	Desc          string
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
