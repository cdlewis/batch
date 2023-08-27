package panera

import (
	"context"
)

type ListNode[T any] interface {
	Node[[]T]
}

type listNodeImpl[T any] struct {
	children []Node[T]
	results  []T
}

func NewListNode[T any](children []Node[T]) ListNode[T] {
	return &listNodeImpl[T]{
		children: children,
	}
}

func (l *listNodeImpl[T]) GetValue(ctx context.Context, id int) []T {
	nodeState := NodeStateFromContext(ctx)
	return nodeState.GetResolvedValue(id).([]T)
}

func (l *listNodeImpl[T]) IsResolved(ctx context.Context, id int) bool {
	nodeState := NodeStateFromContext(ctx)
	return nodeState.GetIsResolved(id)
}

func (l *listNodeImpl[T]) GetChildren() []AnyNode {
	results := make([]AnyNode, 0, len(l.children))
	for _, i := range l.children {
		results = append(results, i.(AnyNode))
	}
	return results
}

func (l *listNodeImpl[T]) Run(ctx context.Context, id int) any {
	nodeState := NodeStateFromContext(ctx)
	childIDs := nodeState.GetChildren(id)

	results := make([]T, 0, len(l.children))
	for idx, c := range l.children {
		results = append(results, c.GetValue(ctx, childIDs[idx]))
	}

	nodeState.SetResolvedValue(id, results)
	nodeState.SetIsResolved(id, true)

	return results
}

func (l *listNodeImpl[T]) Result() any {
	return l.results
}
