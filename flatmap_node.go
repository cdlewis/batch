package main

import (
	"context"
	"fmt"
)

type AnyFlatMap interface {
	FlatMapSentinalFunction()
	FlatMapFullyResolved(context.Context, int) bool
}

type FlatMapNode[T, U any] interface {
	Node[U]
	AnyFlatMap
}

type flatMapNodeImpl[T, U any] struct {
	child Node[T]
	fn    func(T) Node[U]
}

func NewFlatMapNode[T, U any](
	node Node[T],
	transformer func(T) Node[U],
) FlatMapNode[T, U] {
	return &flatMapNodeImpl[T, U]{
		child: node,
		fn:    transformer,
	}
}

func (f *flatMapNodeImpl[T, U]) IsResolved(ctx context.Context, id int) bool {
	nodeState := NodeStateFromContext(ctx)
	return nodeState.GetIsResolved(id)
}

func (f *flatMapNodeImpl[T, U]) Run(ctx context.Context, id int) any {
	nodeState := NodeStateFromContext(ctx)

	if _, ok := nodeState.GetResolvedValue(id).(Node[U]); ok {
		panic("Invariant violation: flatMap run multiple times")
	}

	children := nodeState.GetChildren(id)
	fmt.Println("children", children)
	result := f.fn(f.child.GetValue(ctx, children[0]))
	nodeState.SetResolvedValue(id, result)

	return result
}

func (f *flatMapNodeImpl[T, U]) GetValue(ctx context.Context, id int) U {
	nodeState := NodeStateFromContext(ctx)

	grandChild, ok := nodeState.GetResolvedValue(id).(Node[U])
	if !ok {
		panic("Invariant violation: node state does not contain a valid node")
	}

	children := nodeState.GetChildren(id)
	if len(children) != 2 {
		panic("Invariant violation: node should contain exactly two children")
	}

	return grandChild.GetValue(ctx, children[1])
}

func (f *flatMapNodeImpl[T, U]) GetAnyResolvables() []AnyNode {
	return []AnyNode{f.child}
}

func (f *flatMapNodeImpl[T, U]) FlatMapFullyResolved(ctx context.Context, id int) bool {
	nodeState := NodeStateFromContext(ctx)

	grandChild, ok := nodeState.GetResolvedValue(id).(Node[U])
	fmt.Println("@@", grandChild, ok)
	if !ok {
		return false
	}

	children := nodeState.GetChildren(id)
	fmt.Println("children", children)
	if len(children) != 2 {
		return false
	}

	return grandChild.IsResolved(ctx, children[1])
}

func (f *flatMapNodeImpl[T, U]) FlatMapSentinalFunction() {
	return
}
