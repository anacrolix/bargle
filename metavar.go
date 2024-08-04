package bargle

// A Metavar is an Arg that has a name. This is used for example with positional arguments that
// can't derive an obvious name from their matching forms.
type Metavar interface {
	Metavar() string
}
