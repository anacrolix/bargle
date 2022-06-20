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
		targetReflectType := fieldValue.Addr().Type()
		target := fieldValue.Addr().Interface()
		structField := type_.Field(i)
		argTag := structField.Tag.Get("arg")
		if argTag == "-" {
			continue
		}
		arity := structField.Tag.Get("arity")
		var param Param
		if argTag == "positional" {
			param = &Positional[any]{
				Name:  fmt.Sprintf("%v.%v", type_.Name(), structField.Name),
				Desc:  structField.Tag.Get("help"),
				Value: target,
				U:     mustGetUnaryUnmarshaler(targetReflectType),
			}
			cmd.Positionals = append(cmd.Positionals, param)
		} else {
			longs := []string{xstrings.ToKebabCase(structField.Name)}
			switch typedTarget := target.(type) {
			case *bool:
				param = &Flag{
					Value: typedTarget,
					Longs: longs,
				}
				cmd.Options = append(cmd.Options)
			default:
				option := &UnaryOption[any]{
					Value:       target,
					Unmarshaler: mustGetUnaryUnmarshaler(targetReflectType),
					Longs:       longs,
				}
				if arity == "+" {
					option.Required = true
				}
				param = option
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
