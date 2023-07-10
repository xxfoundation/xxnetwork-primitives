////////////////////////////////////////////////////////////////////////////////
// Copyright © 2022 xx foundation                                             //
//                                                                            //
// Use of this source code is governed by a license that can be found in the  //
// LICENSE file.                                                              //
////////////////////////////////////////////////////////////////////////////////

package utils

import (
	"fmt"
	"strings"
	"testing"
)

/*
// Tests that GetIP returns the expected IP address for each valid address.
func TestGetIP_ValidAddress(t *testing.T) {
	testValues := []struct {
		addr     string
		expected net.IP
	}{
		{"10.40.210.253", net.IPv4(10, 40, 210, 253)},
		{"192.168.0.1", net.IPv4(192, 168, 0, 1)},
		{"192.168.0.1:80", net.IPv4(192, 168, 0, 1)},
		{"::FFFF:127.0.0.1", net.IPv4(127, 0, 0, 1)},
		{"::FFFF:C0A8:1", net.IPv4(192, 168, 0, 1)},
		{"::FFFF:C0A8:0001", net.IPv4(192, 168, 0, 1)},
		{"0000:0000:0000:0000:0000:FFFF:C0A8:1", net.IPv4(192, 168, 0, 1)},
		{"::FFFF:192.168.0.1", net.IPv4(192, 168, 0, 1)},
		{"[::FFFF:C0A8:1]:80", net.IPv4(192, 168, 0, 1)},
		{"::FFFF:C0A8:1:2", net.IP{0, 0, 0, 0, 0, 0, 0, 0, 0xFF, 0xFF, 0xC0, 0xA8, 0, 0x1, 0, 0x2}},
		{"[::FFFF:C0A8:1:2]:80", net.IP{0, 0, 0, 0, 0, 0, 0, 0, 0xFF, 0xFF, 0xC0, 0xA8, 0, 0x1, 0, 0x2}},
		{"2001:0db8:85a3:0000:0000:8a2e:0370:7334", net.IP{0x20, 0x1, 0xd, 0xb8, 0x85, 0xa3, 0x0, 0, 0, 0, 0x8a, 0x2e, 0x3, 0x70, 0x73, 0x34}},
		{"google.com", net.IPv4(173, 194, 215, 138)},
		{"xx.network", net.IPv4(172, 67, 68, 159)},
		{"elixxir.io", net.IPv4(104, 26, 4, 131)},
		{"permissioning.prod.cmix.rip", net.IPv4(18, 157, 190, 101)},
	}

	for i, val := range testValues {
		ip, err := GetIP(val.addr)
		if err != nil {
			t.Errorf("Failed to get IP for %q (%d): %+v", val.addr, i, err)
		}
		if !ip.Equal(val.expected) {
			t.Errorf("Failed to get expected IP for %q (%d)."+
				"\nexpected: %+v\nreceived: %+v", val.addr, i, val.expected, ip)
		}
	}
}

// Error path: tests that GetIP returns the expected error for each invalid
// address.
func TestGetIP_InvalidAddress(t *testing.T) {
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

		"localhost:80",
		"a.bc:80",
		"localhost.local:80",
		"xn--80ak6aa92e.com:80",

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

	for i, address := range addresses {
		t.Log(address)
		ip, err := GetIP(address)
		if err != nil {
			t.Errorf("Failed to get IP for address %q (%d): %+v", address, i, err)
		}
		t.Log(ip)
	}
}
*/

/*
// Happy path.
func TestIsAddress_ValidAddress(t *testing.T) {
	// TODO
}

// Error path.
func TestIsAddress_InvalidAddress(t *testing.T) {
	// TODO
}
*/

// Tests that IsPublicAddress returns nil for valid public addresses.
func TestIsPublicAddress_ValidAddress(t *testing.T) {
	addresses := []string{
		"localhost",
		"a.bc",
		"localhost.local",
		"localhost.localdomain.intern",
		"l.local.intern",
		"ru.link.n.svpncloud.com",
		"permissioning.prod.cmix.rip",

		"example.com",
		"foo.example.com",
		"bar.foo.example.com",
		"exa-mple.co.uk",
		"xn--80ak6aa92e.com",
		"9gag.com",

		"localhost:80",
		"a.bc:80",
		"localhost.local:80",
		"xn--80ak6aa92e.com:80",

		"128.0.0.1:9000",
		"128.0.0.1",
		"192.169.255.255",
		"9.255.0.255",
		"172.32.255.255",
		"::2",
		"fec0::1",
		"feff::1",
		"0100::0001:0000:0000:0000:0000",

		"one.one.one.one",
		"rfree1.blue-shield.at",
		"dns01.prd.kista.ovpn.com",
		"ns.belltele.in",
		"cache0300.ns.eu.uu.net",
		"dns.quad9.net",
		"nsx.euroweb.ro",
		"dns11.quad9.net",
		"ns1.solwaycomms.net",
		"lookup1.resolver.lax-noc.com",
		"resolver5.freedns.zone.powered.by.ihost24.com",
	}

	for i, address := range addresses {
		err := IsPublicAddress(address)
		if err != nil {
			t.Errorf("Address %q incorrectly determined to not be valid "+
				"public address (%d): %+v", address, i, err)
		}
	}
}

