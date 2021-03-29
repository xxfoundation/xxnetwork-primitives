////////////////////////////////////////////////////////////////////////////////
// Copyright © 2020 xx network SEZC                                           //
//                                                                            //
// Use of this source code is governed by a license that can be found in the  //
// LICENSE file                                                               //
////////////////////////////////////////////////////////////////////////////////

package utils

import (
	"strings"
	"testing"
)

// Happy path.
func TestIsAddress_ValidAddress(t *testing.T) {
	// TODO
}

// Error path.
func TestIsAddress_InvalidAddress(t *testing.T) {
	// TODO
}

// Happy path.
func TestIsPublicAddress_ValidAddress(t *testing.T) {
	// TODO
}

// Error path.
func TestIsPublicAddress_InvalidAddress(t *testing.T) {
	// TODO
}

// Happy path.
// TODO
func TestIsHost_ValidHosts(t *testing.T) {
	addresses := []string{
		"localhost",
		"a.bc",
		"localhost.local",
		"localhost.localdomain.intern",
		"l.local.intern",
		"ru.link.n.svpncloud.com",

		"example.com",
		"foo.example.com",
		"bar.foo.example.com",
		"exa-mple.co.uk",
		"xn--80ak6aa92e.com",
		"9gag.com",

		"localhost._localdomain",
		"localhost.localdomain._int",
		"_localhost",
		"a.b.",
		"__",
	}

	for i, address := range addresses {
		err := IsHost(address)
		if err != nil {
			t.Errorf("Address %s incorrectly determined to not be host (%d): %+v",
				address, i, err)
		}
	}
}

// Error path.
func TestIsHost_InvalidHosts(t *testing.T) {
	// TODO
	testValues := []struct {
		addr string
		err  string
	}{
		{"a.b..", "label cannot begin with a period"},
		{"-localhost", "begins with a hyphen"},
		{"localhost.-localdomain", "begins with a hyphen"},
		{"localhost.localdomain.-int", "begins with a hyphen"},
		{"lÖcalhost", "invalid character"},
		{"localhost.lÖcaldomain", "invalid character"},
		{"localhost.localdomain.üntern", "invalid character"},
		{"localhost/", "invalid character"},
		{"127.0.0.1", "begins with a digit"},
		{"[::1]", "invalid character"},
		{"50.50.50.50", "begins with a digit"},
		{"localhost.localdomain.intern:65535", "invalid character"},
		{"漢字汉字", "invalid character"},
		{"www.jubfvq1v3p38i51622y0dvmdk1mymowjyeu26gbtw9andgynj1gg8z3msb1kl5z6" +
			"906k846pj3sulm4kiyk82ln5teqj9nsht59opr0cs5ssltx78lfyvml19lfq1wp4u" +
			"sbl0o36cmiykch1vywbttcus1p9yu0669h8fj4ll7a6bmop505908s1m83q2ec2qr" +
			"9nbvql2589adma3xsq2o38os2z3dmfh2tth4is4ixyfasasasefqwe4t2ub2fz1rm" +
			"e.de", "cannot exceed 255"},
	}

	for i, val := range testValues {
		err := IsHost(val.addr)
		if err == nil || !strings.Contains(err.Error(), val.err) {
			t.Errorf("Address %s incorrectly determined to be host (%d)."+
				"\nexpected: %s\nreceived: %v", val.addr, i, val.err, err)
		}
	}
}

// Happy path.
func TestIsIPv4_ValidIPs(t *testing.T) {
	addresses := []string{
		"10.40.210.253",
		"192.168.0.1",
		"192.168.0.1:80",
		"::FFFF:127.0.0.1",
		"::FFFF:C0A8:1",
		"::FFFF:C0A8:0001",
		"0000:0000:0000:0000:0000:FFFF:C0A8:1",
		"::FFFF:192.168.0.1",
		"[::FFFF:C0A8:1]:80",
	}

	for i, ip := range addresses {
		if !IsIPv4(ip) {
			t.Errorf("Address %s incorrectly determined to not be IPv4 (%d).", ip, i)
		}
	}
}

