package main

type TaskManager struct {
	counter      int
	tasks        map[int]AnyNode
	dependencies map[int]map[int]struct{}
}

func NewTaskManager() *TaskManager {
	return &TaskManager{
		tasks:        map[int]AnyNode{},
		dependencies: map[int]map[int]struct{}{},
	}
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
