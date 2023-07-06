////////////////////////////////////////////////////////////////////////////////
// Copyright © 2022 xx foundation                                             //
//                                                                            //
// Use of this source code is governed by a license that can be found in the  //
// LICENSE file.                                                              //
////////////////////////////////////////////////////////////////////////////////

package region

import "gitlab.com/xx_network/primitives/id"

// permute.go contains the implementation of Heap's algorithm, used to generate
// all possible permutations of n objects

// Permute is based off of Heap's algorithm found here:
// https://en.wikipedia.org/wiki/Heap%27s_algorithm.
//
// It runs n! time, but in place in terms of space. As of writing, we use this
// for permuting all orders of a team, of which team size is small, justifying
// the high complexity.
func Permute(items []*id.ID) [][]*id.ID {
	var helper func([]*id.ID, int)
	var output [][]*id.ID

	// Place inline to make appending output easier
	helper = func(items []*id.ID, numItems int) {
		if numItems == 1 {
			// Create a copy and append the copy to the output
			ourCopy := make([]*id.ID, len(items))
			copy(ourCopy, items)
			output = append(output, ourCopy)
		} else {
			for i := 0; i < numItems; i++ {
				helper(items, numItems-1)
				// Swap choice dependent on parity of k (even or odd)
				if numItems%2 == 1 {
					// Swap the values
					items[i], items[numItems-1] = items[numItems-1], items[i]

				} else {
					// Swap the values
					items[0], items[numItems-1] = items[numItems-1], items[0]

				}
			}
		}
	}

	// Initialize recursive function
	helper(items, len(items))
	return output
}
