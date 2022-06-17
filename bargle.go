package bargle

type Parser interface {
	Parse(ctx Context) error
}

type Matcher interface {
	// Should this allow a nil return for no match?
	Match(args Args) MatchResult
}
