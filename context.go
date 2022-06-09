package bargle

import (
	"fmt"
)

type context struct {
	args     Args
	actions  []func() error
	deferred []func()
}

type Context = *context

type ContextFunc func(Context)

func NewContext(args []string) Context {
	return &context{
		args: NewArgs(args),
	}
}

func (me *context) Args() Args {
	return me.args
}

func (me *context) Run(p Parser) (err error) {
	defer recoverType(func(ce controlError) {
		err = ce
	})
	return p.Parse(me)
}

func (me *context) Parse(p Parser) {
	args := me.args.Clone()
	err := p.Parse(me)
	if err != nil {
		panic(controlError{fmt.Errorf("parsing %q: %w", args.Pop(), err)})
	}
}

func (me *context) Match(p Parser) bool {
	args := me.args.Clone()
	err := p.Parse(me)
	switch err {
	case noMatch:
		me.args = args
		return false
	case nil:
		return true
	default:
		panic(err)
	}
}

func (me *context) Done() bool {
	return me.args.Len() == 0
}

func (me *context) Defer(f func()) {
	me.deferred = append(me.deferred, f)
}

func (me *context) AfterParse(f func() error) {
	me.actions = append(me.actions, f)
}

func (me *context) Unhandled() {
	panic(controlError{unhandledErr{me.args.Pop()}})
}

func (me *context) Fail(err error) {
	panic(controlError{userError(err)})
}

func (me *context) ParseUntilDone(ps ...Parser) {
	for !me.Done() {
		for _, p := range ps {
			if me.Match(p) {
				continue
			}
		}
		me.Unhandled()
	}
}
