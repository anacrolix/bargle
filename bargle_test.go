package bargle

import (
	"testing"

	qt "github.com/frankban/quicktest"
)

func TestParseFlagNoArgs(t *testing.T) {
	ctx := NewContext(nil)
	f := Flag{}
	f.AddLong("debug").AddShort('s')
	err := ctx.Run(func(ctx Context) {
		ctx.Try(&f)
	})
	qt.Assert(t, err, qt.IsNil)
}
