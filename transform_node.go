package main

import "context"

type TransformNode[T, U any] interface {
	Node[U]
}

type transformNodeImpl[T, U any] struct {
	isResolved bool
	child      Node[T]
	fn         func(T) U
}

func NewTransformNode[T, U any](node Node[T], transformer func(T) U) TransformNode[T, U] {
	return &transformNodeImpl[T, U]{
		child: node,
		fn:    transformer,
	}
}

func (l *transformNodeImpl[T, U]) GetValue(ctx context.Context, id int) U {
	nodeState := NodeStateFromContext(ctx)
	return nodeState.GetResolvedValue(id).(U)
}

func (l *transformNodeImpl[T, U]) IsResolved(ctx context.Context, id int) bool {
	nodeState := NodeStateFromContext(ctx)
	return nodeState.GetIsResolved(id)
}

func (l *transformNodeImpl[T, U]) GetAnyResolvables() []AnyNode {
	return []AnyNode{l.child}
}

func (l *transformNodeImpl[T, U]) Run(ctx context.Context, id int) any {
	nodeState := NodeStateFromContext(ctx)
	children := nodeState.GetChildren(id)

	result := l.fn(l.child.GetValue(ctx, children[0]))

	nodeState.SetResolvedValue(id, result)
	nodeState.SetIsResolved(id, true)

	return result
}
