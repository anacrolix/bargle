package bargle

import (
	"fmt"
	"net"
	"net/url"
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
		return signedUnmarshaler(t, 0)
	case *int8:
		return signedUnmarshaler(t, 8)
	case *int16:
		return signedUnmarshaler(t, 16)
	case *int32:
		return signedUnmarshaler(t, 32)
	case *int64:
		return signedUnmarshaler(t, 64)
	case *uint:
		return unsignedUnmarshaler(t, 0)
	case *uint8:
		return unsignedUnmarshaler(t, 8)
	case *uint16:
		return unsignedUnmarshaler(t, 16)
	case *uint32:
		return unsignedUnmarshaler(t, 32)
	case *uint64:
		return unsignedUnmarshaler(t, 64)
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
	string | *url.URL | net.IP | time.Duration | bool |
		int | int8 | int16 | int32 | int64 |
		uint | uint8 | uint16 | uint32 | uint64 |
		float64 | float32
}

type Builtin[T BuiltinUnmarshalerType] struct {
	Value T
}

func (b *Builtin[T]) Unmarshal(ctx UnmarshalContext) error {
	return BuiltinUnmarshaler(&b.Value).Unmarshal(ctx)
}