// Error path.
func TestIsIPv4_InvalidIPs(t *testing.T) {
	addresses := []string{
		"::FFFF:C0A8:1:2",
		"[::FFFF:C0A8:1:2]:80",
		"2001:0db8:85a3:0000:0000:8a2e:0370:7334:3445",
		"2001:0db8:85a3:0000:0000:8a2e:0370:7334",
		"1000.40.210.253",
	}

	for i, ip := range addresses {
		if IsIPv4(ip) {
			t.Errorf("Address %s incorrectly determined to be IPv4 (%d).", ip, i)
		}
	}
}

// Happy path.
func TestIsIPv6_ValidIPs(t *testing.T) {
	addresses := []string{
		"::FFFF:C0A8:1:2",
		"[::FFFF:C0A8:1:2]:80",
		"2001:0db8:85a3:0000:0000:8a2e:0370:7334",
	}

	for i, ip := range addresses {
		if !IsIPv6(ip) {
			t.Errorf("Address %s incorrectly determined to not be IPv6 (%d).", ip, i)
		}
	}
}

// Error path.
func TestIsIPv6_InvalidIPs(t *testing.T) {
	addresses := []string{
		"[::FFFF:C0A8:1]:80",
		"10.40.210.253",
		"192.168.0.1",
		"192.168.0.1:80",
		"::FFFF:127.0.0.1",
		"::FFFF:C0A8:1",
		"::FFFF:C0A8:0001",
		"0000:0000:0000:0000:0000:FFFF:C0A8:1",
		"::FFFF:192.168.0.1",
		"2001:0db8:85a3:0000:0000:8a2e:0370:7334:3445",
	}

	for i, ip := range addresses {
		if IsIPv6(ip) {
			t.Errorf("Address %s incorrectly determined to be IPv6 (%d).", ip, i)
		}
	}
}

// Happy path.
func TestIsIP_ValidIPs(t *testing.T) {
	addresses := []string{
		"128.0.0.1:9000",
		"128.0.0.1",
		"192.169.255.255",
		"9.255.0.255",
		"172.32.255.255",
		"::2",
		"fec0::1",
		"feff::1",
		"0100::0001:0000:0000:0000:0000",
	}

	for i, ip := range addresses {
		if !IsIP(ip) {
			t.Errorf("String %s incorrectly determined not to be IP (%d).", ip, i)
		}
	}
}

// Happy path.
func TestIsIP_InvalidIPs(t *testing.T) {
	addresses := []string{
		"255.255.255.255.0",
		"255.255.255.",
		"0.0.0.256",
	}

	for i, ip := range addresses {
		if IsIP(ip) {
			t.Errorf("String %s incorrectly determined to be IP (%d).", ip, i)
		}
	}
}

// Happy path.
func TestIsPublicIP_ValidIPs(t *testing.T) {
	addresses := []string{
		"128.0.0.1:9000",
		"128.0.0.1",
		"192.169.255.255",
		"9.255.0.255",
		"172.32.255.255",
		"::2",
		"fec0::1",
		"feff::1",
		"0100::0001:0000:0000:0000:0000",
	}

	for i, ip := range addresses {
		err := IsPublicIP(ip)
		if err != nil {
			t.Errorf("Received error for valid IP address %s (%d): %s", ip, i, err)
		}
	}
}

