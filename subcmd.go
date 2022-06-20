package bargle

import (
	"errors"

	"github.com/anacrolix/generics"
)

type Subcommand struct {
	optionDefaults
	Name string
	Command
	Desc string
}

func (me Subcommand) Parse(args Args) error {
	return errors.New("subcommands do not take values")
}

func (me Subcommand) Help() ParamHelp {
	return ParamHelp{
		Forms:       []string{me.Name},
		Description: me.Desc,
	}
}

func (me Subcommand) Match(args Args) MatchResult {
	if args.Len() == 0 {
		return noMatch
	}
	name := args.Pop()
	if name != me.Name {
		return noMatch
	}
	return matchedNoParse{baseMatchResult{
		match: name,
		param: me,
		args:  args,
	}}
}

func (me Subcommand) Subcommand() generics.Option[Command] {
	return generics.Some(me.Command)
}
