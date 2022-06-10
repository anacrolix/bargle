package bargle

import (
	"fmt"
	"io"
	"os"
	"strings"
)

type Help struct {
	params []ParamHelper
}

func (me *Help) Parse(ctx Context) error {
	if ctx.Match(LongParser{Long: "help"}) {
		me.Print(os.Stdout)
		ctx.Success()
	}
	return noMatch
}

type ParamHelp struct {
	Forms       []string
	Description string
	Values      string
	Options     []ParamHelp
}

func (me ParamHelp) Write(w HelpWriter) {
	s := strings.Join(me.Forms, ", ")
	if me.Values != "" {
		s += ": " + me.Values
	}
	w.WriteLine(s)
	if me.Description != "" {
		w.Indented().WriteLine(me.Description)
	}
	for _, o := range me.Options {
		o.Write(w.Indented())
	}
}

type helpFormatter struct {
	Options  []ParamHelp
	Commands []ParamHelp
}

func (me helpFormatter) Write(w io.Writer) {
	hw := HelpWriter{w: w}
	hw.WriteLine("usage:")
	hw = hw.Indented()
	for _, ph := range me.Options {
		ph.Write(hw)
	}
	for _, ph := range me.Commands {
		ph.Write(hw)
	}
}

type HelpFormatter = *helpFormatter

func (me *helpFormatter) AddCommand(name string) {
	me.Commands = append(me.Commands, ParamHelp{
		Forms: []string{name},
	})
}

func (me *helpFormatter) AddOption(ph ParamHelp) {
	me.Options = append(me.Options, ph)
}

func (ph ParamHelp) Print(w HelpWriter) {
	w.WriteLine(strings.Join(ph.Forms, ", "))
	w.Indented().WriteLine(ph.Description)
}

func (me *Help) AddParams(params ...ParamHelper) {
	me.params = append(me.params, params...)
}

func (me Help) Print(w io.Writer) {
	f := helpFormatter{}
	for _, p := range me.params {
		p.Help(&f)
	}
	f.Write(w)
}

type ParamHelper interface {
	Help(f HelpFormatter)
}

type HelpWriter struct {
	indent int
	w      io.Writer
}

func (me HelpWriter) WriteLine(s string) {
	fmt.Fprintf(me.w, "%s%s\n", strings.Repeat("  ", me.indent), s)
}

func (me HelpWriter) Indented() HelpWriter {
	me.indent++
	return me
}
