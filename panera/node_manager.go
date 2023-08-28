package panera

import (
	"fmt"
	"sync"
)

// NodeManager manages Nodes for the lifetime of a single request. It is
// responsible for identifying and tracking dependencies. It can be queried
// by the executor to find new executable nodes.
type NodeManager interface {
	GetNodeByID(NodeID) AnyNode
	AttachNode(NodeID, AnyNode)
	RemoveNodeAsDependency(NodeID)
	GetRunnableNodes() []NodeID
	GetRootNode() AnyNode
	PrintDependencyTree()
}

type nodeManagerImpl struct {
	tasks        map[NodeID]AnyNode
	dependencies map[NodeID]map[NodeID]struct{}
	hasRun       map[NodeID]bool
	rootTask     AnyNode
	// Not particularly efficient. There is a lot of room for improvement
	// in terms of reducing lock contention but a giant lock works for a
	// proof-of-concept.
	mutex sync.RWMutex
}

func NewNodeManager(rootTask AnyNode) NodeManager {
	tm := &nodeManagerImpl{
		tasks:        map[NodeID]AnyNode{},
		dependencies: map[NodeID]map[NodeID]struct{}{},
		hasRun:       map[NodeID]bool{},
		rootTask:     rootTask,
		mutex:        sync.RWMutex{},
	}

	tm.exploreNodeGraph(rootTask, nil)

	return tm
}

func (t *nodeManagerImpl) AttachNode(id NodeID, node AnyNode) {
	t.mutex.Lock()
	defer t.mutex.Unlock()

	t.exploreNodeGraph(node, id)
}

func (t *nodeManagerImpl) GetNodeByID(id NodeID) AnyNode {
	if id == nil {
		panic("Attempted to retrieve invalid task")
	}

	t.mutex.RLock()
	defer t.mutex.RUnlock()

	return t.tasks[id]
}

func (t *nodeManagerImpl) RemoveNodeAsDependency(id NodeID) {
	t.mutex.Lock()
	defer t.mutex.Unlock()

	t.hasRun[id] = true

	for parentTaskID, childTasks := range t.dependencies {
		if _, ok := childTasks[id]; ok {
			delete(t.dependencies[parentTaskID], id)
		}
	}
}

func (t *nodeManagerImpl) GetRunnableNodes() []NodeID {
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

func (t *nodeManagerImpl) GetRootNode() AnyNode {
	return t.rootTask
}

func (t *nodeManagerImpl) PrintDependencyTree() {
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

func (t *nodeManagerImpl) addDependency(parentID, childID NodeID) {
	if t.dependencies[parentID] == nil {
		t.dependencies[parentID] = map[NodeID]struct{}{}
	}
	t.dependencies[parentID][childID] = struct{}{}
}

type NodeParentPair struct {
	ParentID NodeID
	Node     AnyNode
}

func (t *nodeManagerImpl) exploreNodeGraph(root AnyNode, parentID NodeID) {
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
