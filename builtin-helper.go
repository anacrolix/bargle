package args

import (
	"fmt"
	"sync"

	"github.com/anacrolix/generics"
)

type builtinHelper struct {
	helping        bool
	globalArgs     map[Arg]struct{}
	unmatchedArgs  []Arg
	parser         Arg
	initParserOnce sync.Once
}

func (b *builtinHelper) initParser() {
	b.initParserOnce.Do(func() {
		b.parser = Long("help", BuiltinUnmarshalerFromAny(&b.helping))
	})
}

func (b *builtinHelper) ArgInfo() ArgInfo {
	b.initParser()
	return b.parser.ArgInfo()
}

func (b *builtinHelper) printArg(arg Arg) {
	fmt.Println(generics.ConvertToSliceOfAny(arg.ArgInfo().MatchingForms)...)
}

func (b *builtinHelper) DoHelp() {
	for arg := range b.globalArgs {
		b.printArg(arg)
	}
	for _, arg := range b.unmatchedArgs {
		b.printArg(arg)
	}
}

func (b *builtinHelper) Parse(ctx ParseContext) bool {
	b.initParser()
	return b.parser.Parse(ctx)
}

func (b *builtinHelper) Parsed(attempt ParseAttempt) {
	if attempt.Arg.ArgInfo().Global {
		if b.globalArgs == nil {
			b.globalArgs = make(map[Arg]struct{})
		}
		b.globalArgs[attempt.Arg] = struct{}{}
	} else {
		b.unmatchedArgs = append(b.unmatchedArgs, attempt.Arg)
	}
	if attempt.Matched {
		b.unmatchedArgs = b.unmatchedArgs[:0]
	}
}

func (b *builtinHelper) Helping() bool {
	return b.helping
}

var _ Helper = (*builtinHelper)(nil)
