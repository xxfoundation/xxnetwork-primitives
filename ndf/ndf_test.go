////////////////////////////////////////////////////////////////////////////////////////////
// Copyright © 2020 xx network SEZC                                                       //
//                                                                                        //
// Use of this source code is governed by a license that can be found in the LICENSE file //
////////////////////////////////////////////////////////////////////////////////////////////

package ndf

import (
	"bytes"
	"encoding/json"
	"gitlab.com/xx_network/primitives/utils"
	"os"
	"reflect"
	"testing"
)

var (
	ExampleNDF = ""
	TestDef    = &NetworkDefinition{}
)

func TestMain(m *testing.M) {
	data, err := utils.ReadFile("ndf.json")
	if err != nil {
		panic(err)
	}
	ExampleNDF = string(data)
	err = json.Unmarshal(data, TestDef)
	if err != nil {
		panic(err)
	}
	os.Exit(m.Run())
}

// Test that DecodeNDF() produces a NetworkDefinition for the supplied NDF.
func TestNetworkDefinition_Marshal_Unmarshal(t *testing.T) {
	def := &NetworkDefinition{}
	err := json.Unmarshal([]byte(ExampleNDF), def)
	if err != nil {
		t.Errorf("Failed to Unmarshal() NDF: %+v", err)
	}

	data, err := json.Marshal(def)
	if err != nil {
		t.Errorf("Failed to Marshal() NDF: %+v", err)
	}
	if !bytes.Equal([]byte(ExampleNDF), data) {
		t.Errorf("Mashaled data does not match expected NDF."+
			"\n\texpected: %s\n\treceived: %s", ExampleNDF, data)
	}
}

// Happy path.
func TestStripNdf(t *testing.T) {
	strippedNDF := TestDef.StripNdf()

	if reflect.DeepEqual(TestDef, strippedNDF) {
		t.Errorf("Stripped NDF matches normal NDF."+
			"\n\texpected: %+v\n\treceived: %+v", TestDef, strippedNDF)
	}

	// Check that the address and cert fields are empty
	for i, node := range strippedNDF.Nodes {
		expectedNode := Node{ID: TestDef.Nodes[i].ID}
		if !reflect.DeepEqual(expectedNode, node) {
			t.Errorf("StripNdf() did not modify node %d correctly."+
				"\n\texpected: %v\n\trecieved: %v", i, expectedNode, node)
		}
	}
}

// Happy path: this finds an id in the hardcoded/global ExampleNDF
func TestGetNodeId(t *testing.T) {
	testID, err := TestDef.Nodes[0].GetNodeId()
	if err != nil {
		t.Errorf("GetNodeId() produced an error: %+v", err)
	}

	if !bytes.Equal(TestDef.Nodes[0].ID, testID.Bytes()) {
		t.Errorf("GetNodeId() produced an unexpected ID."+
			"\n\texpected: %s\n\trecieved: %s", TestDef.Nodes[0].ID, testID)
	}
}

// Happy path.
func TestGetGatewayId(t *testing.T) {
	testID, err := TestDef.Gateways[0].GetGatewayId()
	if err != nil {
		t.Errorf("GetGatewayId() produced an error: %+v", err)
	}

	if !bytes.Equal(TestDef.Gateways[0].ID, testID.Bytes()) {
		t.Errorf("GetGatewayId() produced an unexpected ID."+
			"\n\texpected: %s\n\trecieved: %s", TestDef.Gateways[0].ID, testID)
	}
}

// Happy path.
func TestNewTestNDF(t *testing.T) {
	def := NewTestNDF(t)

	if !reflect.DeepEqual(TestDef, def) {
		t.Errorf("NewTestNDF() did not return the expected NDF."+
			"\n\texpected: %+v\n\trecieved: %+v", TestDef, def)
	}
}
