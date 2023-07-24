////////////////////////////////////////////////////////////////////////////////
// Copyright © 2022 xx foundation                                             //
//                                                                            //
// Use of this source code is governed by a license that can be found in the  //
// LICENSE file.                                                              //
////////////////////////////////////////////////////////////////////////////////

package region

import (
	"encoding/binary"
	"io"
	"math"

	"github.com/pkg/errors"
	jww "github.com/spf13/jwalterweatherman"

	"gitlab.com/xx_network/primitives/id"
)

func OrderNodeTeam(nodes []*id.ID, countries map[id.ID]string,
	countryToBins map[string]GeoBin, distanceLatency [12][12]int,
	rng io.Reader) ([]*id.ID, int, error) {
	// Make all permutations of nodePermutation
	permutations := Permute(nodes)
	jww.DEBUG.Printf("Looking for most efficient teaming order")
	optimalLatency := math.MaxInt32
	var optimalTeams [][]*id.ID

	// TODO: consider a way to do this more efficiently? As of now,
	//  for larger teams of 10 or greater it takes >2 seconds for round creation
	//  but it runs in the microsecond range with 4 nodePermutation.
	//  Since our use case is smaller teams, we deem this sufficient for now
	for np := range permutations {
		nodePermutation := permutations[np]
		totalLatency := 0
		for i := range nodePermutation {
			thisNode := nodePermutation[i]
			// Get the ordering for the current node
			thisCounty, ok := countries[*thisNode]
			if !ok {
				return nil, 0, errors.Errorf("Unable to locate country for "+
					"node %s: %v", thisNode, countries)
			}
			thisRegion, ok := countryToBins[thisCounty]
			if !ok {
				return nil, 0, errors.Errorf("Unable to locate bin for "+
					"node %s at country %s", thisNode, thisCounty)
			}

			// Get the ordering of the next node, circling back if at the last
			// node
			nextNode := nodePermutation[(i+1)%len(nodePermutation)]
			nextCounty, ok := countries[*nextNode]
			if !ok {
				return nil, 0, errors.Errorf("Unable to locate country for "+
					"node %s: %v", nextCounty, countries)
			}
			nextRegion, ok := countryToBins[nextCounty]
			if !ok {
				return nil, 0, errors.Errorf("Unable to locate bin for "+
					"node %s at country %s", nextNode, nextCounty)
			}
			// Calculate the distance and pull the latency from the table
			totalLatency += distanceLatency[thisRegion][nextRegion]

		}

		// Replace with the best time and order found thus far
		if totalLatency < optimalLatency {
			optimalTeams = make([][]*id.ID, 0)
			optimalTeams = append(optimalTeams, nodePermutation)
			optimalLatency = totalLatency
		} else if totalLatency == optimalLatency {
			optimalTeams = append(optimalTeams, nodePermutation)
		}

	}

	numBytes := make([]byte, 8)
	_, err := rng.Read(numBytes)
	if err != nil {
		return nil, 0, errors.WithMessagef(err, "failed to generate ordering")
	}
	index := binary.BigEndian.Uint64(numBytes) % uint64(len(optimalTeams))

	return optimalTeams[index], optimalLatency, nil
}

