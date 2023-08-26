package main

type AnyNode interface {
	IsResolved() bool
	GetAnyResolvables() []AnyNode
	Run() any
}

type Node[T any] interface {
	AnyNode

	GetValue() T
}

type BatchableNode interface {
	AnyNode

	ResolverID() string
}
