package bargle

import (
	"fmt"
	"io"
	"os"
	"sync"

	g "github.com/anacrolix/generics"
)

type builtinHelper struct {
	helping        bool
	globalArgs     map[Arg]struct{}
	unmatchedArgs  map[ArgType][]Arg
	parser         Arg
	initParserOnce sync.Once
	writer         io.Writer
	helpedCount    int
}

func (b *builtinHelper) initParser() {
	b.initParserOnce.Do(func() {
		b.parser = Long("help", BuiltinUnmarshaler(&b.helping))
	})
}

func (b *builtinHelper) ArgInfo() ArgInfo {
	b.initParser()
	return b.parser.ArgInfo()
}

func (b *builtinHelper) printArg(arg Arg) {
	if mv, ok := arg.(Metavar); ok {
		fmt.Fprintf(b.writer, "%s: ", mv.Metavar())
	}
	fmt.Fprintln(b.writer, g.ConvertToSliceOfAny(arg.ArgInfo().MatchingForms)...)
	descer, ok := arg.(ArgDescer)
	if ok {
		fmt.Fprintf(b.writer, "    %s\n", descer.ArgDesc())
	}
}

func (b *builtinHelper) globalArgsSlice() (slice []Arg) {
	for arg := range b.globalArgs {
		slice = append(slice, arg)
	}
	return
}

const noArgumentsExpectedHelp = "No arguments expected.\n"

type PrintHelpOpts struct {
	// Don't print the usage string which includes the program basename. Helpful for testing or
	// temporary binaries.
	NoPrintUsage bool
}

func (b *builtinHelper) DoHelp(opts PrintHelpOpts) {
	b.helpedCount++
	printedSomething := false
	if !opts.NoPrintUsage {
		fmt.Fprintf(b.writer, "Usage for %v:\n\n", os.Args[0])
	}
	printedSomething = b.printArgBlock(ArgTypeEnvVar, "Environment variables:", b.globalArgsSlice()) || printedSomething
	printedSomething = b.printArgBlock(ArgTypeSwitch, "Switches:", b.unmatchedArgs[ArgTypeSwitch]) || printedSomething
	printedSomething = b.printArgBlock(ArgTypeSwitch, "Positional:", b.unmatchedArgs[ArgTypePos]) || printedSomething
	if !printedSomething {
		fmt.Fprint(b.writer, noArgumentsExpectedHelp)
	}
}

func (b *builtinHelper) printArgBlock(argType ArgType, header string, args []Arg) bool {
	if len(args) == 0 {
		return false
	}
	fmt.Fprintln(b.writer, header)
	for _, arg := range args {
		fmt.Fprint(b.writer, "  ")
		b.printArg(arg)
	}
	return true
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
		g.MakeMapIfNil(&b.unmatchedArgs)
		b.unmatchedArgs[argType] = append(b.unmatchedArgs[argType], arg)
	}
	if attempt.Matched {
		b.unmatchedArgs = nil
	}
}

func (b *builtinHelper) Helping() bool {
	return b.helping
}

func (b *builtinHelper) SetHelping() {
	b.helping = true
}

var _ Helper = (*builtinHelper)(nil)
