package ephemeral

import (
	"bytes"
	"encoding/binary"
	"gitlab.com/xx_network/crypto/csprng"
	"gitlab.com/xx_network/primitives/id"
	_ "golang.org/x/crypto/blake2b"
	"testing"
	"time"
)

// Unit test for GetId
func TestGetId(t *testing.T) {
	testId := id.NewIdFromString("zezima", id.User, t)
	eid, err := GetId(testId, 99, uint64(time.Now().UnixNano()))
	if err == nil {
		t.Error("Should error with size > 64")
	}
	eid, err = GetId(testId, 16, uint64(time.Now().Unix()))
	if err != nil {
		t.Errorf("Failed to create ephemeral ID: %+v", err)
	}
	t.Log(eid)
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
	eid, err := GetIdFromIntermediary(iid, 16, uint64(time.Now().UnixNano()))
	if err != nil {
		t.Errorf("Failed to get id from intermediary: %+v", err)
	}
	if eid[2] != 0 && eid[3] != 0 && eid[4] != 0 && eid[5] != 0 && eid[6] != 0 && eid[7] != 0 {
		t.Errorf("Id was not cleared to proper size: %+v", eid)
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
	newEid, err := eid.Fill(uint(64), csprng.NewSystemRNG())
	if err != nil {
		t.Errorf("Failed to fill ID: %+v", err)
	}
	for i, r := range newEid[:] {
		if r != eid[i] {
			t.Errorf("Fill changed bits in max size ID.  Original: %+v, New: %+v", eid, newEid)
		}
	}

	eid = eid.Clear(16)
	newEid, err = eid.Fill(16, csprng.NewSystemRNG())
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
	salt1 := getRotationSalt(idHash, uint64(ts))
	ts += (12 * time.Hour).Nanoseconds()
	salt2 := getRotationSalt(idHash, uint64(ts))
	ts += (12 * time.Hour).Nanoseconds()
	salt3 := getRotationSalt(idHash, uint64(ts))
	if bytes.Compare(salt1, salt2) == 0 && bytes.Compare(salt2, salt3) == 0 {
		t.Error("Salt did not change as timestamp increased w/ period of one day")
	}
	t.Logf("First: %+v\tSecond: %+v\nThird: %+v\n", salt1, salt2, salt3)
}

// Unit test for UInt64 method on ephemeral ID
func TestId_UInt64(t *testing.T) {
	testId := id.NewIdFromString("zezima", id.User, t)
	eid, err := GetId(testId, 16, uint64(time.Now().Unix()))
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
