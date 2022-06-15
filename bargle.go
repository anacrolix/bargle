package bargle

type Parser interface {
	Parse(ctx Context) error
}

type ParamParser interface {
	Parser
	ParamHelper
}
