package bargle

import (
	"encoding"
	"fmt"
	"reflect"
	"strconv"

	"golang.org/x/exp/constraints"
)

type UnaryUnmarshaler interface {
	UnaryUnmarshal(s string) error
	//CurrentValue() string
	TargetHelp() string
	Matching() bool
}

type basicBuiltinUnaryUnmarshalTarget interface {
	~string | ~int16 | ~uint16
}

type builtinUnaryUnmarshalTarget interface {
	basicBuiltinUnaryUnmarshalTarget
}

type BuiltinUnaryUnmarshaler[T builtinUnaryUnmarshalTarget] struct{}

func parseIntType[R, I constraints.Integer](s string, bits int, f strconvParseIntFunc[I]) (ret R, err error) {
	var i I
	i, err = f(s, 0, bits)
	ret = R(i)
	return
}

type Slice[T any] struct {
	value       *[]T
	Unmarshaler UnaryUnmarshaler
	elemTarget  *T
}

func (s2 Slice[T]) TargetHelp() string {
	return s2.Unmarshaler.TargetHelp() + "..."
}

func (sl Slice[T]) UnaryUnmarshal(s string) error {
	err := sl.Unmarshaler.UnaryUnmarshal(s)
	if err != nil {
		return err
	}
	*sl.value = append(*sl.value, *sl.elemTarget)
	return nil
}

type String struct {
	Target *string
	Ok     bool
}

var _ UnaryUnmarshaler = (*String)(nil)

func (me String) Value() string {
	return *me.Target
}

func (s2 String) TargetHelp() string {
	return "(string)"
}

func (me *String) UnaryUnmarshal(arg string) error {
	*me.Target = arg
	me.Ok = true
	return nil
}

func (me *String) Matching() bool {
	return !me.Ok
}

func NewString() *String {
	return &String{}
}

func makeSimpleAnyUnaryUnmarshalFromFunc[T any](target T, u func(string) error) anyUnaryUnmarshaler {
	matched := false
	return typedToAnyUnaryUnmarshalerWrapper[T]{
		unaryUnmarshalerFunc[T]{
			u: func(s string) error {
				matched = true
				return u(s)
			},
			target: target,
			matching: func() bool {
				return !matched
			},
			help: fmt.Sprintf("(%v)", reflect.TypeOf(target).Elem().String()),
		},
	}
}

type strconvParseIntFunc[I any] func(string, int, int) (I, error)

func makeIntUnaryUnmarshaler[T, I constraints.Integer](p *T, bits int, intFunc strconvParseIntFunc[I]) (anyUnaryUnmarshaler, error) {
	return makeSimpleAnyUnaryUnmarshalFromFunc(p, func(s string) (err error) {
		*p, err = parseIntType[T](s, bits, intFunc)
		return
	}), nil
}

// Does a unary unmarshal, trying to infer a default unmarshaler if necessary.
func makeAnyUnaryUnmarshalerViaReflection(target any) (anyUnaryUnmarshaler, error) {
	if tu, ok := target.(encoding.TextUnmarshaler); ok {
		return makeSimpleAnyUnaryUnmarshalFromFunc(target, func(s string) error {
			return tu.UnmarshalText([]byte(s))
		}), nil
	}
	switch p := target.(type) {
	case *string:
		return typedToAnyUnaryUnmarshalerWrapper[string]{&String{Target: p}}, nil
	case *uint16:
		return makeIntUnaryUnmarshaler(p, 16, strconv.ParseUint)
	case *int64:
		return makeIntUnaryUnmarshaler(p, 64, strconv.ParseInt)
	case *int:
		return makeIntUnaryUnmarshaler(p, 0, strconv.ParseInt)
	}
	targetPtrValue := reflect.ValueOf(target)
	targetValue := targetPtrValue.Elem()
	targetType := targetPtrValue.Type().Elem()
	switch targetType.Kind() {
	case reflect.Slice:
		elemTarget := reflect.New(targetType.Elem())
		eu, err := makeAnyUnaryUnmarshalerViaReflection(elemTarget.Interface())
		if err != nil {
			return nil, fmt.Errorf("getting unmarshaller for slice elem: %w", err)
		}
		return unaryUnmarshalerFunc[any]{
			u: func(s string) error {
				err := eu.UnaryUnmarshal(s)
				if err != nil {
					return err
				}
				targetValue.Set(reflect.Append(targetValue, elemTarget.Elem()))
				return nil
			},
			target: target,
			matching: func() bool {
				return true
			},
			help: fmt.Sprintf("(%v...)", reflect.TypeOf(target).Elem().Elem().String()),
		}, nil
	case reflect.Ptr:
		subTarget := targetValue
		if subTarget.IsNil() {
			subTarget = reflect.New(targetType.Elem())
		}
		elemUnmarshaler, err := makeAnyUnaryUnmarshalerViaReflection(subTarget.Interface())
		if err != nil {
			return nil, fmt.Errorf("getting unmarshaller for target elem: %w", err)
		}
		return unaryUnmarshalerWithUnmarshalFunc[any]{
			func(s string) error {
				err := elemUnmarshaler.UnaryUnmarshal(s)
				if err != nil {
					return err
				}
				targetValue.Set(subTarget)
				return nil
			},
			elemUnmarshaler,
		}, nil
	default:
		return nil, fmt.Errorf("unhandled unary unmarshaler reflection type: %T", target)
	}
}

type unaryUnmarshalerWithUnmarshalFunc[T any] struct {
	uf func(string) error
	UnaryUnmarshaler
}

func (me unaryUnmarshalerWithUnmarshalFunc[T]) UnaryUnmarshal(s string) error {
	return me.uf(s)
}

type anyUnaryUnmarshaler = UnaryUnmarshaler

type typedToAnyUnaryUnmarshalerWrapper[T any] struct {
	UnaryUnmarshaler
}

type unaryUnmarshalerWrapperAnyToTyped[T any] struct {
	anyUnaryUnmarshaler
}

func (me unaryUnmarshalerWrapperAnyToTyped[T]) UnaryUnmarshal(s string) error {
	return me.anyUnaryUnmarshaler.UnaryUnmarshal(s)
}

func AutoUnmarshaler[T any](t *T) (u UnaryUnmarshaler) {
	err := initNilUnmarshalerUsingReflect(&u, t)
	if err != nil {
		panic(err)
	}
	return
}

func initNilUnmarshalerUsingReflect[T any](u *UnaryUnmarshaler, t *T) error {
	if *u != nil {
		return nil
	}
	if t == nil {
		t = new(T)
	}
	unmarshaler, err := makeAnyUnaryUnmarshalerViaReflection(t)
	if err != nil {
		return err
	}
	*u = unaryUnmarshalerWrapperAnyToTyped[T]{unmarshaler}
	return nil
}

type unaryUnmarshalerFunc[T any] struct {
	u        func(string) error
	target   T
	help     string
	matching func() bool
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

func (me unaryUnmarshalerFunc[T]) Matching() bool {
	return me.matching()
}
