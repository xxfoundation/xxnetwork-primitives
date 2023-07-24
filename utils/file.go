////////////////////////////////////////////////////////////////////////////////
// Copyright © 2022 xx foundation                                             //
//                                                                            //
// Use of this source code is governed by a license that can be found in the  //
// LICENSE file.                                                              //
////////////////////////////////////////////////////////////////////////////////

// Package utils contains general utility functions used by our system.
// They are generic and perform basic tasks. As of writing, it mostly contains
// file IO functions to make our system be able to file IO independent of
// platform as well as domain and IP validation.

package utils

import (
	"os"
	"path/filepath"
	"time"

	"github.com/mitchellh/go-homedir"
	"github.com/pkg/errors"
)

const (
	// FilePerms is the default permissions for new files
	FilePerms = os.FileMode(0644)

	// DirPerms is the default permissions for new directory
	DirPerms = os.ModePerm
)

// ExpandPath replaces the '~' character with the user's home directory and
// cleans the path using the following rules:
//  1. Replace multiple Separator elements with a single one.
//  2. Eliminate each . path name element (the current directory).
//  3. Eliminate each inner .. path name element (the parent directory)
//     along with the non-.. element that precedes it.
//  4. Eliminate .. elements that begin a rooted path: that is, replace
//     "/.." by "/" at the beginning of a path, assuming Separator is '/'.
//  5. The returned path ends in a slash only if it represents a root
//     directory.
//  6. Any occurrences of slash are replaced by Separator.
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

// mkdirAll creates all the folders in a path that do not exist. If the path
// already exists, then nothing is done and nil is returned.
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

// WriteFile creates any directories in the path that do not exist and write
// the specified data to the file.
func WriteFile(path string, data []byte, filePerm, dirPerm os.FileMode) error {
	// Expand '~' to user's home directory and clean the path
	path, err := ExpandPath(path)
	if err != nil {
		return err
	}

	// Make directories in the path that do not already exist
	err = mkdirAll(path, dirPerm)
	if err != nil {
		return err
	}

	// Write to the specified file
	err = os.WriteFile(path, data, filePerm)
	return err
}

// WriteFileDef creates any directories in the path that do not exist and write
// the specified data to the file using the default file and directory
// permissions.
func WriteFileDef(path string, data []byte) error {
	return WriteFile(path, data, FilePerms, DirPerms)
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
	return os.ReadFile(path)
}

// Exists checks if a file or directory exists at the specified path.
func Exists(path string) bool {
	// Check if a file or directory exists at the path
	_, exists := exists(path)

	return exists
}

// FileExists checks if the file at the path exists. It returns false if the
// file does not exist or if it is a directory.
func FileExists(path string) bool {
	// Get file description information and if the file exists
	info, exists := exists(path)

	isFile := false
	if info != nil {
		isFile = !info.IsDir()
	}

	// Check if the file is a directory
	return exists && isFile
}

// DirExists checks if the directory at the path exists. It returns false if the
// directory does not exist or if it is a file.
func DirExists(path string) bool {
	// Get file description information and if the directory exists
	info, exists := exists(path)

	// Check if the file is a directory
	return exists && info.IsDir()
}

// GetLastModified returns the time the file was last modified.
func GetLastModified(path string) (time.Time, error) {
	// Get file description information and path errors
	info, err := os.Stat(path)
	if err != nil {
		return time.Time{}, err
	}

	return info.ModTime(), nil
}

// ReadDir reads the named directory, returning all its directory entries
// sorted by filename.
func ReadDir(path string) ([]string, error) {
	entries, err := os.ReadDir(path)
	if err != nil {
		return nil, err
	}

	files := make([]string, 0)
	for _, entry := range entries {
		if !entry.IsDir() {
			files = append(files, entry.Name())
		}
	}

	return files, nil
}

// exists checks if a file or directory exists at the specified path and also
// returns the file's FileInfo.
func exists(path string) (os.FileInfo, bool) {
	// Expand '~' to user's home directory and clean the path
	path, err := ExpandPath(path)
	if err != nil {
		return nil, false
	}

	// Get file description information and path errors
	info, err := os.Stat(path)

	// Check if a file or directory exists at the path
	return info, !os.IsNotExist(err)
}

// SearchDefaultLocations searches for a file path in a default directory in
// a number of hard-coded paths, including the user's home folder and /etc/. If
// the file is found, then its full path is returned. Otherwise, the path is
// blank and an error is returned.
//
// Note that defaultDirectory MUST be a relative path. By default, when checking
// the home directory, a "." is prepended to the defaultDirectory.
func SearchDefaultLocations(
	defaultFileName string, defaultDirectory string) (string, error) {
	// Get the user's home directory
	defaultDirs, err := getDefaultSearchDirs(defaultDirectory)
	if err != nil {
		return "", errors.Errorf("Could not get home directory: %+v", err)
	}

	// Search the directories for the file
	for _, dir := range defaultDirs {
		// Format the path and check for errors
		path := dir + "/" + defaultFileName
		foundFilePath, err := ExpandPath(path)
		if err != nil {
			return "", errors.Errorf("Error expanding path %s: %v", path, err)
		}

		// If the file exists, return its path
		if FileExists(foundFilePath) {
			return foundFilePath, nil
		}
	}

	return "", errors.Errorf("Could not find %s in any of the directories: %v",
		defaultFileName, defaultDirs)
}

// getDefaultSearchDirs retrieves the list of default directories to search for
// configuration files in. Note that defaultDirectory MUST be a relative path.
func getDefaultSearchDirs(defaultDirectory string) ([]string, error) {
	var searchDirs []string

	// Get the user's home directory
	home, err := homedir.Dir()
	if err != nil {
		return nil, errors.Errorf("Could not get home directory: %+v", err)
	}

	// Add the home directory to the search
	searchDirs = append(searchDirs, filepath.Clean(home+"/."+defaultDirectory+"/"))

	// Add /opt/ to the search
	searchDirs = append(searchDirs, filepath.Clean("/opt/"+defaultDirectory+"/"))

	// Add /etc/ to the search
	searchDirs = append(searchDirs, filepath.Clean("/etc/"+defaultDirectory+"/"))

	return searchDirs, nil
}
