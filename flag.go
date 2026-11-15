package bargle

import (
	"fmt"
	"strconv"
	"strings"

	g "github.com/anacrolix/generics"
)

// Returns a boolean switch. Unlike a Long wrapping a bool Unmarshaler, a flag has a negative form
// ("--no-key"), so switches that default to true can be turned off.
func Flag(key string, target *bool) flag {
	return flag{
		key: key,
		set: func(value bool) {
			*target = value
		},
		value: func() any {
			return *target
		},
	}
}

// Returns a boolean switch that records whether it was given at all, for switches where the
// unset state is distinct from false.
func OptionFlag(key string, target *g.Option[bool]) flag {
	return flag{
		key: key,
		set: func(value bool) {
			target.Set(value)
		},
		value: func() any {
			if !target.Ok {
				return nil
			}
			return target.Value
		},
	}
}

// See Flag. The key is joined from elements, like LongElems.
func FlagElems(target *bool, firstElem string, elems ...string) flag {
	return Flag(strings.Join(append(g.Singleton(firstElem), elems...), "-"), target)
}

type flag struct {
	key   string
	set   func(value bool)
	value func() any
}

var _ interface {
	Arg
	ArgValuer
} = flag{}

func (me flag) ArgInfo() ArgInfo {
	return ArgInfo{
		MatchingForms: g.Singleton(fmt.Sprintf("--[no-]%[1]s, --[no-]%[1]s=bool", me.key)),
		ArgType:       ArgTypeSwitch,
	}
}

func (me flag) Value() any {
	return me.value()
}

func (me flag) Parse(ctx ParseContext) bool {
	arg, ok := ctx.Pop()
	if !ok {
		return false
	}
	key, ok := strings.CutPrefix(arg, "--")
	if !ok {
		return false
	}
	key, value, haveValue := strings.Cut(key, "=")
	var negate bool
	switch key {
	case me.key:
	case "no-" + me.key:
		negate = true
	default:
		return false
	}
	u := flagUnmarshaler{me.set, negate}
	if haveValue {
		return ctx.UnmarshalArg(u, value)
	}
	return ctx.Unmarshal(u)
}

// Sets a bool, optionally inverting an explicit value. Giving a flag with no value always means
// the affirmative sense of the form used, so "--no-key" sets false and "--no-key=false" sets true.
type flagUnmarshaler struct {
	set    func(value bool)
	negate bool
}

func (me flagUnmarshaler) ArgTypes() []string {
	return g.Singleton("?bool")
}

func (me flagUnmarshaler) Unmarshal(ctx UnmarshalContext) error {
	value := true
	if ctx.HaveExplicitValue() {
		arg, err := ctx.Pop()
		if err != nil {
			return err
		}
		value, err = strconv.ParseBool(arg)
		if err != nil {
			return err
		}
	}
	me.set(value != me.negate)
	return nil
}