/*
// Tests that IsPublicAddress returns nil for valid public addresses.
func TestIsPublicAddress_ValidAddressFile(t *testing.T) {
	file, err := ReadFile("domain_list.txt")
	if err != nil {
		t.Errorf("Failed to read file: %+v", err)
	}
	addresses := strings.Split(strings.TrimSpace(string(file)), "\n")

	for i, address := range addresses {
		err := IsPublicAddress(address)
		if err != nil {
			t.Errorf("Address %q incorrectly determined to not be valid "+
				"public address (%d): %+v", address, i, err)
		}
	}
}
*/

// Error path: tests that IsPublicAddress returns the expected error for invalid
// and private addresses.
func TestIsPublicAddress_InvalidAddress(t *testing.T) {
	testValues := []struct {
		addr string
		err  string
	}{
		{"2001:0db8:85a3:0000:0000:8a2e:0370:7334", "2001:db8::/32"},
		{"[::FFFF:C0A8:1]:80", "192.168.0.0/16"},

		{"10.40.210.253", "10.0.0.0/8"},
		{"192.168.0.1", "192.168.0.0/16"},
		{"192.168.0.1:80", "192.168.0.0/16"},
		{"::FFFF:127.0.0.1", "127.0.0.0/8"},
		{"::FFFF:C0A8:1", "192.168.0.0/16"},
		{"::FFFF:C0A8:0001", "192.168.0.0/16"},
		{"0000:0000:0000:0000:0000:FFFF:C0A8:1", "192.168.0.0/16"},
		{"::FFFF:192.168.0.1", "192.168.0.0/16"},

		{"10.0.0.0", "10.0.0.0/8"},
		{"10.255.0.3", "10.0.0.0/8"},
		{"10.255.255.255", "10.0.0.0/8"},

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
		err := IsPublicAddress(val.addr)
		if err == nil || !strings.Contains(err.Error(), val.err) {
			t.Errorf("Address %q incorrectly determined to not be valid "+
				"public address (%d).\nexpected: %s\nreceived: %v",
				val.addr, i, val.err, err)
		}
	}
}

// Tests that IsValidPublicAddress returns nil for valid public addresses.
func TestIsValidPublicAddress_ValidAddress(t *testing.T) {
	addresses := []string{
		"localhost",
		"a.bc",
		"localhost.local",
		"localhost.localdomain.intern",
		"l.local.intern",
		"ru.link.n.svpncloud.com",
		"permissioning.prod.cmix.rip",

		"example.com",
		"foo.example.com",
		"bar.foo.example.com",
		"exa-mple.co.uk",
		"xn--80ak6aa92e.com",
		"9gag.com",

		"localhost:80",
		"a.bc:80",
		"localhost.local:80",
		"xn--80ak6aa92e.com:80",

		"128.0.0.1:9000",
		"128.0.0.1",
		"192.169.255.255",
		"9.255.0.255",
		"172.32.255.255",
		"::2",
		"fec0::1",
		"feff::1",
		"0100::0001:0000:0000:0000:0000",

		"one.one.one.one",
		"rfree1.blue-shield.at",
		"dns01.prd.kista.ovpn.com",
		"ns.belltele.in",
		"cache0300.ns.eu.uu.net",
		"dns.quad9.net",
		"nsx.euroweb.ro",
		"dns11.quad9.net",
		"ns1.solwaycomms.net",
		"lookup1.resolver.lax-noc.com",
		"resolver5.freedns.zone.powered.by.ihost24.com",
	}

	for i, address := range addresses {
		err := IsValidPublicAddress(address)
		if err != nil {
			t.Errorf("Address %q incorrectly determined to not be valid "+
				"public address (%d): %+v", address, i, err)
		}
	}
}

