////////////////////////////////////////////////////////////////////////////////////////////
// Copyright © 2020 xx network SEZC                                                       //
//                                                                                        //
// Use of this source code is governed by a license that can be found in the LICENSE file //
////////////////////////////////////////////////////////////////////////////////////////////

// Package ndf contains the structure for our network definition file. This object is used by
// various users, including our cMix nodes and clients of the xx Messenger, among others.
// It also includes functions to unmarshal an NDF from a JSON file, separate the signature
// from the actual NDF content, and serialize the NDF structure into a byte slice

package ndf

import (
	"encoding/base64"
	"encoding/json"
	"github.com/pkg/errors"
	jww "github.com/spf13/jwalterweatherman"
	"gitlab.com/xx_network/primitives/id"
	"strings"
	"time"
)

// This constant string is to be used by our users that request NDFs from permissioning.
// Those that request include cMix nodes, gateways, notification bot and clients.
// Permissioning builds and provides the ndf to those that request it. However, depending on
// the status of the cMix network, it might not have the the ndf ready upon request.
// Permissioning in this case tells the requester that it is not ready with an error message.
// The requester checks if the error message contains this string, and thus knows it needs to ask
// again.
const NO_NDF = "Contacted server does not have an ndf to give"

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
	ID             []byte `json:"Id"`
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
	ID      []byte `json:"Id"`
	Cert    string `json:"Cert"`
	Address string `json:"Address"`
}

// Group is the structure for a group in the JSON file; it is used for the E2E
// and CMIX objects.
type Group struct {
	Prime      string
	SmallPrime string `json:"Small_prime"`
	Generator  string
}

func (g *Group) String() (string, error) {
	data, err := json.Marshal(g)
	if err != nil {
		return "", errors.Errorf("Unable to marshal group: %+v", err)
	}

	return string(data), nil
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

// Returns a stripped down copy of the NDF object to be used by Clients
func (ndf *NetworkDefinition) StripNdf() *NetworkDefinition {
	// Strip down nodes slice of addresses and certs
	var strippedNodes []Node
	for _, node := range ndf.Nodes {
		newNode := Node{
			ID: node.ID,
		}
		strippedNodes = append(strippedNodes, newNode)
	}

	// Create a new Ndf with the stripped information
	return &NetworkDefinition{
		Timestamp:    ndf.Timestamp,
		Gateways:     ndf.Gateways,
		Nodes:        strippedNodes,
		Registration: ndf.Registration,
		Notification: ndf.Notification,
		UDB:          ndf.UDB,
		E2E:          ndf.E2E,
		CMIX:         ndf.CMIX,
	}
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
		b = append(b, val.ID...)
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
	b = append(b, []byte(ndf.UDB.Cert)...)

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

// Marshal returns a json marshal of the ndf
func (ndf *NetworkDefinition) Marshal() ([]byte, error) {
	ndfBytes, err := json.Marshal(ndf)
	if err != nil {
		return nil, err
	}

	return ndfBytes, nil
}

// GetNodeId marshals the node id into the ID type. Returns an error if Marshal
// fails.
func (n *Node) GetNodeId() (*id.ID, error) {
	newID, err := id.Unmarshal(n.ID)
	if err != nil {
		return nil, err
	}

	return newID, nil

}

// GetGatewayId formats the gateway id into the id format specified in the id package of this repo
func (n *Gateway) GetGatewayId() (*id.ID, error) {
	newID, err := id.Unmarshal(n.ID)
	if err != nil {
		return nil, err
	}

	return newID, nil
}
