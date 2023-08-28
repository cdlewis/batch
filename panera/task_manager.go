package panera

import (
	"context"
	"fmt"
	"sync"
)

const _parentIDNotParent = -1

type TaskManager interface {
	GetTask(NodeID) AnyNode
	UpdateTask(context.Context, NodeID, AnyNode)
	FinishTask(NodeID)
	GetRunnableTasksIDs() []NodeID
	GetRootTask() AnyNode
	PrintDependencyTree()
}

type taskManagerImpl struct {
	tasks        map[NodeID]AnyNode
	dependencies map[NodeID]map[NodeID]struct{}
	hasRun       map[NodeID]bool
	rootTask     AnyNode
	mutex        sync.RWMutex
}

func NewTaskManager(ctx context.Context, rootTask AnyNode) TaskManager {
	tm := &taskManagerImpl{
		tasks:        map[NodeID]AnyNode{},
		dependencies: map[NodeID]map[NodeID]struct{}{},
		hasRun:       map[NodeID]bool{},
		rootTask:     rootTask,
		mutex:        sync.RWMutex{},
	}

	tm.exploreTaskGraph(ctx, rootTask, nil)

	return tm
}

func (t *taskManagerImpl) UpdateTask(ctx context.Context, id NodeID, node AnyNode) {
	t.mutex.Lock()
	defer t.mutex.Unlock()

	t.exploreTaskGraph(ctx, node, id)
}

func (t *taskManagerImpl) GetTask(id NodeID) AnyNode {
	if id == nil {
		panic("Attempted to retrieve invalid task")
	}

	t.mutex.RLock()
	defer t.mutex.RUnlock()

	return t.tasks[id]
}

func (t *taskManagerImpl) FinishTask(id NodeID) {
	t.mutex.Lock()
	defer t.mutex.Unlock()

	t.hasRun[id] = true

	for parentTaskID, childTasks := range t.dependencies {
		if _, ok := childTasks[id]; ok {
			delete(t.dependencies[parentTaskID], id)
		}
	}
}

func (t *taskManagerImpl) GetRunnableTasksIDs() []NodeID {
	t.mutex.RLock()
	defer t.mutex.RUnlock()

	result := []NodeID{}

	for id := range t.tasks {
		if len(t.dependencies[id]) == 0 && !t.hasRun[id] {
			result = append(result, id)
		}
	}

	return result
}

func (t *taskManagerImpl) GetRootTask() AnyNode {
	return t.rootTask
}

func (t *taskManagerImpl) PrintDependencyTree() {
	t.mutex.RLock()
	defer t.mutex.RUnlock()

	fmt.Println("Tasks")
	for taskID, task := range t.tasks {
		fmt.Println("\t", taskID, task.Debug())
	}
	fmt.Println()
	fmt.Println("Dependency tree:")
	for nodeID, deps := range t.dependencies {
		fmt.Println("\t", nodeID, t.tasks[nodeID].Debug())
		for childID := range deps {
			fmt.Println("\t\t", childID, t.tasks[childID].Debug())
		}
	}
}

func (t *taskManagerImpl) addDependency(parentID, childID NodeID) {
	if t.dependencies[parentID] == nil {
		t.dependencies[parentID] = map[NodeID]struct{}{}
	}
	t.dependencies[parentID][childID] = struct{}{}
}

type NodeParentPair struct {
	ParentID NodeID
	Node     AnyNode
}

func (t *taskManagerImpl) exploreTaskGraph(ctx context.Context, root AnyNode, parentID NodeID) {
	stack := []NodeParentPair{
		{
			ParentID: parentID,
			Node:     root,
		},
	}

	for len(stack) > 0 {
		nextStack := []NodeParentPair{}

		for len(stack) > 0 {
			nextNode := stack[len(stack)-1]
			stack = stack[:len(stack)-1]

			t.tasks[nextNode.Node.GetID()] = nextNode.Node

			if nextNode.ParentID.IsSet() {
				t.addDependency(nextNode.ParentID, nextNode.Node.GetID())
			}

			children := nextNode.Node.GetChildren()
			for _, w := range children {
				nextStack = append(nextStack, NodeParentPair{
					ParentID: nextNode.Node.GetID(),
					Node:     w,
				})
			}
		}

		stack = nextStack
	}
}
