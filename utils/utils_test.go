////////////////////////////////////////////////////////////////////////////////
// Copyright © 2019 Privategrity Corporation                                   /
//                                                                             /
// All rights reserved.                                                        /
////////////////////////////////////////////////////////////////////////////////

// Contains utility operations used throughout the repo

package utils

import (
	"bytes"
	"github.com/mitchellh/go-homedir"
	"os"
	"path/filepath"
	"testing"
)

const sep = string(filepath.Separator)

// Tests that ExpandPath() properly expands the the "~" character.
func TestGetFullPath_Happy(t *testing.T) {
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
func TestGetFullPath_Default(t *testing.T) {
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

// Tests MakeDirs() by checking if the directory structure exists.
func TestMakeDirs(t *testing.T) {
	path := "temp/temp2/test.txt"
	err := MakeDirs(path, os.ModePerm)

	if err != nil {
		t.Errorf("MakeDirs() produced an unexpected error:\n\t%v", err)
	}

	if _, err := os.Stat(filepath.Dir(path)); os.IsExist(err) {
		t.Errorf("MakeDirs() did not correctly make the directories:"+
			"\n\t%s", path)
	}

	// Remove the file after testing
	_ = os.RemoveAll("temp")
}

// Tests TestMakeDirsAndFile() by checking if the directory structure and the
// file exists.
func TestMakeDirsAndFile(t *testing.T) {
	path := "temp/temp2/test.txt"
	data := []byte("test data")
	err := WriteFile(path, data, os.ModePerm, os.ModePerm)

	if err != nil {
		t.Errorf("WriteFile() produced an unexpected error:\n\t%v", err)
	}

	if _, err := os.Stat(path); os.IsExist(err) {
		t.Errorf("WriteFile() did not correctly make the directories:"+
			"\n\t%s", path)
	}

	// Remove the file after testing
	_ = os.RemoveAll("temp")
}

// Tests that TestMakeDirsAndFile() returns an error with a malformed path.
func TestMakeDirsAndFile_Error(t *testing.T) {
	path := "~a/temp/temp2/test.txt"
	data := []byte("test data")
	err := WriteFile(path, data, os.ModePerm, os.ModePerm)

	if err == nil {
		t.Error("WriteFile() should have produced an error.")
	}
}

func TestReadFile(t *testing.T) {
	path := "test.txt"
	data := []byte("Test string.")
	err := WriteFile(path, data, FilePerms, DirPerms)

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
