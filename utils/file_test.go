////////////////////////////////////////////////////////////////////////////////
// Copyright © 2022 xx foundation                                             //
//                                                                            //
// Use of this source code is governed by a license that can be found in the  //
// LICENSE file.                                                              //
////////////////////////////////////////////////////////////////////////////////

// Contains utility operations used throughout the repo

package utils

import (
	"bytes"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"sort"
	"testing"
	"time"

	"github.com/mitchellh/go-homedir"
	"github.com/stretchr/testify/require"
)

const sep = string(filepath.Separator)

// Tests that ExpandPath properly expands the "~" character.
func TestExpandPath_Happy(t *testing.T) {
	path := sep + "test123" + sep + "test.txt"
	testPath := "~" + path
	homeDir, _ := homedir.Dir()
	expectPath := homeDir + path
	newPath, err := ExpandPath(testPath)

	if err != nil {
		t.Errorf("ExpandPath produced an unexpected error: %v", err)
	}

	if newPath != expectPath {
		t.Errorf("ExpandPath did not correctly expand the \"~\" character in "+
			"the path %s\nexpected: %s\nreceived: %s",
			testPath, expectPath, newPath)
	}
}

// Tests that the path is unchanged by ExpandPath.
func TestExpandPath_Default(t *testing.T) {
	path := sep + "test123" + sep + "test.txt"
	newPath, err := ExpandPath(path)

	if err != nil {
		t.Errorf("ExpandPath produced an unexpected error: %v", err)
	}

	if newPath != path {
		t.Errorf("ExpandPath unexpectedly modified the path %s"+
			"\nexpected: %s\nreceived: %s", path, path, newPath)
	}
}

// Tests that for an empty path, ExpandPath returns an empty string.
func TestExpandPath_EmptyPath(t *testing.T) {
	path := ""
	newPath, err := ExpandPath(path)

	if err != nil {
		t.Errorf("ExpandPath produced an unexpected error: %v", err)
	}

	if newPath != path {
		t.Errorf("ExpandPath unexpectedly modified the path %s"+
			"\nexpected: %s\nreceived: %s", path, path, newPath)
	}
}

// Tests that ExpandPath returns an error for an invalid path.
func TestExpandPath_PathError(t *testing.T) {
	path := "~a/test/test.txt"
	_, err := ExpandPath(path)

	if err == nil {
		t.Errorf("ExpandPath did not produce error when expected:"+
			"\nexpected: %v\nreceived: %v",
			errors.New("cannot expand user-specific home dir"), err)
	}
}

// Tests that mkdirAll creates the directories in the specified path (includes
// the file name) by checking if the directory structure exists.
func TestMkdirAll(t *testing.T) {
	path := "temp/temp2/test.txt"

	// Delete the test file at the end
	defer func() {
		err := os.RemoveAll(path)
		if err != nil {
			t.Fatalf("Error deleting test file %q: %v", path, err)
		}
	}()

	err := mkdirAll(path, DirPerms)
	if err != nil {
		t.Errorf("mkdirAll produced an unexpected error: %v", err)
	}

	if _, err = os.Stat(filepath.Dir(path)); os.IsExist(err) {
		t.Errorf("mkdirAll did not correctly make the directories: %s", path)
	}
}

// Tests that mkdirAll creates the directories in the specified path (does not
// include the file name) by checking if the directory structure exists.
func TestMkdirAll_DirectoryPath(t *testing.T) {
	path := "temp/temp2/"

	// Delete the test file at the end
	defer func() {
		err := os.RemoveAll(path)
		if err != nil {
			t.Fatalf("Error deleting test file %q: %v", path, err)
		}
	}()

	err := mkdirAll(path, DirPerms)
	if err != nil {
		t.Errorf("mkdirAll produced an unexpected error: %v", err)
	}

	if _, err = os.Stat(filepath.Dir(path)); os.IsExist(err) {
		t.Errorf("mkdirAll did not correctly make the directories: %s", path)
	}
}

// Tests that mkdirAll does nothing for an empty path.
func TestMkdirAll_EmptyPath(t *testing.T) {
	path := ""

	// Delete the test file at the end
	defer func() {
		err := os.RemoveAll(path)
		if err != nil {
			t.Fatalf("Error deleting test file %q: %v", path, err)
		}
	}()

	err := mkdirAll(path, DirPerms)
	if err != nil {
		t.Errorf("mkdirAll produced an unexpected error: %v", err)
	}

	if _, err = os.Stat(filepath.Dir(path)); os.IsExist(err) {
		t.Errorf("mkdirAll did not correctly make the directories: %s", path)
	}
}

