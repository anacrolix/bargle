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

func (me withDesc) UnwrapArg() Arg {
	return me.Arg
}

func WithDesc(desc string, arg Arg) interface {
	Arg
	ArgDescer
	ArgWrapper
} {
	return withDesc{desc, arg}
}

// Implemented by Args that wrap another Arg. Wrapping an Arg hides any interfaces it implements
// beyond Arg itself, so wrappers expose the Arg they wrap to let helpers look through them. See
// UnwrapArgAs.
type ArgWrapper interface {
	UnwrapArg() Arg
}

// Returns the first Arg in a chain of ArgWrappers, starting with arg itself, that implements T.
// Use this instead of a plain type assertion when an Arg might have been wrapped by something like
// WithDesc.
func UnwrapArgAs[T any](arg Arg) (t T, ok bool) {
	for {
		t, ok = any(arg).(T)
		if ok {
			return
		}
		wrapper, isWrapper := arg.(ArgWrapper)
		if !isWrapper {
			return
		}
		arg = wrapper.UnwrapArg()
	}
}
