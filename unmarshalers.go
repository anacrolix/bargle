package bargle

import (
	"encoding"
	"fmt"
	"reflect"
	"strconv"

	"golang.org/x/exp/constraints"
)

type UnaryUnmarshaler[T any] interface {
	UnaryUnmarshal(s string, t *T) error
	TargetHelp() string
}

type basicBuiltinUnaryUnmarshalTarget interface {
	~string | ~int16 | ~uint16
}

type builtinUnaryUnmarshalTarget interface {
	basicBuiltinUnaryUnmarshalTarget
}

type BuiltinUnaryUnmarshaler[T builtinUnaryUnmarshalTarget] struct{}

func parseIntType[R constraints.Integer, I constraints.Integer](s string, bits int, f func(string, int, int) (I, error)) (ret R, err error) {
	var i I
	i, err = f(s, 0, bits)
	ret = R(i)
	return
}

func (me BuiltinUnaryUnmarshaler[T]) UnaryUnmarshal(s string, t *T) (err error) {
	switch p := any(t).(type) {
	case *string:
		*p = s
		return nil
	case *int16:
		*p, err = parseIntType[int16](s, 16, strconv.ParseInt)
	case *uint16:
		*p, err = parseIntType[uint16](s, 16, strconv.ParseUint)
	default:
		panic(fmt.Sprintf("builtin unary unmarshaler: unsupported type %T", *t))
	}
	return
}

func (me BuiltinUnaryUnmarshaler[T]) TargetHelp() string {
	var t T
	return fmt.Sprintf("(%T)", t)
}

type Unmarshaler[T any] interface {
	Unmarshal(args Args, t *T) error
}

type Slice[T any] struct {
	Unmarshaler interface {
		UnaryUnmarshaler[T]
		TargetHelper
	}
}

func (s2 Slice[T]) TargetHelp() string {
	return s2.Unmarshaler.TargetHelp() + "..."
}

func NewSlice[T any](u interface {
	UnaryUnmarshaler[T]
	TargetHelper
}) Slice[T] {
	return Slice[T]{Unmarshaler: u}
}

func (sl Slice[T]) UnaryUnmarshal(s string, slice *[]T) error {
	var t T
	err := sl.Unmarshaler.UnaryUnmarshal(s, &t)
	if err != nil {
		return err
	}
	*slice = append(*slice, t)
	return nil
}

type String struct{}

var _ UnaryUnmarshaler[string] = String{}

func (s2 String) TargetHelp() string {
	return "(string)"
}

func (String) Unmarshal(args Args, s *string) error {
	*s = args.Pop()
	return nil
}

func (String) UnaryUnmarshal(arg string, s *string) error {
	*s = arg
	return nil
}

func NewString() *String {
	return &String{}
}

// Does a unary unmarshal, trying to infer a default unmarshaler if necessary.
func doUnaryUnmarshal[T any](s string, t *T, u UnaryUnmarshaler[T]) error {
	if u != nil {
		return u.UnaryUnmarshal(s, t)
	}
	if tu, ok := any(t).(encoding.TextUnmarshaler); ok {
		return tu.UnmarshalText([]byte(s))
	}
	switch p := any(t).(type) {
	case *string:
		return String{}.UnaryUnmarshal(s, p)
	default:
		panic(fmt.Sprintf("unhandled default unary unmarshaler type %T", *t))
	}
}

// Does a unary unmarshal, trying to infer a default unmarshaler if necessary.
func mustGetUnaryUnmarshaler(target reflect.Type) anyUnaryUnmarshaler {
	return anyUnaryUnmarshalerFunc{
		u: mustGetUnaryUnmarshalAnyFunc(target),
		// TODO: Collapse pointers that get allocated automatically for this type. Pass through some
		// nice examples for builtins when/if they are added.
		help: fmt.Sprintf("(%v)", target.String()),
	}
}

// Does a unary unmarshal, trying to infer a default unmarshaler if necessary.
func mustGetUnaryUnmarshalAnyFunc(target reflect.Type) func(string, any) error {
	if target.Implements(reflect.TypeOf((*encoding.TextUnmarshaler)(nil)).Elem()) {
		return func(s string, a any) error {
			return a.(encoding.TextUnmarshaler).UnmarshalText([]byte(s))
		}
	}
	switch target.Elem().Kind() {
	case reflect.Slice:
		eu := mustGetUnaryUnmarshalAnyFunc(reflect.PtrTo(target.Elem().Elem()))
		return func(s string, a any) error {
			ev := reflect.New(target.Elem().Elem())
			err := eu(s, ev.Interface())
			if err != nil {
				return err
			}
			slice := reflect.ValueOf(a).Elem()
			slice.Set(reflect.Append(slice, ev.Elem()))
			return nil
		}
	case reflect.String:
		return func(s string, a any) error {
			*a.(*string) = s
			return nil
		}
	case reflect.Ptr:
		uf := mustGetUnaryUnmarshalAnyFunc(target.Elem())
		return func(s string, a any) error {
			ev := reflect.New(target.Elem())
			err := uf(s, ev.Interface())
			if err != nil {
				return err
			}
			reflect.ValueOf(a).Elem().Set(ev)
			return nil
		}
	default:
		panic(fmt.Errorf("unhandled unary unmarshaler reflection type: %v", target))
	}

}

type anyUnaryUnmarshalerFunc struct {
	u    func(string, any) error
	help string
}

func (me anyUnaryUnmarshalerFunc) UnaryUnmarshal(s string, a *any) error {
	return me.u(s, a)
}

func (me anyUnaryUnmarshalerFunc) TargetHelp() string {
	return me.help
}

type anyUnaryUnmarshaler = UnaryUnmarshaler[any]

type unaryUnmarshalerAnyWrapper[T any] struct {
	UnaryUnmarshaler[T]
}

func (me unaryUnmarshalerAnyWrapper[T]) UnaryUnmarshal(s string, t *any) error {
	return me.UnaryUnmarshaler.UnaryUnmarshal(s, (*t).(*T))
}