// Tests MakeDirs by checking if the directory structure exists.
func TestMakeDirs(t *testing.T) {
	path := "temp/temp2/test.txt"

	// Delete the test file at the end
	defer func() {
		err := os.RemoveAll(path)
		if err != nil {
			t.Fatalf("Error deleting test file %q: %+v", path, err)
		}
	}()

	err := MakeDirs(path, DirPerms)
	if err != nil {
		t.Errorf("MakeDirs produced an unexpected error: %v", err)
	}

	if _, err = os.Stat(filepath.Dir(path)); os.IsExist(err) {
		t.Errorf("MakeDirs did not correctly make the directories: %s", path)
	}
}

// Tests that MakeDirs produces an error on an invalid path.
func TestMakeDirs_PathError(t *testing.T) {
	path := "~a/test/test.txt"

	// Delete the test file at the end
	defer func() {
		err := os.RemoveAll(path)
		if err != nil {
			t.Fatalf("Error deleting test file %q: %+v", path, err)
		}
	}()

	err := MakeDirs(path, DirPerms)
	if err == nil {
		t.Errorf("MakeDirs did not produce error when expected:"+
			"\nexpected: %v\nreceived: %v",
			errors.New("cannot expand user-specific home dir"), err)
	}
}

// Tests WriteFile by checking if the directory structure and that the file
// exists.
func TestWriteFile(t *testing.T) {
	dirPath := "temp/temp2"
	path := dirPath + "/test.txt"
	data := []byte("test data")

	// Delete the test file at the end
	defer func() {
		err := os.RemoveAll("temp/")
		if err != nil {
			t.Fatalf("Error deleting test file %q: %+v", "temp/", err)
		}
	}()

	err := os.MkdirAll(dirPath, os.ModePerm)
	if err != nil {
		t.Errorf("MkdirAll produced an unexpected error: %v", err)
	}

	err = WriteFile(path, data, DirPerms, FilePerms)
	if err != nil {
		t.Errorf("WriteFile produced an unexpected error: %v", err)
	}

	if _, err = os.Stat(path); os.IsExist(err) {
		t.Errorf("WriteFile did not correctly make the directories: %s", path)
	}
}

// Tests that WriteFile returns an error with a malformed path.
func TestWriteFile_PathError(t *testing.T) {
	path := "~a/temp/temp2/test.txt"
	data := []byte("test data")

	// Delete the test file at the end
	defer func() {
		err := os.RemoveAll(path)
		if err != nil {
			t.Fatalf("Error deleting test file %q: %+v", path, err)
		}
	}()

	err := WriteFile(path, data, DirPerms, FilePerms)
	if err == nil {
		t.Errorf("WriteFile did not produce error when expected:"+
			"\nexpected: %v\nreceived: %v",
			errors.New("cannot expand user-specific home dir"), err)
	}
}

// Tests WriteFileDef by checking if the directory structure and that the file
// exists.
func TestWriteFileDef(t *testing.T) {
	dirPath := "temp/temp2"
	path := dirPath + "/test.txt"
	data := []byte("test data")

	// Delete the test file at the end
	defer func() {
		err := os.RemoveAll("temp/")
		if err != nil {
			t.Fatalf("Error deleting test file %q: %+v", "temp/", err)
		}
	}()

	err := os.MkdirAll(dirPath, os.ModePerm)
	if err != nil {
		t.Errorf("MkdirAll produced an unexpected error: %v", err)
	}

	err = WriteFileDef(path, data)
	if err != nil {
		t.Errorf("WriteFileDef produced an unexpected error: %v", err)
	}

	if _, err = os.Stat(path); os.IsExist(err) {
		t.Errorf("WriteFileDef did not correctly make the directories: %s", path)
	}
}

// Tests that ReadFile properly reads the contents of a file created by
// WriteFile.
func TestReadFile(t *testing.T) {
	path := "test.txt"
	data := []byte("Test string.")

	// Delete the test file at the end
	defer func() {
		err := os.RemoveAll(path)
		if err != nil {
			t.Fatalf("Error deleting test file %q: %+v", path, err)
		}
	}()

	err := WriteFile(path, data, FilePerms, FilePerms)
	if err != nil {
		t.Errorf("WriteFile produced an unexpected error: %v", err)
	}

	testData, err := ReadFile(path)
	if err != nil {
		t.Errorf("ReadFile produced an unexpected error: %v", err)
	}

	if !bytes.Equal(testData, data) {
		t.Errorf("ReadFile did not return the correct data from the file %s"+
			"\nexpected: %s\nreceived: %s", path, data, testData)
	}
}

