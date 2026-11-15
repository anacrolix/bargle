package bargle

// A keyword that selects a command, and the parsing to do for the arguments that follow it.
type Subcommand[T any] struct {
	Name string
	// Optional description for help output.
	Desc string
	// Parses the rest of the command's arguments. What it returns is passed back out of
	// ParseSubcommand, so it can carry whatever the caller needs, such as the work to do once the
	// whole command line has parsed successfully.
	Parse func(p *Parser) T
}

// Returns the Subcommand's keyword, described for help output.
func (me Subcommand[T]) Keyword() Arg {
	keyword := Keyword(me.Name)
	if me.Desc == "" {
		return keyword
	}
	return WithDesc(me.Desc, keyword)
}

// Parses the keyword of the first given Subcommand that matches, and runs its Parse. Fails the
// Parser if none of them match. Returns the zero value if nothing matched, or help was requested,
// so check Parser.Ok before using it.
func ParseSubcommand[T any](p *Parser, subs ...Subcommand[T]) (ret T) {
	for _, sub := range subs {
		if p.Parse(sub.Keyword()) {
			return sub.Parse(p)
		}
	}
	if p.Ok() {
		p.Fail()
	}
	return
}
