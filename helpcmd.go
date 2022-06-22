package bargle

import (
	"errors"
	"os"
)

type HelpCommand struct {
	optionDefaults
}

func (h HelpCommand) Init() error {
	// Could walk an attached param here.
	return nil
}

func (h HelpCommand) Parse(args Args) error {
	return errors.New("help does not take value")
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
	}
	sub.Desc = "help with subcommands"
	recurse := &Flag{
		Value: new(bool),
	}
	recurse.longs = []string{"recurse"}
	recurse.shorts = []rune{'r'}
	sub.Options = append(sub.Options, recurse)
	cmd.Positionals = append(cmd.Positionals, &sub)
	addHelpSubcommands(&sub, cmd, recurse.Value)
}

func addHelpSubcommands(to *Subcommand, from *Command, recurse *bool) {
	for _, pos := range from.Positionals {
		fromSub := pos.Subcommand()
		if fromSub.Ok {
			toSub := Subcommand{
				Name: pos.Help().Forms[0],
			}
			addHelpSubcommands(&toSub, &fromSub.Value, recurse)
			//toSub.DefaultAction = helpCommandAction(&fromSub.Value)
			to.Positionals = append(to.Positionals, toSub)
		}
	}
	to.DefaultAction = helpCommandAction(from, recurse)
}

func printCommandHelp(ch commandHelp, recurse bool) {
	helpFormatter{recurse}.Write(os.Stdout, ch)
}

func helpCommandAction(cmd *Command, recurse *bool) func() error {
	return func() error {
		printCommandHelp(cmd.Help(), *recurse)
		return nil
	}
}
