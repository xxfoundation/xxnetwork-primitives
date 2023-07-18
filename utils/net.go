////////////////////////////////////////////////////////////////////////////////
// Copyright © 2022 xx foundation                                             //
//                                                                            //
// Use of this source code is governed by a license that can be found in the  //
// LICENSE file.                                                              //
////////////////////////////////////////////////////////////////////////////////

package utils

import (
	"net"
	"strconv"
	"unicode/utf8"

	"github.com/pkg/errors"
)

// Maximum and Minimum lengths.
const (
	domainMaxLen   = 255
	labelMaxLen    = 63
	topLevelMaxLen = 63

	minPortNum       = 0
	maxPortNum       = 65535
	minAllowablePort = 1024
	maxAllowablePort = maxPortNum
)

// Error messages.
const (
	// lookupIpErr    = "failed to lookup host: %+v"
	invalidIPErr   = "address %q is an invalid IP address"
	nonGlobalIpErr = "address %q is not globally routable (%s: %s [%s])"

	// Error messages for domain name checking
	dnMinLenErr           = "address cannot be empty"
	dnMaxLenErr           = "address length is %d, cannot exceed %d"
	dnLabelStarPeriodErr  = "invalid character '%c' at offset %d: label cannot begin with a period"
	dnRuneErr             = "invalid rune %q at offset %d"
	dnCharErr             = "invalid character %q at offset %d"
	dnLabelMaxLenErr      = "byte length of label %q is %d, cannot exceed %d"
	dnLabelStartHyphenErr = "label %q at offset %d begins with a hyphen"
	dnLabelEndHyphenErr   = "label %q at offset %d end with a hyphen"
	dnTldEndPeriodErr     = "missing top-level domain, domain name cannot end with a period"
	dnTldMaxLenErr        = "byte length of top-level domain %q is %d, cannot exceed %d"
	dnTldStartHyphenErr   = "top-level domain %q at offset %d begins with a hyphen"
	dnTldEndHyphenErr     = "top-level domain %q at offset %d ends with a hyphen"
	dnTldStartDigitErr    = "top-level domain %q at offset %d begins with a digit"
)

/*
// GetIP returns the address as a net.IP object. Expects a valid IPv4 address,
// IPv6 address, or domain name. Ports are allowed; if a port is present, then
// it is stripped.
func GetIP(address string) (net.IP, error) {
	// Split the address from the port if the port exists
	host, _, err := net.SplitHostPort(address)
	if err == nil {
		address = host
	}

	// If the address is a valid IP, then parse it and return it
	ip := net.ParseIP(address)
	if ip != nil {
		return ip, nil
	}

	// Lookup host using the local resolver to get a list of IP addresses
	ips, err := net.LookupIP(address)
	if err != nil {
		return nil, errors.Errorf(lookupIpErr, err)
	}

	fmt.Printf("IPs: %+v\n", ips)

	// Return the first IPv4 address in the list
	for _, ip := range ips {
		if ip.To4() != nil {
			return ip, nil
		}
	}

	// Returns the first IP address for the host
	return ips[0], nil
}
*/

/*
// IsAddress determines if the given address is a valid IP address or domain
// name. Ports are allowed; if a port is present, then it is stripped.
// TODO: add tests
func IsAddress(address string) bool {
	return IsIP(address) || IsDomainName(address) != nil
}
*/

// IsPublicAddress determines if the given address is a public IP address or
// domain name. Any strings that are not IP addresses are determined to be valid
// domain names; the validity of the domain is not checked. Ports are allowed;
// if a port is present, then it is stripped.
func IsPublicAddress(address string) error {
	if ip := ParseIP(address); ip != nil {
		switch ip.To4() {
		case nil:
			return isPrivateV6(ip)
		default:
			return isPrivateV4(ip)
		}
	}

	return nil
}

// IsValidPublicAddress determines if the given address is a public IP address
// or domain name. It also checks the validity of the domain name. Ports are
// allowed; if a port is present, then it is stripped.
func IsValidPublicAddress(address string) error {
	if ip := ParseIP(address); ip != nil {
		switch ip.To4() {
		case nil:
			return isPrivateV6(ip)
		default:
			return isPrivateV4(ip)
		}
	} else {
		return IsDomainName(address)
	}
}

