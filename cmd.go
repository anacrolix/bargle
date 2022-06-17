package bargle

type Command struct {
	Options     []Param
	Positionals []Param
	AfterParse  func() error
}