// Tests that ReadFile returns an error with a malformed path.
func TestReadFile_PathError(t *testing.T) {
	path := "~a/temp/temp2/test.txt"
	_, err := ReadFile(path)

	if err == nil {
		t.Errorf("ReadFile did not produce error when expected:"+
			"\nexpected: %v\nreceived: %v",
			errors.New("cannot expand user-specific home dir"), err)
	}
}

// Tests that TestExist correctly finds a file that exists.
func TestExist(t *testing.T) {
	path := "test.txt"
	data := []byte("Test string.")

	// Delete the test file at the end
	defer func() {
		err := os.RemoveAll(path)
		if err != nil {
			t.Fatalf("Error deleting test file %q: %+v", path, err)
		}
	}()

	err := WriteFile(path, data, FilePerms, FilePerms)
	if err != nil {
		t.Errorf("WriteFile produced an unexpected error: %v", err)
	}

	exists := Exists(path)
	if !exists {
		t.Errorf("Exists did not find a file that should exist")
	}
}

// Tests that TestExist correctly finds a directory that exists.
func TestExist_Dir(t *testing.T) {
	path := "a/"

	// Delete the test file at the end
	defer func() {
		err := os.RemoveAll(path)
		if err != nil {
			t.Fatalf("Error deleting test file %q: %+v", path, err)
		}
	}()

	err := MakeDirs(path+"d", DirPerms)
	if err != nil {
		t.Errorf("MakeDirs produced an unexpected error: %v", err)
	}

	exists := Exists(path)
	if !exists {
		t.Errorf("Exists did not find a directory that should exist")
	}
}

// Tests that TestExist returns false when a file does not exist.
func TestExist_NoFileError(t *testing.T) {
	path := "test.txt"

	// Delete the test file at the end
	defer func() {
		err := os.RemoveAll(path)
		if err != nil {
			t.Fatalf("Error deleting test file %q: %+v", path, err)
		}
	}()

	exists := Exists(path)
	if exists {
		t.Errorf("Exists found a file when one does not exist")
	}
}

// Tests that FileExists correctly finds a file that exists.
func TestFileExists(t *testing.T) {
	path := "test.txt"
	data := []byte("Test string.")

	// Delete the test file at the end
	defer func() {
		err := os.RemoveAll(path)
		if err != nil {
			t.Fatalf("Error deleting test file %q: %+v", path, err)
		}
	}()

	err := WriteFile(path, data, FilePerms, FilePerms)
	if err != nil {
		t.Errorf("WriteFile produced an unexpected error: %v", err)
	}

	exists := FileExists(path)
	if !exists {
		t.Errorf("FileExists did not find a file that should exist")
	}
}

// Tests that FileExists false when the file is a directory.
func TestFileExists_DirError(t *testing.T) {
	path := "a/d"

	// Delete the test file at the end
	defer func() {
		err := os.RemoveAll(path)
		if err != nil {
			t.Fatalf("Error deleting test file %q: %+v", path, err)
		}
	}()

	err := MakeDirs(path, DirPerms)
	if err != nil {
		t.Errorf("MakeDirs produced an unexpected error: %v", err)
	}

	exists := FileExists(path)
	if exists {
		t.Errorf("FileExists found a directory when it was looking for a file")
	}
}

// Tests that FileExists returns false when a file does not exist.
func TestFileExists_NoFileError(t *testing.T) {
	path := "test.txt"
	// Delete the test file at the end
	defer func() {
		err := os.RemoveAll(path)
		if err != nil {
			t.Fatalf("Error deleting test file %q: %+v", path, err)
		}
	}()

	exists := FileExists(path)
	if exists {
		t.Errorf("FileExists found a file when one does not exist")
	}
}

// Tests that DirExists correctly finds a directory that exists.
func TestDirExists(t *testing.T) {
	path := "a/"

	// Delete the test file at the end
	defer func() {
		err := os.RemoveAll(path)
		if err != nil {
			t.Fatalf("Error deleting test file %q: %+v", path, err)
		}
	}()

	err := MakeDirs(path+"d", DirPerms)
	if err != nil {
		t.Errorf("MakeDirs produced an unexpected error: %v", err)
	}

	exists := DirExists(path)
	if !exists {
		t.Errorf("DirExists did not find a directory that should exist")
	}
}

