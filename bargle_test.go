package bargle

import (
	"testing"

	"github.com/davecgh/go-spew/spew"
	qt "github.com/frankban/quicktest"
)

var errorSpewConfig = spew.NewDefaultConfig()

func init() {
	errorSpewConfig.DisableMethods = true
}

func TestParseFlagNoArgs(t *testing.T) {
	ctx := NewContext(nil)
	f := Flag{}
	f.AddLong("debug").AddShort('s')
	err := ctx.Run(Command{Options: []Param{&f}, DefaultAction: func() error {
		// Do nothing? Hm...
		return nil
	}})
	qt.Assert(t, err, qt.IsNil)
}

func TestUnhandledExitCode(t *testing.T) {
	ctx := NewContext([]string{"unhandled"})
	err := ctx.Run(Command{})
	var exitCoder ExitCoder
	c := qt.New(t)
	c.Assert(err, qt.ErrorAs, &exitCoder)
	c.Check(exitCoder.ExitCode(), qt.Equals, 2)
}

func TestParseFailExitCode(t *testing.T) {
	ctx := NewContext([]string{"unhandled"})
	err := ctx.Run(Command{Positionals: []Param{Subcommand{}}})
	//errorSpewConfig.Dump(err)
	var exitCoder ExitCoder
	c := qt.New(t)
	c.Assert(err, qt.ErrorAs, &exitCoder)
	c.Check(exitCoder.ExitCode(), qt.Equals, 2)
}
