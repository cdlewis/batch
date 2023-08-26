package main

import "fmt"

type TaskManager struct {
	counter      int
	tasks        map[int]AnyNode
	dependencies map[int]map[int]struct{}
	hasRun       map[int]bool
	rootTask     AnyNode
}

func NewTaskManager(rootTask AnyNode) *TaskManager {
	tm := &TaskManager{
		tasks:        map[int]AnyNode{},
		dependencies: map[int]map[int]struct{}{},
		hasRun:       map[int]bool{},
		rootTask:     rootTask,
	}

	tm.exploreTaskGraph()

	return tm
}

func (t *TaskManager) AddTask(node AnyNode) int {
	t.counter++
	t.tasks[t.counter] = node
	return t.counter
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

type NodeParentPair struct {
	ParentID int
	Node     AnyNode
}

func (t *TaskManager) exploreTaskGraph() {
	stack := []NodeParentPair{
		{
			ParentID: -1,
			Node:     t.rootTask,
		},
	}

	for len(stack) > 0 {
		nextStack := []NodeParentPair{}

		for len(stack) > 0 {
			nextNode := stack[len(stack)-1]
			stack = stack[:len(stack)-1]

			fmt.Println("Exploring", nextNode)

			currentNodeID := t.AddTask(nextNode.Node)

			if nextNode.ParentID != -1 {
				t.AddDependency(nextNode.ParentID, currentNodeID)
			}

			blockingWork := nextNode.Node.GetAnyResolvables()
			for _, w := range blockingWork {
				nextStack = append(nextStack, NodeParentPair{
					ParentID: currentNodeID,
					Node:     w,
				})
			}
		}

		stack = nextStack
	}
}
