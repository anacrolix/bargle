package bargle

type pseudoPosOnly struct{}

func (p pseudoPosOnly) Parse(ctx ParseContext) bool {
	arg, ok := ctx.Pop()
	return ok && arg == "--"
}

func (p pseudoPosOnly) ArgInfo() ArgInfo {
	return ArgInfo{
		MatchingForms: []string{"--"},
		ArgType:       ArgTypePos,
		Global:        false,
	}
}
