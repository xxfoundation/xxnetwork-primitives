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

// DomainName is the registered domain to be used for calculation of
// unique Gateway DNS addresses by authorizer, gateway, and client.
const DomainName = "elixxirnode.io"

// GetGatewayDns returns the DNS name for the given marshalled GwId.
func GetGatewayDns(gwId []byte) string {
	return fmt.Sprintf("%s.%s",
		base64.URLEncoding.EncodeToString(gwId)[:20], DomainName)
}
