package panera

import (
	"context"
)

type AnyFlatMap interface {
	FlatMapFullyResolved(context.Context) bool
}

type FlatMapNode[T, U any] interface {
	Node[U]
	AnyFlatMap
}

type flatMapNodeImpl[T, U any] struct {
	id    NodeID
	child Node[T]
	fn    func(T) Node[U]
}

func NewFlatMapNode[T, U any](
	node Node[T],
	transformer func(T) Node[U],
) FlatMapNode[T, U] {
	return &flatMapNodeImpl[T, U]{
		id:    NewNodeID(),
		child: node,
		fn:    transformer,
	}
}

func (f *flatMapNodeImpl[T, U]) IsResolved(ctx context.Context) bool {
	nodeState := NodeStateFromContext(ctx)
	return nodeState.GetIsResolved(f.id)
}

func (f *flatMapNodeImpl[T, U]) Run(ctx context.Context) any {
	nodeState := NodeStateFromContext(ctx)

	if _, ok := nodeState.GetResolvedValue(f.id).(Node[U]); ok {
		panic("Invariant violation: flatMap run multiple times")
	}

	result := f.fn(f.child.GetValue(ctx))
	nodeState.SetResolvedValue(f.id, result)

	return result
}

func (f *flatMapNodeImpl[T, U]) GetValue(ctx context.Context) U {
	nodeState := NodeStateFromContext(ctx)

	grandChild, ok := nodeState.GetResolvedValue(f.id).(Node[U])
	if !ok {
		panic("Invariant violation: node state does not contain a valid node")
	}

	return grandChild.GetValue(ctx)
}

func (f *flatMapNodeImpl[T, U]) GetChildren() []AnyNode {
	return []AnyNode{f.child}
}

func (f *flatMapNodeImpl[T, U]) FlatMapFullyResolved(ctx context.Context) bool {
	nodeState := NodeStateFromContext(ctx)

	grandChild, ok := nodeState.GetResolvedValue(f.id).(Node[U])
	if !ok {
		return false
	}

	return grandChild.IsResolved(ctx)
}

func (f *flatMapNodeImpl[T, U]) GetID() NodeID {
	return f.id
}

func (f *flatMapNodeImpl[T, U]) Debug() string {
	return "FlatMap"
}
