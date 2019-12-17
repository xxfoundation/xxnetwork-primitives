package switchboard

import (
	"gitlab.com/elixxir/primitives/id"
	"reflect"
	"strconv"
	"testing"
)




func TestListenersMap_StoreListener(t *testing.T) {
	sb, mockListener:= OneListenerSetup()
	newListenerRecord := &listenerRecord{
		l:  mockListener,
		id: strconv.Itoa(sb.lastID),
	}

	sb.listenersMap.StoreListener(id.ZeroID, int32(5), newListenerRecord)

	mapKey := listenerMapKey{id.ZeroID, int32(5)}
	records_i, ok := sb.listenersMap.listeners.Load(mapKey)
	records := records_i.([]*listenerRecord)
	if !ok{
		t.Logf("Store listener failed, could not find stored object")
		t.Fail()
	}

	testRecords := []*listenerRecord{newListenerRecord}
	if !reflect.DeepEqual(records, testRecords){
		t.Logf("Record found does not match record initially stored")
		t.Fail()
	}
}



func TestListenersMap_GetListenerRecords(t *testing.T) {
	sb, mockListener:= OneListenerSetup()
	newListenerRecord := &listenerRecord{
		l:  mockListener,
		id: strconv.Itoa(sb.lastID),
	}

	sb.listenersMap.StoreListener(id.ZeroID, int32(5), newListenerRecord)

	records, ok := sb.listenersMap.GetListenerRecords(id.ZeroID, int32(5))
	if !ok{
		t.Logf("Get Listener Records failed, could not find stored object")
		t.Fail()
	}

	testRecords := []*listenerRecord{newListenerRecord}
	if !reflect.DeepEqual(records, testRecords){
		t.Logf("Get Listener Records returned wrong object")
		t.Fail()
	}
}

func TestListenersMap_GetMatches(t *testing.T) {
	sb, mockListener:= OneListenerSetup()
	listenerId := strconv.Itoa(sb.lastID)
	newListenerRecord := &listenerRecord{
		l:  mockListener,
		id: listenerId,
	}

	sb.listenersMap.StoreListener(id.ZeroID, int32(5), newListenerRecord)
	matches := []*listenerRecord{}

	//Test that matches are appended
	matchesTestA := sb.listenersMap.GetMatches(matches, id.ZeroID, int32(5))
	if reflect.DeepEqual(matches, matchesTestA){
		t.Logf("Matches were not found when they should have been")
		t.Fail()
	}

	//Test that if you pass in the same variables duplicates don't occur
	matchesTestB := sb.listenersMap.GetMatches(matchesTestA, id.ZeroID, int32(5))
	if !reflect.DeepEqual(matchesTestB, matchesTestA){
		t.Logf("Duplicates found in matches")
		t.Fail()
	}
	
}


func TestListenersMap_RemoveListener(t *testing.T) {
	sb, mockListener:= OneListenerSetup()
	listenerId := strconv.Itoa(sb.lastID)
	newListenerRecord := &listenerRecord{
		l:  mockListener,
		id: listenerId,
	}

	sb.listenersMap.StoreListener(id.ZeroID, int32(5), newListenerRecord)

	sb.listenersMap.RemoveListener(listenerId)

	records, ok := sb.listenersMap.GetListenerRecords(id.ZeroID, int32(5))
	if ok{
		testRecords := []*listenerRecord{newListenerRecord}
		if reflect.DeepEqual(records, testRecords){
			t.Logf("Failed to remove listener from ListenersMap")
			t.Fail()
		}

	}
}