package panera

import (
	"context"
)

type ListNode[T any] interface {
	Node[[]T]
}

type listNodeImpl[T any] struct {
	id       NodeID
	children []Node[T]
	results  []T
}

func NewListNode[T any](children []Node[T]) ListNode[T] {
	return &listNodeImpl[T]{
		id:       NewNodeID(),
		children: children,
	}
}

func (l *listNodeImpl[T]) GetValue(ctx context.Context) []T {
	nodeState := NodeStateFromContext(ctx)
	return nodeState.GetResolvedValue(l.id).([]T)
}

func (l *listNodeImpl[T]) IsResolved(ctx context.Context) bool {
	nodeState := NodeStateFromContext(ctx)
	return nodeState.GetIsResolved(l.id)
}

func (l *listNodeImpl[T]) GetChildren() []AnyNode {
	results := make([]AnyNode, 0, len(l.children))
	for _, i := range l.children {
		results = append(results, i.(AnyNode))
	}
	return results
}

func (l *listNodeImpl[T]) Run(ctx context.Context) any {
	nodeState := NodeStateFromContext(ctx)

	results := make([]T, 0, len(l.children))
	for _, c := range l.children {
		results = append(results, c.GetValue(ctx))
	}

	nodeState.SetResolvedValue(l.id, results)
	nodeState.SetIsResolved(l.id, true)

	return results
}

func (l *listNodeImpl[T]) Result() any {
	return l.results
}

func (l *listNodeImpl[T]) GetID() NodeID {
	return l.id
}

func (l *listNodeImpl[T]) Debug() string {
	return "List"
}
