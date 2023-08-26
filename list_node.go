package main

import "fmt"

type ListNode[T any] interface {
	Node[[]T]
	AnyNode
}

type listNodeImpl[T any] struct {
	isResolved bool
	children   []AnyNode
	results    []any
}

func NewListNode[T any](children []AnyNode) ListNode[T] {
	return &listNodeImpl[T]{
		children: children,
	}
}

func (l *listNodeImpl[T]) GetValue() []T {
	return []T{}
}

func (l *listNodeImpl[T]) IsResolved() bool {
	return l.isResolved
}

func (l *listNodeImpl[T]) GetAnyResolvables() []AnyNode {
	return l.children
}

func (l *listNodeImpl[T]) Run() any {
	results := make([]any, 0, len(l.children))
	fmt.Println("\t children", l.children)
	for _, c := range l.children {
		fmt.Println("\t list child result", c.Result())
		results = append(results, c.Result())
	}
	l.results = results
	return results
}

func (l *listNodeImpl[T]) InjectResult(r any) {
	return
}

func (l *listNodeImpl[T]) Result() any {
	return l.results
}
