package bargle

type Positional[T any] struct {
	Value T
	U     UnaryUnmarshaler[T]
	Name  string
}

func (me *Positional[T]) Help(f HelpFormatter) {
	f.AddOption(ParamHelp{
		Forms: []string{"<" + me.Name + ">"},
	})
}

func (me *Positional[T]) Parse(ctx Context) error {
	if ctx.Args().Len() == 0 {
		return noMatch
	}
	if !ctx.MatchPos() {
		return noMatch
	}
	return doUnaryUnmarshal(ctx.Args().Pop(), &me.Value, me.U)
}

func NewPositional[T any](u UnaryUnmarshaler[T]) *Positional[T] {
	return &Positional[T]{U: u, Name: "arg"}
}
