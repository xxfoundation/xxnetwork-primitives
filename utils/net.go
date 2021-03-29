////////////////////////////////////////////////////////////////////////////////
// Copyright © 2020 xx network SEZC                                           //
//                                                                            //
// Use of this source code is governed by a license that can be found in the  //
// LICENSE file                                                               //
////////////////////////////////////////////////////////////////////////////////

package utils

import (
	"fmt"
	"github.com/pkg/errors"
	"net"
	"strconv"
	"unicode/utf8"
)

// addrBlock describes a block of addresses that have bee registered for as
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
				IP:   net.IP{0x20, 0x1, 0xd, 0xb8, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
				Mask: net.IPMask{0xff, 0xff, 0xff, 0xff, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
			},
		},
		// 2001::/32 (RFC 4380)
		{
			name: "TEREDO",
			rfc:  "RFC4380",
			net: net.IPNet{
				IP: net.IP{0x20, 0x1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
				Mask: net.IPMask{0xff, 0xff, 0xff, 0xff, 0, 0, 0, 0, 0, 0, 0, 0,
					0, 0, 0, 0},
			},
		},
		// 2001::/23 (RFC 2928)
		{
			name: "IETF Protocol Assignments",
			rfc:  "RFC2928",
			net: net.IPNet{
				IP: net.IP{0x20, 0x1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
				Mask: net.IPMask{0xff, 0xff, 0xfe, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
					0, 0, 0},
			},
		},
		// fc00::/7 (RFC 4193)
		{
			name: "Unique-Local",
			rfc:  "RFC4193",
			net: net.IPNet{
				IP:   net.IP{0xfc, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
				Mask: net.IPMask{0xfe, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
			},
		},
		// fe80::/10 (RFC 4291)
		{
			name: "Link-Local Unicast",
			rfc:  "RFC4291",
			net: net.IPNet{
				IP:   net.IP{0xfe, 0x80, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
				Mask: net.IPMask{0xff, 0xc0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
			},
		},
		// ff00::/8 (RFC 4291)
		{
			name: "Multicast",
			rfc:  "RFC4291",
			net: net.IPNet{
				IP:   net.IP{0xff, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
				Mask: net.IPMask{0xff, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
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
				IP:   net.IP{0x20, 0x2, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
				Mask: net.IPMask{0xff, 0xff, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
			},
		},
	}
)

// IsAddress determines if the given address is a valid hostname or IP address.
func IsAddress(address string) bool {
	return IsIP(address) || IsHost(address) != nil
}

// IsPublicAddress determines if the given address is a valid hostname or public
// IP address.
func IsPublicAddress(address string) error {
	if ip := getIP(address); ip != nil {
		switch ip.To4() {
		case nil:
			return isPrivateV6(ip)
		default:
			return isPrivateV4(ip)
		}
	}
	return nil
}

// IsHost determines if the given string is a valid hostname.
func IsHost(address string) error {
	switch {
	case len(address) == 0:
		return errors.New("address is empty")
	case len(address) > 255:
		return errors.Errorf("address length is %d, cannot exceed 255", len(address))
	}

	return checkDomain(address)
}

const maxHostLength = 255

// checkDomain returns an error if the domain name is not valid
// See https://tools.ietf.org/html/rfc1034#section-3.5 and
// https://tools.ietf.org/html/rfc1123#section-2.
// source: https://gist.github.com/chmike/d4126a3247a6d9a70922fc0e8b4f4013
func checkDomain(address string) error {
	// Verify the address is of the correct length
	switch {
	case len(address) == 0:
		return errors.New("address cannot be empty")
	case len(address) > maxHostLength:
		return errors.Errorf("address length is %d, cannot exceed %d", len(address), maxHostLength)
	}

	var l int
	for i, b := range address {
		if b == '.' {
			// Check domain labels validity
			switch {
			case i == l:
				return errors.Errorf("invalid character '%c' at offset %d: label cannot begin with a period", b, i)
			case i-l > 63:
				return errors.Errorf("byte length of label '%s' is %d, cannot exceed 63", address[l:i], i-l)
			case address[l] == '-':
				return errors.Errorf("label '%s' at offset %d begins with a hyphen", address[l:i], l)
			case address[i-1] == '-':
				return errors.Errorf("label '%s' at offset %d ends with a hyphen", address[l:i], l)
			}
			l = i + 1
			continue
		}

		// Test label character validity, note: tests are ordered by decreasing
		// validity frequency
		if !(b >= 'a' && b <= 'z' || b >= '0' && b <= '9' || b == '-' || b >= 'A' && b <= 'Z') {
			// Show the printable unicode character starting at byte offset i
			c, _ := utf8.DecodeRuneInString(address[i:])
			if c == utf8.RuneError {
				return errors.Errorf("invalid rune at offset %d", i)
			}

			return errors.Errorf("invalid character '%c' at offset %d", c, i)
		}
	}

	// Check top level domain validity
	switch {
	case l == len(address):
		return errors.Errorf("missing top level domain, domain cannot end with a period")
	case len(address)-l > 63:
		return errors.Errorf("byte length of top level domain '%s' is %d, cannot exceed 63", address[l:], len(address)-l)
	case address[l] == '-':
		return errors.Errorf("top level domain '%s' at offset %d begins with a hyphen", address[l:], l)
	case address[len(address)-1] == '-':
		return errors.Errorf("top level domain '%s' at offset %d ends with a hyphen", address[l:], l)
	case address[l] >= '0' && address[l] <= '9':
		return errors.Errorf("top level domain '%s' at offset %d begins with a digit", address[l:], l)
	}

	return nil
}

// IsIPv4 determines if the given string is a valid IPv4 address. The IP address
// address may include a port.
func IsIPv4(address string) bool {
	ip := getIP(address)
	return ip != nil && ip.To4() != nil
}

// IsIPv6 determines if the given string is a valid IPv6 address. The IP address
// address may include a port.
func IsIPv6(address string) bool {
	ip := getIP(address)
	return ip != nil && ip.To4() == nil
}

// IsIP determines if the given string is is a valid IP address. The IP address
// address may include a port.
func IsIP(address string) bool {
	return getIP(address) != nil
}

// IsPublicIP determines if the given string is a valid public IP address. The
// IP address address may include a port. If the IP is invalid, then an error is
// returned specifying the reason. Otherwise, it returns nil.
func IsPublicIP(address string) error {
	// Parse the IP address
	ip := getIP(address)
	if ip == nil {
		return errors.Errorf("address %s is an invalid IP", address)
	}

	switch ip.To4() {
	case nil:
		return isPrivateV6(ip)
	default:
		return isPrivateV4(ip)
	}
}

// IsPort determines if the string is a valid network port.
func IsPort(port string) bool {
	portNum, err := strconv.Atoi(port)

	return err == nil && IsPortNum(portNum)
}

// IsPortNum determines if the integer is a valid network port.
func IsPortNum(port int) bool {
	return port > 0 && port < 65536
}

// https://www.iana.org/assignments/service-names-port-numbers/service-names-port-numbers.xhtml?&page=121

// IsValidPort determines if the port is in the allowable range of 1024 to 49151.
func IsValidPort(port int) bool {
	return port > 1023 && port < 49152
}

// isPrivateV4 returns an error if the IPv4 address is part of a private address
// block.
func isPrivateV4(ip net.IP) error {
	for _, privNet := range privateV4Networks {
		if privNet.net.Contains(ip) {
			return errors.Errorf("address %s is not globally routable (%s: %s [%s])",
				ip, privNet.net.String(), privNet.name, privNet.rfc)
		}
	}

	return nil
}

// isPrivateV6 returns an error if the IPv6 address is part of a private address
// block.
func isPrivateV6(ip net.IP) error {
	for _, privNet := range privateV6Networks {
		if privNet.net.Contains(ip) {
			return errors.Errorf("address %s is not globally routable (%s: %s [%s])",
				ip, privNet.net.String(), privNet.name, privNet.rfc)
		}
	}
	return nil
}

// getIP returns the address as a net.IP. The IP address address may include a
// port.
func getIP(address string) net.IP {
	// Split the address from the port if the port exists
	host, _, err := net.SplitHostPort(address)
	if err == nil {
		address = host
	}

	// Parse the IP address
	return net.ParseIP(address)
}

func PrintPrivateV4NetworksSV(s string) {
	var printString string
	for _, block := range privateV4Networks {
		printString += block.net.String() + s + block.name + s + block.rfc + "\n"
	}
	fmt.Print(printString)
}
