////////////////////////////////////////////////////////////////////////////////
// Copyright © 2019 Privategrity Corporation                                   /
//                                                                             /
// All rights reserved.                                                        /
////////////////////////////////////////////////////////////////////////////////

package utils

import (
	"github.com/mitchellh/go-homedir"
	"io/ioutil"
	"os"
	"path/filepath"
)

const (
	// Permissions for new files/directories
	FilePerms = os.FileMode(0644)
	DirPerms  = os.ModePerm
)

// ExpandPath replaces the '~' character with the user's home directory and
// cleans the path using the following rules:
//	1. Replace multiple Separator elements with a single one.
//	2. Eliminate each . path name element (the current directory).
//	3. Eliminate each inner .. path name element (the parent directory)
//	   along with the non-.. element that precedes it.
//	4. Eliminate .. elements that begin a rooted path: that is, replace
//	   "/.." by "/" at the beginning of a path, assuming Separator is '/'.
//	5. The returned path ends in a slash only if it represents a root
//	   directory.
//	6. Any occurrences of slash are replaced by Separator.
func ExpandPath(path string) (string, error) {
	// If the path is empty, then return nothing
	if path == "" {
		return "", nil
	}

	// Replace the '~' character with the user's home directory
	path, err := homedir.Expand(path)
	if err != nil {
		return "", err
	}

	// Cleans the path using the rules in the function description
	path = filepath.Clean(path)

	return path, nil
}

// mkdirAll creates all the folders in a path that do not exist.
func mkdirAll(path string, perm os.FileMode) error {
	// Strip file name from the path
	dir := filepath.Dir(path)

	// Create the directories
	return os.MkdirAll(dir, perm)
}

// MakeDirs expands and cleans the path and then creates all the folders in a
// path that do not exist.
func MakeDirs(path string, perm os.FileMode) error {
	// Expand '~' to user's home directory and clean the path
	path, err := ExpandPath(path)
	if err != nil {
		return err
	}

	// Create all directories in path, if they do not already exist
	return mkdirAll(path, perm)
}

// WriteFile creates any directories in the path that do not exists and write
// the specified data to the file.
func WriteFile(path string, data []byte, filePerm, dirPerm os.FileMode) error {
	// Expand '~' to user's home directory and clean the path
	path, err := ExpandPath(path)
	if err != nil {
		return err
	}

	// Make an directories in the path that do not already exist
	err = mkdirAll(path, dirPerm)
	if err != nil {
		return err
	}

	// Write to the specified file
	err = ioutil.WriteFile(path, data, filePerm)
	return err
}

// ReadFile expands and cleans the specified path, reads the file, and returns
// its contents.
func ReadFile(path string) ([]byte, error) {
	// Expand '~' to user's home directory and clean the path
	path, err := ExpandPath(path)
	if err != nil {
		return nil, err
	}

	// Read the file and return the contents
	return ioutil.ReadFile(path)
}
