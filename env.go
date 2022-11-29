package args

import "os"

func Env(key string, u Unmarshaler) env {
	return env{key, u}
}

type env struct {
	key string
	u   Unmarshaler
}

func (e env) Parse(ctx ParseContext) bool {
	value, ok := os.LookupEnv(e.key)
	if !ok {
		return false
	}
	return ctx.UnmarshalArg(e.u, value)
}

var _ Arg = env{}
