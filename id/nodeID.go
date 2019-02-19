////////////////////////////////////////////////////////////////////////////////
// Copyright © 2018 Privategrity Corporation                                   /
//                                                                             /
// All rights reserved.                                                        /
////////////////////////////////////////////////////////////////////////////////

package id

import (
	"sync"
)

var IsLastNode bool

var nodeId uint64
var setNodeIdOnce sync.Once

func SetNodeID(newNodeID uint64) {
	setNodeIdOnce.Do(func() {
		nodeId = newNodeID
	})
}

func GetNodeID() uint64 {
	return nodeId
}
