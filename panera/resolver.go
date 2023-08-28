package panera

import "context"

type Resolver interface {
	ID() string

	Resolve(context.Context, map[NodeID]any) map[NodeID]any
}
