package panera

import "context"

// Node represents a unit of computation. 
type Node[T any] interface {
	AnyNode

	GetValue(context.Context) T
}

type AnyNode interface {
	GetID() NodeID
	IsResolved(context.Context) bool
	GetChildren() []AnyNode
	Run(context.Context) any
	Debug() string
}
