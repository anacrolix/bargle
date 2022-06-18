package bargle

import (
	"github.com/anacrolix/generics"
)

type AfterParseParamFunc func(ctx Context) error

type Param interface {
	Satisfied() bool
	Matcher
	Subcommand() generics.Option[Command]
	Help() ParamHelp
	AfterParse(ctx Context) error
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
