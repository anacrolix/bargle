package bargle

import (
	"testing"

	qt "github.com/frankban/quicktest"
)

func TestPositionalMatchesOnce(t *testing.T) {
	c := qt.New(t)
	p := NewParser()
	p.SetArgs("first", "second")
	var value string
	pos := Positional("name", String(&value))
	ParseAll(p, pos)
	c.Check(value, qt.Equals, "first")
	c.Check(pos.ParseCount(), qt.Equals, 1)
	c.Check(p.PopAll(), qt.DeepEquals, []string{"second"})
}

func TestPositionals(t *testing.T) {
	c := qt.New(t)
	p := NewParser()
	p.SetArgs("a", "b", "c")
	var values []string
	pos := Positionals("name", AppendSlice(&values, BuiltinUnmarshaler[string]))
	ParseAll(p, pos)
	p.FailIfArgsRemain()
	c.Assert(p.Err(), qt.IsNil)
	c.Check(values, qt.DeepEquals, []string{"a", "b", "c"})
	c.Check(pos.ParseCount(), qt.Equals, 3)
}

// Positionals ordered after each other are filled in turn, and switches can be interleaved.
func TestPositionalOrderingWithSwitches(t *testing.T) {
	c := qt.New(t)
	p := NewParser()
	p.SetArgs("tracker-url", "--port", "1234", "hash-a", "hash-b")
	var (
		tracker string
		port    int
		hashes  []string
	)
	trackerPos := Positional("tracker", String(&tracker))
	hashesPos := Positionals("info-hash", AppendSlice(&hashes, BuiltinUnmarshaler[string]))
	ParseAll(p, Long("port", BuiltinUnmarshaler(&port)), trackerPos, hashesPos)
	p.FailIfArgsRemain()
	c.Assert(p.Err(), qt.IsNil)
	c.Check(tracker, qt.Equals, "tracker-url")
	c.Check(port, qt.Equals, 1234)
	c.Check(hashes, qt.DeepEquals, []string{"hash-a", "hash-b"})
}

func TestRequireGiven(t *testing.T) {
	c := qt.New(t)
	p := NewParser()
	p.SetArgs("a")
	var value string
	pos := Positional("name", String(&value))
	ParseAll(p, pos)
	p.Require(pos)
	c.Check(p.Err(), qt.IsNil)
}

func TestRequireMissing(t *testing.T) {
	c := qt.New(t)
	p := NewParser()
	p.SetArgs()
	var value string
	pos := Positional("name", String(&value))
	ParseAll(p, pos)
	p.Require(pos)
	c.Assert(p.Err(), qt.IsNotNil)
	c.Check(p.Err().Error(), qt.Equals, `name required and not given`)
}

// An argument we couldn't handle is the likelier explanation for a missing required argument, so
// it's reported instead.
func TestRequireReportsUnhandledArgFirst(t *testing.T) {
	c := qt.New(t)
	p := NewParser()
	p.SetArgs("--bogus")
	var value string
	pos := Positional("name", String(&value))
	ParseAll(p, pos)
	p.Require(pos)
	c.Assert(p.Err(), qt.IsNotNil)
	c.Check(p.Err().Error(), qt.Equals, `unused argument: "--bogus"`)
}

// Requesting help isn't a parse failure, so Require mustn't turn it into one.
func TestRequireWhileHelping(t *testing.T) {
	c := qt.New(t)
	p := NewParser()
	p.SetArgs("--help")
	var value string
	pos := Positional("name", String(&value))
	ParseAll(p, pos)
	p.Require(pos)
	c.Check(p.Err(), qt.IsNil)
	c.Check(p.Ok(), qt.IsFalse)
}
