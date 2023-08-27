package panera

import "context"

type ValueNode[T any] struct {
	Node[T]

	value T
}

func NewValueNode[T any](value T) Node[T] {
	return &ValueNode[T]{value: value}
}

func (v *ValueNode[T]) GetValue(_ context.Context, _ int) T {
	return v.value
}

func (v *ValueNode[T]) IsResolved(_ context.Context, _ int) bool {
	return true
}

func (v *ValueNode[T]) GetChildren() []AnyNode {
	return []AnyNode{}
}

func (v *ValueNode[T]) Run(_ context.Context, _ int) any {
	return v.value
}