// Error path: tests that IsValidPublicAddress returns the expected error for
// invalid and private addresses.
func TestIsValidPublicAddress_InvalidAddress(t *testing.T) {
	testValues := []struct {
		addr string
		err  string
	}{
		{"", dnMinLenErr},
		{"a.b..", fmt.Sprintf(dnLabelStarPeriodErr, '.', 4)},
		{"-localhost", fmt.Sprintf(dnTldStartHyphenErr, "-localhost", 0)},
		{"localhost-", fmt.Sprintf(dnTldEndHyphenErr, "localhost-", 0)},
		{"localhost.-localdomain", fmt.Sprintf(dnTldStartHyphenErr, "-localdomain", 10)},
		{"localhost.localdomain.-int", fmt.Sprintf(dnTldStartHyphenErr, "-int", 22)},
		{"localhost.-test.int", fmt.Sprintf(dnLabelStartHyphenErr, "-test", 10)},
		{"localhost.test-.int", fmt.Sprintf(dnLabelEndHyphenErr, "test-", 10)},
		{"lÖcalhost", fmt.Sprintf(dnCharErr, 'Ö', 1)},
		{"l\uFFFDcalhost", fmt.Sprintf(dnRuneErr, '\uFFFD', 1)},
		{"localhost.lÖcaldomain", fmt.Sprintf(dnCharErr, 'Ö', 1)},
		{"localhost.localdomain.üntern", fmt.Sprintf(dnCharErr, 'ü', 22)},
		{"localhost/", fmt.Sprintf(dnCharErr, '/', 9)},
		{"[::1]", fmt.Sprintf(dnCharErr, '[', 0)},
		{"漢字汉字", fmt.Sprintf(dnCharErr, '漢', 0)},
		{"www.jubfvq1v3p38i51622y0dvmdk1mymowjyeu26gbtw9andgynj1gg8z3msb1kl5z6" +
			"906k846pj3sulm4kiyk82ln5teqj9nsht59opr0cs5ssltx78lfyvml19lfq1wp4u" +
			"sbl0o36cmiykch1vywbttcus1p9yu0669h8fj4ll7a6bmop505908s1m83q2ec2qr" +
			"9nbvql2589adma3xsq2o38os2z3dmfh2tth4is4ixyfasasasefqwe4t2ub2fz1rm" +
			"e.de", fmt.Sprintf(dnMaxLenErr, 267, domainMaxLen)},
		{"jubfvq1v3p38i51622y0dvmdk1mymowjyeu26gbtw9andgynj1gg8z3msb1kl5z6.de",
			fmt.Sprintf(dnLabelMaxLenErr,
				"jubfvq1v3p38i51622y0dvmdk1mymowjyeu26gbtw9andgynj1gg8z3msb1kl5z6",
				64, labelMaxLen)},
		{"www.jubfvq1v3p38i51622y0dvmdk1mymowjyeu26gbtw9andgynj1gg8z3msb1kl5z6",
			fmt.Sprintf(dnTldMaxLenErr,
				"jubfvq1v3p38i51622y0dvmdk1mymowjyeu26gbtw9andgynj1gg8z3msb1kl5z6",
				64, labelMaxLen)},
		{"localhost._localdomain", fmt.Sprintf(dnCharErr, '_', 10)},
		{"localhost.localdomain._int", fmt.Sprintf(dnCharErr, '_', 22)},
		{"_localhost", fmt.Sprintf(dnCharErr, '_', 0)},
		{"a.b.", dnTldEndPeriodErr},
		{"__", fmt.Sprintf(dnCharErr, '_', 0)},

		{"2001:0db8:85a3:0000:0000:8a2e:0370:7334:3445", fmt.Sprintf(dnCharErr, ':', 4)},
		{"1000.40.210.253", fmt.Sprintf(dnTldStartDigitErr, "253", 12)},
		{"2001:0db8:85a3:0000:0000:8a2e:0370:7334:3445", fmt.Sprintf(dnCharErr, ':', 4)},

		{"255.255.255.255.0", fmt.Sprintf(dnTldStartDigitErr, "0", 16)},
		{"255.255.255.", dnTldEndPeriodErr},
		{"0.0.0.256", fmt.Sprintf(dnTldStartDigitErr, "256", 6)},

		{"::/128", fmt.Sprintf(dnCharErr, ':', 0)},
	}

	for i, val := range testValues {
		err := IsValidPublicAddress(val.addr)
		if err == nil || !strings.Contains(err.Error(), val.err) {
			t.Errorf("Address %q incorrectly determined to not be valid "+
				"public address (%d).\nexpected: %s\nreceived: %v",
				val.addr, i, val.err, err)
		}
	}
}

