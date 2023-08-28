package panera

import (
	"context"
)

// ValueNode represent a constant value.
type ValueNode[T any] struct {
	Node[T]

	id    NodeID
	value T
}

func NewValueNode[T any](value T) Node[T] {
	return &ValueNode[T]{
		id:    NewNodeID(),
		value: value,
	}
}

func (v *ValueNode[T]) GetValue(_ context.Context) T {
	return v.value
}

func (v *ValueNode[T]) IsResolved(_ context.Context) bool {
	return true
}

func (v *ValueNode[T]) GetChildren() []AnyNode {
	return []AnyNode{}
}

func (v *ValueNode[T]) Run(_ context.Context) any {
	return v.value
}

func (v *ValueNode[T]) GetID() NodeID {
	return v.id
}

func (v *ValueNode[T]) Debug() string {
	return "Value"
}
