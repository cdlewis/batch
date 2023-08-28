package panera

import (
	"context"
	"fmt"
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
			taskManager.PrintDependencyTree()

			// Can we terminate?
			if taskManager.GetRootTask().IsResolved(ctx) {
				done <- true
			}

			runnableTasks := taskManager.GetRunnableTasksIDs()
			fmt.Println("Runnable tasks", runnableTasks)

			// Group tasks into batch and non-batch
			batchableTasks := map[NodeID]AnyBatchQueryNode{}
			batchableQueries := map[string]map[NodeID]any{}
			regularTasks := []NodeID{}
			for _, id := range runnableTasks {
				node := taskManager.GetTask(id)
				if batchTask, ok := node.(AnyBatchQueryNode); ok {
					resolverID := batchTask.ResolverID()
					if batchableQueries[resolverID] == nil {
						batchableQueries[resolverID] = map[NodeID]any{}
					}

					batchableTasks[id] = batchTask
					batchableQueries[resolverID][id] = batchTask.BuildQuery(ctx)
				} else {
					regularTasks = append(regularTasks, id)
				}
			}

			// Kick off async resolvers for the batch tasks

			for resolverID, taskMap := range batchableQueries {
				resolverID, taskMap := resolverID, taskMap
				go func() {
					resultsMap := resolvers[resolverID].Resolve(ctx, taskMap)
					for id, result := range resultsMap {
						fmt.Println("Resolve", id, result)
						batchableTasks[id].SetResult(ctx, result)
						taskManager.FinishTask(id)
					}

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
					fmt.Println(">> Detected flatmap", taskID)
					if flatMapNode.FlatMapFullyResolved(ctx) {
						taskManager.FinishTask(taskID)
					} else {
						newNode := currentTask.Run(ctx).(AnyNode)
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
					fmt.Println("Executing regular task", taskID)
					currentTask := taskManager.GetTask(taskID)
					currentTask.Debug()

					if !currentTask.IsResolved(ctx) {
						currentTask.Run(ctx)
					}

					taskManager.FinishTask(taskID)

					taskResolved <- struct{}{}
				}()
			}
		}
	}()

	select {
	case <-done:
		return taskManager.GetRootTask().(Node[T]).GetValue(ctx)
	case <-ctx.Done():
		panic("timeout")
	}
}
