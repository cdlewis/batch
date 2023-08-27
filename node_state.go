package main

import "context"

var nodeStateKey struct{}

func ContextWithNodeState(ctx context.Context) context.Context {
	return context.WithValue(ctx, nodeStateKey, &nodeState{
		resolved:      map[int]bool{},
		resolvedValue: map[int]any{},
		children:      map[int][]int{},
	})
}

func NodeStateFromContext(ctx context.Context) *nodeState {
	return ctx.Value(nodeStateKey).(*nodeState)
}

type nodeState struct {
	resolved      map[int]bool
	resolvedValue map[int]any
	children      map[int][]int
}

func (n *nodeState) GetIsResolved(id int) bool {
	return n.resolved[id]
}

func (n *nodeState) SetIsResolved(id int, state bool) {
	n.resolved[id] = state
}

func (n *nodeState) SetResolvedValue(id int, value any) {
	n.resolvedValue[id] = value
}

func (n *nodeState) GetChildren(id int) []int {
	return n.children[id]
}

func (n *nodeState) AddChildren(id int, children []int) {
	n.children[id] = append(n.children[id], children...)
}
