package bargle

import (
	"strings"
	"testing"

	qt "github.com/frankban/quicktest"
)

// Wrapping an Arg for its description mustn't hide the interfaces it implements from the help
// formatter.
func TestWithDescKeepsMetavar(t *testing.T) {
	c := qt.New(t)
	p := NewParser()
	p.SetArgs("--help")
	var helpBuf strings.Builder
	p.SetHelper(&builtinHelper{writer: &helpBuf})
	var (
		name string
		port int
	)
	ParseAll(p,
		WithDesc("the torrent name", Positional("name", String(&name))),
		WithDesc("the port to listen on", Long("port", BuiltinUnmarshaler(&port))),
	)
	p.DoHelpIfHelpingOpts(PrintHelpOpts{NoPrintUsage: true})
	help := helpBuf.String()
	c.Check(help, qt.Contains, "name: string")
	c.Check(help, qt.Contains, "the torrent name")
	// Longs have no metavar, so they mustn't gain an empty one.
	c.Check(help, qt.Contains, "--port=int")
	c.Check(help, qt.Not(qt.Contains), ": --port=int")
}

func TestUnwrapArgAs(t *testing.T) {
	c := qt.New(t)
	pos := Positional("name", String(new(string)))
	wrapped := WithDesc("outer", WithDesc("inner", pos))
	metavar, ok := UnwrapArgAs[Metavar](wrapped)
	c.Assert(ok, qt.IsTrue)
	c.Check(metavar.Metavar(), qt.Equals, "name")
	_, ok = UnwrapArgAs[Metavar](Keyword("nope"))
	c.Check(ok, qt.IsFalse)
}
