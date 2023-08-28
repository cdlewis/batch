package panera

var id int64

// NodeID is used to uniquely identify nodes. Originally I passed around
// integers but I think this is cleaner. The goal is to rely on the pointer
// reference to uniquely identify each Node.
//
// An aliased pointer has some interesting properties. It can be used as a
// map key (unlike an Interface). It should stay unique. And it can be
// recycled if it stops being used. There is some funny stuff with struct{}
// pointers though. I'm not yet convinced this is rock solid.
//
// The id int64 type is just for humans.
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
