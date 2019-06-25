package ndf

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"strings"
	"time"
)

var ErrNDFFile = errors.New("NDF file malformed: expected only two or more lines")

// This struct is currently generated in Terraform and decoded here
// So, if the way it's generated in Terraform changes, we also need to change
// the struct
// TODO Use UnmarshalJSON for user and node IDs and groups, at the least
//  We also need to unmarshal the Timestamp to a time.Time
//  See https://blog.gopheracademy.com/advent-2016/advanced-encoding-decoding/
//  for information about how to do this.
type NetworkDefinition struct {
	Timestamp    time.Time
	Gateways     []Gateway
	Nodes        []Node
	Registration Registration
	Udb          UDB
	E2e          Group
	Cmix         Group
}

// Gateway is the structure for the gateways object in the JSON file.
type Gateway struct {
	Address         string
	Tls_certificate string
}

// Node is the structure for the nodes object in the JSON file.
type Node struct {
	Id              []byte
	Dsa_public_key  string
	Address         string
	Tls_certificate string
}

// Registration is the structure for the registration object in the JSON
// file.
type Registration struct {
	Dsa_public_key  string
	Address         string
	Tls_certificate string
}

// UDB is the structure for the udb object in the JSON file.
type UDB struct {
	Id             []byte
	Dsa_public_key string
}

// UDB is the structure for a group in the JSON file; it is used for the
// E2e and Cmix objects.
type Group struct {
	Prime       string
	Small_prime string
	Generator   string
}

// Returns an error if base64 signature decodes incorrectly
// Returns an error if signature verification fails
// Otherwise, returns an object from the json with the contents of the file
func DecodeNDF(ndf string) (*NetworkDefinition, error) {
	// Get JSON data check if the signature is valid
	jsonString, err := separate(ndf)
	if err != nil {
		return nil, err
	}

	// Unmarshal the JSON string into a structure
	networkDefinition := &NetworkDefinition{}
	err = json.Unmarshal([]byte(jsonString), networkDefinition)
	if err != nil {
		return nil, err
	}

	return networkDefinition, nil
}

// separate splits the JSON data from the signature. The JSON data is returned
// at a string and an error is returned if the signature is invalid.
func separate(ndf string) (string, error) {
	lines := strings.Split(ndf, "\n")

	// Check that the NDF string is at least two lines
	if len(lines) < 2 {
		return "", ErrNDFFile
	}

	// Base64 decode the signature to a byte slice
	_, err := base64.StdEncoding.DecodeString(lines[1])
	if err != nil {
		return "", err
	}

	// TODO: verify the signature and not return nil error
	return lines[0], nil
}
