package bargle

import (
	"fmt"
	"os"
)

// Creates a new Parser bound to the system-level arguments.
func NewParser() *Parser {
	return &Parser{
		args: os.Args[1:],
		helper: &builtinHelper{
			writer: os.Stdout,
		},
	}
}

func NewParserNoArgs() *Parser {
	// TODO: Where helper?
	return &Parser{}
}

// A parser for a sequence of strings.
type Parser struct {
	args   []string
	err    error
	helper Helper

	posOnly bool
	// This is to prevent us adding the pseudo positional arg to help multiple times. It gets
	// cleared when an argument matches, so we can add it again once for the next series of parse
	// attempts.
	triedParsingPosOnly bool
}

// Parse the given parameter, if we're in the right state. Returns true if it matched, and sets an
// error if it matched and failed to unmarshal.
func (p *Parser) Parse(arg Arg) (matched bool) {
	return p.parseAndHelp(arg, true)
}

// Parse the given parameter, if we're in the right state. Returns true if it matched, and sets an
// error if it matched and failed to unmarshal.
func (p *Parser) parseAndHelp(arg Arg, addToHelp bool) (matched bool) {
	if addToHelp {
		defer func() {
			if p.helper != nil {
				p.helper.Parsed(ParseAttempt{
					Arg:     arg,
					Matched: matched,
				})
			}
			if matched {
				p.triedParsingPosOnly = false
			}
		}()
	}
	if p.err != nil {
		return false
	}
	if !p.posOnly || arg.ArgInfo().ArgType != ArgTypeSwitch {
		if p.doArgParse(arg) {
			return true
		}
	}
	p.tryParseHelp()
	return false
}

func (p *Parser) parsePseudoPosOnly() {
	if !p.triedParsingPosOnly {
		// This needs to be set before parsing because it recursively calls through to here again.
		p.triedParsingPosOnly = true
		// Currently disabling the -- in help.
		if p.parseAndHelp(pseudoPosOnly{}, false) {
			p.posOnly = true
		}
	}
}

// This parses without checking for existing Parser error or sending messages.
func (p *Parser) doArgParse(arg Arg) (matched bool) {
	// This is done here so that it's never missed, including when we're exclusively just checking
	// if help is requested.
	p.parsePseudoPosOnly()
	pc := parseContext{
		args:    p.args,
		posOnly: p.posOnly,
	}
	parsed := arg.Parse(&pc)
	if parsed {
		p.args = pc.args
	}
	p.err = pc.err
	return parsed
}

// Return existing Parser error, or set one based on how many arguments remain.
func (p *Parser) Fail() error {
	if p.err == nil {
		if len(p.args) == 0 {
			p.err = ErrExpectedArguments
		} else {
			p.FailIfArgsRemain()
		}
	}
	return p.err
}

func (p *Parser) tryParseHelp() {
	if p.err != nil {
		return
	}
	if p.helper == nil {
		return
	}
	if p.helper.Helping() {
		return
	}
	p.doArgParse(p.helper)
}

func (p *Parser) DoHelpIfHelping() {
	p.DoHelpIfHelpingOpts(PrintHelpOpts{})
}

func (p *Parser) DoHelpIfHelpingOpts(opts PrintHelpOpts) {
	if p.helper.Helping() {
		p.helper.DoHelp(opts)
	}
}

// This asserts that no arguments remain, and if they do sets an appropriate error. You would call
// this when you're ready to start actual work after parsing, and then check Parser.Ok().
func (p *Parser) FailIfArgsRemain() {
	if p.err != nil {
		return
	}
	// Check if help is wanted in case no argument parsing attempts have occurred yet.
	p.tryParseHelp()
	if len(p.args) != 0 {
		p.err = fmt.Errorf("unused argument: %q", p.args[0])
	}
}

// Removes and returns all remaining unused arguments. This might be used to pass handling on to
// something else, or to process the rest of the arguments manually.
func (p *Parser) PopAll() (all []string) {
	if p.err != nil {
		return nil
	}
	all = p.args
	p.args = nil
	return
}

// Returns false if there's an error, or help has been issued. You would normally then return
// Parser.Err(), which may be nil.
func (p *Parser) Ok() bool {
	return !p.helper.Helping() && p.Err() == nil
}

// Returns any error the Parser has encountered. Usually this is the first error and blocks further
// parsing until it's convenient to handle it.
func (p *Parser) Err() error {
	return p.err
}

// Applies the given arguments through the unmarshaller. Returns false if an error occurred. TODO:
// This doesn't look completed.
func (p *Parser) SetDefault(u Unmarshaler, args ...string) bool {
	if p.err != nil {
		return false
	}
	pc := parseContext{
		args:    args,
		posOnly: true,
	}
	return pc.Unmarshal(u)
}

func (p *Parser) SetError(err error) {
	if p.err == nil {
		p.err = err
	}
}

func (p *Parser) SetArgs(args ...string) {
	p.args = args
}

func (p *Parser) SetHelper(helper Helper) {
	p.helper = helper
}

func (p *Parser) SetHelping() {
	p.helper.SetHelping()
}
