package args

type Helper interface {
	Arg
	Parsed(ParseAttempt)
	Helping() bool
	DoHelp()
}

type ParseAttempt struct {
	Arg     Arg
	Matched bool
}