// Tests that IsDomainName does not return an error for valid domain names.
func TestIsDomainName_ValidDomainNames(t *testing.T) {
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

		"localhost:80",
		"a.bc:80",
		"localhost.local:80",
		"xn--80ak6aa92e.com:80",
	}

	for i, address := range addresses {
		err := IsDomainName(address)
		if err != nil {
			t.Errorf("Address %q incorrectly determined to not be valid "+
				"domain name (%d): %+v", address, i, err)
		}
	}
}

// Error path: tests that IsDomainName returns the expected error for each
// invalid domain name.
func TestIsDomainName_InvalidDomainNames(t *testing.T) {
	testValues := []struct {
		addr string
		err  string
	}{
		{"", dnMinLenErr},
		{"a.b..", fmt.Sprintf(dnLabelStarPeriodErr, '.', 4)},
		{"-localhost", fmt.Sprintf(dnTldStartHyphenErr, "-localhost", 0)},
		{"localhost-", fmt.Sprintf(dnTldEndHyphenErr, "localhost-", 0)},
		{"localhost.-localdomain", fmt.Sprintf(dnTldStartHyphenErr, "-localdomain", 10)},
		{"localhost.localdomain.-int", fmt.Sprintf(dnTldStartHyphenErr, "-int", 22)},
		{"localhost.-test.int", fmt.Sprintf(dnLabelStartHyphenErr, "-test", 10)},
		{"localhost.test-.int", fmt.Sprintf(dnLabelEndHyphenErr, "test-", 10)},
		{"lÖcalhost", fmt.Sprintf(dnCharErr, 'Ö', 1)},
		{"l\uFFFDcalhost", fmt.Sprintf(dnRuneErr, '\uFFFD', 1)},
		{"localhost.lÖcaldomain", fmt.Sprintf(dnCharErr, 'Ö', 1)},
		{"localhost.localdomain.üntern", fmt.Sprintf(dnCharErr, 'ü', 22)},
		{"localhost/", fmt.Sprintf(dnCharErr, '/', 9)},
		{"127.0.0.1", fmt.Sprintf(dnTldStartDigitErr, "1", 8)},
		{"[::1]", fmt.Sprintf(dnCharErr, '[', 0)},
		{"50.50.50.50", fmt.Sprintf(dnTldStartDigitErr, "50", 9)},
		{"漢字汉字", fmt.Sprintf(dnCharErr, '漢', 0)},
		{"www.jubfvq1v3p38i51622y0dvmdk1mymowjyeu26gbtw9andgynj1gg8z3msb1kl5z6" +
			"906k846pj3sulm4kiyk82ln5teqj9nsht59opr0cs5ssltx78lfyvml19lfq1wp4u" +
			"sbl0o36cmiykch1vywbttcus1p9yu0669h8fj4ll7a6bmop505908s1m83q2ec2qr" +
			"9nbvql2589adma3xsq2o38os2z3dmfh2tth4is4ixyfasasasefqwe4t2ub2fz1rm" +
			"e.de", fmt.Sprintf(dnMaxLenErr, 267, domainMaxLen)},
		{"jubfvq1v3p38i51622y0dvmdk1mymowjyeu26gbtw9andgynj1gg8z3msb1kl5z6.de",
			fmt.Sprintf(dnLabelMaxLenErr,
				"jubfvq1v3p38i51622y0dvmdk1mymowjyeu26gbtw9andgynj1gg8z3msb1kl5z6",
				64, labelMaxLen)},
		{"www.jubfvq1v3p38i51622y0dvmdk1mymowjyeu26gbtw9andgynj1gg8z3msb1kl5z6",
			fmt.Sprintf(dnTldMaxLenErr,
				"jubfvq1v3p38i51622y0dvmdk1mymowjyeu26gbtw9andgynj1gg8z3msb1kl5z6",
				64, labelMaxLen)},
		{"localhost._localdomain", fmt.Sprintf(dnCharErr, '_', 10)},
		{"localhost.localdomain._int", fmt.Sprintf(dnCharErr, '_', 22)},
		{"_localhost", fmt.Sprintf(dnCharErr, '_', 0)},
		{"a.b.", dnTldEndPeriodErr},
		{"__", fmt.Sprintf(dnCharErr, '_', 0)},
	}

	for i, val := range testValues {
		err := IsDomainName(val.addr)
		if err == nil || !strings.Contains(err.Error(), val.err) {
			t.Errorf("Address %q incorrectly determined to be valid domain "+
				"name (%d).\nexpected: %s\nreceived: %v", val.addr, i, val.err, err)
		}
	}
}

// Tests that IsIPv4 returns true for valid IPv4 addresses and IPv4 addresses
// formatted as IPv6 addresses.
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
			t.Errorf("Address %q incorrectly determined to not be IPv4 (%d).", ip, i)
		}
	}
}

