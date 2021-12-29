///////////////////////////////////////////////////////////////////////////////
// Copyright © 2020 xx network SEZC                                          //
//                                                                           //
// Use of this source code is governed by a license that can be found in the //
// LICENSE file                                                              //
///////////////////////////////////////////////////////////////////////////////

// package Set contains a wrapper for the set object which is thread-safe

package excludedRounds

import (
	"github.com/golang-collections/collections/set"
	"gitlab.com/xx_network/primitives/id"
	"sync"
)

// Set struct contains a set of rounds to be excluded from cmix
type Set struct {
	xr *set.Set
	sync.RWMutex
}

func New() *Set {
	return &Set{
		xr:      set.New(nil),
		RWMutex: sync.RWMutex{},
	}
}

func (e *Set) Has(rid id.Round) bool {
	e.RLock()
	defer e.RUnlock()

	return e.xr.Has(rid)
}

func (e *Set) Insert(rid id.Round) {
	e.Lock()
	defer e.Unlock()

	e.xr.Insert(rid)
}

func (e *Set) Remove(rid id.Round) {
	e.Lock()
	defer e.Unlock()

	e.xr.Remove(rid)
}

func (e *Set) Union(other *set.Set) *set.Set {
	e.RLock()
	defer e.RUnlock()

	return e.xr.Union(other)
}

func (e *Set) Len() int {
	e.RLock()
	defer e.RUnlock()

	return e.xr.Len()
}
