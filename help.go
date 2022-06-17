package bargle

import (
	"fmt"
	"io"
	"strings"
)

type TargetHelper interface {
	TargetHelp() string
}

type Help struct {
	params []ParamHelper
}

type FormHelper interface {
	Help(*ParamHelp)
}

func (me Help) matchers() []interface {
	//Matcher
	FormHelper
} {
	return []interface {
		//Matcher
		FormHelper
	}{&LongParser{Long: "help"}, &ShortParser{Short: 'h'}}
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
	if me.Description != "" {
		s += "\t" + me.Description
	}
	w.WriteLine(s)
	for _, o := range me.Options {
		o.Write(w.Indented())
	}
}

type helpFormatter struct {
	Options    []ParamHelp
	Positional []ParamHelp
	Commands   []ParamHelp
}

func (me helpFormatter) Write(w io.Writer) {
	hw := HelpWriter{w: w}
	hw.WriteLine("options:")
	for _, ph := range me.Options {
		ph.Write(hw.Indented())
	}
	if len(me.Positional) != 0 {
		hw.WriteLine("arguments:")
		for _, ph := range me.Positional {
			ph.Write(hw.Indented())
		}
	}
	if len(me.Commands) != 0 {
		hw.WriteLine("commands:")
		for _, ph := range me.Commands {
			ph.Write(hw.Indented())
		}
	}
}

type HelpFormatter = *helpFormatter

func (me *helpFormatter) AddCommand(name, desc string) {
	me.Commands = append(me.Commands, ParamHelp{
		Forms:       []string{name},
		Description: desc,
	})
}

func (me *helpFormatter) AddPositional(ph ParamHelp) {
	me.Positional = append(me.Positional, ph)
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
	me.Help(&f)
	for _, p := range me.params {
		p.Help(&f)
	}
	f.Write(w)
}

func (me Help) Help(f HelpFormatter) {
	ph := ParamHelp{Description: "help (this message)"}
	for _, h := range me.matchers() {
		h.Help(&ph)
	}
	f.AddOption(ph)
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
