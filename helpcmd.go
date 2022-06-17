package bargle

import (
	"os"
)

type HelpCommand struct {
	optionDefaults
}

func (h HelpCommand) Help() ParamHelp {
	return ParamHelp{
		Forms:       []string{"help"},
		Description: "walk commands and output help",
	}
}

func (h HelpCommand) Match(args Args) MatchResult {
	if args.Len() == 0 {
		return noMatch
	}
	if args.Pop() != "help" {
		return noMatch
	}
	return matchedNoParse{baseMatchResult{
		param: h,
		args:  args,
		match: "help",
	}}
}

func (me HelpCommand) AddToCommand(cmd *Command) {
	sub := Subcommand{
		Name: "help",
		Desc: "help with subcommands",
	}
	recurse := &Flag{
		Longs:  []string{"recurse"},
		Shorts: []rune{'r'},
	}
	sub.Options = append(sub.Options, recurse)
	addHelpSubcommands(&sub, cmd)
	cmd.Positionals = append(cmd.Positionals, sub)
}

func addHelpSubcommands(to *Subcommand, from *Command) {
	for _, pos := range from.Positionals {
		fromSub := pos.Subcommand()
		if fromSub.Ok {
			toSub := Subcommand{
				Name: pos.Help().Forms[0],
			}
			addHelpSubcommands(&toSub, &fromSub.Value)
			//toSub.DefaultAction = helpCommandAction(&fromSub.Value)
			to.Positionals = append(to.Positionals, toSub)
		}
	}
	to.DefaultAction = helpCommandAction(from)
}

func helpCommandAction(cmd *Command) func() error {
	return func() error {
		var hf helpFormatter
		formatCommandHelp(&hf, cmd)
		hf.Write(os.Stdout)
		return nil
	}
}

func formatCommandHelp(hf *helpFormatter, cmd *Command) {
	for _, p := range cmd.Options {
		hf.AddOption(p.Help())
	}
	for _, p := range cmd.Positionals {
		if p.Subcommand().Ok {
			hf.AddCommand(p.Help())
		} else {
			hf.AddPositional(p.Help())
		}
	}
}

//var _ Parser = HelpCommand{}