// IsDomainName returns an error if the domain name is not valid. Ports are
// allowed; if a port is present, then it is stripped.
// See https://tools.ietf.org/html/rfc1034#section-3.5 and
// https://tools.ietf.org/html/rfc1123#section-2.
// source: https://gist.github.com/chmike/d4126a3247a6d9a70922fc0e8b4f4013
func IsDomainName(address string) error {
	// Split the address from the port if the port exists
	host, _, err := net.SplitHostPort(address)
	if err == nil {
		address = host
	}

	// Verify the address is of the correct length
	switch {
	case len(address) == 0:
		return errors.New(dnMinLenErr)
	case len(address) > domainMaxLen:
		return errors.Errorf(dnMaxLenErr, len(address), domainMaxLen)
	}

	var l int
	for i, b := range address {
		if b == '.' {
			// Check domain labels validity
			switch {
			case i == l:
				return errors.Errorf(dnLabelStarPeriodErr, b, i)
			case i-l > 63:
				return errors.Errorf(dnLabelMaxLenErr, address[l:i], i-l, labelMaxLen)
			case address[l] == '-':
				return errors.Errorf(dnLabelStartHyphenErr, address[l:i], l)
			case address[i-1] == '-':
				return errors.Errorf(dnLabelEndHyphenErr, address[l:i], l)
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
				return errors.Errorf(dnRuneErr, c, i)
			}

			return errors.Errorf(dnCharErr, c, i)
		}
	}

	// Check top-level domain validity
	switch {
	case l == len(address):
		return errors.Errorf(dnTldEndPeriodErr)
	case len(address)-l > 63:
		return errors.Errorf(dnTldMaxLenErr,
			address[l:], len(address)-l, topLevelMaxLen)
	case address[l] == '-':
		return errors.Errorf(dnTldStartHyphenErr, address[l:], l)
	case address[len(address)-1] == '-':
		return errors.Errorf(dnTldEndHyphenErr, address[l:], l)
	case address[l] >= '0' && address[l] <= '9':
		return errors.Errorf(dnTldStartDigitErr, address[l:], l)
	}

	return nil
}

// IsIPv4 determines if the given string is a valid IPv4 address. Ports are
// allowed; if a port is present, then it is stripped.
func IsIPv4(address string) bool {
	ip := ParseIP(address)
	return ip != nil && ip.To4() != nil
}

// IsIPv6 determines if the given string is a valid IPv6 address. Ports are
// allowed; if a port is present, then it is stripped.
func IsIPv6(address string) bool {
	ip := ParseIP(address)
	return ip != nil && ip.To4() == nil
}

// IsIP determines if the given string is a valid IP address. Ports are allowed;
// if a port is present, then it is stripped.
func IsIP(address string) bool {
	return ParseIP(address) != nil
}

// IsPublicIP determines if the given string is a valid public IP address. The
// IP address may include a port. If the IP is invalid, then an error is
// returned specifying the reason. Otherwise, it returns nil. Ports are allowed;
// if a port is present, then it is stripped.
func IsPublicIP(address string) error {
	// Parse the IP address
	ip := ParseIP(address)
	if ip == nil {
		return errors.Errorf(invalidIPErr, address)
	}

	switch ip.To4() {
	case nil:
		return isPrivateV6(ip)
	default:
		return isPrivateV4(ip)
	}
}

// IsPort determines if the integer is a valid network port.
func IsPort(port int) bool {
	return port >= minPortNum && port <= maxPortNum
}

// IsPortString determines if the string is a valid network port.
func IsPortString(port string) bool {
	portNum, err := strconv.Atoi(port)

	return err == nil && IsPort(portNum)
}

// IsEphemeralPort determines if the port is ephemeral. An ephemeral port is any
// unreserved port, which is any value greater than 1024 (RFC 6056). Note that
// some ports in this range are still assigned.
// https://datatracker.ietf.org/doc/html/rfc6056#section-3.2
func IsEphemeralPort(port int) bool {
	return port >= minAllowablePort && port <= maxAllowablePort
}

// IsEphemeralPortString determines if the string is an ephemeral port.
func IsEphemeralPortString(port string) bool {
	portNum, err := strconv.Atoi(port)
	return err == nil && IsEphemeralPort(portNum)
}

// isPrivateV4 determines if the IPv4 address is a valid public IP address by
// checking that the IP is not part of a private address block. If it in a
// private block, then an error is returned specifying the block.
func isPrivateV4(ip net.IP) error {
	for _, privNet := range privateV4Networks {
		if privNet.net.Contains(ip) {
			return errors.Errorf(nonGlobalIpErr,
				ip, privNet.net.String(), privNet.name, privNet.rfc)
		}
	}

	return nil
}

// isPrivateV6 determines if the IPv6 address is a valid public IP address by
// checking that the IP is not part of a private address block. If it in a
// private block, then an error is returned specifying the block.
func isPrivateV6(ip net.IP) error {
	for _, privNet := range privateV6Networks {
		if privNet.net.Contains(ip) {
			return errors.Errorf(nonGlobalIpErr,
				ip, privNet.net.String(), privNet.name, privNet.rfc)
		}
	}
	return nil
}

// ParseIP returns the IP address as a net.IP object. Expects a valid IPv4 or
// IPv6 address. Ports are allowed; if a port is present, then it is stripped.
func ParseIP(address string) net.IP {
	// Split the address from the port if the port exists
	host, _, err := net.SplitHostPort(address)
	if err == nil {
		address = host
	}

	// Parse the IP address
	return net.ParseIP(address)
}
