package bargle

type HelpCommand struct {
	optionDefaults
}

func (h HelpCommand) Match(args Args) MatchResult {
	if args.Len() == 0 {
		return noMatch
	}
	if args.Pop() != "help" {
		return noMatch
	}
	return matchedNoParse{
		param: h,
		args:  args,
	}
}

//var _ Parser = HelpCommand{}
