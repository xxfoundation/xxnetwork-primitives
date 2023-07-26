////////////////////////////////////////////////////////////////////////////////
// Copyright © 2022 xx foundation                                             //
//                                                                            //
// Use of this source code is governed by a license that can be found in the  //
// LICENSE file.                                                              //
////////////////////////////////////////////////////////////////////////////////

// package Set contains a wrapper for the set object which is thread-safe

package excludedRounds

import (
	"sync"

	"github.com/golang-collections/collections/set"

	"gitlab.com/xx_network/primitives/id"
)

// Set struct contains a set of rounds to be excluded from cmix
type Set struct {
	xr *set.Set
	sync.RWMutex
}

func NewSet() *Set {
	return &Set{xr: set.New(nil)}
}

func (s *Set) Has(rid id.Round) bool {
	s.RLock()
	defer s.RUnlock()

	return s.xr.Has(rid)
}

func (s *Set) Insert(rid id.Round) bool {
	s.Lock()
	defer s.Unlock()

	if s.xr.Has(rid) {
		return false
	}

	s.xr.Insert(rid)
	return true
}

func (s *Set) Remove(rid id.Round) {
	s.Lock()
	defer s.Unlock()

	s.xr.Remove(rid)
}

func (s *Set) Len() int {
	s.RLock()
	defer s.RUnlock()

	return s.xr.Len()
}
