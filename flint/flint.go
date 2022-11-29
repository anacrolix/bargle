package flint

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/huandu/xstrings"

	"github.com/anacrolix/bargle"
)

type sub struct {
	key          string
	cmd          parsedStruct
	newStructPtr reflect.Value
	setOnParse   reflect.Value
}

type pos struct {
	name     string
	arg      args.Arg
	required bool
}

type parsedStruct struct {
	options []args.Arg
	pos     []pos
	subs    []sub
}

func (ps parsedStruct) Run(p *args.Parser) {
opts:
	for _, opt := range ps.options {
		if p.Parse(opt) {
			goto opts
		}
	}
	for _, pos := range ps.pos {
		if !p.Parse(pos.arg) && pos.required {
			p.SetError(fmt.Errorf("%q required and not given", pos.name))
			return
		}
	}
	for _, sub := range ps.subs {
		if p.Parse(args.Keyword(sub.key)) {
			sub.cmd.Run(p)
			sub.setOnParse.Set(sub.newStructPtr)
			break
		}
	}
}

func processStruct(s any, p *args.Parser) (ret parsedStruct) {
	v := reflect.ValueOf(s).Elem()
	structType := v.Type()
	for i := 0; i < v.NumField(); i++ {
		fieldValue := v.Field(i)
		fieldType := fieldValue.Type()
		sf := structType.Field(i)
		name := xstrings.ToKebabCase(sf.Name)
		if value, ok := sf.Tag.Lookup("arg"); ok {
			if value == "subcommand" {
				newStruct := reflect.New(fieldType.Elem())
				ret.subs = append(ret.subs, sub{
					key:          name,
					cmd:          processStruct(newStruct.Interface(), p),
					setOnParse:   fieldValue,
					newStructPtr: newStruct,
				})
				continue
			}
		}
		t := v.Field(i).Addr().Interface()
		fieldPtrAny := v.Field(i).Addr().Interface()
		u := args.BuiltinUnmarshalerFromAny(fieldPtrAny)
		if u == nil {
			panic(fmt.Sprintf("unsupported type: %T", t))
		}
		if value, ok := sf.Tag.Lookup("arg"); ok {
			parts := strings.Split(value, ",")
			required := len(parts) > 1 && parts[1] == "required"
			if parts[0] == "positional" {
				ret.pos = append(ret.pos, pos{
					name:     name,
					arg:      args.Positional(u),
					required: required,
				})
			}
		}
		ret.options = append(ret.options, args.Long(name, u))
		if value, ok := sf.Tag.Lookup("default"); ok {
			p.SetDefault(u, value)
		}
	}
	return
}

// Processes defaults and parses stuff in the struct per flint's implementation. I'm not sure how
// you would mix positionals external to the struct with how this does it.
func ParseStruct[T any](p *args.Parser, s *T) {
	ps := processStruct(s, p)
	ps.Run(p)
}