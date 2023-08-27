package main

import "context"

type ValueNode[T any] struct {
	Node[T]

	value T
}

func NewValueNode[T any](value T) Node[T] {
	return &ValueNode[T]{value: value}
}

func (v *ValueNode[T]) GetValue() T {
	return v.value
}

func (v *ValueNode[T]) IsResolved() bool {
	return true
}

func (v *ValueNode[T]) GetAnyResolvables() []AnyNode {
	return []AnyNode{}
}

func (v *ValueNode[T]) Run(_ context.Context) any {
	return v.value
}
