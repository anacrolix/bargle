package bargle

import (
	"fmt"
	"strings"

	g "github.com/anacrolix/generics"
)

// An experiment with maybe allowing separators and key styling to be propagated from a config down
// the track.
func LongElems(u Unmarshaler, firstElem string, elems ...string) long {
	return Long(strings.Join(append(g.Singleton(firstElem), elems...), "-"), u)
}

// Don't include the prefix "--" for now.
func Long(key string, u Unmarshaler) long {
	return long{key, u}
}

type long struct {
	key string
	u   Unmarshaler
}

func (me long) ArgInfo() ArgInfo {
	argTypes := strings.Join(me.u.ArgTypes(), " ")
	return ArgInfo{
		MatchingForms: g.Singleton(fmt.Sprintf("--%[1]s=%[2]v, --%[1]s %[2]v", me.key, argTypes)),
		ArgType:       ArgTypeSwitch,
	}
}

func (me long) Parse(ctx ParseContext) bool {
	arg, ok := ctx.Pop()
	if !ok {
		return false
	}
	// TODO: Use strings.CutPrefix in go1.20+
	if !strings.HasPrefix(arg, "--") {
		return false
	}
	arg = arg[2:]
	i := strings.IndexByte(arg, '=')
	key := arg
	if i != -1 {
		key = key[:i]
	}
	if key != me.key {
		return false
	}
	if i == -1 {
		return ctx.Unmarshal(me.u)
	}
	return ctx.UnmarshalArg(me.u, arg[i+1:])
}
