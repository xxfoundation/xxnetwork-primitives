////////////////////////////////////////////////////////////////////////////////
// Copyright © 2020 Privategrity Corporation                                   /
//                                                                             /
// All rights reserved.                                                        /
////////////////////////////////////////////////////////////////////////////////

// The version for an entity is composed of a major version , minor version, and
// patch. The major and minor version numbers are both integers and dictate the
// compatibility between two versions. The patch provides extra information that
// is not part of the compatibility check, but still must be present.
package version

import (
	"github.com/pkg/errors"
	"strconv"
	"strings"
)

// String version delimiter
const delimiter = "."

// Version structure contains the major, minor, and patch information.
type Version struct {
	major int
	minor int
	patch string
}

// String prints the Version in a string format of the form "major.minor.path".
func (v *Version) String() string {
	return strconv.Itoa(v.major) + delimiter + strconv.Itoa(v.minor) +
		delimiter + v.patch
}

// ParseVersion creates a Version object based on a string. If the passed string
// is invalid, en error is returned. For a string to be valid, it must have a
// major integer, a minor integer, and a patch string separated by a period.
func ParseVersion(versionString string) (Version, error) {
	versions := strings.Split(versionString, delimiter)

	// Check that the string has three parts
	if len(versions) != 3 {
		return Version{}, errors.Errorf("Version string must contain a "+
			"major, minor, and patch version separated by %#v.", delimiter)
	}

	// Check that the major version is an integer
	major, err := strconv.Atoi(versions[0])
	if err != nil {
		return Version{}, errors.New("Major version could not be parsed as " +
			"an integer.")
	}

	// Check that the minor version is an integer
	minor, err := strconv.Atoi(versions[1])
	if err != nil {
		return Version{}, errors.New("Minor version could not be parsed as " +
			"an integer.")
	}

	// Check that the patch is not empty
	patch := versions[2]
	if patch == "" {
		return Version{}, errors.New("Patch cannot be empty.")
	}

	return Version{
		major: major,
		minor: minor,
		patch: patch,
	}, nil
}

// IsCompatible determines if the present version is compatible with the
// required version. They are compatible when the present major version is equal
// to the required version and the present minor version is greater than or
// equal to the required version. The patches do not need to match.
func IsCompatible(required, present Version) bool {
	// Compare major versions
	if required.major != present.major {
		return false
	}

	// Compare minor versions
	if required.minor > present.minor {
		return false
	}

	return true
}
