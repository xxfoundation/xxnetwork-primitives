////////////////////////////////////////////////////////////////////////////////////////////
// Copyright © 2020 xx network SEZC                                                       //
//                                                                                        //
// Use of this source code is governed by a license that can be found in the LICENSE file //
////////////////////////////////////////////////////////////////////////////////////////////

// Package ndf contains the structure for the network definition file. It is
// generated by permissioning and propagates to nodes, gateways, and clients in
// the network.

package ndf

import (
	"encoding/json"
	"github.com/pkg/errors"
	jww "github.com/spf13/jwalterweatherman"
	"gitlab.com/xx_network/primitives/id"
	"time"
)

// NO_NDF is a string that the permissioning server responds with when a member
// of the network requests an NDF from it but the NDF is not yet available.
const NO_NDF = "Contacted server does not have an ndf to give"

// NetworkDefinition structure hold connection and network information. It
// matches the JSON structure generated in Terraform.
type NetworkDefinition struct {
	Timestamp     time.Time
	Gateways      []Gateway
	Nodes         []Node
	Registration  Registration
	Notification  Notification
	UDB           UDB   `json:"Udb"`
	E2E           Group `json:"E2e"`
	CMIX          Group `json:"Cmix"`
	AddressSpace  []AddressSpace
	ClientVersion string
}

// Gateway contains the connection and identity information of a gateway on the
// network.
type Gateway struct {
	ID             []byte `json:"Id"`
	Address        string
	TlsCertificate string `json:"Tls_certificate"`
}

// Node contains the connection and identity information of a node on the
// network.
type Node struct {
	ID             []byte `json:"Id"`
	Address        string
	TlsCertificate string `json:"Tls_certificate"`
}

// Registration contains the connection information for the permissioning
// server.
type Registration struct {
	Address                   string
	ClientRegistrationAddress string
	TlsCertificate            string `json:"Tls_certificate"`
	EllipticPubKey            string
}

// Notification contains the connection information for the notification bot.
type Notification struct {
	Address        string
	TlsCertificate string `json:"Tls_certificate"`
}

// UDB contains the ID, public key in PEM format, address, and DH public key for
// user discovery.
type UDB struct {
	ID       []byte `json:"Id"`
	Cert     string `json:"Cert"`
	Address  string `json:"Address"`
	DhPubKey []byte `json:"DhPubKey"`
}

// Group contains the information used to reconstruct a cyclic group.
type Group struct {
	Prime      string
	SmallPrime string `json:"Small_prime"`
	Generator  string
}

type AddressSpace struct {
	Size      uint8
	Timestamp time.Time
}

func (g *Group) String() (string, error) {
	data, err := json.Marshal(g)
	if err != nil {
		return "", errors.Errorf("Unable to marshal group: %+v", err)
	}

	return string(data), nil
}

// Marshal returns the JSON encoding of the NDF.
func (ndf *NetworkDefinition) Marshal() ([]byte, error) {
	return json.Marshal(ndf)
}

// Unmarshal parses the JSON encoded data and returns the resulting
// NetworkDefinition.
func Unmarshal(data []byte) (*NetworkDefinition, error) {
	ndf := &NetworkDefinition{}
	err := json.Unmarshal(data, ndf)
	return ndf, err
}

