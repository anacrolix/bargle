package bargle

import (
	"github.com/anacrolix/generics"
)

type Param interface {
	Satisfied() bool
	Matcher
	Subcommand() generics.Option[Command]
}

type paramDefaults struct{}

func (p paramDefaults) Parse(ctx Context) error {
	return nil
}

func (p paramDefaults) Subcommand() (_ generics.Option[Command]) {
	return
}

type optionDefaults struct {
	paramDefaults
}

func (o optionDefaults) Satisfied() bool {
	return true
}

type posDefaults struct {
	paramDefaults
}

func (p posDefaults) Satisfied() bool {
	return false
}
