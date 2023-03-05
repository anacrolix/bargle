package bargle

import (
	"fmt"
	"net"
	"net/url"
	"strconv"
	"time"
)

// Returns an unmarshaler for a builtin type. t must be a pointer to a type in the
// BuiltinUnmarshalerType type set.
func BuiltinUnmarshalerFromAny(t any) Unmarshaler {
	switch t := t.(type) {
	case *string:
		return String(t)
	case *bool:
		return boolUnmarshaler{t}
	case **url.URL:
		return UnaryUnmarshalFunc(t, url.Parse)
	case *time.Duration:
		return UnaryUnmarshalFunc(t, time.ParseDuration)
	case *net.IP:
		return UnaryUnmarshalFunc(t, func(s string) (ip net.IP, err error) {
			ip = net.ParseIP(s)
			if ip == nil {
				err = fmt.Errorf("failed to parse IP from %q", s)
			}
			return
		})
	case *int:
		return intUnmarshaler[int, int64]{
			t:    t,
			f:    strconv.ParseInt,
			bits: 0,
		}
	case *int32:
		return intUnmarshaler[int32, int64]{
			t:    t,
			f:    strconv.ParseInt,
			bits: 0,
		}
	case *uint32:
		return intUnmarshaler[uint32, uint64]{
			t:    t,
			f:    strconv.ParseUint,
			bits: 0,
		}
	case *uint64:
		return intUnmarshaler[uint64, uint64]{
			t:    t,
			f:    strconv.ParseUint,
			bits: 0,
		}
	case *float64:
		return floatUnmarshaler[float64]{
			t:    t,
			bits: 64,
		}
	case *float32:
		return floatUnmarshaler[float32]{
			t:    t,
			bits: 32,
		}
	default:
		return nil
	}
}

// An unmarshaler for any of the types in the BuiltinUnmarshalerType type set.
func BuiltinUnmarshaler[T BuiltinUnmarshalerType](t *T) Unmarshaler {
	u := BuiltinUnmarshalerFromAny(t)
	if u == nil {
		// I expect this shouldn't happen as tne types are enforced by BuiltinUnmarshalerType. We
		// could include a better type error here.
		panic("unreachable")
	}
	return u
}

// A set of types supported by the builtin unmarshaler.
type BuiltinUnmarshalerType interface {
	string | *url.URL | int | net.IP | time.Duration | bool | int32 | uint32 | uint64 | float64 | float32
}

type Builtin[T BuiltinUnmarshalerType] struct {
	Value T
}

func (b *Builtin[T]) Unmarshal(ctx UnmarshalContext) error {
	return BuiltinUnmarshaler(&b.Value).Unmarshal(ctx)
}
