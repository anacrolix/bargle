package bargle

import (
	"testing"

	g "github.com/anacrolix/generics"
	qt "github.com/frankban/quicktest"
)

func TestFlag(t *testing.T) {
	for _, _case := range []struct {
		args     []string
		expected bool
	}{
		{[]string{"--dht"}, true},
		{[]string{"--no-dht"}, false},
		{[]string{"--dht=false"}, false},
		{[]string{"--dht=true"}, true},
		{[]string{"--no-dht=false"}, true},
		{[]string{"--no-dht=true"}, false},
		// The last occurrence wins.
		{[]string{"--no-dht", "--dht"}, true},
		{[]string{"--dht", "--no-dht"}, false},
	} {
		c := qt.New(t)
		p := NewParser()
		p.SetArgs(_case.args...)
		value := !_case.expected
		ParseAll(p, Flag("dht", &value))
		p.FailIfArgsRemain()
		c.Assert(p.Err(), qt.IsNil, qt.Commentf("%q", _case.args))
		c.Check(value, qt.Equals, _case.expected, qt.Commentf("%q", _case.args))
	}
}

// A flag defaulting to true can only be turned off through the negative form.
func TestFlagDefaultTrue(t *testing.T) {
	c := qt.New(t)
	p := NewParser()
	p.SetArgs("--no-progress")
	progress := true
	ParseAll(p, Flag("progress", &progress))
	c.Assert(p.Err(), qt.IsNil)
	c.Check(progress, qt.IsFalse)
}

func TestFlagBadValue(t *testing.T) {
	c := qt.New(t)
	p := NewParser()
	p.SetArgs("--dht=bogus")
	var value bool
	ParseAll(p, Flag("dht", &value))
	c.Check(p.Err(), qt.IsNotNil)
}

// A flag that doesn't match must leave the arguments alone for the next one to try.
func TestFlagNoMatchDoesntConsume(t *testing.T) {
	c := qt.New(t)
	p := NewParser()
	p.SetArgs("--seed")
	var dht, seed bool
	ParseAll(p, Flag("dht", &dht), Flag("seed", &seed))
	p.FailIfArgsRemain()
	c.Assert(p.Err(), qt.IsNil)
	c.Check(dht, qt.IsFalse)
	c.Check(seed, qt.IsTrue)
}

func TestOptionFlag(t *testing.T) {
	c := qt.New(t)
	parse := func(args ...string) g.Option[bool] {
		p := NewParser()
		p.SetArgs(args...)
		var value g.Option[bool]
		ParseAll(p, OptionFlag("private", &value))
		c.Assert(p.Err(), qt.IsNil)
		return value
	}
	c.Check(parse(), qt.Equals, g.None[bool]())
	c.Check(parse("--private"), qt.Equals, g.Some(true))
	c.Check(parse("--no-private"), qt.Equals, g.Some(false))
	c.Check(parse("--private=false"), qt.Equals, g.Some(false))
}

func TestFlagElems(t *testing.T) {
	c := qt.New(t)
	p := NewParser()
	p.SetArgs("--no-port-forward")
	value := true
	ParseAll(p, FlagElems(&value, "port", "forward"))
	c.Assert(p.Err(), qt.IsNil)
	c.Check(value, qt.IsFalse)
}
