////////////////////////////////////////////////////////////////////////////////
// Copyright © 2018 Privategrity Corporation                                   /
//                                                                             /
// All rights reserved.                                                        /
////////////////////////////////////////////////////////////////////////////////

package nodeid

import (
	"testing"
)

func TestNodeID(t *testing.T) {
	// try to set the first time
	actual := uint64(489)
	SetNodeID(actual)
	expected := uint64(489)
	// they should match
	if actual != expected {
		t.Errorf("NodeID: actual (%v) differed from expected (%v)", actual,
			expected)
	}

	// try to set the second time
	changeSecondTime := uint64(55)
	SetNodeID(changeSecondTime)
	expected = uint64(489)
	// they shouldn't match
	if changeSecondTime == expected {
		t.Errorf("NodeID: could set twice: %v, %v", actual,
			expected)
	}
}
