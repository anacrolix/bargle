package bargle

type HelpCommand struct {
}

func (h HelpCommand) Parse(ctx Context) error {
	if ctx.Helping() {
		return noMatch
	}
	if ctx.Args().Len() == 0 {
		return noMatch
	}
	if ctx.Args().Pop() != "help" {
		return noMatch
	}
	ctx.StartHelping()
	return nil
}

var _ Parser = HelpCommand{}
