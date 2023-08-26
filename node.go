package main

type Node[T any] interface {
	AnyNode
	GetValue() T
}

type AnyNode interface {
	IsResolved() bool
	GetAnyResolvables() []AnyNode
	Run() any
}