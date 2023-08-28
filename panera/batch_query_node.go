package panera

import (
	"context"
)

// BatchQueryNode is a special type of node, used as a helper for constructing
// batchable service adapters. As such, users implement a BuildQuery method instead
// of a Run function.
type BatchQueryNode[Q, R any] interface {
	Node[R]
	AnyBatchQueryNode
}

// AnyBatchQueryNode is a generic type param-free version of BatchQueryNode
type AnyBatchQueryNode interface {
	ResolverID() string
	SetResult(context.Context, any)
	BuildQuery(context.Context) any
}

type batchQueryNodeImpl[Q, R any] struct {
	BatchQueryNode[Q, R]

	id         NodeID
	queryFn    func(context.Context) Q
	resolverID string
}

func NewBatchQueryNode[Q, R any](
	queryFn func(context.Context) Q,
	resolverID string,
) BatchQueryNode[Q, R] {
	return &batchQueryNodeImpl[Q, R]{
		id:         NewNodeID(),
		queryFn:    queryFn,
		resolverID: resolverID,
	}
}

func (v *batchQueryNodeImpl[Q, R]) GetValue(ctx context.Context) R {
	nodeState := NodeStateFromContext(ctx)
	return nodeState.GetResolvedValue(v.id).(R)
}

func (v *batchQueryNodeImpl[Q, R]) IsResolved(ctx context.Context) bool {
	nodeState := NodeStateFromContext(ctx)
	return nodeState.GetIsResolved(v.id)
}

func (v *batchQueryNodeImpl[Q, R]) GetChildren() []AnyNode {
	return []AnyNode{}
}

func (v *batchQueryNodeImpl[Q, R]) Run(_ context.Context) any {
	panic("we should batch this -- you screwed up")
}

func (v *batchQueryNodeImpl[Q, R]) ResolverID() string {
	return v.resolverID
}

func (v *batchQueryNodeImpl[Q, R]) BuildQuery(ctx context.Context) any {
	return v.queryFn(ctx)
}

func (v *batchQueryNodeImpl[Q, R]) SetResult(ctx context.Context, result any) {
	nodeState := NodeStateFromContext(ctx)
	nodeState.SetResolvedValue(v.id, result)
	nodeState.SetIsResolved(v.id, true)
}

func (v *batchQueryNodeImpl[Q, R]) GetID() NodeID {
	return v.id
}

func (v *batchQueryNodeImpl[Q, R]) Debug() string {
	return "BatchNode"
}
