package bargle

import (
	"errors"
	"strings"

	"github.com/anacrolix/generics"
	"golang.org/x/exp/maps"
)

type Choice[T any] struct {
	value   T
	Choices map[string]T
}

func (me Choice[T]) Matching() bool {
	return true
}

func (me Choice[T]) Value() T {
	return me.value
}

func (me Choice[T]) TargetHelp() string {
	return strings.Join(maps.Keys(me.Choices), " | ")
}

func (me *Choice[T]) UnaryUnmarshal(choice string) error {
	var ok bool
	me.value, ok = me.Choices[choice]
	if !ok {
		return controlError{errors.New("unknown choice")}
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
