package main

import "fmt"

type TransformNode[T, U any] interface {
	Node[U]
	AnyNode
}

type transformNodeImpl[T, U any] struct {
	isResolved bool
	child      AnyNode
	fn         func(any) any
	result     any
}

func NewTransformNode[T, U any](node Node[T], transformer func(T) U) TransformNode[T, U] {
	return &transformNodeImpl[T, U]{
		child: node,
		fn: func(arg any) any {
			fmt.Println("<<<< adapter running")
			return transformer(arg.(T))
		},
	}
}

func (l *transformNodeImpl[T, U]) GetValue() U {
	return *new(U)
}

func (l *transformNodeImpl[T, U]) IsResolved() bool {
	return l.isResolved
}

func (l *transformNodeImpl[T, U]) GetAnyResolvables() []AnyNode {
	fmt.Println("@@@ res")
	return []AnyNode{l.child}
}

func (l *transformNodeImpl[T, U]) Run() any {
	fmt.Println(">>> running final transform")
	l.result = l.fn(l.child.Result())
	l.isResolved = true
	return l.result
}

func (l *transformNodeImpl[T, U]) InjectResult(r any) {
	return
}

func (l *transformNodeImpl[T, U]) Result() any {
	return l.result
}
