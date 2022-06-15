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
	return me.U.UnaryUnmarshal(ctx.Args().Pop(), &me.Value)
}

func NewPositional[T any](u UnaryUnmarshaler[T]) *Positional[T] {
	return &Positional[T]{U: u, Name: "arg"}
}
