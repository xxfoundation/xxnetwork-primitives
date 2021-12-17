///////////////////////////////////////////////////////////////////////////////
// Copyright © 2020 xx network SEZC                                          //
//                                                                           //
// Use of this source code is governed by a license that can be found in the //
// LICENSE file                                                              //
///////////////////////////////////////////////////////////////////////////////

// package ExcludedRounds contains a wrapper for the set object which is thread-safe

package excludedRounds

import (
	"github.com/golang-collections/collections/set"
	"gitlab.com/xx_network/primitives/id"
	"sync"
)

// ExcludedRounds struct contains a set of rounds to be excluded from cmix
type ExcludedRounds struct {
	xr *set.Set
	sync.RWMutex
}

func New() *ExcludedRounds {
	return &ExcludedRounds{
		xr:      set.New(nil),
		RWMutex: sync.RWMutex{},
	}
}

func (e *ExcludedRounds) Has(rid id.Round) bool {
	e.RLock()
	defer e.RUnlock()

	return e.xr.Has(rid)
}

func (e *ExcludedRounds) Insert(rid id.Round) {
	e.Lock()
	defer e.Unlock()

	e.xr.Insert(rid)
}

func (e *ExcludedRounds) Remove(rid id.Round) {
	e.Lock()
	defer e.Unlock()

	e.xr.Remove(rid)
}

func (e *ExcludedRounds) Union(other *set.Set) *set.Set {
	e.RLock()
	defer e.RUnlock()

	return e.xr.Union(other)
}