// CreateLinkTable creates a latency table that maps different region's
// latencies to all other defined regions. Latency is derived through educated
// guesses right now without any real world data.
//
// TODO: This table needs better real-world accuracy. Once data is collected
//
//	this table can be updated for better accuracy and selection.
func CreateLinkTable() (distanceLatency [12][12]int) {

	// Number of hops on the graph from region.Americas to other regions
	distanceLatency[NorthAmerica][NorthAmerica] = 0
	distanceLatency[NorthAmerica][SouthAndCentralAmerica] = 1
	distanceLatency[NorthAmerica][WesternEurope] = 2
	distanceLatency[NorthAmerica][CentralEurope] = 3
	distanceLatency[NorthAmerica][EasternEurope] = 4
	distanceLatency[NorthAmerica][MiddleEast] = 4
	distanceLatency[NorthAmerica][SouthernAfrica] = 3
	distanceLatency[NorthAmerica][NorthernAfrica] = 3
	distanceLatency[NorthAmerica][Russia] = 4
	distanceLatency[NorthAmerica][EasternAsia] = 2
	distanceLatency[NorthAmerica][WesternAsia] = 3
	distanceLatency[NorthAmerica][Oceania] = 2

	distanceLatency[SouthAndCentralAmerica][NorthAmerica] = 1
	distanceLatency[SouthAndCentralAmerica][SouthAndCentralAmerica] = 0
	distanceLatency[SouthAndCentralAmerica][WesternEurope] = 3
	distanceLatency[SouthAndCentralAmerica][CentralEurope] = 4
	distanceLatency[SouthAndCentralAmerica][EasternEurope] = 5
	distanceLatency[SouthAndCentralAmerica][MiddleEast] = 5
	distanceLatency[SouthAndCentralAmerica][SouthernAfrica] = 4
	distanceLatency[SouthAndCentralAmerica][NorthernAfrica] = 4
	distanceLatency[SouthAndCentralAmerica][Russia] = 5
	distanceLatency[SouthAndCentralAmerica][EasternAsia] = 3
	distanceLatency[SouthAndCentralAmerica][WesternAsia] = 4
	distanceLatency[SouthAndCentralAmerica][Oceania] = 3

	distanceLatency[WesternEurope][NorthAmerica] = 2
	distanceLatency[WesternEurope][SouthAndCentralAmerica] = 3
	distanceLatency[WesternEurope][WesternEurope] = 0
	distanceLatency[WesternEurope][CentralEurope] = 1
	distanceLatency[WesternEurope][EasternEurope] = 2
	distanceLatency[WesternEurope][MiddleEast] = 2
	distanceLatency[WesternEurope][SouthernAfrica] = 1
	distanceLatency[WesternEurope][NorthernAfrica] = 1
	distanceLatency[WesternEurope][Russia] = 3
	distanceLatency[WesternEurope][EasternAsia] = 4
	distanceLatency[WesternEurope][WesternAsia] = 3
	distanceLatency[WesternEurope][Oceania] = 4

	distanceLatency[CentralEurope][NorthAmerica] = 3
	distanceLatency[CentralEurope][SouthAndCentralAmerica] = 4
	distanceLatency[CentralEurope][WesternEurope] = 1
	distanceLatency[CentralEurope][CentralEurope] = 0
	distanceLatency[CentralEurope][EasternEurope] = 1
	distanceLatency[CentralEurope][MiddleEast] = 1
	distanceLatency[CentralEurope][SouthernAfrica] = 1
	distanceLatency[CentralEurope][NorthernAfrica] = 1
	distanceLatency[CentralEurope][Russia] = 2
	distanceLatency[CentralEurope][EasternAsia] = 3
	distanceLatency[CentralEurope][WesternAsia] = 2
	distanceLatency[CentralEurope][Oceania] = 4

	distanceLatency[EasternEurope][NorthAmerica] = 4
	distanceLatency[EasternEurope][SouthAndCentralAmerica] = 5
	distanceLatency[EasternEurope][WesternEurope] = 2
	distanceLatency[EasternEurope][CentralEurope] = 1
	distanceLatency[EasternEurope][EasternEurope] = 0
	distanceLatency[EasternEurope][MiddleEast] = 1
	distanceLatency[EasternEurope][SouthernAfrica] = 2
	distanceLatency[EasternEurope][NorthernAfrica] = 2
	distanceLatency[EasternEurope][Russia] = 1
	distanceLatency[EasternEurope][EasternAsia] = 3
	distanceLatency[EasternEurope][WesternAsia] = 2
	distanceLatency[EasternEurope][Oceania] = 4

	distanceLatency[MiddleEast][NorthAmerica] = 4
	distanceLatency[MiddleEast][SouthAndCentralAmerica] = 5
	distanceLatency[MiddleEast][WesternEurope] = 2
	distanceLatency[MiddleEast][CentralEurope] = 1
	distanceLatency[MiddleEast][EasternEurope] = 1
	distanceLatency[MiddleEast][MiddleEast] = 0
	distanceLatency[MiddleEast][SouthernAfrica] = 2
	distanceLatency[MiddleEast][NorthernAfrica] = 2
	distanceLatency[MiddleEast][Russia] = 2
	distanceLatency[MiddleEast][EasternAsia] = 2
	distanceLatency[MiddleEast][WesternAsia] = 1
	distanceLatency[MiddleEast][Oceania] = 3

	distanceLatency[NorthernAfrica][NorthAmerica] = 3
	distanceLatency[NorthernAfrica][SouthAndCentralAmerica] = 4
	distanceLatency[NorthernAfrica][WesternEurope] = 1
	distanceLatency[NorthernAfrica][CentralEurope] = 1
	distanceLatency[NorthernAfrica][EasternEurope] = 2
	distanceLatency[NorthernAfrica][MiddleEast] = 2
	distanceLatency[NorthernAfrica][SouthernAfrica] = 2
	distanceLatency[NorthernAfrica][NorthernAfrica] = 0
	distanceLatency[NorthernAfrica][Russia] = 3
	distanceLatency[NorthernAfrica][EasternAsia] = 4
	distanceLatency[NorthernAfrica][WesternAsia] = 3
	distanceLatency[NorthernAfrica][Oceania] = 5

	distanceLatency[SouthernAfrica][NorthAmerica] = 3
	distanceLatency[SouthernAfrica][SouthAndCentralAmerica] = 4
	distanceLatency[SouthernAfrica][WesternEurope] = 1
	distanceLatency[SouthernAfrica][CentralEurope] = 1
	distanceLatency[SouthernAfrica][EasternEurope] = 2
	distanceLatency[SouthernAfrica][MiddleEast] = 2
	distanceLatency[SouthernAfrica][SouthernAfrica] = 0
	distanceLatency[SouthernAfrica][NorthernAfrica] = 2
	distanceLatency[SouthernAfrica][Russia] = 3
	distanceLatency[SouthernAfrica][EasternAsia] = 4
	distanceLatency[SouthernAfrica][WesternAsia] = 3
	distanceLatency[SouthernAfrica][Oceania] = 5

	distanceLatency[Russia][NorthAmerica] = 4
	distanceLatency[Russia][SouthAndCentralAmerica] = 5
	distanceLatency[Russia][WesternEurope] = 3
	distanceLatency[Russia][CentralEurope] = 2
	distanceLatency[Russia][EasternEurope] = 1
	distanceLatency[Russia][MiddleEast] = 2
	distanceLatency[Russia][SouthernAfrica] = 3
	distanceLatency[Russia][NorthernAfrica] = 3
	distanceLatency[Russia][Russia] = 0
	distanceLatency[Russia][EasternAsia] = 2
	distanceLatency[Russia][WesternAsia] = 1
	distanceLatency[Russia][Oceania] = 3

	distanceLatency[EasternAsia][NorthAmerica] = 2
	distanceLatency[EasternAsia][SouthAndCentralAmerica] = 3
	distanceLatency[EasternAsia][WesternEurope] = 4
	distanceLatency[EasternAsia][CentralEurope] = 3
	distanceLatency[EasternAsia][EasternEurope] = 3
	distanceLatency[EasternAsia][MiddleEast] = 2
	distanceLatency[EasternAsia][SouthernAfrica] = 4
	distanceLatency[EasternAsia][NorthernAfrica] = 4
	distanceLatency[EasternAsia][Russia] = 2
	distanceLatency[EasternAsia][EasternAsia] = 0
	distanceLatency[EasternAsia][WesternAsia] = 1
	distanceLatency[EasternAsia][Oceania] = 1

	distanceLatency[WesternAsia][NorthAmerica] = 3
	distanceLatency[WesternAsia][SouthAndCentralAmerica] = 4
	distanceLatency[WesternAsia][WesternEurope] = 3
	distanceLatency[WesternAsia][CentralEurope] = 2
	distanceLatency[WesternAsia][EasternEurope] = 2
	distanceLatency[WesternAsia][MiddleEast] = 1
	distanceLatency[WesternAsia][SouthernAfrica] = 3
	distanceLatency[WesternAsia][NorthernAfrica] = 3
	distanceLatency[WesternAsia][Russia] = 1
	distanceLatency[WesternAsia][EasternAsia] = 1
	distanceLatency[WesternAsia][WesternAsia] = 0
	distanceLatency[WesternAsia][Oceania] = 2

	distanceLatency[Oceania][NorthAmerica] = 2
	distanceLatency[Oceania][SouthAndCentralAmerica] = 3
	distanceLatency[Oceania][WesternEurope] = 4
	distanceLatency[Oceania][CentralEurope] = 4
	distanceLatency[Oceania][EasternEurope] = 4
	distanceLatency[Oceania][MiddleEast] = 3
	distanceLatency[Oceania][SouthernAfrica] = 5
	distanceLatency[Oceania][NorthernAfrica] = 5
	distanceLatency[Oceania][Russia] = 3
	distanceLatency[Oceania][EasternAsia] = 1
	distanceLatency[Oceania][WesternAsia] = 2
	distanceLatency[Oceania][Oceania] = 0

	return
}

func CreateSetLatencyTableWeights(distanceLatency [12][12]int) [12][12]int {
	weights := make([]int, len(distanceLatency))
	weight := 1
	for i := 0; i < len(weights); i++ {
		weights[i] = weight
		weight += 2
	}

	for i := 0; i < len(distanceLatency); i++ {
		for j := 0; j < len(distanceLatency); j++ {
			distanceLatency[i][j] = weights[distanceLatency[i][j]]
		}
	}

	return distanceLatency
}
