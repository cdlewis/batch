package panera

import (
	"context"
	"sync"
)

var nodeStateKey struct{}

func ContextWithNodeState(ctx context.Context) context.Context {
	return context.WithValue(ctx, nodeStateKey, &nodeState{
		resolved:      map[int]bool{},
		resolvedValue: map[int]any{},
		children:      map[int][]int{},
		mutex:         sync.RWMutex{},
	})
}

func NodeStateFromContext(ctx context.Context) *nodeState {
	return ctx.Value(nodeStateKey).(*nodeState)
}

type nodeState struct {
	resolved      map[int]bool
	resolvedValue map[int]any
	children      map[int][]int
	mutex         sync.RWMutex
}

func (n *nodeState) GetIsResolved(id int) bool {
	n.mutex.RLock()
	defer n.mutex.RUnlock()
	return n.resolved[id]
}

func (n *nodeState) SetIsResolved(id int, state bool) {
	n.mutex.Lock()
	defer n.mutex.Unlock()
	n.resolved[id] = state
}

func (n *nodeState) GetResolvedValue(id int) any {
	n.mutex.RLock()
	defer n.mutex.RUnlock()
	return n.resolvedValue[id]
}

func (n *nodeState) SetResolvedValue(id int, value any) {
	n.mutex.Lock()
	defer n.mutex.Unlock()
	n.resolvedValue[id] = value
}

func (n *nodeState) GetChildren(id int) []int {
	n.mutex.RLock()
	defer n.mutex.RUnlock()
	return n.children[id]
}

func (n *nodeState) AddChildren(id int, children []int) {
	n.mutex.Lock()
	defer n.mutex.Unlock()
	n.children[id] = append(n.children[id], children...)
}
