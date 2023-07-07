////////////////////////////////////////////////////////////////////////////////
// Copyright © 2022 xx foundation                                             //
//                                                                            //
// Use of this source code is governed by a license that can be found in the  //
// LICENSE file.                                                              //
////////////////////////////////////////////////////////////////////////////////

package ephemeral

import (
	"bytes"
	"crypto"
	"encoding/binary"
	"encoding/json"
	"math"
	"strconv"
	"testing"
	"time"

	_ "golang.org/x/crypto/blake2b"

	"gitlab.com/xx_network/crypto/csprng"
	"gitlab.com/xx_network/primitives/id"
)

// Unit test for GetId
func TestGetId(t *testing.T) {
	testId := id.NewIdFromString("zezima", id.User, t)
	eid, _, _, err := GetId(testId, 99, time.Now().UnixNano())
	if err == nil {
		t.Error("Should error with size > 64")
	}
	eid, _, _, err = GetId(testId, 16, time.Now().Unix())
	if err != nil {
		t.Errorf("Failed to create ephemeral ID: %+v", err)
	}
	t.Log(eid)
}

// Unit test for GetIdsByRange
func TestGetIdByRange(t *testing.T) {
	testId := id.NewIdFromString("zezima", id.User, t)
	eids, err := GetIdsByRange(testId, 99, time.Now(), 25)
	if err == nil {
		t.Error("Should error with size > 64")
	}
	duration := 7 * 24 * time.Hour

	eids, err = GetIdsByRange(testId, 16, time.Now(), duration)
	if err != nil {
		t.Errorf("Failed to create ephemeral ID: %+v", err)
	}
	expectedLength := int(int64(duration)/Period) + 1

	if len(eids) != expectedLength {
		t.Errorf("Unexpected list of ephemeral IDs."+
			"\nexpected: %d\nreceived: %d", expectedLength, len(eids))
	}

	// Test that the time variances are correct
	for i := 0; i < len(eids)-1; i++ {
		next := i + 1
		if eids[i].End != eids[next].Start {
			t.Errorf("The next identity after %d does not start "+
				"when the current identity ends:\nend: %s\nstart: %s",
				i, eids[i].End, eids[next].Start)
		}
		if int64(eids[i].End.Sub(eids[i].Start)) != Period {
			t.Errorf("Delta between start and end on %d does not equal the "+
				"Period:\nend: %s\nstart: %s",
				i, eids[i].End, eids[next].Start)
		}
	}
}

func TestGetIntermediaryId(t *testing.T) {
	testId := id.NewIdFromString("zezima", id.User, t)
	iid, err := GetIntermediaryId(testId)
	if err != nil {
		t.Errorf("Failed to get intermediary ID: %+v", err)
	}
	if iid == nil || len(iid) == 0 {
		t.Errorf("iid returned with no data: %v", iid)
	}
}

func TestGetIdFromIntermediary(t *testing.T) {
	testId := id.NewIdFromString("zezima", id.User, t)
	iid, err := GetIntermediaryId(testId)
	if err != nil {
		t.Errorf("Failed to get intermediary ID: %+v", err)
	}
	eid, _, _, err := GetIdFromIntermediary(iid, 16, time.Now().UnixNano())
	if err != nil {
		t.Errorf("Failed to get ID from intermediary: %+v", err)
	}
	if !bytes.Equal(eid[:6], make([]byte, 6)) {
		t.Errorf("Id was not cleared to proper size: %v", eid)
	}
}

// Check that given precomputed input that should generate a reserved
// ephemeral ID, GetIdFromIntermediary does not generate a reserved Id
func TestGetIdFromIntermediary_Reserved(t *testing.T) {
	// Hardcoded to ensure a collision with a reserved ID
	hardcodedTimestamp := int64(1614199942358373731)
	size := uint(4)
	testId := id.NewIdFromString(strconv.Itoa(41), id.User, t)

	// Intermediary ID expected to generate a reserved ephemeral ID
	iid, err := GetIntermediaryId(testId)
	if err != nil {
		t.Errorf("Failed to get intermediary ID: %+v", err)
	}
	// Generate an ephemeral Id given the input above. This specific
	// call does not check if the outputted Id is reserved
	salt, _, _ := getRotationSalt(iid, hardcodedTimestamp)
	b2b := crypto.BLAKE2b_256.New()
	expectedReservedEID, err := getIdFromIntermediary(b2b, iid, salt, size)
	if err != nil {
		t.Errorf("Failed to get ID from intermediary: %+v", err)
	}

	// Check that the ephemeral Id generated with hardcoded data is a reserved ID
	if !IsReserved(expectedReservedEID) {
		t.Errorf("Expected reserved eid is no longer reserved; may need to " +
			"find a new ID. Use FindReservedID in this case.")
	}

	// Generate an ephemeral ID which given the same input above with the
	// production facing call
	eid, _, _, err := GetIdFromIntermediary(iid, size, hardcodedTimestamp)
	if err != nil {
		t.Errorf("Failed to get id from intermediary: %+v", err)
	}

	// Check that the ephemeralID generated is not reserved.
	if IsReserved(eid) {
		t.Errorf("Ephemeral ID generated should not be reserved!"+
			"\nReserved IDs: %v"+
			"\nGenerated ID: %v", ReservedIDs, eid)
	}

}

