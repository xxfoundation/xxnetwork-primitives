////////////////////////////////////////////////////////////////////////////////
// Copyright © 2022 xx foundation                                             //
//                                                                            //
// Use of this source code is governed by a license that can be found in the  //
// LICENSE file.                                                              //
////////////////////////////////////////////////////////////////////////////////

package authorizer

import (
	"encoding/hex"
	"fmt"
)

const (
	// DomainName is the registered domain to be used for calculation of
	// unique Gateway DNS addresses by authorizer, gateway, and client.
	DomainName = "xxnode.io"

	// Maximum length of DNS name. Determined by third party service.
	maxLength = 64

	// Maximum number of characters of gateway ID to use. Subtract length of
	// domain plus the additional period from maxLength.
	maxGwIdLength = maxLength - len(DomainName) - 1
)

// GetGatewayDns returns the DNS name for the given marshalled gateway ID.
// Strips the domain name, if it exists, from the encoded ID.
func GetGatewayDns(gwID []byte) string {
	encoded := hex.EncodeToString(gwID)
	if len(encoded) > maxGwIdLength {
		encoded = encoded[:maxGwIdLength]
	}
	return fmt.Sprintf("%s.%s", encoded, DomainName)
}
