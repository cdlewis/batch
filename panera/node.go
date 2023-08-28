package panera

import "context"

type AnyNode interface {
	GetID() NodeID
	IsResolved(context.Context) bool
	GetChildren() []AnyNode
	Run(context.Context) any
	Debug() string
}

type Node[T any] interface {
	AnyNode

	GetValue(context.Context) T
}
