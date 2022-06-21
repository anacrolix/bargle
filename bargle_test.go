package bargle

import (
	"testing"

	"github.com/anacrolix/tagflag"
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

// Checks that we can unmarshal into pointers, and in to pointers to pointers that implement
// encoding.TextUnmarshaler.
func TestUnmarshalPointer(t *testing.T) {
	c := qt.New(t)
	// Test unmarshalling in place
	var b tagflag.Bytes
	uf, err := makeAnyUnaryUnmarshalerViaReflection(&b)
	c.Assert(err, qt.IsNil)
	err = uf.UnaryUnmarshal("1M")
	c.Assert(err, qt.IsNil)
	c.Check(b.Int64(), qt.Equals, int64(1_000_000))
	// Unmarshal into nil pointer
	var pb *tagflag.Bytes
	uf, err = makeAnyUnaryUnmarshalerViaReflection(&pb)
	c.Assert(err, qt.IsNil)
	err = uf.UnaryUnmarshal("2M")
	c.Assert(err, qt.IsNil)
	c.Check(pb.Int64(), qt.Equals, int64(2_000_000))
	// Check that reusing an unmarshal func doesn't adapt to a new pointer
	npb := new(tagflag.Bytes)
	pb = npb
	err = uf.UnaryUnmarshal("3M")
	c.Assert(err, qt.IsNil)
	c.Check(pb.Int64(), qt.Equals, int64(3_000_000))
	c.Check(pb, qt.Not(qt.Equals), npb)
	c.Check(npb.Int64(), qt.Equals, int64(0))
	// Unmarshal into pointer to pointer
	pab := new(tagflag.Bytes)
	uf, err = makeAnyUnaryUnmarshalerViaReflection(pab)
	c.Assert(err, qt.IsNil)
	err = uf.UnaryUnmarshal("4M")
	c.Assert(err, qt.IsNil)
	c.Check(pab.Int64(), qt.Equals, int64(4_000_000))
}
