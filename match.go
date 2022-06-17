package bargle

import (
	"github.com/anacrolix/generics"
)

type MatchResult interface {
	// Could this be represented by a nil interface?
	Matched() generics.Option[string]
	Args() Args
	Parse(ctx Context) error
	// Should the caller of Matcher.Match remember this?
	Param() Param
}

var noMatch noMatchType

type noMatchType struct{}

func (n noMatchType) Matched() (_ generics.Option[string]) {
	return
}

func (n noMatchType) Args() Args {
	panic("unimplemented")
}

func (n noMatchType) Parse(ctx Context) error {
	panic("unimplemented")
}

func (n noMatchType) Param() Param {
	panic("unimplemented")
}

type matchedNoParse struct {
	match string
	param Param
	args  Args
}

func (m matchedNoParse) Matched() generics.Option[string] {
	return generics.Some(m.match)
}

func (m matchedNoParse) Args() Args {
	return m.args
}

func (m matchedNoParse) Parse(ctx Context) error {
	return nil
}

func (m matchedNoParse) Param() Param {
	return m.param
}
