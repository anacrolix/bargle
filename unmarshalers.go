package bargle

import (
	"encoding"
	"fmt"
	"reflect"
	"strconv"

	"golang.org/x/exp/constraints"
)

type UnaryUnmarshaler[T any] interface {
	UnaryUnmarshal(s string) error
	Value() T
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
	value       *[]T
	Unmarshaler UnaryUnmarshaler[T]
}

func (s2 Slice[T]) TargetHelp() string {
	return s2.Unmarshaler.TargetHelp() + "..."
}

func (sl Slice[T]) UnaryUnmarshal(s string) error {
	err := sl.Unmarshaler.UnaryUnmarshal(s)
	if err != nil {
		return err
	}
	*sl.value = append(*sl.value, sl.Unmarshaler.Value())
	return nil
}

type String struct {
	value *string
}

var _ UnaryUnmarshaler[string] = String{}

func (me String) Value() string {
	return *me.value
}

func (s2 String) TargetHelp() string {
	return "(string)"
}

func (String) Unmarshal(args Args, s *string) error {
	*s = args.Pop()
	return nil
}

func (me String) UnaryUnmarshal(arg string) error {
	*me.value = arg
	return nil
}

func NewString() *String {
	return &String{}
}

// Does a unary unmarshal, trying to infer a default unmarshaler if necessary.
func mustGetUnaryUnmarshaler(target any) (_ anyUnaryUnmarshaler, err error) {
	unmarshalFunc, err := mustGetUnaryUnmarshalAnyFunc(target)
	if err != nil {
		err = fmt.Errorf("getting unmarshal func for a %T: %w", target, err)
		return
	}
	return unaryUnmarshalerFunc[any]{
		target: target,
		u:      unmarshalFunc,
		// TODO: Collapse pointers that get allocated automatically for this type. Pass through some
		// nice examples for builtins when/if they are added.
		help: fmt.Sprintf("(%v)", reflect.TypeOf(target).Elem().String()),
	}, nil
}

// Does a unary unmarshal, trying to infer a default unmarshaler if necessary.
func mustGetUnaryUnmarshalAnyFunc(target any) (func(string) error, error) {
	if tu, ok := target.(encoding.TextUnmarshaler); ok {
		return func(s string) error {
			return tu.UnmarshalText([]byte(s))
		}, nil
	}
	switch p := target.(type) {
	case *string:
		return String{value: p}.UnaryUnmarshal, nil
	case *uint16:
		return func(s string) error {
			u64, err := strconv.ParseUint(s, 0, 16)
			*p = uint16(u64)
			return err
		}, nil
	}
	targetPtrValue := reflect.ValueOf(target)
	targetValue := targetPtrValue.Elem()
	targetType := targetPtrValue.Type().Elem()
	switch targetType.Kind() {
	case reflect.Slice:
		elemTarget := reflect.New(targetType.Elem())
		eu, err := mustGetUnaryUnmarshalAnyFunc(elemTarget.Interface())
		if err != nil {
			return nil, fmt.Errorf("getting unmarshaller for slice elem: %w", err)
		}
		return func(s string) error {
			err := eu(s)
			if err != nil {
				return err
			}
			targetValue.Set(reflect.Append(targetValue, elemTarget.Elem()))
			return nil
		}, nil
	//case reflect.Ptr:
	//	uf := mustGetUnaryUnmarshalAnyFunc(target.Elem())
	//	return func(s string, a any) error {
	//		ev := reflect.New(target.Elem())
	//		err := uf(s, ev.Interface())
	//		if err != nil {
	//			return err
	//		}
	//		reflect.ValueOf(a).Elem().Set(ev)
	//		return nil
	//	}
	default:
		return nil, fmt.Errorf("unhandled unary unmarshaler reflection type: %T", target)
	}
}

type anyUnaryUnmarshaler = UnaryUnmarshaler[any]

type unaryUnmarshalerAnyWrapper[T any] struct {
	UnaryUnmarshaler[T]
}

func (me unaryUnmarshalerAnyWrapper[T]) UnaryUnmarshal(s string) error {
	return me.UnaryUnmarshaler.UnaryUnmarshal(s)
}

type unaryUnmarshalerWrapperAnyToTyped[T any] struct {
	anyUnaryUnmarshaler
}

func (me unaryUnmarshalerWrapperAnyToTyped[T]) UnaryUnmarshal(s string) error {
	return me.anyUnaryUnmarshaler.UnaryUnmarshal(s)
}

func (me unaryUnmarshalerWrapperAnyToTyped[T]) Value() T {
	return me.anyUnaryUnmarshaler.Value().(T)
}

func initNilUnmarshalerUsingReflect[T any](u *UnaryUnmarshaler[T]) error {
	if *u != nil {
		return nil
	}
	t := new(T)
	unmarshaler, err := mustGetUnaryUnmarshaler(t)
	if err != nil {
		return err
	}
	*u = unaryUnmarshalerWrapperAnyToTyped[T]{unmarshaler}
	return nil
}

type unaryUnmarshalerFunc[T any] struct {
	u      func(string) error
	target T
	help   string
}

func (me unaryUnmarshalerFunc[T]) UnaryUnmarshal(s string) error {
	return me.u(s)
}

func (me unaryUnmarshalerFunc[T]) TargetHelp() string {
	return me.help
}

func (me unaryUnmarshalerFunc[T]) Value() T {
	targetValue := reflect.ValueOf(me.target)
	if targetValue.IsNil() {
		panic(fmt.Errorf("target is nil: %v", me.target))
	}
	return targetValue.Elem().Interface().(T)
}
