package bargle

import "fmt"

// An Arg that tracks how many times it has matched, so that required arguments can be checked
// after parsing.
type ParseCounter interface {
	Arg
	ParseCount() int
}

// Fails the Parser if the argument hasn't matched. Does nothing if the Parser has already failed,
// or help was requested. Unhandled arguments are reported first, since they're the likelier reason
// something is missing.
func (p *Parser) Require(arg ParseCounter) {
	if !p.Ok() || arg.ParseCount() != 0 {
		return
	}
	p.FailIfArgsRemain()
	if p.Ok() {
		p.SetError(fmt.Errorf("%v required and not given", ArgName(arg)))
	}
}

// A name for an Arg suitable for error messages. Prefers a Metavar, falling back to the first form
// the Arg matches.
func ArgName(arg Arg) string {
	if metavar, ok := UnwrapArgAs[Metavar](arg); ok {
		return metavar.Metavar()
	}
	forms := arg.ArgInfo().MatchingForms
	if len(forms) != 0 {
		return forms[0]
	}
	return fmt.Sprintf("%T", arg)
}
