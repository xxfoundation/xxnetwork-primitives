////////////////////////////////////////////////////////////////////////////////
// Copyright © 2022 xx foundation                                             //
//                                                                            //
// Use of this source code is governed by a license that can be found in the  //
// LICENSE file.                                                              //
////////////////////////////////////////////////////////////////////////////////

package region

import (
	"math/rand"
	"testing"

	"gitlab.com/xx_network/primitives/id"
)

// Happy path
func TestPermute(t *testing.T) {
	const totalNodes = 3
	nodeList := make([]*id.ID, totalNodes)
	prng := rand.New(rand.NewSource(42))

	// Build node states with unique ordering
	for i := range nodeList {
		// Make a node state and place it in the list
		nodeList[i] = id.NewRandomTestID(prng, id.Node, t)
	}

	// Permute the nodes
	permutations := Permute(nodeList)
	expectedLen := factorial(totalNodes, t)

	// Verify that the amount of permutations is
	// factorial of the original amount of nodes
	if len(permutations) != expectedLen {
		t.Errorf("Permutations did not produce the expected amount of "+
			"permutations (factorial of amount of nodes)!"+
			"\nexpected: %d\nreceived: %d", expectedLen, len(permutations))
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

func factorial(n int, t testing.TB) int {
	factVal := 1
	if n < 0 {
		t.Errorf("Factorial of negative number doesn't exist: %d", n)
	} else {
		for i := 1; i <= n; i++ {
			factVal *= i
		}
	}

	return factVal
}
