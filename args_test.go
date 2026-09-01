package bargle

import (
	"strings"
	"testing"

	qt "github.com/frankban/quicktest"
)

// Wrapping an Arg for its description mustn't hide the interfaces it implements from the help
// formatter.
func TestWithDescKeepsMetavarAndValue(t *testing.T) {
	c := qt.New(t)
	p := NewParser()
	p.SetArgs("--help")
	var helpBuf strings.Builder
	p.SetHelper(&builtinHelper{writer: &helpBuf})
	var (
		name  string
		files []string
		debug bool
	)
	ParseAll(p,
		WithDesc("the torrent name", Positional("name", String(&name))),
		WithDesc("files to add", Positionals("file", AppendSlice(&files, BuiltinUnmarshaler[string]))),
		WithDesc("enable debug logging", Flag("debug", &debug)),
	)
	p.DoHelpIfHelpingOpts(PrintHelpOpts{NoPrintUsage: true})
	help := helpBuf.String()
	c.Check(help, qt.Contains, "name: string")
	c.Check(help, qt.Contains, "the torrent name")
	c.Check(help, qt.Contains, "file: string...")
	// Flags have no metavar, so they mustn't gain an empty one.
	c.Check(help, qt.Contains, "--[no-]debug")
	c.Check(help, qt.Not(qt.Contains), ": --[no-]debug")
}

func TestUnwrapArgAs(t *testing.T) {
	c := qt.New(t)
	pos := Positional("name", String(new(string)))
	wrapped := WithDesc("outer", WithDesc("inner", pos))
	metavar, ok := UnwrapArgAs[Metavar](wrapped)
	c.Assert(ok, qt.IsTrue)
	c.Check(metavar.Metavar(), qt.Equals, "name")
	_, ok = UnwrapArgAs[ParseCounter](Keyword("nope"))
	c.Check(ok, qt.IsFalse)
}

func TestArgName(t *testing.T) {
	c := qt.New(t)
	c.Check(ArgName(WithDesc("d", Positional("torrent file", String(new(string))))), qt.Equals, "torrent file")
	c.Check(ArgName(Flag("debug", new(bool))), qt.Equals, "--[no-]debug, --[no-]debug=bool")
}
