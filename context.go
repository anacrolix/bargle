package bargle

import (
	"fmt"
	"strings"

	"github.com/anacrolix/generics"
)

type context struct {
	args     Args
	actions  *[]func() error
	deferred *[]func()
	tried    []ParamHelper
}

type Context = *context

type ContextFunc func(ctx Context)

func NewContext(args []string) Context {
	ctx := &context{
		args: NewArgs(args),
	}
	generics.InitNew(&ctx.actions)
	generics.InitNew(&ctx.deferred)
	return ctx
}

func (me *context) Args() Args {
	return me.args
}

func (ctx *context) Run(f ContextFunc) (err error) {
	defer recoverType(func(ce controlError) {
		err = ce
	})
	defer recoverType(func(success) {})
	defer func() {
		for i := range *ctx.deferred {
			(*ctx.deferred)[len(*ctx.deferred)-1-i]()
		}
	}()
	for {
		again := false
		func() {
			defer recoverType(func(tried) {
				again = true
			})
			f(ctx)
		}()
		if !again {
			break
		}
		ctx.tried = nil
	}
	if ctx.args.Len() > 0 {
		ctx.implicitHelp()
		err = unhandledErr{ctx.args.Pop()}
		return
	}
	for _, f := range *ctx.actions {
		err = f()
		if err != nil {
			break
		}
	}
	return
}

func (me *context) implicitHelp() bool {
	help := Help{}
	help.AddParams(me.tried...)
	return me.Match(&help)
}

func (me *context) addTry(p Parser) {
	if ph, ok := p.(ParamHelper); ok {
		me.tried = append(me.tried, ph)
	}
}

func (me *context) Parse(p Parser) {
	args := me.args.Clone()
	me.addTry(p)
	err := p.Parse(me)
	if err == noMatch {
		me.implicitHelp()
	}
	if err != nil {
		panic(controlError{fmt.Errorf("parsing %q: %w", args.Pop(), err)})
	}
}

func (me *context) Match(p Parser) bool {
	args := me.args.Clone()
	me.addTry(p)
	err := p.Parse(me)
	switch err {
	case noMatch:
		me.args = args
		return false
	case nil:
		return true
	default:
		panic(controlError{err})
	}
}

func (me *context) MustParseOne(params ...Parser) {
	for _, p := range params {
		if me.Match(p) {
			return
		}
	}
	me.implicitHelp()
	me.Unhandled()
}

func (me *context) Done() bool {
	return me.args.Len() == 0
}

func (me *context) Defer(f func()) {
	*me.deferred = append(*me.deferred, f)
}

func (me *context) AfterParse(f func() error) {
	*me.actions = append(*me.actions, f)
}

func (me *context) Unhandled() {
	if me.args.Len() == 0 {
		panic(controlError{parseFailure{}})
	}
	panic(controlError{unhandledErr{me.args.Pop()}})
}

func (me *context) Fail(err error) {
	panic(controlError{userError(err)})
}

func (me *context) ParseUntilDone(ps ...Parser) {
start:
	for me.Done() {
		return
	}
	for _, p := range ps {
		if me.Match(p) {
			goto start
		}
	}
	me.Unhandled()
}

func (me *context) MissingArgument(name string) {
	me.Fail(fmt.Errorf("missing argument: %s", name))
}

func (me *context) Success() {
	panic(success{})
}

func (me *context) Try(p Parser) {
	if me.Match(p) {
		panic(tried{})
	}
}

// Returns whether the next arg can be parsed as positional. This could allow to handle -- and drop
// into positional only arguments.
func (me *context) MatchPos() bool {
	if me.args.Len() == 0 {
		return true
	}
	args := me.args.Clone()
	if strings.HasPrefix(args.Pop(), "-") {
		return false
	}
	return true
}

func (me *context) NewChild() Context {
	child := *me
	child.tried = nil
	return &child
}
