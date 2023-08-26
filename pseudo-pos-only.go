package bargle

type pseudoPosOnly struct{}

var _ interface {
	ArgDescer
} = pseudoPosOnly{}

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

func (p pseudoPosOnly) ArgDesc() string {
	return "starts parsing positional arguments only"
}
