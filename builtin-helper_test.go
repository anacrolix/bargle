package bargle

import (
	"strings"
	"testing"

	"github.com/go-quicktest/qt"
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
	qt.Check(t, qt.Equals(helper.helpedCount, 1))
}

func TestHelpAfterPosOnlyNoArgumentsExpected(t *testing.T) {
	p := NewParser()
	p.SetArgs("--", "--help")
	var helpBuf strings.Builder
	helper := builtinHelper{writer: &helpBuf}
	p.SetHelper(&helper)
	p.FailIfArgsRemain()
	p.DoHelpIfHelpingOpts(PrintHelpOpts{NoPrintUsage: true})
	qt.Check(t, qt.Equals(helper.helpedCount, 1))
	qt.Assert(t, qt.Equals(helpBuf.String(), noArgumentsExpectedHelp))
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
	qt.Assert(t, qt.Equals(helper.helpedCount, 1))
}
