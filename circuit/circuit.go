package circuit

import (
	jww "github.com/spf13/jwalterweatherman"
	"gitlab.com/elixxir/primitives/id"
)

type Circuit struct {
	nodes       []*id.Node
	nodeIndexes map[id.Node]int
}

// New makes a list of node addresses for use.  It finds
// the passed "myId" and denotes it internally for use with
// utility functions.  The nodeID are copied instead of linked
// to ensure any modification of them does not change the
// Circuit structure.  Will panic if the length of the passed
// list is zero.
func New(list []*id.Node) *Circuit {
	c := Circuit{
		nodes:       make([]*id.Node, 0),
		nodeIndexes: make(map[id.Node]int),
	}

	if len(list) == 0 {
		jww.FATAL.Panicf("Cannot build a Circuit of len 0")
	}

	for index, nid := range list {
		if _, ok := c.nodeIndexes[*nid]; ok {
			jww.FATAL.Panicf("NodeIDs must be unique for the circuit.Circuit, "+
				"%s passed multiple times", nid)
		}

		c.nodeIndexes[*nid] = index
		c.nodes = append(c.nodes, nid.DeepCopy())
	}

	return &c
}

// GetNodeLocation returns the location of the passed node in the list.
// Returns -1 if the node is not in the list
func (c *Circuit) GetNodeLocation(node *id.Node) int {
	nodeLoc, ok := c.nodeIndexes[*node]

	if !ok {
		return -1
	}

	return nodeLoc
}

// GetNodeAtIndex returns the node at the given index.  Panics
// if the index does not exist within the circuit
func (c *Circuit) GetNodeAtIndex(index int) *id.Node {
	if index < 0 || index >= len(c.nodes) {
		jww.FATAL.Panicf("Cannot get an index %v which is outside"+
			" the Circut (len=%v)", index, len(c.nodes))
	}
	return c.nodes[index].DeepCopy()
}

// Len returns the number of nodes in the circuit
func (c *Circuit) Len() int {
	return len(c.nodes)
}

// GetNextNode gets the node following the passed node in
// the list. It wraps around to the beginning of the list
// if the passed node is the last node.
func (c *Circuit) GetNextNode(from *id.Node) *id.Node {
	loc := c.GetNodeLocation(from)

	if loc == -1 {
		jww.FATAL.Panicf("Cannot get the next node in the circuit.Circut"+
			"for node %s which is not present", from)
	}

	return c.nodes[(loc+1)%len(c.nodes)].DeepCopy()
}

// GetNextNode gets the node preceding the passed node in
// the list. It wraps around to the end of the list
// if the passed node is the first node.
func (c *Circuit) GetPrevNode(from *id.Node) *id.Node {
	loc := c.GetNodeLocation(from)

	if loc == -1 {
		jww.FATAL.Panicf("Cannot get the previous node in the circuit.Circut"+
			"for node %s which is not present", from)
	}

	var prevLoc int
	if loc == 0 {
		prevLoc = len(c.nodes) - 1
	} else {
		prevLoc = loc - 1
	}

	return c.nodes[prevLoc].DeepCopy()
}

// IsFirstNode returns true if the passed node is the
// first node, false otherwise
func (c *Circuit) IsFirstNode(node *id.Node) bool {
	return c.GetNodeLocation(node) == 0
}

// IsLastNode returns true if the passed node is the
// last node, false otherwise
func (c *Circuit) IsLastNode(node *id.Node) bool {
	return c.GetNodeLocation(node) == c.Len()-1
}
