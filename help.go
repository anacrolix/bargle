package bargle

import (
	"fmt"
	"io"
	"strings"
)

type TargetHelper interface {
	TargetHelp() string
}

type FormHelper interface {
	Help(*ParamHelp)
}

type ParamHelp struct {
	Forms       []string
	Description string
	Values      string
	Options     []ParamHelp
	Subcommand  commandHelp
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

type commandHelp struct {
	Options    []ParamHelp
	Positional []ParamHelp
	Commands   []ParamHelp
	Desc       string
}

type helpFormatter struct {
	recurse bool
}

func (me helpFormatter) writeParam(hw HelpWriter, ph ParamHelp) {
	ph.Write(hw)
	if me.recurse {
		me.writeCommand(hw.Indented(), ph.Subcommand)
	}
}

func (me helpFormatter) writeCommand(hw HelpWriter, ch commandHelp) {
	if ch.Desc != "" {
		hw.WriteLine(ch.Desc)
	}
	if len(ch.Options) != 0 {
		hw.WriteLine("options:")
		for _, ph := range ch.Options {
			me.writeParam(hw.Indented(), ph)
		}
	}
	if len(ch.Positional) != 0 {
		hw.WriteLine("arguments:")
		for _, ph := range ch.Positional {
			me.writeParam(hw.Indented(), ph)
		}
	}
	if len(ch.Commands) != 0 {
		hw.WriteLine("commands:")
		for _, ph := range ch.Commands {
			me.writeParam(hw.Indented(), ph)
		}
	}
}

func (me helpFormatter) Write(w io.Writer, ch commandHelp) {
	hw := HelpWriter{w: w}
	me.writeCommand(hw, ch)
}

func (me *commandHelp) AddCommand(ph ParamHelp) {
	me.Commands = append(me.Commands, ph)
}

func (me *commandHelp) AddPositional(ph ParamHelp) {
	me.Positional = append(me.Positional, ph)
}

func (me *commandHelp) AddOption(ph ParamHelp) {
	me.Options = append(me.Options, ph)
}

func (ph ParamHelp) Print(w HelpWriter) {
	w.WriteLine(strings.Join(ph.Forms, ", "))
	w.Indented().WriteLine(ph.Description)
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
