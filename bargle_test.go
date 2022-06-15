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
	err := ctx.Run(func(ctx Context) {
		ctx.Try(&f)
	})
	qt.Assert(t, err, qt.IsNil)
}

func TestUnhandledExitCode(t *testing.T) {
	ctx := NewContext([]string{"unhandled"})
	err := ctx.Run(func(ctx Context) {
		ctx.Unhandled()
	})
	var exitCoder ExitCoder
	c := qt.New(t)
	c.Assert(err, qt.ErrorAs, &exitCoder)
	c.Check(exitCoder.ExitCode(), qt.Equals, 2)
}

func TestParseFailExitCode(t *testing.T) {
	ctx := NewContext([]string{"unhandled"})
	err := ctx.Run(func(ctx Context) {
		ctx.Parse(Command{})
	})
	//errorSpewConfig.Dump(err)
	var exitCoder ExitCoder
	c := qt.New(t)
	c.Assert(err, qt.ErrorAs, &exitCoder)
	c.Check(exitCoder.ExitCode(), qt.Equals, 2)
}
