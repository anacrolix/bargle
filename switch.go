package bargle

import (
	"errors"
)

// A switch is an option that takes no values and has no negative form.
type Switch struct {
	optionDefaults
	Longs          []string
	Shorts         []rune
	Desc           string
	AfterParseFunc AfterParseParamFunc
}

func (f Switch) Parse(args Args) error {
	return errors.New("switches do not take values")
}

func (f Switch) Match(args Args) (mr MatchResult) {
	for _, l := range f.Longs {
		_args := args.Clone()
		p := LongParser{Long: l, CanUnary: false}
		if p.Match(_args) {
			return f.matchResult(_args, args.Pop())
		}
	}
	for _, l := range f.Shorts {
		_args := args.Clone()
		p := ShortParser{Short: l, CanUnary: false}
		if p.Match(_args) {
			return f.matchResult(_args, args.Pop())
		}
	}
	return noMatch
}

func (f Switch) matchResult(args Args, matchedArg string) MatchResult {
	mr := matchedNoParse{
		baseMatchResult: baseMatchResult{
			args:  args,
			param: f,
			match: matchedArg,
		},
	}
	return mr
}

func (f Switch) Help() ParamHelp {
	ph := ParamHelp{
		Description: f.Desc,
	}
	for _, l := range f.Longs {
		ph.Forms = append(ph.Forms, "--"+l)
	}
	for _, s := range f.Shorts {
		ph.Forms = append(ph.Forms, "-"+string(s))
	}
	return ph
}

func (f Switch) AfterParse(ctx Context) error {
	if f.AfterParseFunc == nil {
		return f.optionDefaults.AfterParse(ctx)
	}
	return f.AfterParseFunc(ctx)
}