// Error path: tests that IsIPv4 returns false for IPv6 and invalid IPv4
// addresses.
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

// Tests that IsIPv6 returns true for valid IPv6 addresses.
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

// Error path: tests that IsIPv6 returns false for IPv4 and invalid IPv6
// addresses.
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

// Tests that IsIP returns true for valid IP addresses.
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

// Error path: tests that IsIP returns false for invalid IP addresses.
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

// Tests that IsPublicIP returns true for a list of different valid public IP
// addresses.
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

// Error path: tests that IsPublicIP returns an error for IPs in the start,
// middle, and end of each private address block.
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

// Tests that IsPort return true for ports greater than 0 and less than 65536.
func TestIsPort_ValidPorts(t *testing.T) {
	ports := []int{
		minPortNum, maxPortNum, minPortNum + 1, maxPortNum - 1, 23443,
	}

	for i, port := range ports {
		if !IsPort(port) {
			t.Errorf("Integer %d incorrectly determined to be a port (%d).",
				port, i)
		}
	}
}

// Error path: tests that IsPort returns false for ports less than 1 and greater
// than 65536.
func TestIsPort_InvalidPorts(t *testing.T) {
	ports := []int{
		minPortNum - 1, maxPortNum + 1, -50,
	}

	for i, port := range ports {
		if IsPort(port) {
			t.Errorf("Integer %d incorrectly determined not to be a port (%d).",
				port, i)
		}
	}
}

// Tests that IsPortString returns true for valid port strings.
func TestIsPortString_ValidPorts(t *testing.T) {
	ports := []string{
		"1", "65535", "23443",
	}

	for i, port := range ports {
		if !IsPortString(port) {
			t.Errorf("String %q incorrectly determined to be a port (%d).",
				port, i)
		}
	}
}

// Error path: tests that IsPortString returns false for invalid port strings.
func TestIsPortString_InvalidPorts(t *testing.T) {
	ports := []string{
		"-1", "65536", "-50", "hello",
	}

	for i, port := range ports {
		if IsPortString(port) {
			t.Errorf("String %q incorrectly determined not to be a port (%d).",
				port, i)
		}
	}
}

// Tests that IsEphemeralPort return true for valid ephemeral ports.
func TestIsIsEphemeralPort_ValidPorts(t *testing.T) {
	ports := []int{
		minAllowablePort, minAllowablePort + 1, 49151, 23443, maxAllowablePort,
		maxAllowablePort - 1,
	}

	for i, port := range ports {
		if !IsEphemeralPort(port) {
			t.Errorf("Integer %d incorrectly determined to be a port (%d).",
				port, i)
		}
	}
}

// Error path: tests that IsEphemeralPort returns false for ports that are
// invalid or are reserved.
func TestIsEphemeralPort_InvalidPorts(t *testing.T) {
	ports := []int{
		0, minAllowablePort - 1, 65536, -50, maxAllowablePort + 1,
	}

	for i, port := range ports {
		if IsEphemeralPort(port) {
			t.Errorf("Integer %d incorrectly determined not to be a port (%d).",
				port, i)
		}
	}
}

// Tests that IsEphemeralPortString return true for valid ephemeral ports.
func TestIsEphemeralPortString_ValidPorts(t *testing.T) {
	ports := []string{
		"1024", "1025", "49151", "23443", "65535", "65534",
	}

	for i, port := range ports {
		if !IsEphemeralPortString(port) {
			t.Errorf("String %q incorrectly determined to be a port (%d).",
				port, i)
		}
	}
}

// Error path: tests that IsEphemeralPort returns false for strings that are
// invalid ports or are reserved.
func TestIsEphemeralPortString_InvalidPorts(t *testing.T) {
	ports := []string{
		"0", "1023", "65536", "-50", "65536",
	}

	for i, port := range ports {
		if IsEphemeralPortString(port) {
			t.Errorf("String %q incorrectly determined not to be a port (%d).",
				port, i)
		}
	}
}

// Tests that ParseIP does not return an error for a list of valid IP addresses.
func Test_ParseIP(t *testing.T) {
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
		if ParseIP(address) == nil {
			t.Errorf("Failed to parse %q (%d).", address, i)
		}
	}
}

// Error path: check that ParseIP returns nil for invalid IPs.
func Test_ParseIP_InvalidIPs(t *testing.T) {
	addresses := []string{
		"2001:0db8:85a3:0000:0000:8a2e:0370:7334:3445",
		"1000.40.210.253",
		"",
	}

	for i, address := range addresses {
		if ParseIP(address) != nil {
			t.Errorf("Parsed invalid IP %q (%d).", address, i)
		}
	}
}
