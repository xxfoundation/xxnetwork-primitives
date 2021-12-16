///////////////////////////////////////////////////////////////////////////////
// Copyright © 2020 xx network SEZC                                          //
//                                                                           //
// Use of this source code is governed by a license that can be found in the //
// LICENSE file                                                              //
///////////////////////////////////////////////////////////////////////////////

// package excludedRounds contains a wrapper for the set object which is thread-safe

package excludedRounds

import (
	"github.com/golang-collections/collections/set"
	"sync"
)

// ExcludedRounds struct contains a set of rounds to be excluded from cmix
type ExcludedRounds struct {
	xr *set.Set
	sync.RWMutex
}

func (e *ExcludedRounds) Has(element interface{}) bool {
	e.RLock()
	defer e.RUnlock()

	return e.xr.Has(element)
}

func (e *ExcludedRounds) Insert(element interface{}) {
	e.Lock()
	defer e.Unlock()

	e.xr.Insert(element)
}
