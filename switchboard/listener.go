////////////////////////////////////////////////////////////////////////////////////////////
// Copyright © 2020 xx network SEZC                                                       //
//                                                                                        //
// Use of this source code is governed by a license that can be found in the LICENSE file //
////////////////////////////////////////////////////////////////////////////////////////////

package switchboard

import (
	jww "github.com/spf13/jwalterweatherman"
	"gitlab.com/elixxir/primitives/id"
)

const AnyType = int32(0)

type Item interface {
	GetSender() *id.ID
	GetMessageType() int32
}

//interface records adhere to
type Listener interface {
	Hear(item Item)
	Name() string
}

// This function type defines callbacks that get passed when the listener is
// listened to. It will always be called in its own goroutine. It may be called
// multiple times simultaneously
type ListenerFunc func(item Item)

// id object returned when a listener is created and is used to delete it from
// the system
type ListenerID struct {
	userID      *id.ID
	messageType int32
	listener    Listener
}

//getter for userID
func (lid ListenerID) GetUserID() *id.ID {
	return lid.userID
}

//getter for message type
func (lid ListenerID) GetMessageType() int32 {
	return lid.messageType
}

//getter for name
func (lid ListenerID) GetName() string {
	return lid.listener.Name()
}

//internal listener implementations
type funcListener struct {
	listener ListenerFunc
	name     string
}

func newFuncListener(listener ListenerFunc, name string) *funcListener {
	return &funcListener{
		listener: listener,
		name:     name,
	}
}

func (fl *funcListener) Hear(item Item) {
	fl.listener(item)
}

func (fl *funcListener) Name() string {
	return fl.name
}

type chanListener struct {
	listener chan Item
	name     string
}

func newChanListener(listener chan Item, name string) *chanListener {
	return &chanListener{
		listener: listener,
		name:     name,
	}
}

func (cl *chanListener) Hear(item Item) {
	select {
	case cl.listener <- item:
	default:
		jww.WARN.Printf("Switchboard failed to speak on channel "+
			"listener %s", cl.name)
	}
}

func (cl *chanListener) Name() string {
	return cl.name
}
