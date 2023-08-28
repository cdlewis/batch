package panera

import (
	"context"
)

func ExecuteGraph[T any](
	parentCtx context.Context,
	node AnyNode,
	resolvers map[string]Resolver,
) (T, error) {
	ctx := ContextWithNodeState(parentCtx)

	nodeManager := NewNodeManager(node)

	nodeResolved := make(chan struct{}, 1)
	nodeResolved <- struct{}{}
	done := make(chan bool)

	go func() {
		// This is our main loop. We receive notifications about completed nodes
		// through the 'nodeResolved' channel, which triggers a re-evaluate if
		// any new work is executable.
		for range nodeResolved {

			// Can we terminate?
			if nodeManager.GetRootNode().IsResolved(ctx) {
				done <- true
			}

			runnableNodes := nodeManager.GetRunnableNodes()

			// Group nodes into batch and non-batch
			batchableNodes := map[NodeID]AnyBatchQueryNode{}
			batchableQueries := map[string]map[NodeID]any{}
			regularNodes := []NodeID{}
			for _, id := range runnableNodes {
				node := nodeManager.GetNodeByID(id)
				if batchNode, ok := node.(AnyBatchQueryNode); ok {
					resolverID := batchNode.ResolverID()
					if batchableQueries[resolverID] == nil {
						batchableQueries[resolverID] = map[NodeID]any{}
					}

					batchableNodes[id] = batchNode
					batchableQueries[resolverID][id] = batchNode.BuildQuery(ctx)
				} else {
					regularNodes = append(regularNodes, id)
				}
			}

			// Kick off async resolvers for the batch nodes

			for resolverID, nodeMap := range batchableQueries {
				resolverID, nodeMap := resolverID, nodeMap
				go func() {
					resultsMap := resolvers[resolverID].Resolve(ctx, nodeMap)
					for id, result := range resultsMap {
						batchableNodes[id].SetResult(ctx, result)
						nodeManager.RemoveNodeAsDependency(id)
					}

					nodeResolved <- struct{}{}
				}()
			}

			// Kick off individual nodes

			for _, nodeID := range regularNodes {
				nodeID := nodeID
				currentNode := nodeManager.GetNodeByID(nodeID)

				// FlatMaps are a special case because they are the only node that can
				// produce new nodes. These new nodes then have to be added to the
				// node manager.
				//
				// If we do this in parallel we risk re-evaluating node candidates before
				// the new node has been scheduled (a race condition). This is fixable
				// but requires (more) special handling so for the moment we just evaluate the
				// flatMap's transform function syncronously.
				if flatMapNode, isFlatMap := currentNode.(AnyFlatMap); isFlatMap {
					if flatMapNode.FlatMapFullyResolved(ctx) {
						nodeManager.RemoveNodeAsDependency(nodeID)
					} else {
						newNode := currentNode.Run(ctx).(AnyNode)
						nodeManager.AttachNode(nodeID, newNode)
					}

					// Technically we don't know if a node has been resolved. We 'kick' the
					// scheduler here to avoid deadlocks.
					go func() {
						nodeResolved <- struct{}{}
					}()
					continue
				}

				go func() {
					currentNode := nodeManager.GetNodeByID(nodeID)
					currentNode.Debug()

					if !currentNode.IsResolved(ctx) {
						currentNode.Run(ctx)
					}

					nodeManager.RemoveNodeAsDependency(nodeID)

					nodeResolved <- struct{}{}
				}()
			}
		}
	}()

	// Wait until all nodes are resolved or we time out
	select {
	case <-done:
		return nodeManager.GetRootNode().(Node[T]).GetValue(ctx), nil
	case <-ctx.Done():
		return *new(T), ctx.Err()
	}
}
