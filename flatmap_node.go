package main

/*
type FlatMapNode[T, U any] interface {
	Node[U]
}

type flatMapNodeImpl[T, U any] struct {
	isResolved bool
	child      Node[T]
	fn         func(T) Node[U]
	result     U
}

func NewFlatMapNode[T, U any](
	node Node[T],
	transformer func(T) Node[U],
) FlatMapNode[T, U] {
	return flatMapNodeImpl[T, U]{}
}

func (f flatMapNodeImpl[T, U]) GetValue() U {
	return f.result
}

func (f flatMapNodeImpl[T, U]) IsResolved() bool {
	return f.isResolved
}
*/
