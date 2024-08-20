package bargle

type Arg interface {
	Parse(ctx ParseContext) bool
	ArgInfo() ArgInfo
}

type ArgValuer interface {
	Value() any
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

type withDesc struct {
	desc string
	Arg
}

func (me withDesc) ArgDesc() string {
	return me.desc
}

func WithDesc(desc string, arg Arg) interface {
	Arg
	ArgDescer
} {
	return withDesc{desc, arg}
}
