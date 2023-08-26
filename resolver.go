package main

type Resolver interface {
	ID() string

	Resolve([]int, *TaskManager)
}
