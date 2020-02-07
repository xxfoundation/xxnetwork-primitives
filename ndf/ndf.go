////////////////////////////////////////////////////////////////////////////////
// Copyright © 2020 Privategrity Corporation                                   /
//                                                                             /
// All rights reserved.                                                        /
////////////////////////////////////////////////////////////////////////////////

package ndf

import (
	"encoding/base64"
	"encoding/json"
	jww "github.com/spf13/jwalterweatherman"
	"strings"
	"time"
)

const NO_NDF = "Permissioning server does not have an ndf to give to client"

// NetworkDefinition structure matches the JSON structure generated in
// Terraform, which allows it to be decoded to Go.
type NetworkDefinition struct {
	Timestamp    time.Time
	Gateways     []Gateway
	Nodes        []Node
	Registration Registration
	Notification Notification
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
	Address        string
	TlsCertificate string `json:"Tls_certificate"`
}

// Registration is the structure for the registration object in the JSON file.
type Registration struct {
	Address        string
	TlsCertificate string `json:"Tls_certificate"`
}

// Notifications is the structure for the registration object in the JSON file.
type Notification struct {
	Address        string
	TlsCertificate string `json:"Tls_certificate"`
}

// UDB is the structure for the UDB object in the JSON file.
type UDB struct {
	ID []byte `json:"Id"`
}

// Group is the structure for a group in the JSON file; it is used for the E2E
// and CMIX objects.
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
	jsonData, signature := separate(ndf)

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

// StripNdf returns a stripped down version of the ndf for clients
func StripNdf(ndfJson []byte) ([]byte, error) {
	definition, _, err := DecodeNDF(string(ndfJson))
	if err != nil {
		return nil, err
	}

	// Strip down nodes slice of addresses and certs
	var strippedNodes []Node
	for _, node := range definition.Nodes {
		newNode := Node{
			ID: node.ID,
		}
		strippedNodes = append(strippedNodes, newNode)
	}

	// Create a new Ndf with the stripped information
	strippedNdf := &NetworkDefinition{
		Timestamp:    definition.Timestamp,
		Gateways:     definition.Gateways,
		Nodes:        strippedNodes,
		Registration: definition.Registration,
		Notification: definition.Notification,
		UDB:          definition.UDB,
		E2E:          definition.E2E,
		CMIX:         definition.CMIX,
	}

	data, err := json.Marshal(strippedNdf)
	if err != nil {
		return nil, err
	}

	return data, nil
}

// separate splits the JSON data from the signature. The NDF string is expected
// to have the JSON data starting on line 1 and its signature on the last line.
// Returns JSON data and signature as separate strings. If the signature is not
// present, it is returned as an empty string.
func separate(ndf string) (string, string) {
	var jsonLineEnd int
	var signature string
	lines := strings.Split(ndf, "\n")

	// Determine which line the JSON ends and which line the signature is on
	for i := len(lines) - 1; i >= 0; i-- {
		line := strings.TrimSpace(lines[i])
		if line != "" {
			if strings.HasSuffix(line, "}") {
				jsonLineEnd = i
				break
			} else {
				signature = line
			}
		}
	}

	return strings.Join(lines[0:jsonLineEnd+1], "\n"), signature
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
		b = append(b, []byte(val.Address)...)
		b = append(b, []byte(val.TlsCertificate)...)
	}

	// Convert Registration to byte slice
	b = append(b, []byte(ndf.Registration.Address)...)
	b = append(b, []byte(ndf.Registration.TlsCertificate)...)

	// Convert UDB to byte slice
	b = append(b, ndf.UDB.ID...)

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
