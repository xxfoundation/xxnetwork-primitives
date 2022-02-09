package ephemeral

import (
	"bytes"
	"crypto"
	"encoding/binary"
	"gitlab.com/xx_network/primitives/id"
	_ "golang.org/x/crypto/blake2b"
	"math"
	"math/rand"
	"strconv"
	"testing"
	"time"
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
			"\n\tExpected: %d"+
			"\n\tReceived: %d", expectedLength, len(eids))
	}

	// test that the time variances are correct
	for i := 0; i < len(eids)-1; i++ {
		next := i + 1
		if eids[i].End != eids[next].Start {
			t.Errorf("The next identity after %d does not start "+
				"when the current identity ends: \n\t end: %s \n\t start: %s",
				i, eids[i].End, eids[next].Start)
		}
		if int64(eids[i].End.Sub(eids[i].Start)) != Period {
			t.Errorf("Delta between start and end on %d does not equal the "+
				"Period: \n\t end: %s \n\t start: %s",
				i, eids[i].End, eids[next].Start)
		}
	}
}

func TestGetIntermediaryId(t *testing.T) {
	testId := id.NewIdFromString("zezima", id.User, t)
	iid, err := GetIntermediaryId(testId)
	if err != nil {
		t.Errorf("Failed to get intermediary id: %+v", err)
	}
	if iid == nil || len(iid) == 0 {
		t.Errorf("iid returned with no data: %+v", iid)
	}
}

func TestGetIdFromIntermediary(t *testing.T) {
	testId := id.NewIdFromString("zezima", id.User, t)
	iid, err := GetIntermediaryId(testId)
	if err != nil {
		t.Errorf("Failed to get intermediary id: %+v", err)
	}
	eid, _, _, err := GetIdFromIntermediary(iid, 16, time.Now().UnixNano())
	if err != nil {
		t.Errorf("Failed to get id from intermediary: %+v", err)
	}
	if eid[2] != 0 && eid[3] != 0 && eid[4] != 0 && eid[5] != 0 && eid[6] != 0 && eid[7] != 0 {
		t.Errorf("Id was not cleared to proper size: %+v", eid)
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
		t.Errorf("Failed to get intermediary id: %+v", err)
	}
	// Generate an ephemeral Id given the input above. This specific
	// call does not check if the outputted Id is reserved
	salt, _, _ := getRotationSalt(iid, hardcodedTimestamp)
	b2b := crypto.BLAKE2b_256.New()
	expectedReservedEID, err := getIdFromIntermediaryHelper(b2b, iid, salt, size)
	if err != nil {
		t.Errorf("Failed to get id from intermediary: %+v", err)
	}

	// Check that the ephemeral Id generated with hardcoded data is a reserved ID
	if !IsReserved(expectedReservedEID) {
		t.Errorf("Expected reserved eid is no longer reserved, " +
			"\n\tmay need to find a new ID. Use FindReservedID in this case.")
	}

	// Generate an ephemeral ID which given the same input above with the production facing call
	eid, _, _, err := GetIdFromIntermediary(iid, size, hardcodedTimestamp)
	if err != nil {
		t.Errorf("Failed to get id from intermediary: %+v", err)
	}

	// Check that the ephemeralID generated is not reserved.
	if IsReserved(eid) {
		t.Errorf("Ephemeral ID generated should not be reserved!"+
			"\n\tReserved IDs: %v"+
			"\n\tGenerated ID: %v", ReservedIDs, eid)
	}

}

// Will find a reserved ephemeral ID and returns the
// associated intermediary ID
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
		eid, err := getIdFromIntermediaryHelper(b2b, iid, salt, size)
		if err != nil {
			t.Errorf("Failed to get id from intermediary: %+v", err)
		}

		// Check if ephemeral ID is reserved exit
		if IsReserved(eid) {
			t.Logf("Found input which generates a reserved id. Input as follows."+
				"\n\tSize: %d"+
				"\n\tTimestamp: %d"+
				"\n\tTestID: %v"+
				"\n\tTestID generated using the following line of code: "+
				"\n\t\ttestId := id.NewIdFromString(strconv.Itoa(%d), id.User, t)",
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
		t.Errorf("Bytes were cleared from max size id: %+v", newEid)
	}

	newEid = eid.Clear(16)
	if newEid[0] != 0 || newEid[1] != 0 || newEid[2] != 0 || newEid[3] != 0 || newEid[4] != 0 || newEid[5] != 0 {
		t.Errorf("Proper bits were not cleared from size 16 id: %+v", newEid)
	}
	if eid[0] == 0 && eid[1] == 0 && eid[2] == 0 && eid[3] == 0 && eid[4] == 0 && eid[5] == 0 {
		t.Errorf("Bits were cleared from original id: %+v", eid)
	}
	if newEid[6] != eid[6] && newEid[7] != eid[7] {
		t.Errorf("Proper bits do not match in ids.  Original: %+v, cleared: %+v", eid, newEid)
	}
}

func TestId_Fill(t *testing.T) {
	eid := Id{}
	dummyData := []byte{201, 99, 103, 45, 68, 2, 56, 7}
	copy(eid[:], dummyData)

	eid = eid.Clear(uint(64))
	prng := rand.New(rand.NewSource(42))
	newEid, err := eid.Fill(uint(64), rand.New(prng))
	if err != nil {
		t.Errorf("Failed to fill ID: %+v", err)
	}
	for i, r := range newEid[:] {
		if r != eid[i] {
			t.Errorf("Fill changed bits in max size ID.  Original: %+v, New: %+v", eid, newEid)
		}
	}

	eid = eid.Clear(16)
	newEid, err = eid.Fill(16, prng)
	if err != nil {
		t.Errorf("Failed to fill ID: %+v", err)
	}
	if newEid[0] == eid[0] || newEid[1] == eid[1] || newEid[2] == eid[2] || newEid[3] == eid[3] ||
		newEid[4] == eid[4] || newEid[5] == eid[5] {
		t.Errorf("Proper bits were not filled from size 16 id: %+v", newEid)
	}
	if newEid[6] != eid[6] && newEid[7] != eid[7] {
		t.Errorf("Proper bits do not match in ids.  Original: %+v, cleared: %+v", eid, newEid)
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
	t.Logf("First: %+v\tSecond: %+v\nThird: %+v\n", salt1, salt2, salt3)
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
	eid, err = eid.Fill(16, rand.New(rand.NewSource(42)))
	if err != nil {
		t.Errorf("Failed to fill ephemeral ID: %+v", err)
	}
	maxuint64Id := Id{}
	binary.BigEndian.PutUint64(maxuint64Id[:], math.MaxUint64)
	if maxuint64Id.Int64() != math.MinInt64 {
		t.Error("Did not properly convert from uint to int")
		t.Error(maxuint64Id.Int64())
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
	eid, err = eid.Fill(16, rand.New(rand.NewSource(42)))
	if err != nil {
		t.Errorf("Failed to fill ephemeral ID: %+v", err)
	}
	eid2, err := Marshal(eid[:])
	if err != nil {
		t.Errorf("Failed to marshal id from bytes")
	}
	if bytes.Compare(eid[:], eid2[:]) != 0 {
		t.Errorf("Failed to load ephermeral ID from bytes.  Original: %+v, Loaded: %+v", eid, eid2)
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