// Error path: tests that an IP in the start, middle, and end of each private
// address block triggers an error.
func TestIsPublicIP_InvalidIPs(t *testing.T) {
	testValues := []struct {
		ip  string
		err string
	}{
		{"255.255.255.255.0", "invalid IP address"},
		{"255.255.255.", "invalid IP address"},
		{"0.0.0.256", "invalid IP address"},

		{"10.0.0.0", "10.0.0.0/8"},
		{"10.255.0.3", "10.0.0.0/8"},
		{"10.255.255.255", "10.0.0.0/8"},

		{"172.16.0.0", "172.16.0.0/12"},
		{"172.16.255.255", "172.16.0.0/12"},
		{"172.31.255.255", "172.16.0.0/12"},

		{"192.168.0.0", "192.168.0.0/16"},
		{"192.168.254.254", "192.168.0.0/16"},
		{"192.168.255.255", "192.168.0.0/16"},

		{"127.0.0.0", "127.0.0.0/8"},
		{"127.0.0.1", "127.0.0.0/8"},
		{"127.255.255.255", "127.0.0.0/8"},

		{"0.0.0.0", "0.0.0.0/8"},
		{"0.0.0.1", "0.0.0.0/8"},
		{"0.255.255.255", "0.0.0.0/8"},

		{"169.254.0.0", "169.254.0.0/16"},
		{"169.254.42.55", "169.254.0.0/16"},
		{"169.254.255.255", "169.254.0.0/16"},

		{"192.0.0.0", "192.0.0.0/24"},
		{"192.0.0.197", "192.0.0.0/24"},
		{"192.0.0.255", "192.0.0.0/24"},

		{"192.0.2.0", "192.0.2.0/24"},
		{"192.0.2.133", "192.0.2.0/24"},
		{"192.0.2.255", "192.0.2.0/24"},

		{"198.51.100.0", "198.51.100.0/24"},
		{"198.51.100.102", "198.51.100.0/24"},
		{"198.51.100.255", "198.51.100.0/24"},

		{"203.0.113.0", "203.0.113.0/24"},
		{"203.0.113.180", "203.0.113.0/24"},
		{"203.0.113.255", "203.0.113.0/24"},

		{"192.88.99.0", "192.88.99.0/24"},
		{"192.88.99.239", "192.88.99.0/24"},
		{"192.88.99.255", "192.88.99.0/24"},

		{"192.18.0.0", "192.18.0.0/15"},
		{"192.18.199.93", "192.18.0.0/15"},
		{"192.19.255.255", "192.18.0.0/15"},

		{"224.0.0.0", "224.0.0.0/4"},
		{"224.14.211.255", "224.0.0.0/4"},
		{"239.255.255.255", "224.0.0.0/4"},

		{"240.0.0.0", "240.0.0.0/4"},
		{"240.147.219.211", "240.0.0.0/4"},
		{"255.255.255.254", "240.0.0.0/4"},

		{"255.255.255.255", "255.255.255.255/32"},

		{"100.64.0.0", "100.64.0.0/10"},
		{"100.92.11.255", "100.64.0.0/10"},
		{"100.127.255.255", "100.64.0.0/10"},

		{"::/128", "::/128"},

		// {"::ffff:0:0", "::ffff:0:0/96"},
		// {"::ffff:459b:d6eb", "::ffff:0:0/96"},
		// {"::ffff:ffff:ffff", "::ffff:0:0/96"},

		{"100::", "100::/64"},
		{"100::54CE:37FA:E888:9FD7", "100::/64"},
		{"100::FFFF:FFFF:FFFF:FFFF", "100::/64"},

		// {"2001::", "2001::/23"},
		{"2001:00E8:2FCF:D96E:965E:6AC4:D050:907A", "2001::/23"},
		{"2001:1ff:ffff:ffff:ffff:ffff:ffff:ffff", "2001::/23"},

		{"2001:2::", "2001:2::/48"},
		{"2001:2::10f4:a868:917f:b819:1158", "2001:2::/48"},
		{"2001:2::ffff:ffff:ffff:ffff:ffff", "2001:2::/48"},

		{"2001:db8::", "2001:db8::/32"},
		{"2001:db8:4dd8:f94:2121:d712:4ffa:e8e0", "2001:db8::/32"},
		{"2001:db8:ffff:ffff:ffff:ffff:ffff:ffff", "2001:db8::/32"},

		{"2001::", "2001::/32"},
		{"2001::5b01:c5d9:d151:3c69:d0a9:6db7", "2001::/32"},
		{"2001::ffff:ffff:ffff:ffff:ffff:ffff", "2001::/32"},

		{"fc00::", "fc00::/7"},
		{"fc75:dcc6:c94a:81ed:252c:681a:b200:8088", "fc00::/7"},
		{"fdff:ffff:ffff:ffff:ffff:ffff:ffff:ffff", "fc00::/7"},

		{"fe80::", "fe80::/10"},
		{"fe92:3c00:e510:b69a:3314:d5e3:b3cf:182b", "fe80::/10"},
		{"febf:ffff:ffff:ffff:ffff:ffff:ffff:ffff", "fe80::/10"},

		{"ff00::", "ff00::/8"},
		{"ffb0:f622:b362:29f6:4dd2:5ee0:bf97:551b", "ff00::/8"},
		{"ffff:ffff:ffff:ffff:ffff:ffff:ffff:ffff", "ff00::/8"},

		{"2002::", "2002::/16"},
		{"2002:e05d:9290:fbaf:b69e:94ef:4d2e:6c0b", "2002::/16"},
		{"2002:ffff:ffff:ffff:ffff:ffff:ffff:ffff", "2002::/16"},

		{"0100::", "100::/64"},
		{"0100::0000:ffff:ffff:ffff:ffff", "100::/64"},

		{"fe80::1", "fe80::/10"},
		{"febf::1", "fe80::/10"},
		{"ff00::1", "ff00::/8"},
		{"ff10::1", "ff00::/8"},
		{"ffff:ffff:ffff:ffff:ffff:ffff:ffff:ffff", "ff00::/8"},
		{"2002::", "2002::/16"},
		{"2002:ffff:ffff:ffff:ffff:ffff:ffff:ffff", "2002::/16"},
	}

	for i, val := range testValues {
		err := IsPublicIP(val.ip)
		if err == nil || !strings.Contains(err.Error(), val.err) {
			t.Errorf("Failed to error for invalid address %s (%d)."+
				"\nexpected: %s\nreceived: %v", val.ip, i, val.err, err)
		}
	}
}

