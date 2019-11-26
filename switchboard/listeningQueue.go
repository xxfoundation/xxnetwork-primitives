////////////////////////////////////////////////////////////////////////////////
// Copyright © 2018 Privategrity Corporation                                   /
//                                                                             /
// All rights reserved.                                                        /
////////////////////////////////////////////////////////////////////////////////

package switchboard

import (
	"gitlab.com/elixxir/primitives/id"
)

type ListeningQueue chan Item

// Set up a listening queue and add it to the switchboard
func (s *Switchboard) ListenChannel(
	messageType int32, sender *id.User, channelBufferSize int) (id string,
	messageQueue ListeningQueue) {
	messageQueue = make(ListeningQueue, channelBufferSize)
	id = s.Register(sender, messageType, messageQueue)
	return id, messageQueue
}

// TODO What happens if you use pointer receiver? Should test whether it still works correctly.
// Multiple threads can write to this buffer simultaneously through the
// switchboard using this method, although because the writes are to adjacent
// elements, performance is likely to be suboptimal
func (l ListeningQueue) Hear(item Item, isHeardElsewhere bool, i ...interface{}) {
	l <- item
}