// Will find a reserved ephemeral ID and returns the associated intermediary ID.
func FindReservedID(size uint, timestamp int64, t *testing.T) []byte {
	b2b := crypto.BLAKE2b_256.New()

	// Loops through until a reserved ID is found
	counter := 0
	for {
		testId := id.NewIdFromString(strconv.Itoa(counter), id.User, t)
		iid, err := GetIntermediaryId(testId)
		if err != nil {
			t.Errorf("Failed to get intermediary id: %+v", err)
		}

		// Generate an ephemeral ID
		salt, _, _ := getRotationSalt(iid, timestamp)
		eid, err := getIdFromIntermediary(b2b, iid, salt, size)
		if err != nil {
			t.Errorf("Failed to get id from intermediary: %+v", err)
		}

		// Check if ephemeral ID is reserved exit
		if IsReserved(eid) {
			t.Logf("Found input which generates a reserved id. Input as follows."+
				"\nSize: %d"+
				"\nTimestamp: %d"+
				"\nTestID: %v"+
				"\nTestID generated using the following line of code: "+
				"\ntestId := id.NewIdFromString(strconv.Itoa(%d), id.User, t)",
				size, timestamp, testId, counter)
			return iid
		}
		// Increment the counter
		counter++
	}
}

func TestId_Clear(t *testing.T) {
	eid := Id{}
	dummyData := []byte{201, 99, 103, 45, 68, 2, 56, 7}
	copy(eid[:], dummyData)

	newEid := eid.Clear(uint(64))
	var ok bool
	if bytes.Map(func(r rune) rune { ok = ok || r == 0; return r }, eid[:]); ok {
		t.Errorf("Bytes were cleared from max size ID: %v", newEid)
	}

	newEid = eid.Clear(16)
	if !bytes.Equal(newEid[:6], make([]byte, 6)) {
		t.Errorf("Proper bits were not cleared from size 16 ID: %v", newEid)
	}
	if bytes.Equal(eid[:6], make([]byte, 6)) {
		t.Errorf("Bits were cleared from original ID: %v", eid)
	}
	if !bytes.Equal(eid[6:8], newEid[6:8]) {
		t.Errorf("Proper bits do not match in IDs.\noriginal: %v\ncleared:  %v",
			eid, newEid)
	}
}

func TestId_Fill(t *testing.T) {
	var eid Id
	dummyData := []byte{201, 99, 103, 45, 68, 2, 56, 7}
	copy(eid[:], dummyData)

	eid = eid.Clear(uint(64))
	newEid, err := eid.Fill(uint(64), csprng.NewSystemRNG())
	if err != nil {
		t.Errorf("Failed to fill ID: %+v", err)
	}
	for i, r := range newEid[:] {
		if r != eid[i] {
			t.Errorf("Fill changed bits in max size ID (%d)."+
				"\noriginal: %v\nnew:      %v", i, eid, newEid)
		}
	}

	eid = eid.Clear(16)
	newEid, err = eid.Fill(16, csprng.NewSystemRNG())
	if err != nil {
		t.Errorf("Failed to fill ID: %+v", err)
	}
	if bytes.Equal(newEid[:6], eid[:6]) {
		t.Errorf("Proper bits were not filled from size 16 ID %v", newEid)
	}
	if !bytes.Equal(newEid[6:8], eid[6:8]) {
		t.Errorf("Proper bits do not match in IDs.\noriginal: %v\ncleared:  %v",
			eid, newEid)
	}
}