// Happy path.
func TestIsPort_ValidPorts(t *testing.T) {
	ports := []string{
		"1", "65535", "23443",
	}

	for i, port := range ports {
		if !IsPort(port) {
			t.Errorf("String %s incorrectly determined to be a port (%d).", port, i)
		}
	}
}

// Error path.
func TestIsPort_InvalidPorts(t *testing.T) {
	ports := []string{
		"0", "65536", "-50",
	}

	for i, port := range ports {
		if IsPort(port) {
			t.Errorf("String %s incorrectly determined not to be a port (%d).", port, i)
		}
	}
}

// Happy path.
func TestIsPortNum_ValidPorts(t *testing.T) {
	ports := []int{
		1, 65535, 23443,
	}

	for i, port := range ports {
		if !IsPortNum(port) {
			t.Errorf("String %d incorrectly determined to be a port (%d).", port, i)
		}
	}
}

// Error path.
func TestIsPortNum_InvalidPorts(t *testing.T) {
	ports := []int{
		0, 65536, -50,
	}

	for i, port := range ports {
		if IsPortNum(port) {
			t.Errorf("String %d incorrectly determined not to be a port (%d).", port, i)
		}
	}
}

// Happy path.
func TestIsIsValidPort_ValidPorts(t *testing.T) {
	ports := []int{
		1024, 49151, 23443,
	}

	for i, port := range ports {
		if !IsValidPort(port) {
			t.Errorf("String %d incorrectly determined to be a port (%d).", port, i)
		}
	}
}

// Error path.
func TestIsValidPort_InvalidPorts(t *testing.T) {
	ports := []int{
		1023, 49152, -50,
	}

	for i, port := range ports {
		if IsValidPort(port) {
			t.Errorf("String %d incorrectly determined not to be a port (%d).", port, i)
		}
	}
}

// Happy path.
func Test_getIP(t *testing.T) {
	addresses := []string{
		"192.168.0.1",
		"192.168.0.1",
		"192.168.0.1:80",
		"::FFFF:127.0.0.1",
		"::FFFF:C0A8:1",
		"::FFFF:C0A8:0001",
		"0000:0000:0000:0000:0000:FFFF:C0A8:1",
		"::FFFF:192.168.0.1",
		"[::FFFF:C0A8:1]:80",
		"2001:0db8:85a3:0000:0000:8a2e:0370:7334",
	}

	for i, address := range addresses {
		if getIP(address) == nil {
			t.Errorf("Failed to parse %s (%d).", address, i)
		}
	}
}

// Error path: check that nil is returned for invalid IPs.
func Test_getIP_InvalidIPs(t *testing.T) {
	addresses := []string{
		"2001:0db8:85a3:0000:0000:8a2e:0370:7334:3445",
		"1000.40.210.253",
		"",
	}

	for i, address := range addresses {
		if getIP(address) != nil {
			t.Errorf("Parsed invalid IP %s (%d).", address, i)
		}
	}
}

func TestPrintPrivateV4NetworksCSV(t *testing.T) {
	PrintPrivateV4NetworksSV("\t")
}
