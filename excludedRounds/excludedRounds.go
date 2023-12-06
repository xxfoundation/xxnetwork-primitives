////////////////////////////////////////////////////////////////////////////////
// Copyright © 2024 xx foundation                                             //
//                                                                            //
// Use of this source code is governed by a license that can be found in the  //
// LICENSE file.                                                              //
////////////////////////////////////////////////////////////////////////////////

package excludedRounds

import (
	"gitlab.com/xx_network/primitives/id"
)

// ExcludedRounds is a list of rounds that are excluded from sending on.
type ExcludedRounds interface {
	// Has indicates if the round is in the list.
	Has(rid id.Round) bool

	// Insert adds the round to the list. Returns true if the round was added.
	Insert(rid id.Round) bool

	// Remove deletes the round from the list.
	Remove(rid id.Round)

	// Len returns the number of rounds in the list.
	Len() int
}
