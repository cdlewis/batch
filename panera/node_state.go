package panera

import (
	"context"
	"fmt"
	"sync"
)

var nodeStateKey struct{}

func ContextWithNodeState(ctx context.Context) context.Context {
	nodeState := &nodeState{
		resolved:      map[NodeID]bool{},
		resolvedValue: map[NodeID]any{},
		children:      map[NodeID][]NodeID{},
		mutex:         sync.RWMutex{},
	}

	return context.WithValue(ctx, nodeStateKey, nodeState)
}

func NodeStateFromContext(ctx context.Context) *nodeState {
	return ctx.Value(nodeStateKey).(*nodeState)
}

type nodeState struct {
	resolved      map[NodeID]bool
	resolvedValue map[NodeID]any
	children      map[NodeID][]NodeID
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
	fmt.Println("Getting value for %v", id, n.resolvedValue[id])
	return n.resolvedValue[id]
}

func (n *nodeState) SetResolvedValue(id NodeID, value any) {
	n.mutex.Lock()
	defer n.mutex.Unlock()
	fmt.Println("Storing value", id, value)
	n.resolvedValue[id] = value
}

func (n *nodeState) GetChildren(id NodeID) []NodeID {
	n.mutex.RLock()
	defer n.mutex.RUnlock()
	return n.children[id]
}

func (n *nodeState) AddChildren(id NodeID, children []NodeID) {
	n.mutex.Lock()
	defer n.mutex.Unlock()
	n.children[id] = append(n.children[id], children...)
}
