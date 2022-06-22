package bargle

import (
	"fmt"
	"log"
	"testing"

	"github.com/anacrolix/tagflag"
	"github.com/davecgh/go-spew/spew"
	qt "github.com/frankban/quicktest"
)

var errorSpewConfig = spew.NewDefaultConfig()

func init() {
	errorSpewConfig.DisableMethods = true
	log.SetFlags(log.Flags() | log.Lshortfile)
}

func TestParseFlagNoArgs(t *testing.T) {
	ctx := NewContext(nil)
	f := NewFlag(nil)
	f.AddLong("debug")
	f.AddShort('s')
	err := ctx.Run(Command{Options: []Param{f.Make()}, DefaultAction: func() error {
		// Do nothing? Hm...
		return nil
	}})
	qt.Assert(t, err, qt.IsNil)
}

func TestParseFlagNoTarget(t *testing.T) {
	ctx := NewContext([]string{"-s"})
	fm := NewFlag(nil)
	fm.AddLong("debug")
	fm.AddShort('s')
	f := fm.Make()
	err := ctx.Run(Command{Options: []Param{f}, DefaultAction: func() error {
		// Do nothing? Hm...
		return nil
	}})
	c := qt.New(t)
	c.Assert(err, qt.IsNil)
	c.Check(*f.Value, qt.IsTrue)
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

func TestFromStructDefaults(t *testing.T) {
	c := qt.New(t)
	var struct_ struct {
		DefaultTrue bool `default:"true"`
		NoDefault   bool
		Default420  string   `default:"420"`
		SetManually []string `default:"world"`
	}
	struct_.SetManually = append(struct_.SetManually, "hello")
	FromStruct(&struct_)
	c.Check(struct_.DefaultTrue, qt.IsTrue)
	c.Check(struct_.NoDefault, qt.IsFalse)
	c.Check(struct_.Default420, qt.Equals, "420")
	c.Check(struct_.SetManually, qt.DeepEquals, []string{"hello", "world"})
}

type trigram string

func (me *trigram) UnmarshalText(b []byte) error {
	if len(b) != 3 {
		return fmt.Errorf("expected 3 chars, got %v", len(b))
	}
	*me = trigram(string(b))
	return nil
}

func TestUnmarshalStructSliceTextUnmarshaler(t *testing.T) {
	var struct_ struct {
		Nope  string    `arg:"positional"`
		Hello []trigram `arg:"positional" arity:"+"`
	}
	withCmd := func(f func(cmd Command)) {
		struct_.Nope = ""
		struct_.Hello = nil
		cmd := FromStruct(&struct_)
		cmd.DefaultAction = func() error { return nil }
		f(cmd)
	}
	c := qt.New(t)

	withCmd(func(cmd Command) {
		ctx := NewContext(nil)
		err := ctx.Run(cmd)
		var up unsatisfiedParam
		if c.Check(err, qt.ErrorAs, &up) {
			c.Check(up.p, qt.Equals, cmd.Positionals[0])
		}
	})

	withCmd(func(cmd Command) {
		ctx := NewContext([]string{"nope"})
		err := ctx.Run(cmd)
		var up unsatisfiedParam
		if c.Check(err, qt.ErrorAs, &up) {
			c.Check(up.p, qt.Equals, cmd.Positionals[1])
		}
	})

	withCmd(func(cmd Command) {
		ctx := NewContext([]string{"nope", "abc"})
		err := ctx.Run(cmd)
		c.Assert(err, qt.IsNil)
		c.Check(struct_.Hello, qt.DeepEquals, []trigram{"abc"})
	})

	withCmd(func(cmd Command) {
		ctx := NewContext([]string{"nope", "abc", "def"})
		err := ctx.Run(cmd)
		c.Assert(err, qt.IsNil)
		c.Check(struct_.Hello, qt.DeepEquals, []trigram{"abc", "def"})
	})

	withCmd(func(cmd Command) {
		ctx := NewContext([]string{"nope", "abc", "herp", "def"})
		err := ctx.Run(cmd)
		c.Assert(err, qt.IsNotNil)
		c.Check(struct_.Hello, qt.DeepEquals, []trigram{"abc"})
	})
}
