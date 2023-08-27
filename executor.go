package main

import (
	"context"
	"fmt"
)

func ExecuteGraph[T any](parentCtx context.Context, node AnyNode, resolvers map[string]Resolver) T {
	ctx := ContextWithNodeState(parentCtx)

	taskManager := NewTaskManager(ctx, node)
	taskManager.PrintDependencyTree()

	runNext := taskManager.GetRunnableTasksIDs()

	for len(runNext) > 0 {
		batchableTasks := map[string][]int{}
		regularTasks := []int{}

		for _, id := range runNext {
			node := taskManager.GetTask(id)
			if bachTask, ok := node.(BatchableNode); ok {
				batchableTasks[bachTask.ResolverID()] = append(batchableTasks[bachTask.ResolverID()], id)
			} else {
				regularTasks = append(regularTasks, id)
			}
		}

		for resolverID, taskIDs := range batchableTasks {
			resolvers[resolverID].Resolve(ctx, taskIDs, taskManager)
		}

		for _, taskID := range regularTasks {
			currentTask := taskManager.GetTask(taskID)
			fmt.Println("Getting current task", taskID, currentTask)

			if flatMapNode, isFlatMap := currentTask.(AnyFlatMap); isFlatMap {
				if flatMapNode.FlatMapFullyResolved(ctx, taskID) {
					taskManager.FinishTask(taskID)
					break
				}
				fmt.Println("Detected", taskID, "is flatmap")
				newNode := currentTask.Run(ctx, taskID).(AnyNode)
				fmt.Println("Re-running deps")
				id := taskManager.UpdateTask(ctx, taskID, newNode)
				fmt.Println("new root task with", id)
				taskManager.PrintDependencyTree()
				continue
			}

			if !currentTask.IsResolved(ctx, taskID) {
				currentTask.Run(ctx, taskID)
			}

			fmt.Println("!!!! FINISHED", taskID)
			taskManager.FinishTask(taskID)
		}
		fmt.Println("fetching new task list")
		taskManager.PrintDependencyTree()
		runNext = taskManager.GetRunnableTasksIDs()
	}

	return taskManager.GetRootTask().(Node[T]).GetValue(ctx, 1)
}
