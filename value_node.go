package main

type ValueNode[T any] struct {
	Node[T]
	AnyNode

	value T
}

func NewValueNode[T any](value T) ValueNode[T] {
	return ValueNode[T]{value: value}
}

func (v ValueNode[T]) GetValue() T {
	return v.value
}

func (v ValueNode[T]) IsResolved() bool {
	return true
}

func (v ValueNode[T]) GetAnyResolvables() []AnyNode {
	return []AnyNode{}
}

func (v ValueNode[T]) Run() any {
	return v.value
}

func (v ValueNode[T]) InjectResult(r any) {
	return
}

func (v ValueNode[T]) Result() any {
	return v.value
}
