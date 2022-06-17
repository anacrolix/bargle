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
	baseMatchResult
}

func (m matchedNoParse) Parse(ctx Context) error {
	return nil
}

type baseMatchResult struct {
	args  Args
	param Param
	match string
}

func (me baseMatchResult) Matched() generics.Option[string] {
	return generics.Some(me.match)
}

func (me baseMatchResult) Args() Args {
	return me.args
}

func (me baseMatchResult) Param() Param {
	return me.param
}
