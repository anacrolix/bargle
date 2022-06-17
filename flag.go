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

func (f *Flag) matchResult(no bool, us UnarySwitch, args Args, matchedArg string) MatchResult {
	mr := flagMatchResult{
		baseMatchResult: baseMatchResult{
			args:  args,
			param: f,
			match: matchedArg,
		},
		target: &f.Value,
		no:     no,
	}
	if us.GotValue() {
		mr.value = generics.Some(args.Pop())
	}
	return mr
}

type flagMatchResult struct {
	baseMatchResult
	value  generics.Option[string]
	target *bool
	no     bool
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

var _ MatchResult = flagMatchResult{}

func (f *Flag) Match(args Args) (mr MatchResult) {
	for _, l := range f.Longs {
		_args := args.Clone()
		p := LongParser{Long: l, CanUnary: true}
		if p.Match(_args) {
			return f.matchResult(false, p, _args, args.Pop())
		}
		_args = args.Clone()
		p.Long = "no-" + l
		if p.Match(_args) {
			return f.matchResult(true, p, _args, args.Pop())
		}
	}
	for _, l := range f.Shorts {
		_args := args.Clone()
		p := ShortParser{Short: l, CanUnary: true}
		if p.Match(_args) {
			return f.matchResult(false, p, _args, args.Pop())
		}
		_args = args.Clone()
		p.Prefix = '+'
		if p.Match(_args) {
			return f.matchResult(true, p, _args, args.Pop())
		}
	}
	return noMatch
}
