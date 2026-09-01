package bargle

import (
	"strings"
	"testing"

	g "github.com/anacrolix/generics"
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

// Printing an Arg's current value must cope with Unmarshalers that don't track one, and with
// values that aren't strings.
func TestHelpArgValues(t *testing.T) {
	c := qt.New(t)
	p := NewParser()
	p.SetArgs("--help")
	var helpBuf strings.Builder
	p.SetHelper(&builtinHelper{writer: &helpBuf})
	var (
		name     string
		port     int
		debug    bool
		private  g.Option[bool]
		trackers []string
		event    textUnmarshalerValue
	)
	ParseAll(p,
		Long("name", String(&name)),
		Long("port", BuiltinUnmarshaler(&port)),
		Flag("debug", &debug),
		OptionFlag("private", &private),
		// An accumulating Unmarshaler has no single current value.
		Long("tracker", AppendSlice(&trackers, BuiltinUnmarshaler[string])),
		// Neither does one that only knows how to consume text.
		Long("event", TextUnmarshaler(&event)),
	)
	p.DoHelpIfHelpingOpts(PrintHelpOpts{NoPrintUsage: true})
	help := helpBuf.String()
	c.Check(help, qt.Contains, `--name=string, --name string [current value: ""]`)
	c.Check(help, qt.Contains, `--port=int, --port int [current value: 0]`)
	c.Check(help, qt.Contains, `--[no-]debug, --[no-]debug=bool [current value: false]`)
	// Not given, so there's nothing to report.
	c.Check(help, qt.Contains, "--[no-]private, --[no-]private=bool\n")
	c.Check(help, qt.Contains, "--tracker=string..., --tracker string...\n")
	c.Check(help, qt.Contains, "--event=string, --event string\n")
	c.Check(help, qt.Not(qt.Contains), "%!q")
}

type textUnmarshalerValue struct{}

func (textUnmarshalerValue) UnmarshalText([]byte) error { return nil }
