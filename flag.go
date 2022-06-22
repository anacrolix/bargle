package bargle

import (
	"strconv"

	"github.com/anacrolix/generics"
)

// A flag is an optional argument that results in a boolean value. It has a negative form, and can
// parse from a bound value in the same argument (with '=').
type Flag struct {
	optionDefaults
	Value *bool
	switchesOpts
}

func NewFlag(target *bool) *flagMaker {
	if target == nil {
		target = new(bool)
	}
	return &flagMaker{
		target: target,
	}
}

func (f Flag) Init() error {
	return nil
}

func (f Flag) Parse(args Args) (err error) {
	*f.Value, err = strconv.ParseBool(args.Pop())
	return
}

func (f Flag) Help() ParamHelp {
	ph := ParamHelp{
		Values: "[true|1|false|0]",
	}
	for _, l := range f.longs {
		ph.Forms = append(ph.Forms, "--"+l, "--no-"+l)
	}
	for _, s := range f.shorts {
		ph.Forms = append(ph.Forms, "-"+string(s), "+"+string(s))
	}
	return ph
}

func (f *Flag) matchResult(no bool, us UnarySwitch, args Args, matchedArg string) MatchResult {
	mr := flagMatchResult{
		baseMatchResult: baseMatchResult{
			args:  args,
			param: f,
			match: matchedArg,
		},
		target: f.Value,
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

func (f flagMatchResult) Parse(Args) (err error) {
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

func (f Flag) Match(args Args) (mr MatchResult) {
	for _, l := range f.longs {
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
	for _, l := range f.shorts {
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

type flagMaker struct {
	switchesMaker
	target *bool
}

func (m flagMaker) Make() Flag {
	return Flag{
		Value:        m.target,
		switchesOpts: m.switchesMaker.switchesOpts,
	}
}
