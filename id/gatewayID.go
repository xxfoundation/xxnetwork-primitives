////////////////////////////////////////////////////////////////////////////////////////////
// Copyright © 2020 xx network SEZC                                                       //
//                                                                                        //
// Use of this source code is governed by a license that can be found in the LICENSE file //
////////////////////////////////////////////////////////////////////////////////////////////

// Package id contains identification structures for gateways, nodes and users.
// The ID structures include structures for the ID of gateways, nodes, and users.
// Each ID is a byte slice with a constant length. These structures each have functions
// allowing IDs to be created, compared, copied, serialised, and be converted to strings.

package id

import (
	"encoding/base64"
)

// Length of gateway IDs in bytes
const GatewayIdLen = 32

// Gateway ID array
type Gateway [GatewayIdLen]byte

// Used as a temporary gateway id untill we come up with a better solution for generating gateway ID's
func NewTmpGateway() *Gateway {
	gateway := new(Gateway)
	copy(gateway[:], "tmp")
	return gateway
}

// NewGateway returns a new Gateway ID from a Node ID if the byte slice has the
// correct length. Otherwise, it returns a Gateway ID that is all zeroes.
func (n *Node) NewGateway() *Gateway {
	gateway := new(Gateway)
	if len(n) == GatewayIdLen {
		copy(gateway[:], n[:])
	}
	return gateway
}

// Bytes returns a copy of a Gateway ID as a byte slice.
func (g *Gateway) Bytes() []byte {
	bytes := make([]byte, GatewayIdLen)
	copy(bytes, g[:])

	return bytes
}

// Cmp determines whether two Gateway IDs are the same.
func (g *Gateway) Cmp(y *Gateway) bool {
	return *g == *y
}

// DeepCopy creates a completely new copy of the Gateway ID.
func (g *Gateway) DeepCopy() *Gateway {
	if g == nil {
		return nil
	}

	var newGateway Gateway
	copy(newGateway[:], (*g)[:])

	return &newGateway
}

// String implements Stringer and allows Gateway IDs to be used as connection
// IDs. Currently, it just base64 encodes the Gateway ID and appends "-Gateway".
func (g *Gateway) String() string {
	return base64.StdEncoding.EncodeToString(g.Bytes()) + "-Gateway"
}
