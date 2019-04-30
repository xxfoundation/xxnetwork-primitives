////////////////////////////////////////////////////////////////////////////////
// Copyright © 2018 Privategrity Corporation                                   /
//                                                                             /
// All rights reserved.                                                        /
////////////////////////////////////////////////////////////////////////////////

package id

// Length of node IDs in bytes
const NodeIdLen = 32

// Node ID array
type NodeID [NodeIdLen]byte

// SetBytes sets the bytes of the node ID to the provided byte slice and returns
// it if the byte slice has the correct length. Otherwise, returns a user ID
// that is all zeroes.
func (n *NodeID) SetBytes(data []byte) *NodeID {
	if len(data) != NodeIdLen {
		return new(NodeID)
	} else {
		copy(n[:], data)
		return n
	}
}

// Bytes converts a node ID to a byte slice.
func (n *NodeID) Bytes() []byte {
	return n[:]
}

// Equals determines whether two node IDs are the same.
func (n *NodeID) Cmp(y *NodeID) bool {
	return *n == *y
}

// DeepCopy creates a completely new copy of the node ID.
func (n *NodeID) DeepCopy() *NodeID {
	if n == nil {
		return nil
	}

	var newNode NodeID
	copy(newNode[:], (*n)[:])

	return &newNode
}
