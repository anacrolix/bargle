package bargle

import (
	"fmt"
	"reflect"

	"github.com/huandu/xstrings"
)

func FromStruct(target interface{}) (cmd Command) {
	value := reflect.ValueOf(target).Elem()
	type_ := value.Type()
	for i := 0; i < value.NumField(); i++ {
		fieldValue := value.Field(i)
		target := fieldValue.Addr().Interface()
		structField := type_.Field(i)
		argTag := structField.Tag.Get("arg")
		if argTag == "-" {
			continue
		}
		arity := structField.Tag.Get("arity")
		var param Param
		getUnmarshaler := func() anyUnaryUnmarshaler {
			unmarshaler, err := makeAnyUnaryUnmarshalerViaReflection(target)
			if err != nil {
				panic(fmt.Errorf("getting unmarshaler for %v: %w", structField, err))
			}
			return unmarshaler
		}
		if argTag == "positional" {
			param = &Positional{
				Name:  fmt.Sprintf("%v.%v", type_.Name(), structField.Name),
				Desc:  structField.Tag.Get("help"),
				Value: getUnmarshaler(),
			}
			// TODO: Handle required/not-required.
			cmd.Positionals = append(cmd.Positionals, param)
		} else {
			longs := []string{xstrings.ToKebabCase(structField.Name)}
			switch typedTarget := target.(type) {
			case *bool:
				flag := NewFlag(typedTarget)
				flag.AddLongs(longs...)
				param = flag.Make()
			default:
				option := NewUnaryOption(getUnmarshaler())
				option.AddLongs(longs...)
				if arity == "+" {
					option.SetRequired()
				}
				param = option.Make()
			}
			cmd.Options = append(cmd.Options, param)
		}
		default_ := structField.Tag.Get("default")
		if default_ != "" {
			err := param.Parse(NewArgs([]string{default_}))
			if err != nil {
				panic(fmt.Errorf("setting default %q: %w", default_, err))
			}
		}
	}
	return
}
