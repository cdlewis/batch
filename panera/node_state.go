package panera

import (
	"context"
	"sync"
)

var nodeStateKey struct{}

func ContextWithNodeState(ctx context.Context) context.Context {
	nodeState := &nodeState{
		resolved:      map[NodeID]bool{},
		resolvedValue: map[NodeID]any{},

		// Not particularly efficient. There is a lot of room for improvement
		// in terms of reducing lock contention but a giant lock works for a
		// proof-of-concept.
		mutex: sync.RWMutex{},
	}

	return context.WithValue(ctx, nodeStateKey, nodeState)
}

func NodeStateFromContext(ctx context.Context) *nodeState {
	return ctx.Value(nodeStateKey).(*nodeState)
}

type nodeState struct {
	resolved      map[NodeID]bool
	resolvedValue map[NodeID]any
	mutex         sync.RWMutex
}

func (n *nodeState) GetIsResolved(id NodeID) bool {
	n.mutex.RLock()
	defer n.mutex.RUnlock()
	return n.resolved[id]
}

func (n *nodeState) SetIsResolved(id NodeID, state bool) {
	n.mutex.Lock()
	defer n.mutex.Unlock()
	n.resolved[id] = state
}

func (n *nodeState) GetResolvedValue(id NodeID) any {
	n.mutex.RLock()
	defer n.mutex.RUnlock()

	return n.resolvedValue[id]
}

func (n *nodeState) SetResolvedValue(id NodeID, value any) {
	n.mutex.Lock()
	defer n.mutex.Unlock()

	n.resolvedValue[id] = value
}
