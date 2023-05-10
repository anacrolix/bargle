package bargle

import (
	qt "github.com/frankban/quicktest"
	"strings"
	"testing"
)

func TestHelpAfterPosOnly(t *testing.T) {
	p := NewParser()
	p.SetArgs("--", "--help")
	var helpBuf strings.Builder
	helper := builtinHelper{writer: &helpBuf}
	p.SetHelper(&helper)
	var et string
	ParseLongBuiltin(p, &et, "phone", "home")
	p.FailIfArgsRemain()
	c := qt.New(t)
	c.Assert(helper.helpedCount, qt.Equals, 1)
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
