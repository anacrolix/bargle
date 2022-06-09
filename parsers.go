package bargle

import (
	"fmt"
	"strings"

	"github.com/anacrolix/generics"
)

type LongParser struct {
	Long     string
	CanUnary bool
}

func (me LongParser) Parse(ctx Context) error {
	args := ctx.Args()
	next := args.Pop()
	if !strings.HasPrefix(next, "--") {
		return noMatch
	}
	before, after, found := strings.Cut(next[2:], "=")
	if me.CanUnary && found {
		args.Push(after)
	}
	if before != me.Long {
		return noMatch
	}
	return nil
}

type Subcommand struct {
	Commands map[string]ContextFunc
}

func (me Subcommand) Parse(ctx Context) error {
	cmd := ctx.Args().Pop()
	defer recoverType(func(err controlError) {
		panic(controlError{fmt.Errorf("%s: %w", cmd, err)})
	})
	me.Commands[cmd](ctx)
	return nil
}

type Choice[T any] struct {
	Choices map[string]T
	Target  *T
}

func (me Choice[T]) Parse(ctx Context) error {
	args := ctx.Args()
	var ok bool
	choice := args.Pop()
	*me.Target, ok = me.Choices[choice]
	if !ok {
		return noMatch
	}
	return nil
}

func (me Choice[T]) Add(name string, value T) {
	generics.MakeMapIfNil(&me.Choices)
	me.Choices[name] = value
}

func (me Choice[T]) SetDefault(arg string) {
	var ok bool
	*me.Target, ok = me.Choices[arg]
	if !ok {
		panic(fmt.Sprintf("no such choice: %q", arg))
	}
}

type positionalParserTarget interface {
	*string
}

type positionalParser[T positionalParserTarget] struct {
	target T
}

func (me positionalParser[T]) Parse(ctx Context) error {
	args := ctx.Args()
	next := args.Pop()
	if strings.HasPrefix(next, "-") {
		return noMatch
	}
	*me.target = next
	return nil
}

func Positional[T positionalParserTarget](target T) positionalParser[T] {
	return positionalParser[T]{target}
}
