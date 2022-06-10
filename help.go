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
	Options     []ParamHelp
}

func (me ParamHelp) Write(w io.Writer) {
	fmt.Fprintf(w, "  %s\n", strings.Join(me.Forms, ", "))
	if me.Description != "" {
		fmt.Fprintf(w, "    %s\n", me.Description)
	}
	for _, o := range me.Options {
		o.Write(w)
	}
}

type helpFormatter struct {
	Options  []ParamHelp
	Commands []ParamHelp
}

func (me helpFormatter) Write(w io.Writer) {
	for _, ph := range me.Options {
		ph.Write(w)
	}
	for _, ph := range me.Commands {
		ph.Write(w)
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

func (ph ParamHelp) Print(w io.Writer) {
	fmt.Fprintf(w, "  %s\n", strings.Join(ph.Forms, ", "))
	fmt.Fprintf(w, "    %s\n", ph.Description)
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
