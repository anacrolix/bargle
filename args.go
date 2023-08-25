package bargle

type Arg interface {
	Parse(ctx ParseContext) bool
	ArgInfo() ArgInfo
}

type Input struct {
	args []string
}

type ArgType int

const (
	ArgTypeSwitch = iota + 1
	ArgTypeEnvVar
	ArgTypePos
)

type ArgInfo struct {
	MatchingForms []string
	ArgType       ArgType
	// Whether the argument is set at a global level and so always relevant to a parsing scope.
	// Environment variables for example.
	Global bool
}

type ArgDescer interface {
	ArgDesc() string
}
