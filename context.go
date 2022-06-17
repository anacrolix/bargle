package bargle

import (
	"errors"
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

func (ctx *context) Run(cmd Command) (err error) {
	defer func() {
		for i := range *ctx.deferred {
			(*ctx.deferred)[len(*ctx.deferred)-1-i]()
		}
	}()
	err = ctx.runCommand(cmd)
	if err != nil {
		return
	}
	if ctx.args.Len() != 0 {
		return fmt.Errorf("%v unused args, starting with %q", ctx.args.Len(), ctx.args.Pop())
	}
	for _, f := range *ctx.actions {
		err = f()
		if err != nil {
			break
		}
	}
	return
}

func (ctx *context) runCommand(cmd Command) error {
options:
	matches := make([]MatchResult, 0, len(cmd.Options))
	for _, opt := range cmd.Options {
		mr := ctx.Match(opt)
		if mr.Matched().Ok {
			matches = append(matches, mr)
		}
	}
	switch len(matches) {
	case 1:
		err := ctx.runMatchResult(matches[0])
		if err != nil {
			return err
		}
		goto options
	default:
		return errors.New("matched multiple options")
	case 0:
	}
	if !ctx.Done() {
		for _, pos := range cmd.Positionals {
			mr := ctx.Match(pos)
			if !mr.Matched().Ok {
				return fmt.Errorf("%v is next but couldn't match", pos)
			}
			err := ctx.runMatchResult(mr)
			if err != nil {
				return err
			}
			goto options
		}
	}
	*ctx.actions = append(*ctx.actions, cmd.AfterParse)
	return nil
}

func (ctx *context) runMatchResult(mr MatchResult) error {
	ctx.args = mr.Args()
	err := mr.Parse(ctx)
	p := mr.Param()
	if err != nil {
		return fmt.Errorf("parsing %q: %w", mr.Matched().Unwrap(), err)
	}
	sub := p.Subcommand()
	if sub.Ok {
		err := ctx.runCommand(sub.Value)
		if err != nil {
			return err
		}
	}
	return nil
}

func (ctx *context) doHelpCommand() {
	if ctx.Helping() {
		var help Help
		help.AddParams(ctx.tried...)
		help.Print(os.Stdout)
		ctx.Success()
	}
}

func (me *context) Match(m Matcher) (ret MatchResult) {
	return m.Match(me.args.Clone())
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

func (me *context) MissingArgument(name string) {
	me.Fail(fmt.Errorf("missing argument: %s", name))
}

func (me *context) Success() {
	panic(success{})
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
