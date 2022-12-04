package args

import (
	"fmt"
	"os"

	. "github.com/anacrolix/generics"
)

func Env(key string, u Unmarshaler) *env {
	return &env{key, u, ArgInfo{
		MatchingForms: Singleton(fmt.Sprintf(`%s="$value"`, key)),
		ArgType:       ArgTypeEnvVar,
		Global:        true,
	}}
}

type env struct {
	key  string
	u    Unmarshaler
	info ArgInfo
}

func (e env) ArgInfo() ArgInfo {
	return e.info
}

func (e env) Parse(ctx ParseContext) bool {
	value, ok := os.LookupEnv(e.key)
	if !ok {
		return false
	}
	return ctx.UnmarshalArg(e.u, value)
}

var _ Arg = env{}
