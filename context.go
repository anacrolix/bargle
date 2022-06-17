package bargle

import (
	"fmt"
	"os"
	"strings"

	"github.com/anacrolix/generics"
)

type context struct {
	args     Args
	actions  *[]func() error
	deferred *[]func()
	tried    []ParamHelper
	helping  bool
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
	ctx.doHelpCommand()
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

func (ctx *context) doHelpCommand() {
	if ctx.Helping() {
		var help Help
		help.AddParams(ctx.tried...)
		help.Print(os.Stdout)
		ctx.Success()
	}
}

func (me *context) implicitHelp() bool {
	if me.Helping() {
		return false
	}
	help := Help{}
	help.AddParams(me.tried...)
	return me.matchAddTry(&help, false)
}

func (me *context) addTry(p Parser) {
	if ph, ok := p.(ParamHelper); ok {
		me.tried = append(me.tried, ph)
	}
}

func (me *context) Parse(p Parser) {
	args := me.args.Clone()
	me.addTry(p)
	if me.Helping() {
		return
	}
	err := p.Parse(me)
	if err == noMatch {
		if me.implicitHelp() {
			return
		}
	}
	if err != nil {
		var arg generics.Option[string]
		if args.Len() != 0 {
			arg.Set(args.Pop())
		}
		panic(controlError{parseError{
			inner: err,
			arg:   arg,
			param: p,
		}})
	}
}

func (me *context) matchAddTry(p Parser, addTry bool) bool {
	args := me.args.Clone()
	if addTry {
		me.addTry(p)
	}
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

func (me *context) Match(p Parser) bool {
	return me.matchAddTry(p, true)
}

func (me *context) MustParseOne(params ...Parser) {
	for _, p := range params {
		if me.Match(p) /*&& !me.Helping()*/ {
			return
		}
	}
	if me.Helping() {
		return
	}
	if me.implicitHelp() {
		return
	}
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
	for _, p := range ps {
		if me.Match(p) {
			goto start
		}
	}
	me.implicitHelp()
	if me.Helping() {
		return
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

func (me *context) Helping() bool {
	return me.helping
}

func (me *context) StartHelping() {
	me.helping = true
}
