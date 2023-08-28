package panera

import (
	"context"
	"sync"
)

// NodeState tracks the state of each node in the graph. Our goal is to avoid
// state existing on the nodes themselves since this wouldn't make it safe to
// run the same graph multiple times.
//
// Ultimately there may be better ways of accomplishing this goal. Forcing nodes
// to all read/write through a single mutex may create unnecessary lock contention
// vs having them, e.g. write to context directly.
//
// But for a proof-of-concept, this is fine.
type NodeState interface {
	GetIsResolved(NodeID) bool
	SetIsResolved(NodeID, bool)
	GetResolvedValue(NodeID) any
	SetResolvedValue(NodeID, any)
}

var nodeStateKey struct{}

func ContextWithNodeState(ctx context.Context) context.Context {
	nodeState := &nodeStateImpl{
		resolved:      map[NodeID]bool{},
		resolvedValue: map[NodeID]any{},

		// Not particularly efficient. There is a lot of room for improvement
		// in terms of reducing lock contention but a giant lock works for a
		// proof-of-concept.
		mutex: sync.RWMutex{},
	}

	return context.WithValue(ctx, nodeStateKey, nodeState)
}

func NodeStateFromContext(ctx context.Context) NodeState {
	return ctx.Value(nodeStateKey).(NodeState)
}

type nodeStateImpl struct {
	resolved      map[NodeID]bool
	resolvedValue map[NodeID]any
	mutex         sync.RWMutex
}

func (n *nodeStateImpl) GetIsResolved(id NodeID) bool {
	n.mutex.RLock()
	defer n.mutex.RUnlock()
	return n.resolved[id]
}

func (n *nodeStateImpl) SetIsResolved(id NodeID, state bool) {
	n.mutex.Lock()
	defer n.mutex.Unlock()
	n.resolved[id] = state
}

func (n *nodeStateImpl) GetResolvedValue(id NodeID) any {
	n.mutex.RLock()
	defer n.mutex.RUnlock()

	return n.resolvedValue[id]
}

func (n *nodeStateImpl) SetResolvedValue(id NodeID, value any) {
	n.mutex.Lock()
	defer n.mutex.Unlock()

	n.resolvedValue[id] = value
}
