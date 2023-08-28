package panera

import "context"

// Resolver represents our batching interface. It allows you to provide
// a resolution strategy for a group of nodes. It accepts a map of NodeID->Request
// and expects you to return a map of NodeID->Response.
type Resolver interface {
	ID() string

	Resolve(context.Context, map[NodeID]any) map[NodeID]any
}
