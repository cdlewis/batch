package main

import "context"

type AnyNode interface {
	IsResolved(context.Context, int) bool
	GetAnyResolvables() []AnyNode
	Run(context.Context, int) any
}

type Node[T any] interface {
	AnyNode

	GetValue(context.Context, int) T
}

type BatchableNode interface {
	AnyNode

	ResolverID() string
}
