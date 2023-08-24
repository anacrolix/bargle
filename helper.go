package bargle

type Helper interface {
	Arg
	Parsed(ParseAttempt)
	Helping() bool
	DoHelp()
	SetHelping()
}

type ParseAttempt struct {
	Arg     Arg
	Matched bool
}
