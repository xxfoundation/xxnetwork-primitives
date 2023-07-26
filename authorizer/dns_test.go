////////////////////////////////////////////////////////////////////////////////
// Copyright © 2022 xx foundation                                             //
//                                                                            //
// Use of this source code is governed by a license that can be found in the  //
// LICENSE file.                                                              //
////////////////////////////////////////////////////////////////////////////////

package authorizer

import (
	"math/rand"
	"testing"
)

// Consistency test of GetGatewayDns.
func TestGetGatewayDns(t *testing.T) {
	prng := rand.New(rand.NewSource(389257))

	expectedDnsNames := []string{
		"883968c2360cb4.xxnode.io",
		"72a1a31943d4282c700b614ab911769b316e064bf001081d3484aa.xxnode.io",
		"3f346ab53235279738c9355b97adde74f8703b7f8a55f100617ad0.xxnode.io",
		"df75c745cd8f0e3226cb8b48368a7eefac708b63da8b793e747d6e.xxnode.io",
		"de5e5e0ddfaa873b04d1a6e38d8b000edb.xxnode.io",
		"ebfb75cd082cacc0a24106ff35f818d5e1355a9cf2c57541c8515f.xxnode.io",
		"265d6ed02d478816e9455572a0a292d8240d1bdb4db854e508cdef.xxnode.io",
		"847d2b5077f2afb1cb519e0297cc8b66d287f29020d68d312853e7.xxnode.io",
		"19441cb0394816046e7d76cd410dd75944255f79a37f0677a7b125.xxnode.io",
		"b09d84ed9bf54d0e8831e6c11def92bceb25eb3bd4166106ea23de.xxnode.io",
		"d81e9c84aa04d9fb848ec5de6333a0b4818735e5f7604b60b848f2.xxnode.io",
		"abed7c32d04c12100092856e2f524ed5cb9f2764ece34f1c15f512.xxnode.io",
		"936acb388ff7ab670ab83c68bf0fbb6e6abc4121edff348a11ea56.xxnode.io",
		"082de3dea0097c7bec94e0f91990a2f1b4f17a973873bd6819c076.xxnode.io",
		"5515a6895c3c129fdb7d36925fb164c882d378cd92d672424bec5c.xxnode.io",
		"541d092e3000108fbf5f24cdde601c30a04d9014bc070ae5f2c26f.xxnode.io",
		"5d244632decacbc4f08d2da5fae3a8f2854865d6e1f0380cd5ae2e.xxnode.io",
		"1d2e6ff5023016518c59b7d0b71b77924033fb624cd22a97b2d97c.xxnode.io",
		"57a8402f3ac28f1a0901c2d23819356ec6c0c9b38b52484e5e2327.xxnode.io",
		"3ece16fbf31224624fd657df55bb9f363048dcbaffa379df2ecd03.xxnode.io",
		"2649315e.xxnode.io",
		"39f0119b3114b56dd810183063fe1bc92157efb046e5bde0b98673.xxnode.io",
		"777e17f50624354acf1cb6527f9dc3e1bc87d2e8162424ddb5d5fa.xxnode.io",
		"b95e85de4ea45f6a9e4dacf16a58686d9632764da4043715f51254.xxnode.io",
		"43a01ed5093752390e7c4cf17a222d2ee321f38889516c6970ed8a.xxnode.io",
	}

	for i, expected := range expectedDnsNames {
		gwID := make([]byte, prng.Intn(maxGwIdLength*2))
		prng.Read(gwID)
		dnsName := GetGatewayDns(gwID)

		if expected != dnsName {
			t.Errorf("Unexpected DNS name for gateway ID %X (%d)."+
				"\nexpected: %q\nreceived: %q", gwID, i, expected, dnsName)
		}
	}
}
