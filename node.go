package main

import "context"

type AnyNode interface {
	IsResolved() bool
	GetAnyResolvables() []AnyNode
	Run(context.Context) any
}

type Node[T any] interface {
	AnyNode

	GetValue() T
}

type BatchableNode interface {
	AnyNode

	ResolverID() string
}
