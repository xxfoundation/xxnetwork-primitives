////////////////////////////////////////////////////////////////////////////////
// Copyright © 2022 xx foundation                                             //
//                                                                            //
// Use of this source code is governed by a license that can be found in the  //
// LICENSE file.                                                              //
////////////////////////////////////////////////////////////////////////////////

package keyMutex

import (
	"testing"
)

// TestKeyMutexSmoke adds 2 elements to the array, then locks
// and unlocks them a few times.
func TestKeyMutexSmoke(t *testing.T) {
	km := New()
	for i := 0; i < 10; i++ {
		a := km.Lock("lockA")
		b := km.Lock("lockB")
		a.Unlock()
		b.Unlock()
	}

	cnt := 0
	km.m.Range(func(k, v interface{}) bool {
		cnt += 1
		return true
	})
	if cnt != 2 {
		t.Errorf("invalid count, expected 2, got %d", cnt)
	}

	km.Delete("lockA")
	cnt = 0
	km.m.Range(func(k, v interface{}) bool {
		cnt += 1
		return true
	})
	if cnt != 1 {
		t.Errorf("invalid count, expected 1, got %d", cnt)
	}

	km.Delete("lockB")
	cnt = 0
	km.m.Range(func(k, v interface{}) bool {
		cnt += 1
		return true
	})
	if cnt != 0 {
		t.Errorf("invalid count, expected 0, got %d", cnt)
	}
}
