package panera

import (
	"context"
)

// TransformNode represents a 'map' operation on the Node type, transforming
// the contents of a node from type T to type U.
type TransformNode[T, U any] interface {
	Node[U]
}

type transformNodeImpl[T, U any] struct {
	id    NodeID
	child Node[T]
	fn    func(T) U
}

func NewTransformNode[T, U any](node Node[T], transformer func(T) U) TransformNode[T, U] {
	return &transformNodeImpl[T, U]{
		id:    NewNodeID(),
		child: node,
		fn:    transformer,
	}
}

func (l *transformNodeImpl[T, U]) GetID() NodeID {
	return l.id
}

func (l *transformNodeImpl[T, U]) GetValue(ctx context.Context) U {
	nodeState := NodeStateFromContext(ctx)
	return nodeState.GetResolvedValue(l.id).(U)
}

func (l *transformNodeImpl[T, U]) IsResolved(ctx context.Context) bool {
	nodeState := NodeStateFromContext(ctx)
	return nodeState.GetIsResolved(l.id)
}

func (l *transformNodeImpl[T, U]) GetChildren() []AnyNode {
	return []AnyNode{l.child}
}

func (l *transformNodeImpl[T, U]) Run(ctx context.Context) any {
	nodeState := NodeStateFromContext(ctx)
	result := l.fn(l.child.GetValue(ctx))

	nodeState.SetResolvedValue(l.id, result)
	nodeState.SetIsResolved(l.id, true)

	return result
}

func (l *transformNodeImpl[T, U]) Debug() string {
	return "Transform"
}
