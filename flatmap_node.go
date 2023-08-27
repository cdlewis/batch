package main

type AnyFlatMap interface {
	FlatMapSentinalFunction()
	FlatMapFullyResolved() bool
}

type FlatMapNode[T, U any] interface {
	Node[U]
	AnyFlatMap
}

type flatMapNodeImpl[T, U any] struct {
	child      Node[T]
	grandChild Node[U]
	fn         func(T) Node[U]
	isResolved bool
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

func (f *flatMapNodeImpl[T, U]) IsResolved() bool {
	return f.isResolved
}

func (f *flatMapNodeImpl[T, U]) Run() any {
	f.grandChild = f.fn(f.child.GetValue())
	return f.grandChild
}

func (f *flatMapNodeImpl[T, U]) GetValue() U {
	if grandChild := f.grandChild; grandChild != nil {
		return grandChild.GetValue()
	}

	panic("Unexpected access of node value")
}

func (f *flatMapNodeImpl[T, U]) GetAnyResolvables() []AnyNode {
	return []AnyNode{f.child}
}

func (f *flatMapNodeImpl[T, U]) FlatMapFullyResolved() bool {
	if grandChild := f.grandChild; grandChild != nil {
		return grandChild.IsResolved()
	}

	return false
}

func (f *flatMapNodeImpl[T, U]) FlatMapSentinalFunction() {
	return
}
