package main

import (
	"fmt"
)

func main() {
	// userResolver := Resolver[string, string]{
	// 	id: "user",
	// 	resolve: func(requests []ResolvableValue[string]) []string {
	// 		results := []string{}
	// 		for _, r := range requests {
	// 			results = append(results, "resolved: "+r.arg.(string))
	// 		}
	// 		return results
	// 	},
	// }

	// resolvers := map[string]Resolver[any, any]{
	// 	"user": userResolver,
	// }

	// firstUser := userResolver.fetch("chris")
	// secondUser := userResolver.fetch("mike")
	users := NewTransformNode[[]string, string](
		NewListNode[string]([]AnyNode{
			NewValueNode("Chris"),
			NewValueNode("Mike"),
		}),
		func(results []string) string {
			fmt.Println("trans running")
			result := ""
			for _, i := range results {
				result += i + ","
			}
			return result
		},
	)

	fmt.Println(resolve[string](users))
}

type ExploreNextJob struct {
	ParentID int
	Node     AnyNode
}

func resolve[T any](node AnyNode) T {
	counter := 0
	tasks := map[int]AnyNode{}
	blocked := map[int][]int{}

	runNext := []int{}

	exploreNext := []ExploreNextJob{
		{
			ParentID: 0,
			Node:     node,
		},
	}
	for len(exploreNext) > 0 {
		nextNode := exploreNext[0]
		fmt.Println("Exploring", nextNode)
		exploreNext = exploreNext[1:]

		counter++
		tasks[counter] = nextNode.Node

		if nextNode.ParentID != 0 {
			blocked[nextNode.ParentID] = append(blocked[nextNode.ParentID], counter)
		}

		blockingWork := nextNode.Node.GetAnyResolvables()
		if len(blockingWork) == 0 {
			runNext = append(runNext, counter)
			continue
		}

		for _, w := range blockingWork {
			exploreNext = append(exploreNext, ExploreNextJob{
				ParentID: counter,
				Node:     w,
			})
		}
	}
	fmt.Println("task list", tasks)

	for len(runNext) > 0 {
		nextRunNext := []int{}

		for len(runNext) > 0 {
			taskID := runNext[0]
			runNext = runNext[1:]

			result := tasks[taskID].Run()
			fmt.Println(taskID, "completed with", result)

			// inject the result into each blocked task and see if it can be run now
			for blockedTaskID, blockingTasks := range blocked {
				canRun := true
				for _, t := range blockingTasks {
					if t != taskID {
						if !tasks[t].IsResolved() {
							canRun = false
						}
						continue
					}
					fmt.Println("Found", taskID, "blocks", blockedTaskID)
				}

				if canRun {
					nextRunNext = append(nextRunNext, blockedTaskID)
				}
			}
		}

		runNext = nextRunNext
	}

	fmt.Println(tasks)
	return tasks[1].Result().(T)
}

type AnyNode interface {
	IsResolved() bool
	GetAnyResolvables() []AnyNode
	Run() any
	InjectResult(any)
	Result() any
}

type AnyResolvable interface{}

func NewResolver[T, U any](id string, resolve func(T) U) Resolver[T, U] {
	return Resolver[T, U]{
		id:      id,
		resolve: resolve,
	}
}

type Resolver[T, U any] struct {
	id      string
	resolve func(T) U
}

func (r Resolver[T, U]) fetch(arg T) ResolvableValue[U] {
	return ResolvableValue[U]{
		key: r.id,
		arg: arg,
	}
}

type ResolvableValue[U any] struct {
	Node[U]

	key string
	arg any
}

func (r ResolvableValue[U]) Key() string {
	return r.key
}

type Transform[T, U any] interface {
	Apply(T) U
}

type Node[T any] interface {
	AnyNode
	GetValue() T
}

type Resolvable[T any] interface {
	Node[T]

	Key() string
}

// type FlatMapNode[T, U any] struct {
// 	base Node[T]
// 	fn   func(Node[T]) Node[U]
// }

/*
type MapNode[T any, U any] struct {
	node Node[T]
	fn   func(T) U
}

func MapNode[T, U any](node Node[T], fn func(T) U) Node[U] {
	return MapNode[T, U]{
		node: node,
		fn:   fn,
	}
}
*/

type nodeSentinal struct{}
type NodeResult interface {
	NodeResultMustBePartOfYourResultStruct(nodeSentinal)
}