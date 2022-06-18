package bargle

type Command struct {
	Options     []Param
	Positionals []Param
	// Action taken if no subcommand is invoked.
	DefaultAction func() error
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
	return
}
