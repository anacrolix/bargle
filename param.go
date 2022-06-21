package bargle

import (
	"github.com/anacrolix/generics"
)

type AfterParseParamFunc func(ctx Context) error

type Param interface {
	Satisfied() bool
	Matcher
	// This should not be used directly, normally, and instead Parse is called on a match result.
	// This is used for parsing defaults for example. TODO: Move to Unmarshaler Value and don't
	// expose here?
	Parse(args Args) error
	Subcommand() generics.Option[Command]
	Help() ParamHelp
	AfterParse(ctx Context) error
	Init() error
}

type paramDefaults struct{}

func (paramDefaults) Subcommand() (_ generics.Option[Command]) {
	return
}

func (paramDefaults) AfterParse(Context) error {
	return nil
}

type optionDefaults struct {
	paramDefaults
}

func (optionDefaults) Satisfied() bool {
	return true
}

type posDefaults struct {
	paramDefaults
}

func (posDefaults) Satisfied() bool {
	return false
}
