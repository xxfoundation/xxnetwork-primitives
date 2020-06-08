////////////////////////////////////////////////////////////////////////////////////////////
// Copyright © 2020 xx network SEZC                                                       //
//                                                                                        //
// Use of this source code is governed by a license that can be found in the LICENSE file //
////////////////////////////////////////////////////////////////////////////////////////////

// Contains utility operations used throughout the repo

package utils

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/mitchellh/go-homedir"
	"os"
	"path/filepath"
	"reflect"
	"testing"
	"time"
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
		t.Errorf("ExpandPath() produced an unexpected error: %v", err)
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
		t.Errorf("ExpandPath() produced an unexpected error: %v", err)
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
		t.Errorf("ExpandPath() produced an unexpected error: %v", err)
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

	// Delete the test file at the end
	defer func() {
		err := os.RemoveAll(path)
		if err != nil {
			t.Fatalf("Error deleting test file %#v:\n%v", path, err)
		}
	}()

	err := mkdirAll(path, DirPerms)
	if err != nil {
		t.Errorf("mkdirAll() produced an unexpected error: %v", err)
	}

	if _, err = os.Stat(filepath.Dir(path)); os.IsExist(err) {
		t.Errorf("mkdirAll() did not correctly make the directories:"+
			"\n\t%s", path)
	}
}

// Tests that mkdirAll() creates the directories in the specified path (does not
// include the file name) by checking if the directory structure exists.
func TestMkdirAll_DirectoryPath(t *testing.T) {
	path := "temp/temp2/"

	// Delete the test file at the end
	defer func() {
		err := os.RemoveAll(path)
		if err != nil {
			t.Fatalf("Error deleting test file %#v:\n%v", path, err)
		}
	}()

	err := mkdirAll(path, DirPerms)
	if err != nil {
		t.Errorf("mkdirAll() produced an unexpected error: %v", err)
	}

	if _, err = os.Stat(filepath.Dir(path)); os.IsExist(err) {
		t.Errorf("mkdirAll() did not correctly make the directories:"+
			"\n\t%s", path)
	}
}

// Tests that mkdirAll() does nothing for an empty path.
func TestMkdirAll_EmptyPath(t *testing.T) {
	path := ""

	// Delete the test file at the end
	defer func() {
		err := os.RemoveAll(path)
		if err != nil {
			t.Fatalf("Error deleting test file %#v:\n%v", path, err)
		}
	}()

	err := mkdirAll(path, DirPerms)
	if err != nil {
		t.Errorf("mkdirAll() produced an unexpected error: %v", err)
	}

	if _, err = os.Stat(filepath.Dir(path)); os.IsExist(err) {
		t.Errorf("mkdirAll() did not correctly make the directories:"+
			"\n\t%s", path)
	}
}

// Tests MakeDirs() by checking if the directory structure exists.
func TestMakeDirs(t *testing.T) {
	path := "temp/temp2/test.txt"

	// Delete the test file at the end
	defer func() {
		err := os.RemoveAll(path)
		if err != nil {
			t.Fatalf("Error deleting test file %#v:\n%v", path, err)
		}
	}()

	err := MakeDirs(path, DirPerms)
	if err != nil {
		t.Errorf("MakeDirs() produced an unexpected error: %v", err)
	}

	if _, err = os.Stat(filepath.Dir(path)); os.IsExist(err) {
		t.Errorf("MakeDirs() did not correctly make the directories:"+
			"\n\t%s", path)
	}
}