// Tests that DirExists false when the file is a directory.
func TestDirExists_FileError(t *testing.T) {
	path := "test.txt"
	data := []byte("Test string.")

	// Delete the test file at the end
	defer func() {
		err := os.RemoveAll(path)
		if err != nil {
			t.Fatalf("Error deleting test file %q: %+v", path, err)
		}
	}()

	err := WriteFile(path, data, FilePerms, FilePerms)
	if err != nil {
		t.Errorf("WriteFile produced an unexpected error: %v", err)
	}

	exists := DirExists(path)
	if exists {
		t.Errorf("DirExists found a file when it was looking for a directory")
	}
}

// Tests that DirExists returns false when a file does not exist.
func TestDirExists_NoDirError(t *testing.T) {
	path := "a/b/c/"

	// Delete the test file at the end
	defer func() {
		err := os.RemoveAll(path)
		if err != nil {
			t.Fatalf("Error deleting test file %q: %+v", path, err)
		}
	}()

	exists := FileExists(path)
	if exists {
		t.Errorf("DirExists found a directroy when one does not exist")
	}
}

// Tests that GetLastModified will return an accurate last modified timestamp.
func TestGetLastModified(t *testing.T) {
	path := "test.txt"
	data := []byte("Test string.")

	// Delete the test file at the end
	defer func() {
		require.NoError(t, os.RemoveAll(path))

	}()

	// Record approximately when we are writing to file
	firstWriteTimestamp := time.Now()

	// Write to file
	require.NoError(t, WriteFile(path, data, FilePerms, FilePerms))

	// Retrieve the last modification of the file
	lastModified, err := GetLastModified(path)
	require.NoError(t, err)

	// The last modified timestamp should not differ by more than a few
	// milliseconds from the timestamp taken before the write operation took
	// place.
	require.True(t, lastModified.Sub(firstWriteTimestamp) < 2*time.Millisecond ||
		lastModified.Sub(firstWriteTimestamp) > 2*time.Millisecond)

	// Retrieve modified timestamp again
	newLastModified, err := GetLastModified(path)
	require.NoError(t, err)

	// Ensure last modified does not change arbitrarily
	require.Equal(t, newLastModified, lastModified)
}

// ReadDir unit test.
func TestReadDir(t *testing.T) {
	files, err := ReadDir("./")
	require.NoError(t, err)

	// NOTE: This test uses the files in the utils package as expected values.
	//       If at any point files are added or moved, refactor this list
	//       accordingly.
	var expectedFiles = []string{"gen.go", "net.go", "net_test.go",
		"privNet.go", "file.go", "file_test.go"}
	sort.Strings(expectedFiles)

	require.Equal(t, expectedFiles, files)
}

// Tests that GetLastModified will update after a write operation to a file.
func TestGetLastModified_Update(t *testing.T) {

	path := "test.txt"
	data := []byte("Test string.")

	// Delete the test file at the end
	defer func() {
		require.NoError(t, os.RemoveAll(path))

	}()

	// Write to file
	require.NoError(t, WriteFile(path, data, FilePerms, FilePerms))

	// Retrieve the last modification of the file
	lastModified, err := GetLastModified(path)
	require.NoError(t, err)

	time.Sleep(50 * time.Millisecond)

	// Record timestamp of second write
	secondWriteTimestamp := time.Now()

	// Write again to the same file path
	newData := []byte("New data")
	require.NoError(t, WriteFile(path, newData, FilePerms, FilePerms))

	// Retrieve last modified after re-writing to file
	newLastModified, err := GetLastModified(path)
	require.NoError(t, err)

	// Ensure last modified has been updated, and is not returning an old value
	require.NotEqual(t, newLastModified, lastModified)

	// The last modified timestamp should not differ by more than a few
	// milliseconds from the timestamp taken before the write operation took
	// place.
	require.True(t, lastModified.Sub(secondWriteTimestamp) < 2*time.Millisecond ||
		lastModified.Sub(secondWriteTimestamp) > 2*time.Millisecond)

}

