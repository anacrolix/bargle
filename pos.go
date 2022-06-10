package bargle

type Positional[T any] struct {
	Value T
	U     Unmarshaler[T]
	Name  string
}

func (me *Positional[T]) Help(f HelpFormatter) {
	f.AddOption(ParamHelp{
		Forms: []string{"<" + me.Name + ">"},
	})
}

func (me *Positional[T]) Parse(ctx Context) error {
	return me.U.Unmarshal(ctx.Args(), &me.Value)
}

func NewPositional[T any](u Unmarshaler[T]) *Positional[T] {
	return &Positional[T]{U: u, Name: "arg"}
}
