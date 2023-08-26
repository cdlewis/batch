package main

type TransformNode[T, U any] interface {
	Node[U]
	AnyNode
}

type transformNodeImpl[T, U any] struct {
	isResolved bool
	child      Node[T]
	fn         func(T) U
	result     U
}

func NewTransformNode[T, U any](node Node[T], transformer func(T) U) TransformNode[T, U] {
	return &transformNodeImpl[T, U]{
		child: node,
		fn:    transformer,
	}
}

func (l *transformNodeImpl[T, U]) GetValue() U {
	return l.result
}

func (l *transformNodeImpl[T, U]) IsResolved() bool {
	return l.isResolved
}

func (l *transformNodeImpl[T, U]) GetAnyResolvables() []AnyNode {
	return []AnyNode{l.child}
}

func (l *transformNodeImpl[T, U]) Run() any {
	l.result = l.fn(l.child.GetValue())
	l.isResolved = true
	return l.result
}
