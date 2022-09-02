////////////////////////////////////////////////////////////////////////////////
// Copyright © 2022 xx foundation                                             //
//                                                                            //
// Use of this source code is governed by a license that can be found in the  //
// LICENSE file.                                                              //
////////////////////////////////////////////////////////////////////////////////

package region

import (
	"fmt"
	"gitlab.com/xx_network/primitives/id"
	"math/rand"
	"testing"
)

// Happy path
func TestPermute(t *testing.T) {

	const totalNodes = 3
	nodeList := make([]*id.ID, totalNodes)

	prng := rand.New(rand.NewSource(42))

	// Build node states with unique ordering
	for i := 0; i < totalNodes; i++ {
		// Make a node state

		newNode, _ := id.NewRandomID(prng, id.Node)

		// Place new node in list
		nodeList[i] = newNode

	}

	// Permute the nodes
	permutations := Permute(nodeList)
	expectedLen := factorial(totalNodes)

	// Verify that the amount of permutations is
	// factorial of the original amount of nodes
	if len(permutations) != expectedLen {
		t.Errorf("Permutations did not produce the expected amount of permutations "+
			"(factorial of amount of nodes)!"+
			"\n\tExpected: %d"+
			"\n\tReceived: %d", expectedLen, len(permutations))
	}

	expectedPermutations := make(map[string]bool)

	// Iterate through all the permutations to ensure uniqueness between orderings
	for _, permutation := range permutations {
		var concatenatedOrdering string
		// Concatenate orderings into a single string
		for _, ourNode := range permutation {
			concatenatedOrdering += ourNode.String()
		}
		// If that ordering has been encountered before, error
		if expectedPermutations[concatenatedOrdering] {
			t.Errorf("Permutation %s has occurred more than once!", concatenatedOrdering)
		}

		// Mark permutation as seen
		expectedPermutations[concatenatedOrdering] = true

	}

}

func factorial(n int) int {
	factVal := 1
	if n < 0 {
		fmt.Println("Factorial of negative number doesn't exist.")
	} else {
		for i := 1; i <= n; i++ {
			factVal *= i
		}

	}
	return factVal
}
