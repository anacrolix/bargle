package bargle

import (
	"encoding"
	g "github.com/anacrolix/generics"
)

type Unmarshaler interface {
	Unmarshal(ctx UnmarshalContext) error
	ArgTypes() []string
}

type UnmarshalerValuer interface {
	Value() any
}

func String(s *string) Unmarshaler {
	return stringUnmarshaler{s}
}

type stringUnmarshaler struct {
	s *string
}

func (me stringUnmarshaler) Value() any {
	return *me.s
}

func (me stringUnmarshaler) ArgTypes() []string {
	return g.Singleton("string")
}

func (me stringUnmarshaler) Unmarshal(ctx UnmarshalContext) (err error) {
	*me.s, err = ctx.Pop()
	return err
}

func TextUnmarshaler(tu encoding.TextUnmarshaler) Unmarshaler {
	return textUnmarshaler{tu}
}

type textUnmarshaler struct {
	inner encoding.TextUnmarshaler
}

func (t textUnmarshaler) Unmarshal(ctx UnmarshalContext) (err error) {
	s, err := ctx.Pop()
	if err != nil {
		return
	}
	err = t.inner.UnmarshalText([]byte(s))
	return
}

func (t textUnmarshaler) ArgTypes() []string {
	return []string{"string"}
}
