package main

import "context"

type ListNode[T any] interface {
	Node[[]T]
	AnyNode
}

type listNodeImpl[T any] struct {
	isResolved bool
	children   []Node[T]
	results    []T
}

func NewListNode[T any](children []Node[T]) ListNode[T] {
	return &listNodeImpl[T]{
		children: children,
	}
}

func (l *listNodeImpl[T]) GetValue() []T {
	return l.results
}

func (l *listNodeImpl[T]) IsResolved() bool {
	return l.isResolved
}

func (l *listNodeImpl[T]) GetAnyResolvables() []AnyNode {
	results := make([]AnyNode, 0, len(l.children))
	for _, i := range l.children {
		results = append(results, i.(AnyNode))
	}
	return results
}

func (l *listNodeImpl[T]) Run(ctx context.Context) any {
	results := make([]T, 0, len(l.children))
	for _, c := range l.children {
		results = append(results, c.GetValue())
	}
	l.results = results
	return results
}

func (l *listNodeImpl[T]) Result() any {
	return l.results
}