// Tests that MakeDirs() produces an error on an invalid path.
func TestMakeDirs_PathError(t *testing.T) {
	path := "~a/test/test.txt"

	// Delete the test file at the end
	defer func() {
		err := os.RemoveAll(path)
		if err != nil {
			t.Fatalf("Error deleting test file %#v:\n%v", path, err)
		}
	}()

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
	dirPath := "temp/temp2"
	path := dirPath + "/test.txt"
	data := []byte("test data")

	// Delete the test file at the end
	defer func() {
		err := os.RemoveAll("temp/")
		if err != nil {
			t.Fatalf("Error deleting test file %#v:\n%v", "temp/", err)
		}
	}()

	err := os.MkdirAll(dirPath, os.ModePerm)
	if err != nil {
		t.Errorf("MkdirAll() produced an unexpected error: %v", err)
	}

	err = WriteFile(path, data, DirPerms, FilePerms)
	if err != nil {
		t.Errorf("WriteFile() produced an unexpected error: %v", err)
	}

	if _, err = os.Stat(path); os.IsExist(err) {
		t.Errorf("WriteFile() did not correctly make the directories:"+
			"\n\t%s", path)
	}
}

// Tests that WriteFile() returns an error with a malformed path.
func TestWriteFile_PathError(t *testing.T) {
	path := "~a/temp/temp2/test.txt"
	data := []byte("test data")

	// Delete the test file at the end
	defer func() {
		err := os.RemoveAll(path)
		if err != nil {
			t.Fatalf("Error deleting test file %#v:\n%v", path, err)
		}
	}()

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

	// Delete the test file at the end
	defer func() {
		err := os.RemoveAll(path)
		if err != nil {
			t.Fatalf("Error deleting test file %#v:\n%v", path, err)
		}
	}()

	err := WriteFile(path, data, FilePerms, FilePerms)
	if err != nil {
		t.Errorf("WriteFile() produced an unexpected error: %v", err)
	}

	testData, err := ReadFile(path)
	if err != nil {
		t.Errorf("ReadFile() produced an unexpected error: %v", err)
	}

	if !bytes.Equal(testData, data) {
		t.Errorf("ReadFile() did not return the correct data from the file %s"+
			"\n\texpected: %s\n\treceived: %s", path, data, testData)
	}
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

// Tests that TestExist() correctly finds a file that exists.
func TestExist(t *testing.T) {
	path := "test.txt"
	data := []byte("Test string.")

	// Delete the test file at the end
	defer func() {
		err := os.RemoveAll(path)
		if err != nil {
			t.Fatalf("Error deleting test file %#v:\n%v", path, err)
		}
	}()

	err := WriteFile(path, data, FilePerms, FilePerms)
	if err != nil {
		t.Errorf("WriteFile() produced an unexpected error: %v", err)
	}

	exists := Exists(path)
	if !exists {
		t.Errorf("Exists() did not find a file that should exist")
	}
}

// Tests that TestExist() correctly finds a directory that exists.
func TestExist_Dir(t *testing.T) {
	path := "a/"

	// Delete the test file at the end
	defer func() {
		err := os.RemoveAll(path)
		if err != nil {
			t.Fatalf("Error deleting test file %#v:\n%v", path, err)
		}
	}()

	err := MakeDirs(path+"d", DirPerms)
	if err != nil {
		t.Errorf("MakeDirs() produced an unexpected error: %v", err)
	}

	exists := Exists(path)
	if !exists {
		t.Errorf("Exists() did not find a directory that should exist")
	}
}

// Tests that TestExist() returns false when a file does not exist.
func TestExist_NoFileError(t *testing.T) {
	path := "test.txt"

	// Delete the test file at the end
	defer func() {
		err := os.RemoveAll(path)
		if err != nil {
			t.Fatalf("Error deleting test file %#v:\n%v", path, err)
		}
	}()

	exists := Exists(path)
	if exists {
		t.Errorf("Exists() found a file when one does not exist")
	}
}

// Tests that FileExists() correctly finds a file that exists.
func TestFileExists(t *testing.T) {
	path := "test.txt"
	data := []byte("Test string.")

	// Delete the test file at the end
	defer func() {
		err := os.RemoveAll(path)
		if err != nil {
			t.Fatalf("Error deleting test file %#v:\n%v", path, err)
		}
	}()

	err := WriteFile(path, data, FilePerms, FilePerms)
	if err != nil {
		t.Errorf("WriteFile() produced an unexpected error: %v", err)
	}

	exists := FileExists(path)
	if !exists {
		t.Errorf("FileExists() did not find a file that should exist")
	}
}

// Tests that FileExists() false when the file is a directory.
func TestFileExists_DirError(t *testing.T) {
	path := "a/d"

	// Delete the test file at the end
	defer func() {
		err := os.RemoveAll(path)
		if err != nil {
			t.Fatalf("Error deleting test file %#v:\n%v", path, err)
		}
	}()

	err := MakeDirs(path, DirPerms)
	if err != nil {
		t.Errorf("MakeDirs() produced an unexpected error: %v", err)
	}

	exists := FileExists(path)
	if exists {
		t.Errorf("FileExists() found a directory when it was looking for a file")
	}
}

// Tests that FileExists() returns false when a file does not exist.
func TestFileExists_NoFileError(t *testing.T) {
	path := "test.txt"
	// Delete the test file at the end
	defer func() {
		err := os.RemoveAll(path)
		if err != nil {
			t.Fatalf("Error deleting test file %#v:\n%v", path, err)
		}
	}()

	exists := FileExists(path)
	if exists {
		t.Errorf("FileExists() found a file when one does not exist")
	}
}

// Tests that DirExists() correctly finds a directory that exists.
func TestDirExists(t *testing.T) {
	path := "a/"

	// Delete the test file at the end
	defer func() {
		err := os.RemoveAll(path)
		if err != nil {
			t.Fatalf("Error deleting test file %#v:\n%v", path, err)
		}
	}()

	err := MakeDirs(path+"d", DirPerms)
	if err != nil {
		t.Errorf("MakeDirs() produced an unexpected error: %v", err)
	}

	exists := DirExists(path)
	if !exists {
		t.Errorf("DirExists() did not find a directory that should exist")
	}
}

// Tests that DirExists() false when the file is a directory.
func TestDirExists_FileError(t *testing.T) {
	path := "test.txt"
	data := []byte("Test string.")

	// Delete the test file at the end
	defer func() {
		err := os.RemoveAll(path)
		if err != nil {
			t.Fatalf("Error deleting test file %#v:\n%v", path, err)
		}
	}()

	err := WriteFile(path, data, FilePerms, FilePerms)
	if err != nil {
		t.Errorf("WriteFile() produced an unexpected error: %v", err)
	}

	exists := DirExists(path)
	if exists {
		t.Errorf("DirExists() found a file when it was looking for a directory")
	}
}

// Tests that DirExists() returns false when a file does not exist.
func TestDirExists_NoDirError(t *testing.T) {
	path := "a/b/c/"

	// Delete the test file at the end
	defer func() {
		err := os.RemoveAll(path)
		if err != nil {
			t.Fatalf("Error deleting test file %#v:\n%v", path, err)
		}
	}()

	exists := FileExists(path)
	if exists {
		t.Errorf("DirExists() found a directroy when one does not exist")
	}
}

// Tests that Test_exist() correctly finds a file that exists and returns the
// correct FileInfo.
func Test_exist(t *testing.T) {
	path := "test.txt"
	data := []byte("Test string.")

	// Delete the test file at the end
	defer func() {
		err := os.RemoveAll(path)
		if err != nil {
			t.Fatalf("Error deleting test file %#v:\n%v", path, err)
		}
	}()

	err := WriteFile(path, data, FilePerms, FilePerms)
	if err != nil {
		t.Errorf("WriteFile() produced an unexpected error: %v", err)
	}

	info, exists := exists(path)
	expectedInfo, err := os.Stat(path)

	if !exists && err != nil {
		t.Errorf("exists() did not find a file that should exist:"+
			"\n\t%v", err)
	} else if !exists {
		t.Errorf("exists() did not find a file that should exist")
	}

	if !reflect.DeepEqual(info, expectedInfo) {
		t.Errorf("exists() did not return the expected FileInfo."+
			"\n\texpected: %v\n\treceived: %v", expectedInfo, info)
	}
}

// Tests that Test_exist() returns false when a file does not exist. and returns
// a nil FileInfo.
func Test_exist_NoFileError(t *testing.T) {
	path := "test.txt"

	info, exists := exists(path)

	if exists {
		t.Errorf("exists() found a file when one does not exist")
	}

	if info != nil {
		t.Errorf("exists() unexpectedly returned a non-nil FileInfo."+
			"\n\texpected: %v\n\treceived: %v", nil, info)
	}
}

// Tests that SearchDefaultLocations() finds the specified file in the user's
// home directory
func TestSearchDefaultLocations(t *testing.T) {
	testDir := fmt.Sprintf("testDir-%d/", time.Now().Nanosecond())
	testFile := fmt.Sprintf("testFile-%d.txt", time.Now().Nanosecond())
	testPath := testDir + testFile
	expectedPath, err := ExpandPath("~/." + testPath)
	if err != nil {
		t.Fatalf("ExpandPath() failed to exapnd the path %s: %+v", testPath, err)
	}
	expectedDir, err := ExpandPath("~/" + testDir)
	if err != nil {
		t.Fatalf("ExpandPath() failed to exapnd the path %s: %+v", testPath, err)
	}

	// Delete the test file at the end
	defer func() {
		err := os.RemoveAll(expectedDir)
		if err != nil {
			t.Fatalf("Error deleting test file %#v:\n%v", expectedDir, err)
		}
	}()

	err = WriteFile(expectedPath, []byte("TEST"), FilePerms, DirPerms)
	if err != nil {
		t.Fatalf("WriteFile() failed to create file %s: %+v", testPath, err)
	}

	foundPath, err := SearchDefaultLocations(testFile, testDir)
	if err != nil {
		t.Errorf("SearchDefaultLocations() produced an unexpected error: %+v",
			err)
	}

	if foundPath != expectedPath {
		t.Errorf("SearchDefaultLocations() did not find the correct file."+
			"\n\texpected: %s\n\treceived: %s", expectedPath, foundPath)
	}
}

// Tests that SearchDefaultLocations() return an error when the file does not
// exist.
func TestSearchDefaultLocations_NotFoundError(t *testing.T) {
	testDir := fmt.Sprintf(".testDir-%d/", time.Now().Nanosecond())
	testFile := fmt.Sprintf("testFile-%d.txt", time.Now().Nanosecond())

	foundPath, err := SearchDefaultLocations(testFile, testDir)
	if err == nil {
		t.Errorf("SearchDefaultLocations() did not error when expected.")
	}

	if foundPath != "" {
		t.Errorf("SearchDefaultLocations() did not return an empty path on error."+
			"\n\texpected: %s\n\treceived: %s", "", foundPath)
	}
}

// Tests that getDefaultSearchDirs generates the correct list of default paths.
func TestGetDefaultSearchDirs(t *testing.T) {
	testDir := "xxnetwork"
	expectedDir0, err := ExpandPath("~/." + testDir + "/")
	expectedDir1, err := ExpandPath("/etc/" + testDir + "/")

	testDirs, err := getDefaultSearchDirs(testDir)
	if err != nil {
		t.Errorf("getDefaultSearchDirs() produced an unxpected error: %+v", err)
	}

	if testDirs[0] != expectedDir0 {
		t.Errorf("getDefaultSearchDirs() did not return the correct path for "+
			"home.\n\texpected: %s\n\treceived: %s", expectedDir0, testDirs[0])
	}

	if testDirs[1] != expectedDir1 {
		t.Errorf("getDefaultSearchDirs() did not return the correct path for "+
			"/etc/.\n\texpected: %s\n\treceived: %s", expectedDir1, testDirs[1])
	}
}
