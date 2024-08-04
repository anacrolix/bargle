package bargle

import (
	"strings"
	"testing"

	qt "github.com/frankban/quicktest"
)

func TestHelpAfterPosOnly(t *testing.T) {
	p := NewParser()
	p.SetArgs("--", "--help")
	var helpBuf strings.Builder
	helper := builtinHelper{writer: &helpBuf}
	p.SetHelper(&helper)
	var et string
	ParseLongBuiltin(p, &et, "phone", "home")
	p.DoHelpIfHelpingOpts(PrintHelpOpts{NoPrintUsage: true})
	c := qt.New(t)
	c.Check(helper.helpedCount, qt.Equals, 1)
}

func TestHelpAfterPosOnlyNoArgumentsExpected(t *testing.T) {
	p := NewParser()
	p.SetArgs("--", "--help")
	var helpBuf strings.Builder
	helper := builtinHelper{writer: &helpBuf}
	p.SetHelper(&helper)
	p.FailIfArgsRemain()
	p.DoHelpIfHelpingOpts(PrintHelpOpts{NoPrintUsage: true})
	c := qt.New(t)
	c.Check(helper.helpedCount, qt.Equals, 1)
	c.Assert(helpBuf.String(), qt.Equals, noArgumentsExpectedHelp)
}

func TestSolitaryHelp(t *testing.T) {
	p := NewParser()
	p.SetArgs("--help")
	var helpBuf strings.Builder
	helper := builtinHelper{writer: &helpBuf}
	p.SetHelper(&helper)
	var et string
	ParseLongBuiltin(p, &et, "phone", "home")
	p.DoHelpIfHelping()
	c := qt.New(t)
	c.Assert(helper.helpedCount, qt.Equals, 1)
}
