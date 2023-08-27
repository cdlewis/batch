package main

import (
	"context"
)

func ExecuteGraph[T any](parentCtx context.Context, node AnyNode, resolvers map[string]Resolver) T {
	ctx := ContextWithNodeState(parentCtx)

	taskManager := NewTaskManager(ctx, node)

	taskResolved := make(chan struct{}, 1)
	taskResolved <- struct{}{}
	done := make(chan bool)

	go func() {
		// This is our main loop. We receive notifications about completed tasks
		// through the 'taskResolved' channel, which triggers a re-evaluate if
		// any new work is executable.
		for range taskResolved {

			// Can we terminate?
			if taskManager.rootTask.IsResolved(ctx, 1) {
				done <- true
			}

			runnableTasks := taskManager.GetRunnableTasksIDs()

			// Group tasks into batch and non-batch

			batchableTasks := map[string][]int{}
			regularTasks := []int{}
			for _, id := range runnableTasks {
				node := taskManager.GetTask(id)
				if bachTask, ok := node.(BatchableNode); ok {
					batchableTasks[bachTask.ResolverID()] = append(batchableTasks[bachTask.ResolverID()], id)
				} else {
					regularTasks = append(regularTasks, id)
				}
			}

			// Kick off async resolvers for the batch tasks

			for resolverID, taskIDs := range batchableTasks {
				resolverID, taskIDs := resolverID, taskIDs
				go func() {
					resolvers[resolverID].Resolve(ctx, taskIDs, taskManager)
					taskResolved <- struct{}{}
				}()
			}

			// Kick off individual tasks

			for _, taskID := range regularTasks {
				taskID := taskID
				currentTask := taskManager.GetTask(taskID)

				// FlatMaps are a special case because they are the only node that can
				// produce new nodes. These new nodes then have to be added to the
				// task manager.
				//
				// If we do this in parallel we risk re-evaluating task candidates before
				// the new task has been scheduled (a race condition). This is fixable
				// but requires special handling so for the moment we just evaluate the
				// flatMap's transform function syncronously.
				if flatMapNode, isFlatMap := currentTask.(AnyFlatMap); isFlatMap {
					if flatMapNode.FlatMapFullyResolved(ctx, taskID) {
						taskManager.FinishTask(taskID)
					} else {
						newNode := currentTask.Run(ctx, taskID).(AnyNode)
						taskManager.UpdateTask(ctx, taskID, newNode)
					}

					// Technically we don't know if a task has been resolved. We 'kick' the
					// scheduler here to avoid deadlocks.
					go func() {
						taskResolved <- struct{}{}
					}()
					continue
				}

				go func() {
					currentTask := taskManager.GetTask(taskID)

					if !currentTask.IsResolved(ctx, taskID) {
						currentTask.Run(ctx, taskID)
					}

					taskManager.FinishTask(taskID)

					taskResolved <- struct{}{}
				}()
			}
		}
	}()

	select {
	case <-done:
		return taskManager.GetRootTask().(Node[T]).GetValue(ctx, 1)
	case <-ctx.Done():
		panic("timeout")
	}
}
