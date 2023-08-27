package panera

import "context"

type Resolver interface {
	ID() string

	Resolve(context.Context, map[int]any) map[int]any
}
