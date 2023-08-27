package main

import (
	"context"
	"fmt"
)

const _parentIDNotParent = -1

type TaskManager struct {
	counter      int
	tasks        map[int]AnyNode
	dependencies map[int]map[int]struct{}
	hasRun       map[int]bool
	rootTask     AnyNode
}

func NewTaskManager(ctx context.Context, rootTask AnyNode) *TaskManager {
	tm := &TaskManager{
		tasks:        map[int]AnyNode{},
		dependencies: map[int]map[int]struct{}{},
		hasRun:       map[int]bool{},
		rootTask:     rootTask,
	}

	tm.exploreTaskGraph(ctx, rootTask, _parentIDNotParent)

	return tm
}

func (t *TaskManager) UpdateTask(ctx context.Context, id int, node AnyNode) int {
	return t.exploreTaskGraph(ctx, node, id)
}

func (t *TaskManager) GetTask(id int) AnyNode {
	return t.tasks[id]
}

func (t *TaskManager) AddDependency(parentID, childID int) {
	if t.dependencies[parentID] == nil {
		t.dependencies[parentID] = map[int]struct{}{}
	}
	t.dependencies[parentID][childID] = struct{}{}
}

func (t *TaskManager) FinishTask(id int) []int {
	t.hasRun[id] = true
	unblockedTasks := []int{}

	for parentTaskID, childTasks := range t.dependencies {
		if _, ok := childTasks[id]; ok {
			delete(t.dependencies[parentTaskID], id)

			if len(t.dependencies[id]) == 0 {
				unblockedTasks = append(unblockedTasks, parentTaskID)
			}
		}
	}

	return unblockedTasks
}

func (t *TaskManager) GetRunnableTasksIDs() []int {
	result := []int{}

	for id := range t.tasks {
		if len(t.dependencies[id]) == 0 && !t.hasRun[id] {
			result = append(result, id)
		}
	}

	return result
}

func (t *TaskManager) GetRootTask() AnyNode {
	return t.rootTask
}

func (t *TaskManager) PrintDependencyTree() {
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

type NodeParentPair struct {
	NodeID   int
	ParentID int
	Node     AnyNode
}

func (t *TaskManager) exploreTaskGraph(ctx context.Context, root AnyNode, parentID int) int {
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

			fmt.Println("Exploring", nextNode)

			if nextNode.ParentID != -1 {
				newRoot = nextNode.NodeID
				t.AddDependency(nextNode.ParentID, nextNode.NodeID)
			}

			children := nextNode.Node.GetAnyResolvables()
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

			fmt.Println("Adding children", nextNode.NodeID, childNodeIDs)
			nodeState.AddChildren(nextNode.NodeID, childNodeIDs)
		}

		stack = nextStack
	}

	return newRoot
}
