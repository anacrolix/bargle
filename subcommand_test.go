package bargle

import (
	"strings"
	"testing"

	qt "github.com/frankban/quicktest"
)

func testSubcommands(p *Parser) string {
	return ParseSubcommand(p,
		Subcommand[string]{
			Name: "magnet",
			Desc: "print a magnet link",
			Parse: func(p *Parser) string {
				return "magnet"
			},
		},
		Subcommand[string]{
			Name: "pprint",
			Desc: "pretty print",
			Parse: func(p *Parser) string {
				var files bool
				ParseAll(p, Flag("files", &files))
				if files {
					return "pprint files"
				}
				return "pprint"
			},
		},
	)
}

func TestParseSubcommand(t *testing.T) {
	c := qt.New(t)
	parse := func(args ...string) (string, error) {
		p := NewParser()
		p.SetArgs(args...)
		ret := testSubcommands(p)
		p.FailIfArgsRemain()
		return ret, p.Err()
	}

	ret, err := parse("magnet")
	c.Assert(err, qt.IsNil)
	c.Check(ret, qt.Equals, "magnet")

	ret, err = parse("pprint", "--files")
	c.Assert(err, qt.IsNil)
	c.Check(ret, qt.Equals, "pprint files")

	_, err = parse("bogus")
	c.Assert(err, qt.IsNotNil)
	c.Check(err.Error(), qt.Equals, `unused argument: "bogus"`)

	_, err = parse()
	c.Check(err, qt.Equals, ErrExpectedArguments)
}

// Nothing matches while help is being requested, and that mustn't be reported as a failure.
func TestParseSubcommandHelp(t *testing.T) {
	c := qt.New(t)
	p := NewParser()
	p.SetArgs("--help")
	var helpBuf strings.Builder
	p.SetHelper(&builtinHelper{writer: &helpBuf})
	testSubcommands(p)
	p.FailIfArgsRemain()
	c.Assert(p.Err(), qt.IsNil)
	c.Check(p.Ok(), qt.IsFalse)
	p.DoHelpIfHelpingOpts(PrintHelpOpts{NoPrintUsage: true})
	c.Check(helpBuf.String(), qt.Contains, "print a magnet link")
	c.Check(helpBuf.String(), qt.Contains, "pretty print")
}
