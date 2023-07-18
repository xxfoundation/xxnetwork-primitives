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
	"time"

	"gitlab.com/xx_network/primitives/id"
)

func TestCreateLatencyTable(t *testing.T) {
	o := CreateLinkTable()

	for i := 0; i < len(o); i++ {
		for j := 0; j < len(o); j++ {
			if o[i][j] != o[j][i] {
				t.Errorf("Orders of %s and %s did not have the same distance "+
					"(%s - %d vs %s -  %d)", GeoBin(i), GeoBin(j),
					GeoBin(i), o[i][j], GeoBin(j), o[j][i])
			}
		}
	}
}

// Test that a team of 8 nodes, each in a different region
// is assembled into a round with an efficient order
func TestCreateRound_EfficientTeam_AllRegions(t *testing.T) {
	const teamSize = 9

	// Build the nodes
	nodeList := make([]*id.ID, teamSize)

	rng := rand.New(rand.NewSource(42))

	// Craft regions for nodes
	countries := []string{
		"CA", "BZ", "GB", "DE", "UA", "IL", "EG", "ZA", "RU", "BD", "IN", "AU"}
	countryMap := make(map[id.ID]string)

	for i := uint64(0); i < uint64(len(nodeList)); i++ {
		nid := id.NewIdFromUInt(i, id.Node, t)
		nodeList[i] = nid
		countryMap[*nid] = countries[i]
	}

	start := time.Now()

	latencyTable := CreateSetLatencyTableWeights(CreateLinkTable())

	bestOrdering, weight, err :=
		OrderNodeTeam(nodeList, countryMap, GetCountryBins(), latencyTable, rng)

	duration := time.Now().Sub(start)
	t.Logf("CreateRound took: %v\n", duration)

	if err != nil {
		t.Fatalf("Failed to get best ordering: %+v", err)
	}

	expectedDuration := 500 * time.Millisecond

	if duration > expectedDuration {
		t.Logf("Warning, creating round for a team of %d took longer than expected."+
			"\nexpected: ~%s\nreceived: %s", teamSize, expectedDuration, duration)
	}

	var regionOrder []GeoBin
	var regionOrderStr []string
	for _, n := range bestOrdering {
		order, _ := GetCountryBin(countryMap[*n])
		regionOrder = append(regionOrder, order)
		regionOrderStr = append(regionOrderStr, order.String())
	}

	t.Logf("Team order outputted by CreateRound with weight %d: %v",
		weight, regionOrderStr)

	// Go through the regions, checking for any long jumps
	validRegionTransitions := newTransitions()
	longTransitions := uint32(0)
	for i, thisRegion := range regionOrder {
		// Get the next region to  see if it's a long distant jump
		nextRegion := regionOrder[(i+1)%len(regionOrder)]
		if !validRegionTransitions.isValidTransition(thisRegion, nextRegion) {
			longTransitions++
		}

	}

	t.Logf("Amount of long distant jumps: %v", longTransitions)

	// Check that the long jumps does not exceed over half the jumps
	if longTransitions > teamSize/2+1 {
		t.Errorf("Number of long distant transitions beyond acceptable amount!"+
			"\nAcceptable long distance transitions: %v"+
			"\nReceived long distance transitions: %v", teamSize/2+1, longTransitions)
	}

}

// Test that a team of 8 nodes, each in a different region
// is assembled into a round with an efficient order
func TestCreateRound_EfficientTeam_CloseAndFar(t *testing.T) {
	const teamSize = 5

	// Build the nodes
	nodeList := make([]*id.ID, teamSize)

	rng := rand.New(rand.NewSource(42))

	// Craft regions for nodes
	countries := []string{"AX", "RU", "KZ", "CN", "AU"}
	countryMap := make(map[id.ID]string)

	for i := uint64(0); i < uint64(len(nodeList)); i++ {
		nid := id.NewIdFromUInt(i, id.Node, t)
		nodeList[i] = nid
		countryMap[*nid] = countries[i]
	}

	start := time.Now()

	latencyTable := CreateSetLatencyTableWeights(CreateLinkTable())

	_, _, err :=
		OrderNodeTeam(nodeList, countryMap, GetCountryBins(), latencyTable, rng)

	duration := time.Now().Sub(start)
	t.Logf("CreateRound took: %v\n", duration)

	if err != nil {
		t.Fatalf("Failed to get best ordering: %+v", err)
	}

	expectedDuration := 60 * time.Millisecond

	if duration > expectedDuration {
		t.Errorf("Warning, creating round for a team of %d took longer than "+
			"expected.\nexpected: ~%s\nreceived: %s",
			teamSize, expectedDuration, duration)
	}

	// var regionOrder []GeoBin
	// var regionOrderStr []string
	// for _, n := range bestOrdering {
	//	order, _ := GetCountryBin(countryMap[*n])
	//	regionOrder = append(regionOrder, order)
	//	regionOrderStr = append(regionOrderStr, order.String())
	// }
	//
	//
	// for i := 0; i < len(bestOrdering)-1; i++ {
	//	x :=  latencyTable[(regionOrder[i])][regionOrder[i+1]]
	//	if x > 2 {
	//		t.Fatalf("blah")
	//	}
	// }
}