// Tests that Test_exist correctly finds a file that exists and returns the
// correct FileInfo.
func Test_exist(t *testing.T) {
	path := "test.txt"
	data := []byte("Test string.")

	// Delete the test file at the end
	defer func() {
		err := os.RemoveAll(path)
		if err != nil {
			t.Fatalf("Error deleting test file %q: %+v", path, err)
		}
	}()

	err := WriteFile(path, data, FilePerms, FilePerms)
	if err != nil {
		t.Errorf("WriteFile produced an unexpected error: %v", err)
	}

	info, exists := exists(path)
	expectedInfo, err := os.Stat(path)

	if !exists && err != nil {
		t.Errorf("exists did not find a file that should exist: %v", err)
	} else if !exists {
		t.Errorf("exists did not find a file that should exist")
	}

	if !reflect.DeepEqual(info, expectedInfo) {
		t.Errorf("exists did not return the expected FileInfo."+
			"\nexpected: %v\nreceived: %v", expectedInfo, info)
	}
}

// Tests that Test_exist returns false when a file does not exist. and returns
// a nil FileInfo.
func Test_exist_NoFileError(t *testing.T) {
	path := "test.txt"

	info, exists := exists(path)

	if exists {
		t.Errorf("exists found a file when one does not exist")
	}

	if info != nil {
		t.Errorf("exists unexpectedly returned a non-nil FileInfo."+
			"\nexpected: %v\nreceived: %v", nil, info)
	}
}

// Tests that SearchDefaultLocations finds the specified file in the user's
// home directory
func TestSearchDefaultLocations(t *testing.T) {
	testDir := fmt.Sprintf("testDir-%d/", time.Now().Nanosecond())
	testFile := fmt.Sprintf("testFile-%d.txt", time.Now().Nanosecond())
	testPath := testDir + testFile
	expectedPath, err := ExpandPath("~/." + testPath)
	if err != nil {
		t.Fatalf("ExpandPath failed to exapnd the path %s: %+v", testPath, err)
	}
	expectedDir, err := ExpandPath("~/" + testDir)
	if err != nil {
		t.Fatalf("ExpandPath failed to exapnd the path %s: %+v", testPath, err)
	}

	// Delete the test file at the end
	defer func() {
		err := os.RemoveAll(expectedDir)
		if err != nil {
			t.Fatalf("Error deleting test file %q: %+v", expectedDir, err)
		}
	}()

	err = WriteFile(expectedPath, []byte("TEST"), FilePerms, DirPerms)
	if err != nil {
		t.Fatalf("WriteFile failed to create file %s: %+v", testPath, err)
	}

	foundPath, err := SearchDefaultLocations(testFile, testDir)
	if err != nil {
		t.Errorf("SearchDefaultLocations produced an unexpected error: %+v",
			err)
	}

	if foundPath != expectedPath {
		t.Errorf("SearchDefaultLocations did not find the correct file."+
			"\nexpected: %s\nreceived: %s", expectedPath, foundPath)
	}
}

// Tests that SearchDefaultLocations return an error when the file does not
// exist.
func TestSearchDefaultLocations_NotFoundError(t *testing.T) {
	testDir := fmt.Sprintf(".testDir-%d/", time.Now().Nanosecond())
	testFile := fmt.Sprintf("testFile-%d.txt", time.Now().Nanosecond())

	foundPath, err := SearchDefaultLocations(testFile, testDir)
	if err == nil {
		t.Errorf("SearchDefaultLocations did not error when expected.")
	}

	if foundPath != "" {
		t.Errorf("SearchDefaultLocations did not return an empty path on error."+
			"\nexpected: %s\nreceived: %s", "", foundPath)
	}
}

// Tests that getDefaultSearchDirs generates the correct list of default paths.
func TestGetDefaultSearchDirs(t *testing.T) {
	testDir := "xxnetwork"
	expectedDir0, err := ExpandPath("~/." + testDir + "/")
	expectedDir1, err := ExpandPath("/opt/" + testDir + "/")
	expectedDir2, err := ExpandPath("/etc/" + testDir + "/")

	testDirs, err := getDefaultSearchDirs(testDir)
	if err != nil {
		t.Errorf("getDefaultSearchDirs produced an unxpected error: %+v", err)
	}

	if testDirs[0] != expectedDir0 {
		t.Errorf("getDefaultSearchDirs did not return the correct path for "+
			"home.\nexpected: %s\nreceived: %s", expectedDir0, testDirs[0])
	}

	if testDirs[1] != expectedDir1 {
		t.Errorf("getDefaultSearchDirs did not return the correct path for "+
			"/etc/.\nexpected: %s\nreceived: %s", expectedDir1, testDirs[1])
	}

	if testDirs[2] != expectedDir2 {
		t.Errorf("getDefaultSearchDirs did not return the correct path for "+
			"/etc/.\nexpected: %s\nreceived: %s", expectedDir2, testDirs[2])
	}
}
