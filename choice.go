package bargle

import (
	"fmt"

	"github.com/anacrolix/generics"
)

type Choice[T any] struct {
	Choices map[string]T
}

func (me Choice[T]) Help(ph *ParamHelp) {
	for choice := range me.Choices {
		ph.Options = append(ph.Options, ParamHelp{
			Forms: []string{choice},
			//Description: fmt.Sprintf("%v", v),
		})
	}
}

func (me Choice[T]) Unmarshal(choice string, t *T) error {
	var ok bool
	*t, ok = me.Choices[choice]
	if !ok {
		return controlError{fmt.Errorf("unknown choice: %q", choice)}
	}
	return nil
}

func (me Choice[T]) Add(name string, value T) {
	generics.MakeMapIfNil(&me.Choices)
	me.Choices[name] = value
}

func (me Choice[T]) Get(key string) T {
	return me.Choices[key]
}

func NewChoice[T any](choices map[string]T) *Choice[T] {
	return &Choice[T]{Choices: choices}
}

type Choices[T any] map[string]T
