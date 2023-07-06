////////////////////////////////////////////////////////////////////////////////
// Copyright © 2022 xx foundation                                             //
//                                                                            //
// Use of this source code is governed by a license that can be found in the  //
// LICENSE file.                                                              //
////////////////////////////////////////////////////////////////////////////////

package utils

import (
	"net"
)

// addrBlock describes a block of addresses that have been registered for as
// specific purpose.
type addrBlock struct {
	name string    // Name for the purpose of the block
	rfc  string    // The standard describing block
	net  net.IPNet // The IPNet that describes the address block
}

// A list of IPv4 and IPv6 addresses block that are globally unreachable.
// List sourced from:
// https://www.iana.org/assignments/iana-ipv4-special-registry/iana-ipv4-special-registry.xhtml
// https://www.iana.org/assignments/iana-ipv6-special-registry/iana-ipv6-special-registry.xhtml
// https://github.com/letsencrypt/boulder/blob/master/bdns/dns.go
var (
	privateV4Networks = []addrBlock{
		// 10.0.0.0/8 (RFC 1918)
		{
			name: "Private-Use",
			rfc:  "RFC1918",
			net: net.IPNet{
				IP:   []byte{10, 0, 0, 0},
				Mask: []byte{255, 0, 0, 0},
			},
		},
		// 172.16.0.0/12 (RFC 1918)
		{
			name: "Private-Use",
			rfc:  "RFC1918",
			net: net.IPNet{
				IP:   []byte{172, 16, 0, 0},
				Mask: []byte{255, 240, 0, 0},
			},
		},
		// 192.168.0.0/16 (RFC 1918)
		{
			name: "Private-Use",
			rfc:  "RFC1918",
			net: net.IPNet{
				IP:   []byte{192, 168, 0, 0},
				Mask: []byte{255, 255, 0, 0},
			},
		},
		// 127.0.0.0/8 (RFC 5735, Section 3.2.1.3)
		{
			name: "Loopback",
			rfc:  "RFC1122",
			net: net.IPNet{
				IP:   []byte{127, 0, 0, 0},
				Mask: []byte{255, 0, 0, 0},
			},
		},
		// 0.0.0.0/8 (RFC 1122, Section 3.2)
		{
			name: "\"This network\"",
			rfc:  "RFC1122",
			net: net.IPNet{
				IP:   []byte{0, 0, 0, 0},
				Mask: []byte{255, 0, 0, 0},
			},
		},
		// 169.254.0.0/16 (RFC 3927)
		{
			name: "Link Local",
			rfc:  "RFC3927",
			net: net.IPNet{
				IP:   []byte{169, 254, 0, 0},
				Mask: []byte{255, 255, 0, 0},
			},
		},
		// 192.0.0.0/24 (RFC 5736, Section 2.1)
		{
			name: "IETF Protocol Assignments",
			rfc:  "RFC3927",
			net: net.IPNet{
				IP:   []byte{192, 0, 0, 0},
				Mask: []byte{255, 255, 255, 0},
			},
		},
		// 192.0.2.0/24 (RFC 5737)
		{
			name: "Documentation (TEST-NET-1)",
			rfc:  "RFC5737",
			net: net.IPNet{
				IP:   []byte{192, 0, 2, 0},
				Mask: []byte{255, 255, 255, 0},
			},
		},
		// 198.51.100.0/24 (RFC 5737)
		{
			name: "Documentation (TEST-NET-2)",
			rfc:  "RFC5737",
			net: net.IPNet{
				IP:   []byte{198, 51, 100, 0},
				Mask: []byte{255, 255, 255, 0},
			},
		},
		// 203.0.113.0/24 (RFC 5737)
		{
			name: "Documentation (TEST-NET-3)",
			rfc:  "RFC5737",
			net: net.IPNet{
				IP:   []byte{203, 0, 113, 0},
				Mask: []byte{255, 255, 255, 0},
			},
		},
		// 192.88.99.0/24 (RFC 7526)
		{
			name: "6to4 Relay Anycast",
			rfc:  "RFC7526",
			net: net.IPNet{
				IP:   []byte{192, 88, 99, 0},
				Mask: []byte{255, 255, 255, 0},
			},
		},
		// 192.18.0.0/15 (RFC 2544)
		{
			name: "Benchmarking",
			rfc:  "RFC2544",
			net: net.IPNet{
				IP:   []byte{192, 18, 0, 0},
				Mask: []byte{255, 254, 0, 0},
			},
		},
		// 224.0.0.0/4 (RFC 3171)
		{
			name: "Multicast",
			rfc:  "RFC3171",
			net: net.IPNet{
				IP:   []byte{224, 0, 0, 0},
				Mask: []byte{240, 0, 0, 0},
			},
		},
		// 255.255.255.255/32 (RFC 919, Section 7)
		{
			name: "Limited Broadcast",
			rfc:  "RFC919",
			net: net.IPNet{
				IP:   []byte{255, 255, 255, 255},
				Mask: []byte{255, 255, 255, 255},
			},
		},
		// 240.0.0.0/4 (RFC 1112, Section 4)
		{
			name: "Reserved",
			rfc:  "RFC1112",
			net: net.IPNet{
				IP:   []byte{240, 0, 0, 0},
				Mask: []byte{240, 0, 0, 0},
			},
		},
		//
		// 100.64.0.0/10 (RFC 6598)
		{
			name: "Shared Address Space",
			rfc:  "RFC6598",
			net: net.IPNet{
				IP:   []byte{100, 64, 0, 0},
				Mask: []byte{255, 192, 0, 0},
			},
		},
	}
	privateV6Networks = []addrBlock{
		// ::/128 (RFC 4291)
		{
			name: "Unspecified Address",
			rfc:  "RFC4291",
			net: net.IPNet{
				IP: net.IP{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
				Mask: net.IPMask{0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff,
					0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff},
			},
		},
		// ::1/128 (RFC 4291)
		{
			name: "Loopback Address",
			rfc:  "RFC4291",
			net: net.IPNet{
				IP: net.IP{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0x1},
				Mask: net.IPMask{0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff,
					0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff},
			},
		},
		// ::ffff:0:0/96 (RFC 4291)
		{
			name: "IPv4-mapped Address",
			rfc:  "RFC4291",
			net: net.IPNet{
				IP: net.IP{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0xff, 0xff, 0, 0, 0, 0},
				Mask: net.IPMask{0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff,
					0xff, 0xff, 0xff, 0xff, 0, 0, 0, 0},
			},
		},
		// 100::/64 (RFC 6666)
		{
			name: "Discard-Only Address Block",
			rfc:  "RFC6666",
			net: net.IPNet{
				IP: net.IP{0x1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
				Mask: net.IPMask{0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff,
					0, 0, 0, 0, 0, 0, 0, 0},
			},
		},
		// 2001:2::/48 (RFC 5180)
		{
			name: "Benchmarking",
			rfc:  "RFC5180",
			net: net.IPNet{
				IP: net.IP{0x20, 0x1, 0, 0x2, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
				Mask: net.IPMask{0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0, 0, 0, 0,
					0, 0, 0, 0, 0, 0},
			},
		},
		// 2001:db8::/32 (RFC 3849)
		{
			name: "Documentation",
			rfc:  "RFC3849",
			net: net.IPNet{
				IP: net.IP{
					0x20, 0x1, 0xd, 0xb8, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
				Mask: net.IPMask{
					0xff, 0xff, 0xff, 0xff, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
			},
		},
		// 2001::/32 (RFC 4380)
		{
			name: "TEREDO",
			rfc:  "RFC4380",
			net: net.IPNet{
				IP: net.IP{0x20, 0x1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
				Mask: net.IPMask{
					0xff, 0xff, 0xff, 0xff, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
			},
		},
		// 2001::/23 (RFC 2928)
		{
			name: "IETF Protocol Assignments",
			rfc:  "RFC2928",
			net: net.IPNet{
				IP: net.IP{0x20, 0x1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
				Mask: net.IPMask{
					0xff, 0xff, 0xfe, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
			},
		},
		// fc00::/7 (RFC 4193)
		{
			name: "Unique-Local",
			rfc:  "RFC4193",
			net: net.IPNet{
				IP: net.IP{0xfc, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
				Mask: net.IPMask{
					0xfe, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
			},
		},
		// fe80::/10 (RFC 4291)
		{
			name: "Link-Local Unicast",
			rfc:  "RFC4291",
			net: net.IPNet{
				IP: net.IP{
					0xfe, 0x80, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
				Mask: net.IPMask{
					0xff, 0xc0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
			},
		},
		// ff00::/8 (RFC 4291)
		{
			name: "Multicast",
			rfc:  "RFC4291",
			net: net.IPNet{
				IP: net.IP{
					0xff, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
				Mask: net.IPMask{
					0xff, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
			},
		},
		// 2002::/16 (RFC 7526)
		// We disable validations to IPs under the 6to4 anycase prefix because
		// there's too much risk of a malicious actor advertising the prefix and
		// answering validations for a 6to4 host they do not control.
		// https://community.letsencrypt.org/t/problems-validating-ipv6-against-host-running-6to4/18312/9
		{
			name: "6to4",
			rfc:  "RFC7526",
			net: net.IPNet{
				IP: net.IP{
					0x20, 0x2, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
				Mask: net.IPMask{
					0xff, 0xff, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
			},
		},
	}
)
