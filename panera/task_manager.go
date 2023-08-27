package panera

import (
	"context"
	"fmt"
	"sync"
)

const _parentIDNotParent = -1

type TaskManager interface {
	GetTask(int) AnyNode
	UpdateTask(context.Context, int, AnyNode)
	FinishTask(int)
	GetRunnableTasksIDs() []int
	GetRootTask() AnyNode
	PrintDependencyTree()
}

type taskManagerImpl struct {
	counter      int
	tasks        map[int]AnyNode
	dependencies map[int]map[int]struct{}
	hasRun       map[int]bool
	rootTask     AnyNode
	mutex        sync.RWMutex
}

func NewTaskManager(ctx context.Context, rootTask AnyNode) TaskManager {
	tm := &taskManagerImpl{
		tasks:        map[int]AnyNode{},
		dependencies: map[int]map[int]struct{}{},
		hasRun:       map[int]bool{},
		rootTask:     rootTask,
		mutex:        sync.RWMutex{},
	}

	tm.exploreTaskGraph(ctx, rootTask, _parentIDNotParent)

	return tm
}

func (t *taskManagerImpl) UpdateTask(ctx context.Context, id int, node AnyNode) {
	t.mutex.Lock()
	defer t.mutex.Unlock()

	t.exploreTaskGraph(ctx, node, id)
}

func (t *taskManagerImpl) GetTask(id int) AnyNode {
	t.mutex.RLock()
	defer t.mutex.RUnlock()

	return t.tasks[id]
}

func (t *taskManagerImpl) FinishTask(id int) {
	t.mutex.Lock()
	defer t.mutex.Unlock()

	t.hasRun[id] = true

	for parentTaskID, childTasks := range t.dependencies {
		if _, ok := childTasks[id]; ok {
			delete(t.dependencies[parentTaskID], id)
		}
	}
}

func (t *taskManagerImpl) GetRunnableTasksIDs() []int {
	t.mutex.RLock()
	defer t.mutex.RUnlock()

	result := []int{}

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
		fmt.Println("\t", taskID, task)
	}
	fmt.Println()
	fmt.Println("Dependency tree:")
	for nodeID, deps := range t.dependencies {
		fmt.Println("\t", nodeID, deps)
	}
}

func (t *taskManagerImpl) addDependency(parentID, childID int) {
	if t.dependencies[parentID] == nil {
		t.dependencies[parentID] = map[int]struct{}{}
	}
	t.dependencies[parentID][childID] = struct{}{}
}

type NodeParentPair struct {
	NodeID   int
	ParentID int
	Node     AnyNode
}

func (t *taskManagerImpl) exploreTaskGraph(ctx context.Context, root AnyNode, parentID int) int {
	nodeState := NodeStateFromContext(ctx)

	t.counter++
	newRoot := t.counter
	stack := []NodeParentPair{
		{
			NodeID:   newRoot,
			ParentID: parentID,
			Node:     root,
		},
	}
	t.tasks[newRoot] = root

	if parentID != _parentIDNotParent {
		nodeState.AddChildren(parentID, []int{newRoot})
	}

	for len(stack) > 0 {
		nextStack := []NodeParentPair{}

		for len(stack) > 0 {
			nextNode := stack[len(stack)-1]
			stack = stack[:len(stack)-1]

			if nextNode.ParentID != -1 {
				newRoot = nextNode.NodeID
				t.addDependency(nextNode.ParentID, nextNode.NodeID)
			}

			children := nextNode.Node.GetChildren()
			childNodeIDs := make([]int, 0, len(children))
			for _, w := range children {
				t.counter++
				id := t.counter
				t.tasks[id] = w

				nextStack = append(nextStack, NodeParentPair{
					ParentID: nextNode.NodeID,
					NodeID:   id,
					Node:     w,
				})

				childNodeIDs = append(childNodeIDs, id)
			}

			nodeState.AddChildren(nextNode.NodeID, childNodeIDs)
		}

		stack = nextStack
	}

	return newRoot
}