/*
// Test that a team of 8 nodes from random regions,
// is assembled into a round with an efficient order
func TestCreateRound_EfficientTeam_RandomRegions(t *testing.T) {
	testPool := NewWaitingPool()

	// Build scheduling params
	testParams := Params{
		TeamSize:            8,
		BatchSize:           32,
		Threshold:           2,
		NodeCleanUpInterval: 3,
	}

	// Build network state
	privKey, _ := rsa.GenerateKey(rand.Reader, 2048)

	testState, err := storage.NewState(privKey, 8, "", region.GetCountryBins())
	if err != nil {
		t.Errorf("Failed to create test state: %v", err)
		t.FailNow()
	}

	// Build the nodes
	nodeList := make([]*id.ID, testParams.TeamSize*2)
	nodeStateList := make([]*node.State, testParams.TeamSize*2)

	// Craft regions for nodes
	regions := []string{"CR", "GB", "SK",
		"HR", "IQ", "BF", "RU", "CX"}

	// Populate the pool with 2x the team size
	for i := uint64(0); i < uint64(len(nodeList)); i++ {
		// Randomize the regions of the nodes
		index := mathRand.Intn(8)

		// Generate a test ID
		nid := id.NewIdFromUInt(i, id.Node, t)
		nodeList[i] = nid

		// Add the node to that node map
		// Place the node in a random region
		err := testState.GetNodeMap().AddNode(nodeList[i], regions[index], "", "", 0)
		if err != nil {
			t.Errorf("Couldn't add node: %v", err)
			t.FailNow()
		}

		// Add the node to the pool
		nodeState := testState.GetNodeMap().GetNode(nid)
		nodeStateList[i] = nodeState
		testPool.Add(nodeState)
	}

	roundID, err := testState.IncrementRoundID()
	if err != nil {
		t.Errorf("IncrementRoundID failed: %+v", err)
	}

	// Create the protoRound
	start := time.Now()
	testProtoRound, err :=
        createSecureRound(testParams, testPool, roundID, testState)
	if err != nil {
		t.Errorf("Error in happy path: %v", err)
	}

	duration := time.Now().Sub(start)
	expectedDuration := int64(45)

	// Check that it did not take an excessive amount of time
	// to create the round
	if duration.Milliseconds() > expectedDuration {
		t.Errorf("Warning, creating round for a team of 8 took longer than expected."+
			"\nexpected: ~%v ms\nreceived: %v ms", expectedDuration, duration)
	}

	// Parse the order of the regions
	// one for testing and one for logging
	var regionOrder []region.GeoBin
	var regionOrderStr []string
	for _, n := range testProtoRound.NodeStateList {
		order, _ := region.GetCountryBin(n.GetOrdering())
		regionOrder = append(regionOrder, order)
		regionOrderStr = append(regionOrderStr, order.String())
	}

	// Output the teaming order to the log in human-readable format
	t.Log("Team order outputted by CreateRound: ", regionOrderStr)

	// Measure the amount of longer than necessary jumps
	validRegionTransitions := newTransitions()
	longTransitions := uint32(0)
	for i, thisRegion := range regionOrder {
		// Get the next region to  see if it's a long distant jump
		nextRegion := regionOrder[(i+1)%len(regionOrder)]
		if !validRegionTransitions.isValidTransition(thisRegion, nextRegion) {
			longTransitions++
		}

	}

	t.Logf("Amount of long distant jumps: %v", longTransitions)

	// Check that the long distant jumps do not exceed half the jumps
	if longTransitions > testParams.TeamSize/2+1 {
		t.Errorf("Number of long distant transitions beyond acceptable amount!"+
			"\nAcceptable long distance transitions: %v"+
			"\nReceived long distance transitions: %v",
            testParams.TeamSize/2+1, longTransitions)
	}

}*/

// Based on the control state logic used for rounds. Based on the map
// discerned from internet cable maps
type regionTransition [12]regionTransitionValidation

// Transitional information used for each region
type regionTransitionValidation struct {
	from [12]bool
}

// Create the valid jumps for each region
func newRegionTransitionValidation(from ...GeoBin) regionTransitionValidation {
	tv := regionTransitionValidation{}

	for _, f := range from {
		tv.from[f] = true
	}

	return tv
}

// Valid transitions are defined as region jumps that are not long distant.
// Long distant is defined by internet cable maps. It was defined in an
// undirected graph of which are good internet connections.
func newTransitions() regionTransition {
	t := regionTransition{}

	latencyTable := CreateLinkTable()

	for i := 0; i < len(latencyTable); i++ {
		acceptable := make([]GeoBin, 0)
		for j := 0; j < len(latencyTable[i]); j++ {
			if latencyTable[i][j] <= 3 {
				acceptable = append(acceptable, GeoBin(j))
			}
		}

		t[i] = newRegionTransitionValidation(acceptable...)
	}

	return t
}

// IsValidTransition checks the transitionValidation to see if
//
//	the attempted transition is valid
func (r regionTransition) isValidTransition(from, to GeoBin) bool {
	return r[to].from[from]
}
