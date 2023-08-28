package panera

var id int64

type NodeID = *nodeIDImpl

type nodeIDImpl struct {
	// Note this is purely for debugging purposes. Task execution
	// does not make use of it.
	id int64
}

func NewNodeID() NodeID {
	id++
	return &nodeIDImpl{
		id: id,
	}
}

func (n *nodeIDImpl) IsSet() bool {
	return n != nil
}
