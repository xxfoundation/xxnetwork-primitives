////////////////////////////////////////////////////////////////////////////////
// Copyright © 2022 xx foundation                                             //
//                                                                            //
// Use of this source code is governed by a license that can be found in the  //
// LICENSE file.                                                              //
////////////////////////////////////////////////////////////////////////////////

package region

import (
	"fmt"
	"testing"
)

func TestWeights(t *testing.T) {
	to := 7

	weights := make([]int, to+1)

	weights[1] = 1

	for current := 2; current < to; current++ {
		var items []int
		for i := 1; i < current; i++ {
			for j := 0; j < (current+i-1)/i; j++ {
				items = append(items, i)
			}
		}

		permutations := All(items)
		fmt.Printf("Permuted %d got %d results\n", current, len(permutations))
		maxWeight := 0
		for _, permutation := range permutations {
			total := 0
			weight := 0
			for _, element := range permutation {
				total += element
				weight += weights[element]
			}
			if total != current {
				continue
			}
			if weight > maxWeight {
				maxWeight = weight
			}
		}

		weights[current] = maxWeight + 1
		fmt.Printf("Current:%d, weight:%d\n", current, weights[current])
	}

	fmt.Printf("FinalList:%v,\n", weights)

}

// All returns all combinations for a given string array.
// This is essentially a powerset of the given set except that the empty set is
// disregarded.
func All(set []int) (subsets [][]int) {
	length := uint(len(set))

	// Go through all possible combinations of objects
	// from 1 (only first object in subset) to 2^length (all objects in subset)
	for subsetBits := 1; subsetBits < (1 << length); subsetBits++ {
		var subset []int

		for object := uint(0); object < length; object++ {
			// checks if object is contained in subset
			// by checking if bit 'object' is set in subsetBits
			if (subsetBits>>object)&1 == 1 {
				// add object to subset
				subset = append(subset, set[object])
			}
		}
		// add subset to subsets
		subsets = append(subsets, subset)
	}
	return subsets
}
