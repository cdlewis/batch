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
		NewListNode([]Node[string]{
			NewValueNode("Chris"),
			NewValueNode("Mike"),
		}),
		func(results []string) string {
			result := ""
			for _, i := range results {
				result += (i + ",")
			}
			fmt.Println("trans running", results)
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
	taskManager := NewTaskManager()

	runNext := []int{}

	exploreNext := []ExploreNextJob{
		{
			ParentID: -1,
			Node:     node,
		},
	}
	for len(exploreNext) > 0 {
		nextNode := exploreNext[0]
		fmt.Println("Exploring", nextNode)
		exploreNext = exploreNext[1:]

		currentNodeID := taskManager.AddTask(nextNode.Node)

		if nextNode.ParentID != -1 {
			taskManager.AddDependency(nextNode.ParentID, currentNodeID)
		}

		blockingWork := nextNode.Node.GetAnyResolvables()
		if len(blockingWork) == 0 {
			runNext = append(runNext, currentNodeID)
			continue
		}

		for _, w := range blockingWork {
			exploreNext = append(exploreNext, ExploreNextJob{
				ParentID: currentNodeID,
				Node:     w,
			})
		}
	}

	fmt.Println(taskManager)
	for len(runNext) > 0 {
		nextRunNext := []int{}

		for len(runNext) > 0 {
			taskID := runNext[0]
			runNext = runNext[1:]
			fmt.Println(taskID)
			currentTask := taskManager.GetTask(taskID)
			if !currentTask.IsResolved() {
				currentTask.Run()
			}
			fmt.Println(taskID, "completed with", currentTask.Result())

			unblockedTasks := taskManager.FinishTask(taskID)
			nextRunNext = append(nextRunNext, unblockedTasks...)
		}

		runNext = nextRunNext
	}

	return taskManager.GetTask(1).Result().(T)
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
