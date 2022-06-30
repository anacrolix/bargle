package bargle

import (
	"fmt"
)

type Positional struct {
	posDefaults
	Value          UnaryUnmarshaler
	Name           string
	Desc           string
	ok             bool
	AfterParseFunc AfterParseParamFunc
}

func (me *Positional) Init() error {
	return nil
}

func (me *Positional) Parse(args Args) error {
	return me.Value.UnaryUnmarshal(args.Pop())
}

func (me *Positional) Match(args Args) MatchResult {
	if !me.Value.Matching() {
		return noMatch
	}
	mr := unaryMatchResult{
		u: me.Value,
	}
	mr.args = args.Clone()
	mr.param = me
	mr.match = args.Pop()
	return mr
}

func (me *Positional) AfterParse(ctx Context) error {
	me.ok = true
	if me.AfterParseFunc != nil {
		return me.AfterParseFunc(ctx)
	}
	return nil
}

func (me *Positional) Satisfied() bool {
	return me.ok
}

func (me *Positional) Help() ParamHelp {
	return ParamHelp{
		Forms:       []string{fmt.Sprintf("<%v>", me.Name)},
		Values:      me.Value.TargetHelp(),
		Description: me.Desc,
	}
}

//func (me *Positional[T]) Parse(ctx Context) error {
//	if ctx.Args().Len() == 0 {
//		return missingArgument
//	}
//	if !ctx.MatchPos() {
//		return noMatch
//	}
//	return doUnaryUnmarshal(ctx.Args().Pop(), &me.Value, me.Value)
//}

func NewPositional(u UnaryUnmarshaler) *Positional {
	return &Positional{Value: u, Name: "arg"}
}
