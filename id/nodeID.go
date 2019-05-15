////////////////////////////////////////////////////////////////////////////////
// Copyright © 2018 Privategrity Corporation                                   /
//                                                                             /
// All rights reserved.                                                        /
////////////////////////////////////////////////////////////////////////////////

package id

import "encoding/base64"

// Length of node IDs in bytes
const NodeIdLen = 32

// Node ID array
type Node [NodeIdLen]byte

// NewNodeFromBytes returns a new Node ID from bytes slice if
// the byte slice has the correct length. Otherwise, it returns
// a user ID that is all zeroes.
func NewNodeFromBytes(data []byte) *Node {
	node := new(Node)
	if len(data) != NodeIdLen {
		return node
	} else {
		copy(node[:], data)
		return node
	}
}

// Bytes converts a node ID to a byte slice.
func (n *Node) Bytes() []byte {
	return n[:]
}

// Equals determines whether two node IDs are the same.
func (n *Node) Cmp(y *Node) bool {
	return *n == *y
}

// DeepCopy creates a completely new copy of the node ID.
func (n *Node) DeepCopy() *Node {
	if n == nil {
		return nil
	}

	var newNode Node
	copy(newNode[:], (*n)[:])

	return &newNode
}

// String() implements Stringer, and allows node IDs to be used as connection IDs
// Currently, it just base64 encodes the node ID
func (n *Node) String() string {
	return base64.StdEncoding.EncodeToString(n.Bytes())
}