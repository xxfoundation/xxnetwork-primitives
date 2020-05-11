////////////////////////////////////////////////////////////////////////////////////////////
// Copyright © 2020 xx network SEZC                                                       //
//                                                                                        //
// Use of this source code is governed by a license that can be found in the LICENSE file //
////////////////////////////////////////////////////////////////////////////////////////////

package switchboard

import (
	"gitlab.com/elixxir/primitives/id"
	"sync"
	"testing"
	"time"
)

// Demonstrates that all the messages can be heard when multiple threads are
// producing items
func TestListeningQueue_Hear(t *testing.T) {
	numItems := 2000
	numThreads := 8
	var wg sync.WaitGroup
	wg.Add(numThreads * numItems)

	s := NewSwitchboard()
	_, queue := s.ListenChannel(0, id.ZeroID, 12)

	var items []Item

	user := id.NewUserFromUints(&[4]uint64{0, 0, 0, 3})
	// Hopefully this would be enough to cause a race condition
	for j := 0; j < numThreads; j++ {
		go func() {
			for i := 0; i < numItems; i++ {
				s.Speak(&Message{
					Contents:    []byte{},
					Sender:      user,
					MessageType: 5,
				})
				wg.Done()
				time.Sleep(time.Millisecond)
			}
		}()
	}
	// Listen to the heard messages
	// If there aren't enough items, this will block forever instead of failing
	// the test
	for len(items) < numThreads*numItems {
		items = append(items, <-queue)
	}
	// Check that all items are represented
	wg.Wait()
	time.Sleep(50 * time.Millisecond)
	if len(items) != numThreads*numItems {
		t.Error("Didn't get the expected number of items on the channel")
	}
	// Make sure there isn't anything else available on the channel: there
	// should be exactly the right number of items available
	select {
	case <-queue:
		t.Error("There was another item on the channel that shouldn't have" +
			" been there")
	default:
	}
}
