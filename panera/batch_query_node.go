package panera

import (
	"context"
)

type AnyBatchQueryNode interface {
	ResolverID() string
	SetResult(context.Context, int, any)
	BuildQuery(context.Context) any
}

type BatchQueryNode[Q, R any] interface {
	Node[R]
	AnyBatchQueryNode
}

type batchQueryNodeImpl[Q, R any] struct {
	BatchQueryNode[Q, R]

	queryFn    func(context.Context) Q
	resolverID string
}

func NewBatchQueryNode[Q, R any](
	queryFn func(context.Context) Q,
	resolverID string,
) BatchQueryNode[Q, R] {
	return &batchQueryNodeImpl[Q, R]{
		queryFn:    queryFn,
		resolverID: resolverID,
	}
}

func (v *batchQueryNodeImpl[Q, R]) GetValue(ctx context.Context, id int) R {
	nodeState := NodeStateFromContext(ctx)
	return nodeState.GetResolvedValue(id).(R)
}

func (v *batchQueryNodeImpl[Q, R]) IsResolved(ctx context.Context, id int) bool {
	nodeState := NodeStateFromContext(ctx)
	return nodeState.GetIsResolved(id)
}

func (v *batchQueryNodeImpl[Q, R]) GetChildren() []AnyNode {
	return []AnyNode{}
}

func (v *batchQueryNodeImpl[Q, R]) Run(_ context.Context, id int) any {
	panic("we should batch this -- you screwed up")
}

func (v *batchQueryNodeImpl[Q, R]) ResolverID() string {
	return v.resolverID
}

func (v *batchQueryNodeImpl[Q, R]) BuildQuery(ctx context.Context) any {
	return v.queryFn(ctx)
}

func (v *batchQueryNodeImpl[Q, R]) SetResult(ctx context.Context, id int, result any) {
	nodeState := NodeStateFromContext(ctx)
	nodeState.SetResolvedValue(id, result)
	nodeState.SetIsResolved(id, true)
}
