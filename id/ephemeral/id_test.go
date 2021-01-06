package ephemeral

import (
	"gitlab.com/xx_network/primitives/id"
	_ "golang.org/x/crypto/blake2b"
	"testing"
)

// Unit test for GetId
func TestGetId(t *testing.T) {
	testId := id.NewIdFromString("zezima", id.User, t)
	eid, err := GetId(testId, 99)
	if err == nil {
		t.Error("Should error with size > 64")
	}
	eid, err = GetId(testId, 28)
	if err != nil {
		t.Errorf("Failed to create ephemeral ID: %+v", err)
	}
	t.Log(eid)
}
