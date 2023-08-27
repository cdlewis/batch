package main

import "context"

type Resolver interface {
	ID() string

	Resolve(context.Context, []int, *TaskManager)
}
