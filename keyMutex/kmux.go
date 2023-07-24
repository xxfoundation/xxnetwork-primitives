////////////////////////////////////////////////////////////////////////////////
// Copyright © 2022 xx foundation                                             //
//                                                                            //
// Use of this source code is governed by a license that can be found in the  //
// LICENSE file.                                                              //
////////////////////////////////////////////////////////////////////////////////

// Package keyMutex implements a keyed mutex with a sync map.
// This allows you to Lock with a generic interface key and unlock.
package keyMutex

import "sync"

// KeyMutex is a keyed mutex map.
type KeyMutex struct {
	m *sync.Map
}

// New creates a new KeyMutex
func New() *KeyMutex {
	return &KeyMutex{
		m: &sync.Map{},
	}
}

// Lock returns a locked mutex for the given key. If the key does not
// exist, a new sync.Mutex is created for that key. The sync.Mutex is
// returned in the "Lock"ed state, so the caller is responsible for
// calling Unlock to prevent deadlocks.
func (km *KeyMutex) Lock(key interface{}) *sync.Mutex {
	l, ok := km.m.Load(key)
	if !ok {
		lck := &sync.Mutex{}
		km.m.Store(key, lck)
		lck.Lock()
		return lck
	}
	lck := l.(*sync.Mutex)
	lck.Lock()
	return lck
}

// Delete removes the mutex from the KeyMutex
func (km *KeyMutex) Delete(key interface{}) {
	km.m.Delete(key)
}