// DeepCopy returns a deep copy of the NDF. Note that this function is
// intentionally verbose so that tests will fail if this function is not updated
// when the NetworkDefinition is modified.
func (ndf *NetworkDefinition) DeepCopy() *NetworkDefinition {
	newNDF := &NetworkDefinition{
		Gateways:     make([]Gateway, len(ndf.Gateways)),
		Nodes:        make([]Node, len(ndf.Nodes)),
		AddressSpace: make([]AddressSpace, len(ndf.AddressSpace)),
	}

	// Copy timestamp
	newNDF.Timestamp = ndf.Timestamp

	// Copy Gateways
	copy(newNDF.Gateways, ndf.Gateways)

	// Copy Nodes
	copy(newNDF.Nodes, ndf.Nodes)

	// Copy Registration
	newNDF.Registration = Registration{
		Address:        ndf.Registration.Address,
		TlsCertificate: ndf.Registration.TlsCertificate,
		EllipticPubKey: ndf.Registration.EllipticPubKey,
	}

	// Copy Notification
	newNDF.Notification = Notification{
		Address:        ndf.Notification.Address,
		TlsCertificate: ndf.Notification.TlsCertificate,
	}

	// Copy UD
	newNDF.UDB = UDB{
		ID:       make([]byte, len(ndf.UDB.ID)),
		Cert:     ndf.UDB.Cert,
		Address:  ndf.UDB.Address,
		DhPubKey: make([]byte, len(ndf.UDB.DhPubKey)),
	}
	copy(newNDF.UDB.ID, ndf.UDB.ID)
	copy(newNDF.UDB.DhPubKey, ndf.UDB.DhPubKey)

	// Copy E2E group
	newNDF.E2E = Group{
		Prime:      ndf.E2E.Prime,
		SmallPrime: ndf.E2E.SmallPrime,
		Generator:  ndf.E2E.Generator,
	}

	// Copy CMIX group
	newNDF.CMIX = Group{
		Prime:      ndf.CMIX.Prime,
		SmallPrime: ndf.CMIX.SmallPrime,
		Generator:  ndf.CMIX.Generator,
	}

	// Copy AddressSpace
	copy(newNDF.AddressSpace, ndf.AddressSpace)

	// Copy ClientVersion
	newNDF.ClientVersion = ndf.ClientVersion

	return newNDF
}

// StripNdf returns a stripped down copy of the NetworkDefinition to be used by
// Clients.
func (ndf *NetworkDefinition) StripNdf() *NetworkDefinition {
	// Remove address and TLS cert for every node.
	var strippedNodes []Node
	for _, node := range ndf.Nodes {
		strippedNodes = append(strippedNodes, Node{ID: node.ID})
	}

	// Create a new NetworkDefinition with the stripped information
	return &NetworkDefinition{
		Timestamp:    ndf.Timestamp,
		Gateways:     ndf.Gateways,
		Nodes:        strippedNodes,
		Registration: ndf.Registration,
		Notification: ndf.Notification,
		UDB:          ndf.UDB,
		E2E:          ndf.E2E,
		CMIX:         ndf.CMIX,
		AddressSpace: ndf.AddressSpace,
	}
}

// Serialize serializes the NetworkDefinition into a byte slice.
func (ndf *NetworkDefinition) Serialize() []byte {
	var b []byte

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
	b = append(b, []byte(ndf.Registration.EllipticPubKey)...)

	// Convert notification bot to byte slice
	b = append(b, []byte(ndf.Notification.Address)...)
	b = append(b, []byte(ndf.Notification.TlsCertificate)...)

	// Convert UDB to byte slice
	b = append(b, ndf.UDB.ID...)
	b = append(b, []byte(ndf.UDB.Cert)...)
	b = append(b, ndf.UDB.Address...)
	b = append(b, ndf.UDB.DhPubKey...)

	// Convert E2E to byte slice
	b = append(b, []byte(ndf.E2E.Prime)...)
	b = append(b, []byte(ndf.E2E.Generator)...)
	b = append(b, []byte(ndf.E2E.SmallPrime)...)

	// Convert CMIX to byte slice
	b = append(b, []byte(ndf.CMIX.Prime)...)
	b = append(b, []byte(ndf.CMIX.Generator)...)
	b = append(b, []byte(ndf.CMIX.SmallPrime)...)

	// Convert AddressSpace to byte slice
	for _, val := range ndf.AddressSpace {
		b = append(b, val.Size)

		timeBytes, err := val.Timestamp.MarshalBinary()
		if err != nil {
			jww.FATAL.Panicf("Failed to marshal NetworkDefinition "+
				"AddressSpace timestamp: %v", err)
		}

		b = append(b, timeBytes...)
	}

	return b
}

// GetNodeId unmarshalls the Node's ID bytes into an id.ID and returns it.
func (n *Node) GetNodeId() (*id.ID, error) {
	return id.Unmarshal(n.ID)
}

// GetGatewayId unmarshalls the Gateway's ID bytes into an id.ID and returns it.
func (g *Gateway) GetGatewayId() (*id.ID, error) {
	return id.Unmarshal(g.ID)
}
