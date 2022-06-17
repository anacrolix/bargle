package bargle

import (
	"strconv"

	"github.com/anacrolix/generics"
)

type Flag struct {
	optionDefaults
	Value  bool
	Longs  []string
	Shorts []rune
}

func (f *Flag) Help() ParamHelp {
	ph := ParamHelp{
		Values: "[true|1|false|0]",
	}
	for _, l := range f.Longs {
		ph.Forms = append(ph.Forms, "--"+l, "--no-"+l)
	}
	for _, s := range f.Shorts {
		ph.Forms = append(ph.Forms, "-"+string(s), "+"+string(s))
	}
	return ph
}

func (f *Flag) AddLong(l string) *Flag {
	f.Longs = append(f.Longs, l)
	return f
}

func (f *Flag) AddShort(s rune) *Flag {
	f.Shorts = append(f.Shorts, s)
	return f
}

func (f *Flag) matchResult(no bool, us UnarySwitch, args Args) MatchResult {
	mr := flagMatchResult{
		args:   args,
		target: &f.Value,
		no:     no,
	}
	if us.GotValue() {
		mr.value = generics.Some(args.Pop())
	}
	return mr
}

type flagMatchResult struct {
	matched string
	args    Args
	value   generics.Option[string]
	target  *bool
	no      bool
}

func (f flagMatchResult) Matched() generics.Option[string] {
	return generics.Some(f.matched)
}

func (f flagMatchResult) Args() Args {
	return f.args
}

func (f flagMatchResult) Parse(ctx Context) (err error) {
	if f.value.Ok {
		*f.target, err = strconv.ParseBool(f.value.Value)
	} else {
		*f.target = true
	}
	if f.no {
		*f.target = !*f.target
	}
	return

}

func (f flagMatchResult) Param() Param {
	//TODO implement me
	panic("implement me")
}

var _ MatchResult = flagMatchResult{}

func (f *Flag) Match(args Args) (mr MatchResult) {
	for _, l := range f.Longs {
		_args := args.Clone()
		p := LongParser{Long: l, CanUnary: true}
		if p.Match(_args) {
			return f.matchResult(false, p, _args)
		}
		_args = args.Clone()
		p.Long = "no-" + l
		if p.Match(_args) {
			return f.matchResult(true, p, _args)
		}
	}
	for _, l := range f.Shorts {
		_args := args.Clone()
		p := ShortParser{Short: l, CanUnary: true}
		if p.Match(_args) {
			return f.matchResult(false, p, _args)
		}
		_args = args.Clone()
		p.Prefix = '+'
		if p.Match(_args) {
			return f.matchResult(true, p, _args)
		}
	}
	return noMatch
}
