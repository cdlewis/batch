package main

import "fmt"

func ExecuteGraph[T any](node AnyNode, resolvers map[string]Resolver) T {
	taskManager := NewTaskManager(node)

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
			resolvers[resolverID].Resolve(taskIDs, taskManager)
		}

		for _, taskID := range regularTasks {
			currentTask := taskManager.GetTask(taskID)

			if flatMapNode, isFlatMap := currentTask.(AnyFlatMap); isFlatMap {
				if flatMapNode.FlatMapFullyResolved() {
					taskManager.FinishTask(taskID)
					break
				}
				fmt.Println("Detected", taskID, "is flatmap")
				newNode := currentTask.Run().(AnyNode)
				fmt.Println("Re-running deps")
				id := taskManager.UpdateTask(taskID, newNode)
				fmt.Println("new root task with", id)
				taskManager.PrintDependencyTree()
				continue
			}

			if !currentTask.IsResolved() {
				currentTask.Run()
			}

			fmt.Println("!!!! FINISHED", taskID)
			taskManager.FinishTask(taskID)
		}
		fmt.Println("fetching new task list")
		taskManager.PrintDependencyTree()
		runNext = taskManager.GetRunnableTasksIDs()
	}

	return taskManager.GetRootTask().(Node[T]).GetValue()
}
