package bargle

type Parser interface {
	Parse(ctx Context) error
}
