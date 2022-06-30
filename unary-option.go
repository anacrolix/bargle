package bargle

import (
	"fmt"

	"github.com/anacrolix/generics"
)

type UnaryOption struct {
	optionDefaults
	unaryOptionOpts
	switchesOpts
	parsed bool
	desc   string
}

type unaryOptionOpts struct {
	Required bool
	Value    UnaryUnmarshaler
}

func (me *unaryOptionMaker) SetRequired() *unaryOptionMaker {
	me.Required = true
	return me
}

func (me *UnaryOption) Init() error {
	return nil
}

func (me *UnaryOption) switchForms() (ret []string) {
	for _, l := range me.longs {
		ret = append(ret, "--"+l)
	}
	for _, s := range me.shorts {
		ret = append(ret, "-"+string(s))
	}
	return
}

func (me *UnaryOption) Help() ParamHelp {
	return ParamHelp{
		Forms:       me.switchForms(),
		Values:      me.Value.TargetHelp(),
		Description: me.desc,
	}
}

func (me *UnaryOption) AfterParse(Context) error {
	me.parsed = true
	return nil
}

func (me *UnaryOption) Satisfied() bool {
	return !me.Required || me.parsed
}

func (me *UnaryOption) Match(args Args) MatchResult {
	return me.matchSwitch(args)
}

type unaryMatchResult struct {
	baseMatchResult
	u UnaryUnmarshaler
}

func (me unaryMatchResult) Parse(args Args) error {
	if args.Len() == 0 {
		return missingArgument
	}
	arg := args.Pop()
	err := me.u.UnaryUnmarshal(arg)
	if err != nil {
		err = fmt.Errorf("unmarshalling %q: %w", arg, err)
	}
	return err
}

func (me *UnaryOption) matchSwitch(args Args) MatchResult {
	for _, l := range me.longs {
		_args := args.Clone()
		gv := &LongParser{Long: l, CanUnary: true}
		if gv.Match(_args) {
			return unaryMatchResult{baseMatchResult{_args, me, args.Clone().Pop()}, me.Value}
		}
	}
	for _, s := range me.shorts {
		_args := args.Clone()
		gv := &ShortParser{Short: s, CanUnary: true}
		if gv.Match(_args) {
			return unaryMatchResult{baseMatchResult{_args, me, args.Clone().Pop()}, me.Value}
		}
	}
	return noMatch
}

func (me *UnaryOption) Parse(args Args) error {
	return me.Value.UnaryUnmarshal(args.Pop())
}

type unaryOptionMaker struct {
	unaryOptionOpts
	switchesMaker
	default_ generics.Option[string]
	desc     string
}

func (me *unaryOptionMaker) SetDefault(default_ string) {
	me.default_ = generics.Some(default_)
}

func NewUnaryOption(u UnaryUnmarshaler) *unaryOptionMaker {
	if u == nil {
		panic("nil UnaryUnmarshaler")
	}
	ret := &unaryOptionMaker{}
	ret.Value = u
	return ret
}

func (me *unaryOptionMaker) Make() *UnaryOption {
	if me.default_.Ok {
		err := me.Value.UnaryUnmarshal(me.default_.Value)
		if err != nil {
			err = fmt.Errorf("unmarshaling default: %w", err)
			panic(err)
		}
	}
	return &UnaryOption{
		unaryOptionOpts: me.unaryOptionOpts,
		switchesOpts:    me.switchesOpts,
		desc:            me.desc,
	}
}

func (me *unaryOptionMaker) Description(desc string) {
	me.desc = desc
}
