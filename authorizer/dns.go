////////////////////////////////////////////////////////////////////////////////
// Copyright © 2022 xx foundation                                             //
//                                                                            //
// Use of this source code is governed by a license that can be found in the  //
// LICENSE file.                                                              //
////////////////////////////////////////////////////////////////////////////////

package authorizer

import (
	"encoding/base64"
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

// GetGatewayDns returns the DNS name for the given marshalled GwId.
// Truncates
func GetGatewayDns(gwId []byte) string {
	return fmt.Sprintf("%s.%s",
		base64.URLEncoding.EncodeToString(gwId)[:maxGwIdLength], DomainName)
}
