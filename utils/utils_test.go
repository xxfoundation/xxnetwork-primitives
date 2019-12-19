////////////////////////////////////////////////////////////////////////////////
// Copyright © 2020 Privategrity Corporation                                   /
//                                                                             /
// All rights reserved.                                                        /
////////////////////////////////////////////////////////////////////////////////

// Contains utility operations used throughout the repo

package utils

import (
	"bytes"
	"errors"
	"github.com/mitchellh/go-homedir"
	"os"
	"path/filepath"
	"testing"
)

const sep = string(filepath.Separator)

// Tests that ExpandPath() properly expands the the "~" character.
func TestExpandPath_Happy(t *testing.T) {
	path := sep + "test123" + sep + "test.txt"
	testPath := "~" + path
	homeDir, _ := homedir.Dir()
	expectPath := homeDir + path
	newPath, err := ExpandPath(testPath)

	if err != nil {
		t.Errorf("ExpandPath() produced an unexpected error:\n\t%v", err)
	}

	if newPath != expectPath {
		t.Errorf("ExpandPath() did not correctly expand the \"~\" character in the path %s"+
			"\n\texpected: %s\n\treceived: %s", testPath, expectPath, newPath)
	}
}

// Tests that the path is unchanged by ExpandPath().
func TestExpandPath_Default(t *testing.T) {
	path := sep + "test123" + sep + "test.txt"
	newPath, err := ExpandPath(path)

	if err != nil {
		t.Errorf("ExpandPath() produced an unexpected error:\n\t%v", err)
	}

	if newPath != path {
		t.Errorf("ExpandPath() unexpectedly modified the path %s"+
			"\n\texpected: %s\n\treceived: %s", path, path, newPath)
	}
}

// Tests that for an empty path, ExpandPath() returns an empty string.
func TestExpandPath_EmptyPath(t *testing.T) {
	path := ""
	newPath, err := ExpandPath(path)

	if err != nil {
		t.Errorf("ExpandPath() produced an unexpected error:\n\t%v", err)
	}

	if newPath != path {
		t.Errorf("ExpandPath() unexpectedly modified the path %s"+
			"\n\texpected: %s\n\treceived: %s", path, path, newPath)
	}
}

// Tests that ExpandPath() returns an error for an invalid path.
func TestExpandPath_PathError(t *testing.T) {
	path := "~a/test/test.txt"
	_, err := ExpandPath(path)

	if err == nil {
		t.Errorf("ExpandPath() did not produce error when expected:"+
			"\n\texpected: %v\n\treceived: %v",
			errors.New("cannot expand user-specific home dir"), err)
	}
}

// Tests that mkdirAll() creates the directories in the specified path (includes
// the file name) by checking if the directory structure exists.
func TestMkdirAll(t *testing.T) {
	path := "temp/temp2/test.txt"
	err := mkdirAll(path, DirPerms)

	if err != nil {
		t.Errorf("mkdirAll() produced an unexpected error:\n\t%v", err)
	}

	if _, err = os.Stat(filepath.Dir(path)); os.IsExist(err) {
		t.Errorf("mkdirAll() did not correctly make the directories:"+
			"\n\t%s", path)
	}

	// Remove the file after testing
	_ = os.RemoveAll("temp")
}

// Tests that mkdirAll() creates the directories in the specified path (does not
// include the file name) by checking if the directory structure exists.
func TestMkdirAll_DirectoryPath(t *testing.T) {
	path := "temp/temp2/"
	err := mkdirAll(path, DirPerms)

	if err != nil {
		t.Errorf("mkdirAll() produced an unexpected error:\n\t%v", err)
	}

	if _, err = os.Stat(filepath.Dir(path)); os.IsExist(err) {
		t.Errorf("mkdirAll() did not correctly make the directories:"+
			"\n\t%s", path)
	}

	// Remove the file after testing
	_ = os.RemoveAll("temp")
}

// Tests that mkdirAll() does nothing for an empty path.
func TestMkdirAll_EmptyPath(t *testing.T) {
	path := ""
	err := mkdirAll(path, DirPerms)

	if err != nil {
		t.Errorf("mkdirAll() produced an unexpected error:\n\t%v", err)
	}

	if _, err = os.Stat(filepath.Dir(path)); os.IsExist(err) {
		t.Errorf("mkdirAll() did not correctly make the directories:"+
			"\n\t%s", path)
	}
}

// Tests MakeDirs() by checking if the directory structure exists.
func TestMakeDirs(t *testing.T) {
	path := "temp/temp2/test.txt"
	err := MakeDirs(path, DirPerms)

	if err != nil {
		t.Errorf("MakeDirs() produced an unexpected error:\n\t%v", err)
	}

	if _, err = os.Stat(filepath.Dir(path)); os.IsExist(err) {
		t.Errorf("MakeDirs() did not correctly make the directories:"+
			"\n\t%s", path)
	}

	// Remove the file after testing
	_ = os.RemoveAll("temp")
}

// Tests that MakeDirs() produces an error on an invalid path.
func TestMakeDirs_PathError(t *testing.T) {
	path := "~a/test/test.txt"
	err := MakeDirs(path, DirPerms)

	if err == nil {
		t.Errorf("MakeDirs() did not produce error when expected:"+
			"\n\texpected: %v\n\treceived: %v",
			errors.New("cannot expand user-specific home dir"), err)
	}
}

// Tests WriteFile() by checking if the directory structure and the
// file exists.
func TestWriteFile(t *testing.T) {
	path := "temp/temp2/test.txt"
	data := []byte("test data")
	err := WriteFile(path, data, DirPerms, FilePerms)

	if err != nil {
		t.Errorf("WriteFile() produced an unexpected error:\n\t%v", err)
	}

	if _, err = os.Stat(path); os.IsExist(err) {
		t.Errorf("WriteFile() did not correctly make the directories:"+
			"\n\t%s", path)
	}

	// Remove the file after testing
	_ = os.RemoveAll("temp")
}

// Tests that WriteFile() returns an error with a malformed path.
func TestWriteFile_PathError(t *testing.T) {
	path := "~a/temp/temp2/test.txt"
	data := []byte("test data")
	err := WriteFile(path, data, DirPerms, FilePerms)

	if err == nil {
		t.Errorf("WriteFile() did not produce error when expected:"+
			"\n\texpected: %v\n\treceived: %v",
			errors.New("cannot expand user-specific home dir"), err)
	}
}

// Tests that ReadFile() properly reads the contents of a file created by
// WriteFile().
func TestReadFile(t *testing.T) {
	path := "test.txt"
	data := []byte("Test string.")
	err := WriteFile(path, data, FilePerms, FilePerms)

	if err != nil {
		t.Errorf("WriteFile() produced an unexpected error:\n\t%v", err)
	}

	testData, err := ReadFile(path)

	if err != nil {
		t.Errorf("ReadFile() produced an unexpected error:\n\t%v", err)
	}

	if !bytes.Equal(testData, data) {
		t.Errorf("ReadFile() did not return the correct data from the file %s"+
			"\n\texpected: %s\n\treceived: %s", path, data, testData)
	}

	// Remove the file after testing
	_ = os.RemoveAll(path)
}

// Tests that ReadFile() returns an error with a malformed path.
func TestReadFile_PathError(t *testing.T) {
	path := "~a/temp/temp2/test.txt"
	_, err := ReadFile(path)

	if err == nil {
		t.Errorf("ReadFile() did not produce error when expected:"+
			"\n\texpected: %v\n\treceived: %v",
			errors.New("cannot expand user-specific home dir"), err)
	}
}
