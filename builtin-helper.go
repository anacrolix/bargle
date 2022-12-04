package args

import (
	"fmt"
	"sync"

	"github.com/anacrolix/generics"
)

type builtinHelper struct {
	helping        bool
	globalArgs     map[Arg]struct{}
	unmatchedArgs  map[ArgType][]Arg
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

func (b *builtinHelper) globalArgsSlice() (slice []Arg) {
	for arg := range b.globalArgs {
		slice = append(slice, arg)
	}
	return
}

func (b *builtinHelper) DoHelp() {
	b.printArgBlock(ArgTypeEnvVar, "Environment variables:", b.globalArgsSlice())
	b.printArgBlock(ArgTypeSwitch, "Switches:", b.unmatchedArgs[ArgTypeSwitch])
	b.printArgBlock(ArgTypeSwitch, "Positional:", b.unmatchedArgs[ArgTypePos])
}

func (b *builtinHelper) printArgBlock(argType ArgType, header string, args []Arg) {
	if len(args) == 0 {
		return
	}
	fmt.Println(header)
	for _, arg := range args {
		fmt.Print("  ")
		b.printArg(arg)
	}
}

func (b *builtinHelper) Parse(ctx ParseContext) bool {
	b.initParser()
	return b.parser.Parse(ctx)
}

func (b *builtinHelper) Parsed(attempt ParseAttempt) {
	arg := attempt.Arg
	argInfo := arg.ArgInfo()
	argType := argInfo.ArgType
	if argInfo.Global {
		if b.globalArgs == nil {
			b.globalArgs = make(map[Arg]struct{})
		}
		b.globalArgs[arg] = struct{}{}
	} else {
		generics.MakeMapIfNil(&b.unmatchedArgs)
		b.unmatchedArgs[argType] = append(b.unmatchedArgs[argType], arg)
	}
	if attempt.Matched {
		b.unmatchedArgs = nil
	}
}

func (b *builtinHelper) Helping() bool {
	return b.helping
}

var _ Helper = (*builtinHelper)(nil)
