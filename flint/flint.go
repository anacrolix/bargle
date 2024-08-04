package flint

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/huandu/xstrings"

	"github.com/anacrolix/bargle/v2"
)

type sub struct {
	key          string
	cmd          parsedStruct
	newStructPtr reflect.Value
	setOnParse   reflect.Value
}

type pos struct {
	name     string
	arg      bargle.Arg
	required bool
}

type parsedStruct struct {
	options []bargle.Arg
	pos     []pos
	subs    []sub
}

func (ps parsedStruct) Run(p *bargle.Parser) {
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
		if p.Parse(bargle.Keyword(sub.key)) {
			sub.cmd.Run(p)
			sub.setOnParse.Set(sub.newStructPtr)
			break
		}
	}
}

func processStruct(s any, p *bargle.Parser) (ret parsedStruct) {
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
		u := bargle.BuiltinUnmarshalerFromAny(fieldPtrAny)
		if u == nil {
			panic(fmt.Sprintf("unsupported type: %T", t))
		}
		if value, ok := sf.Tag.Lookup("arg"); ok {
			parts := strings.Split(value, ",")
			required := len(parts) > 1 && parts[1] == "required"
			if parts[0] == "positional" {
				ret.pos = append(ret.pos, pos{
					name:     name,
					arg:      bargle.Positional(u),
					required: required,
				})
			}
		}
		ret.options = append(ret.options, bargle.Long(name, u))
		if value, ok := sf.Tag.Lookup("default"); ok {
			p.SetDefault(u, value)
		}
	}
	return
}

// Processes defaults and parses stuff in the struct per flint's implementation. I'm not sure how
// you would mix positionals external to the struct with how this does it.
func ParseStruct[T any](p *bargle.Parser, s *T) {
	ps := processStruct(s, p)
	ps.Run(p)
}

// Sets defaults for the type. Probably does other stuff but shouldn't.
func SetDefaults[T any](s *T) {
	ParseStruct(bargle.NewParserNoArgs(), s)
}
