package bargle

import (
	"strconv"

	"github.com/anacrolix/generics"
)

// A flag is an optional argument that results in a boolean value. It has a negative form, and can
// parse from a bound value in the same argument (with '=').
type Flag struct {
	optionDefaults
	target flagTarget
	switchesOpts
}

type boolTarget struct {
	target *bool
}

func (b boolTarget) Set(value bool) {
	*b.target = value
}

type boolPtrTarget struct {
	target **bool
}

func (b boolPtrTarget) Set(value bool) {
	if *b.target == nil {
		*b.target = new(bool)
	}
	**b.target = value
}

func NewFlag(target any) *flagMaker {
	if target == nil {
		target = new(bool)
	}
	return &flagMaker{
		target: func() flagTarget {
			switch typedTarget := target.(type) {
			case *bool:
				return boolTarget{typedTarget}
			case **bool:
				return boolPtrTarget{typedTarget}
			}
			panic("unsupported target type")
		}(),
	}
}

func (f Flag) Init() error {
	return nil
}

func (f Flag) Parse(args Args) (err error) {
	value, err := strconv.ParseBool(args.Pop())
	if err != nil {
		return
	}
	f.target.Set(value)
	return nil
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
		target: f.target,
		no:     no,
	}
	if us.GotValue() {
		mr.value = generics.Some(args.Pop())
	}
	return mr
}

type flagTarget interface {
	Set(bool)
}

type flagMatchResult struct {
	baseMatchResult
	value  generics.Option[string]
	target flagTarget
	no     bool
}

func (f flagMatchResult) Parse(Args) (err error) {
	value, err := func() (value bool, err error) {
		if f.value.Ok {
			return strconv.ParseBool(f.value.Value)
		} else {
			return true, nil
		}
	}()
	if err != nil {
		return
	}
	if f.no {
		value = !value
	}
	f.target.Set(value)
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
	target flagTarget
}

func (m flagMaker) Make() Flag {
	return Flag{
		target:       m.target,
		switchesOpts: m.switchesMaker.switchesOpts,
	}
}
