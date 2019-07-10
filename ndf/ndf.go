////////////////////////////////////////////////////////////////////////////////
// Copyright © 2018 Privategrity Corporation                                   /
//                                                                             /
// All rights reserved.                                                        /
////////////////////////////////////////////////////////////////////////////////

package ndf

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	jww "github.com/spf13/jwalterweatherman"
	"strings"
	"time"
)

// Error for when the NDF file has less than two lines
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

// DecodeNDF decodes the given JSON string into the NetworkDefinition structure
// and decodes the base 64 signature to a byte slice. The NDF string is expected
// to have the JSON data on line 1 and its signature on line 2. Returns an error
// if separating the lines fails or if the JSON unmarshal fails.
func DecodeNDF(ndf string) (*NetworkDefinition, []byte, error) {
	// Get JSON data and check if the separating failed
	jsonData, signature, err := separate(ndf)
	if err != nil {
		return nil, nil, err
	}

	// Decode the signature form base 64 and check for errors
	signatureBytes, err := base64.StdEncoding.DecodeString(signature)
	if err != nil {
		return nil, nil, err
	}

	// Unmarshal the JSON string into a structure
	networkDefinition := &NetworkDefinition{}
	err = json.Unmarshal([]byte(jsonData), networkDefinition)
	if err != nil {
		return nil, nil, err
	}

	return networkDefinition, signatureBytes, nil
}

// separate splits the JSON data from the signature. The NDF string is expected
// to have the JSON data on line 1 and its signature on line 2. Returns JSON
// data and signature as separate strings. Returns an error if there are less
// than two lines in the NDF string.
func separate(ndf string) (string, string, error) {
	lines := strings.Split(ndf, "\n")

	// Check that the NDF string is at least two lines
	if len(lines) < 2 {
		return "", "", ErrNDFFile
	}

	return lines[0], lines[1], nil
}

// Serialize converts the NetworkDefinition into a byte slice.
func (ndf *NetworkDefinition) Serialize() []byte {
	b := make([]byte, 0)

	// Convert timestamp to a byte slice
	timeBytes, err := ndf.Timestamp.MarshalBinary()
	if err != nil {
		jww.FATAL.Panicf("Failed to marshal NetworkDefinition timestamp: %v", err)
	}

	b = append(b, timeBytes...)

	// Convert Gateways slice to byte slice
	for _, val := range ndf.Gateways {
		b = append(b, []byte(val.Address)...)
		b = append(b, []byte(val.TlsCertificate)...)
	}

	// Convert Nodes slice to byte slice
	for _, val := range ndf.Nodes {
		b = append(b, val.ID...)
		b = append(b, []byte(val.DsaPublicKey)...)
		b = append(b, []byte(val.Address)...)
		b = append(b, []byte(val.TlsCertificate)...)
	}

	// Convert Registration to byte slice
	b = append(b, []byte(ndf.Registration.DsaPublicKey)...)
	b = append(b, []byte(ndf.Registration.Address)...)
	b = append(b, []byte(ndf.Registration.TlsCertificate)...)

	// Convert UDB to byte slice
	b = append(b, []byte(ndf.UDB.ID)...)
	b = append(b, []byte(ndf.UDB.DsaPublicKey)...)

	// Convert E2E to byte slice
	b = append(b, []byte(ndf.E2E.Prime)...)
	b = append(b, []byte(ndf.E2E.Generator)...)
	b = append(b, []byte(ndf.E2E.SmallPrime)...)

	// Convert CMIX to byte slice
	b = append(b, []byte(ndf.CMIX.Prime)...)
	b = append(b, []byte(ndf.CMIX.Generator)...)
	b = append(b, []byte(ndf.CMIX.SmallPrime)...)

	return b
}
