package bargle

import (
	"strconv"
)

type Flag struct {
	Value  bool
	Longs  []string
	Shorts []rune
}

func (f *Flag) Help(fm HelpFormatter) {
	ph := ParamHelp{
		Values: "[true|1|false|0]",
	}
	for _, l := range f.Longs {
		ph.Forms = append(ph.Forms, "--"+l, "--no-"+l)
	}
	for _, s := range f.Shorts {
		ph.Forms = append(ph.Forms, "-"+string(s), "+"+string(s))
	}
	fm.AddOption(ph)
}

func (f *Flag) AddLong(l string) *Flag {
	f.Longs = append(f.Longs, l)
	return f
}

func (f *Flag) parseValue(no bool, us UnarySwitch, ctx Context) (err error) {
	if us.GotValue() {
		f.Value, err = strconv.ParseBool(ctx.Args().Pop())
	} else {
		f.Value = true
	}
	if no {
		f.Value = !f.Value
	}
	return
}

func (f *Flag) Parse(ctx Context) (err error) {
	for _, l := range f.Longs {
		p := LongParser{Long: l, CanUnary: true}
		if ctx.Match(&p) {
			return f.parseValue(false, p, ctx)
		}
		p.Long = "no-" + l
		if ctx.Match(&p) {
			return f.parseValue(true, p, ctx)
		}
	}
	for _, l := range f.Shorts {
		p := ShortParser{Short: l, CanUnary: true}
		if ctx.Match(&p) {
			return f.parseValue(false, p, ctx)
		}
		p.Prefix = '+'
		if ctx.Match(&p) {
			return f.parseValue(true, p, ctx)
		}
	}
	return noMatch
}