func TestGetRotationSalt(t *testing.T) {
	ts := time.Now().UnixNano()
	idHash, err := GetIntermediaryId(id.NewIdFromString("zezima", id.User, t))
	if err != nil {
		t.Errorf("Failed to get intermediary id hash: %+v", err)
	}
	salt1, _, _ := getRotationSalt(idHash, ts)
	ts += (12 * time.Hour).Nanoseconds()
	salt2, _, _ := getRotationSalt(idHash, ts)
	ts += (12 * time.Hour).Nanoseconds()
	salt3, _, _ := getRotationSalt(idHash, ts)
	if bytes.Compare(salt1, salt2) == 0 && bytes.Compare(salt2, salt3) == 0 {
		t.Error("Salt did not change as timestamp increased w/ Period of one day")
	}
	t.Logf("First:  %v\tSecond: %v\nThird:  %v\n", salt1, salt2, salt3)
}

// Unit test for UInt64 method on ephemeral ID
func TestId_UInt64(t *testing.T) {
	testId := id.NewIdFromString("zezima", id.User, t)
	eid, _, _, err := GetId(testId, 16, time.Now().Unix())
	if err != nil {
		t.Errorf("Failed to create ephemeral ID: %+v", err)
	}
	ueid := eid.UInt64()
	var b [8]byte
	binary.BigEndian.PutUint64(b[:], ueid)
	if bytes.Compare(b[:], eid[:]) != 0 {
		t.Error("UInt64 conversion is wrong")
	}
}

// Test the int64 conversion from ephemeral ID
func TestId_Int64(t *testing.T) {
	testId := id.NewIdFromString("zezima", id.User, t)
	eid, _, _, err := GetId(testId, 16, time.Now().Unix())
	if err != nil {
		t.Errorf("Failed to create ephemeral ID: %+v", err)
	}
	eid, err = eid.Fill(16, csprng.NewSystemRNG())
	if err != nil {
		t.Errorf("Failed to fill ephemeral ID: %+v", err)
	}
	var maxUint64Id Id
	binary.BigEndian.PutUint64(maxUint64Id[:], math.MaxUint64)
	if maxUint64Id.Int64() != math.MinInt64 {
		t.Error("Did not properly convert from uint to int")
		t.Error(maxUint64Id.Int64())
	}

	zerouint64Id := Id{}
	binary.BigEndian.PutUint64(zerouint64Id[:], 0)
	if zerouint64Id.Int64() != 0 {
		t.Error("Did not properly convert a zero id to id and back")
		t.Error(zerouint64Id.Int64())
	}
}

// Unit test for ephemeral ID load function
func TestMarshal(t *testing.T) {
	testId := id.NewIdFromString("zezima", id.User, t)
	eid, _, _, err := GetId(testId, 16, time.Now().Unix())
	if err != nil {
		t.Errorf("Failed to create ephemeral ID: %+v", err)
	}
	eid, err = eid.Fill(16, csprng.NewSystemRNG())
	if err != nil {
		t.Errorf("Failed to fill ephemeral ID: %+v", err)
	}
	eid2, err := Marshal(eid[:])
	if err != nil {
		t.Errorf("Failed to marshal id from bytes")
	}
	if bytes.Compare(eid[:], eid2[:]) != 0 {
		t.Errorf("Failed to load ephermeral ID from bytes."+
			"\noriginal: %v\nloaded:   %v", eid, eid2)
	}

	_, err = Marshal(nil)
	if err == nil {
		t.Error("nil data should return an error when marshaled")
	}

	_, err = Marshal([]byte("Test"))
	if err == nil {
		t.Error("Data < size 8 should return an error when marshalled")
	}
}

// Tests that an Id can be JSON marshalled and unmarshalled.
func TestId_JSONMarshalUnmarshal(t *testing.T) {
	testID := id.NewIdFromString("zezima", id.User, t)
	expected, _, _, err := GetId(testID, 16, time.Now().Unix())

	data, err := json.Marshal(expected)
	if err != nil {
		t.Errorf("Failed to JSON marshal %T: %+v", expected, err)
	}

	var eid Id
	err = json.Unmarshal(data, &eid)
	if err != nil {
		t.Errorf("Failed to JSON umarshal %T: %+v", eid, err)
	}

	if expected != eid {
		t.Errorf("Marshalled and unamrshalled Id does not match expected."+
			"\nexpected: %s\nreceived: %s", expected, eid)
	}
}
