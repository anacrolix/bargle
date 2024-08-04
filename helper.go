package bargle

type Helper interface {
	Arg
	Parsed(ParseAttempt)
	Helping() bool
	DoHelp(PrintHelpOpts)
	SetHelping()
}

type ParseAttempt struct {
	Arg     Arg
	Matched bool
}
