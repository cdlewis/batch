package panera

import "context"

type AnyNode interface {
	IsResolved(context.Context, int) bool
	GetChildren() []AnyNode
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
