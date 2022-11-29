package args

import "fmt"

type Parser struct {
	args []string
	err  error
}

func (p *Parser) Parse(arg Arg) bool {
	if p.err != nil {
		return false
	}
	pc := parseContext{
		args: p.args,
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
	if p.err != nil {
		return p.err
	}
	if len(p.args) == 0 {
		return ErrExpectedArguments
	}
	p.FailIfArgsRemain()
	return p.err
}

func (p *Parser) FailIfArgsRemain() {
	if p.err != nil {
		return
	}
	if len(p.args) != 0 {
		p.err = fmt.Errorf("unused argument: %q", p.args[0])
	}
}

func (p *Parser) PopAll() (all []string) {
	if p.err != nil {
		return nil
	}
	all = p.args
	p.args = nil
	return
}

func (p *Parser) Err() error {
	return p.err
}

func (p *Parser) SetDefault(u Unmarshaler, args ...string) bool {
	if p.err != nil {
		return false
	}
	pc := parseContext{
		args: args,
	}
	return pc.Unmarshal(u)
}

func (p *Parser) SetError(err error) {
	if p.err == nil {
		p.err = err
	}
}
