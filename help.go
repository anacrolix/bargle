package bargle

import (
	"fmt"
	"io"
)

type Help struct {
	params []Parser
}

func (me *Help) AddParams(params ...Parser) {
	me.params = append(me.params, params...)
}

func (me Help) Print(w io.Writer) {
	for _, p := range me.params {
		fmt.Fprintf(w, "%v\n", p)
	}
}
