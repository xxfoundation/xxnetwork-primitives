////////////////////////////////////////////////////////////////////////////////
// Copyright © 2022 xx foundation                                             //
//                                                                            //
// Use of this source code is governed by a license that can be found in the  //
// LICENSE file.                                                              //
////////////////////////////////////////////////////////////////////////////////

// Package version describes a version used for a repository.
//
// The version for an entity is composed of a major version, minor version, and
// patch. The major and minor version numbers are both integers and dictate the
// compatibility between two versions. The patch provides extra information that
// is not part of the compatibility check, but still must be present.
package version

import (
	"strconv"
	"strings"

	"github.com/pkg/errors"
)

// String version delimiter
const delimiter = "."

// Version structure contains the major, minor, and patch information.
type Version struct {
	major int
	minor int
	patch string
}

// New returns a new Version with the given major, minor and patch.
func New(major, minor int, patch string) Version {
	return Version{
		major: major,
		minor: minor,
		patch: patch,
	}
}

// Major returns the major integer of the Version.
func (v Version) Major() int {
	return v.major
}

// Minor returns the minor integer of the Version.
func (v Version) Minor() int {
	return v.minor
}

// Patch returns the patch string of the Version.
func (v Version) Patch() string {
	return v.patch
}

// String prints the Version in a string format of the form "major.minor.path".
func (v Version) String() string {
	return strconv.Itoa(v.major) + delimiter + strconv.Itoa(v.minor) +
		delimiter + v.patch
}

// ParseVersion parses a string into a Version object. An error is returned for
// invalid version string. To be valid, a string must contain a major integer,
// a minor integer, and a patch string separated by a period.
func ParseVersion(versionString string) (Version, error) {
	// Separate string into the three parts
	versionParts := strings.SplitN(versionString, delimiter, 3)

	// Check that the string has three parts
	if len(versionParts) != 3 {
		return Version{}, errors.Errorf("version string must contain 3 parts: "+
			"a major, minor, and patch separated by \"%s\". Received string "+
			"with %d part(s).", delimiter, len(versionParts))
	}

	version := Version{}
	var err error

	// Check that the major version is an integer
	version.major, err = strconv.Atoi(versionParts[0])
	if err != nil {
		return Version{},
			errors.Errorf("expected integer for major version: %+v", err)
	}

	// Check that the minor version is an integer
	version.minor, err = strconv.Atoi(versionParts[1])
	if err != nil {
		return Version{},
			errors.Errorf("expected integer for minor version: %+v", err)
	}

	// Check that the patch is not empty
	if versionParts[2] == "" {
		return Version{}, errors.New("patch cannot be empty.")
	}
	version.patch = versionParts[2]

	return version, nil
}

// IsCompatible determines if the current Version is compatible with the
// required Version. Version are compatible when the major versions are equal
// and the current minor version is greater than or equal to the required
// version. The patches do not need to match.
func IsCompatible(required, current Version) bool {
	// Compare major versions
	if required.major != current.major {
		return false
	}

	// Compare minor versions
	if required.minor > current.minor {
		return false
	}

	return true
}

// Equal determines if two Version objects are equal. They are equal when major,
// minor, and patch are all the same
func Equal(a, b Version) bool {
	return a.major == b.major && a.minor == b.minor && a.patch == b.patch
}

// Cmp compares two Versions. Return 1 if a is greater than b, -1 if a is less
// than b, and 0 if they are equal. A Version
func Cmp(a, b Version) int {
	if a.major > b.major {
		return 1
	} else if a.major < b.major {
		return -1
	}

	if a.minor > b.minor {
		return 1
	} else if a.minor < b.minor {
		return -1
	}

	return 0
}
