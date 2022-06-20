package bargle

import (
	"bytes"
	"errors"
	"fmt"
	"strings"

	"github.com/anacrolix/generics"
)

type context struct {
	args     Args
	actions  *[]func() error
	deferred *[]func()
	exitCode generics.Option[int]
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
	err = cmd.Init()
	if err != nil {
		return fmt.Errorf("initing command: %w", err)
	}
	err = ctx.runCommand(cmd)
	if err != nil {
		return withExitCode(2, err)
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

func (ctx *context) tryAppendMatch(matches *[]MatchResult, p Param) {
	mr := ctx.Match(p)
	if mr.Matched().Ok {
		*matches = append(*matches, mr)
	}
}

func (ctx *context) runCommand(cmd Command) error {
	ranSubCmd := false
options:
	if ctx.exitCode.Ok {
		return nil
	}
	matches := make([]MatchResult, 0, len(cmd.Options))
	for _, opt := range cmd.Options {
		ctx.tryAppendMatch(&matches, opt)
	}
	if len(matches) == 0 {
		// We only try this if an existing option in the Command hasn't already matched.
		ctx.tryAppendMatch(&matches, Switch{
			Longs:  []string{"help"},
			Shorts: []rune{'h'},
			Desc:   "help/usage",
			AfterParseFunc: func(ctx Context) error {
				printCommandHelp(cmd.Help(), false)
				ctx.Success()
				return nil
			},
		})
	}
	switch len(matches) {
	case 1:
		err := ctx.runMatchResult(matches[0], &ranSubCmd)
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
			mr := ctx.MatchPos(pos)
			if !mr.Matched().Ok {
				continue
			}
			err := ctx.runMatchResult(mr, &ranSubCmd)
			if err != nil {
				return err
			}
			goto options
		}
		return fmt.Errorf("unhandled arg: %q", ctx.args.Pop())
	}
	err := ctx.assertCommandSatisfied(cmd)
	if err != nil {
		return err
	}
	if !ranSubCmd {
		if cmd.DefaultAction != nil {
			*ctx.actions = append(*ctx.actions, cmd.DefaultAction)
		} else {
			return errors.New("no subcommand invoked and no default action")
		}
	}
	return nil
}

func (ctx *context) assertCommandSatisfied(cmd Command) error {
	for _, p := range cmd.AllParams() {
		if !p.Satisfied() {
			var buf bytes.Buffer
			hw := HelpWriter{w: &buf}.Indented()
			p.Help().Write(hw)
			return paramError{fmt.Sprintf("unsatisfied param:\n%s", buf.Bytes())}
		}
	}
	return nil
}

func (ctx *context) runMatchResult(mr MatchResult, ranSubCmd *bool) error {
	// To get here we should have checked that Matched is Ok. Let's unwrap anyway to assert that
	// behaviour.
	matchedArg := mr.Matched().Unwrap()
	ctx.args = mr.Args()
	err := mr.Parse(ctx.args)
	p := mr.Param()
	if err != nil {
		return fmt.Errorf("parsing %q: %w", matchedArg, err)
	}
	err = p.AfterParse(ctx)
	if err != nil {
		return fmt.Errorf("running after parse for %q: %w", matchedArg, err)
	}
	sub := p.Subcommand()
	if sub.Ok {
		err := ctx.runCommand(sub.Value)
		if err != nil {
			return err
		}
		*ranSubCmd = true
	}
	return nil
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
	me.exitCode = generics.Some(0)
}

// Returns whether the next arg can be parsed as positional. This could allow to handle -- and drop
// into positional only arguments.
func (me *context) MatchPos(p Param) MatchResult {
	if me.args.Len() == 0 {
		return noMatch
	}
	args := me.args.Clone()
	if strings.HasPrefix(args.Pop(), "-") {
		return noMatch
	}
	return me.Match(p)
}

func (me *context) ExitCode() int {
	return me.exitCode.UnwrapOr(0)
}
