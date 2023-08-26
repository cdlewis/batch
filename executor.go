package main

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
			if !currentTask.IsResolved() {
				currentTask.Run()
			}

			taskManager.FinishTask(taskID)
		}

		runNext = taskManager.GetRunnableTasksIDs()
	}

	return taskManager.GetRootTask().(Node[T]).GetValue()
}
