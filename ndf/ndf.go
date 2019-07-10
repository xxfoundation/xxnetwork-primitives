// //////////////////////////////////////////////////////////////////////////////
// Copyright © 2018 Privategrity Corporation                                   /
//                                                                             /
// All rights reserved.                                                        /
// //////////////////////////////////////////////////////////////////////////////

package ndf

import (
	"encoding/json"
	"errors"
	"strings"
	"time"
)

var ErrNDFFile = errors.New(
	"NDF file malformed: expected only two or more lines")

// NetworkDefinition structure matches the JSON structure generated in
// Terraform, which allows it to be decoded to Go. If the JSON structure
// changes, then this structure needs to be updated.
type NetworkDefinition struct {
	Timestamp    time.Time
	Gateways     []Gateway
	Nodes        []Node
	Registration Registration
	UDB          UDB   `json:"Udb"`
	E2E          Group `json:"E2e"`
	CMIX         Group `json:"Cmix"`
}

// Gateway is the structure for the gateways object in the JSON file.
type Gateway struct {
	Address        string
	TlsCertificate string `json:"Tls_certificate"`
}

// Node is the structure for the nodes object in the JSON file.
type Node struct {
	ID             []byte `json:"Id"`
	DsaPublicKey   string `json:"Dsa_public_key"`
	Address        string
	TlsCertificate string `json:"Tls_certificate"`
}

// Registration is the structure for the registration object in the JSON
// file.
type Registration struct {
	DsaPublicKey   string `json:"Dsa_public_key"`
	Address        string
	TlsCertificate string `json:"Tls_certificate"`
}

// UDB is the structure for the udb object in the JSON file.
type UDB struct {
	ID           []byte `json:"Id"`
	DsaPublicKey string `json:"Dsa_public_key"`
}

// UDB is the structure for a group in the JSON file; it is used for the E2E and
// CMIX objects.
type Group struct {
	Prime      string
	SmallPrime string `json:"Small_prime"`
	Generator  string
}

// DecodeNDF decodes the given JSON string into the NetworkDefinition structure.
// The JSON string is expected to have the JSON data on line 1 and its signature
// on line 2. Returns an error if separating the lines fails or if the JSON
// unmarshal fails.
func DecodeNDF(ndf string) (*NetworkDefinition, error) {
	// Get JSON data and check if the separating failed
	jsonString, err := separate(ndf)
	if err != nil {
		return nil, err
	}

	// Unmarshal the JSON string into a structure
	networkDefinition := &NetworkDefinition{}
	err = json.Unmarshal([]byte(jsonString), networkDefinition)
	if err != nil {
		return nil, err
	}

	return networkDefinition, nil
}

// separate splits the JSON data from the signature. The JSON string is expected
// to have the JSON data on line 1 and its signature on line 2. Returns JSON
// data as a string. Returns an error if there are less than two lines in the
// NDF string.
func separate(ndf string) (string, error) {
	lines := strings.Split(ndf, "\n")

	// Check that the NDF string is at least two lines
	if len(lines) < 2 {
		return "", ErrNDFFile
	}

	return lines[0], nil
}
