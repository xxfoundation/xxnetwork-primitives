package switchboard

import (
	"gitlab.com/elixxir/primitives/id"
	"sync"
	jww "github.com/spf13/jwalterweatherman"
)

//We use this structure to map these value to a listenerRecord array in the sync.Map
type listenerMapKey struct {
	userId      *id.User
	messageType int32
}

type listenersMap struct{
	// Key = listenerMapKey , value = []*listenerRecords
	listeners sync.Map
	// Key = listenerId String, Value listenerMapKey
	// ListenerIds maps listener ids to listenerMapKeys making the reverse search in unregister o(k).
	listenerKeys sync.Map
}

// Given a user id and message type it will return the slice stored for this user
// and a boolean identifying that the lookup was successful
func (lm *listenersMap) GetListenerRecords(user *id.User, messageType int32) ([]*listenerRecord, bool){
	mapKey := listenerMapKey{user, messageType}
	//sync map returns an interface, so give it a type then return that slice
	listenerRecords_i, ok :=  lm.listeners.Load(mapKey)
	if ok {
		return listenerRecords_i.([]*listenerRecord), ok
	}
	// No record exists so ok will by default be false and we return nil for records
	return nil, ok
}

// Stores a new listener to a sync map with the key of {id.User, MessageType}
func (lm *listenersMap) StoreListener(user *id.User, messageType int32, newListener *listenerRecord){
	mapKey := listenerMapKey{user, messageType}
	newListenerRecordSlice := []*listenerRecord{}

	listenerRecords, ok := lm.GetListenerRecords(user, messageType)
	if ok {
		// Append the new Listener to existing records
		newListenerRecordSlice = append(listenerRecords, newListener)
	} else {
		// Append the new listener to a new empty array, this is the first entry for this mapKey
		newListenerRecordSlice = append(newListenerRecordSlice, newListener)
	}

	// We store here to be able to do an efficient reverse look up of where a listener id is located
	lm.listenerKeys.Store(newListener.id, mapKey)
	// Here we store the actual listener slice appended with the new listener
	lm.listeners.Store(mapKey, newListenerRecordSlice)
	lm.String()

	return
}

// Removes a listener using its listener ID, This method uses a map of listenerIds
// to listenerMapKey objects so we know where the listener object is making the search o(k).
func (lm *listenersMap) RemoveListener(listenerID string) bool{
	// Load the map of the stored listenerId, returns an interface
	unregisterId_i, ok := lm.listenerKeys.Load(listenerID)
	if ok {
		// We need to type the sync.Map output since it returns an interface.
		unregisterMapId := unregisterId_i.(listenerMapKey)
		// Now that we know the location of this listener, we return the relevant slices.
		listeners, ok := lm.GetListenerRecords(unregisterMapId.userId, unregisterMapId.messageType)
		if ok {
			for i := range listeners {
				if listenerID == listeners[i].id {
					//In deleting here is it important to maintain order? quicker solution if not
					newListeners := deleteSliceElem(i, listeners)
					lm.listenerKeys.Delete(listenerID)
					lm.listeners.Store(unregisterMapId, newListeners)
					return true
				}
			}
		} else {
			// Could not be found therefore doesnt exist
			return false
		}
	}
	//ListenerId could not be found
	return false
}

func (lm *listenersMap) String(){
	lm.listeners.Range(func(key interface{}, value interface{}) bool {
		mapKey := key.(listenerMapKey)
		listeners := value.([]*listenerRecord)
		for i := range listeners {
			jww.ERROR.Printf("Listener %v: %v, user %v, "+
				" type %v, ", i, listeners[i].id, mapKey.userId, mapKey.messageType)
			return true
		}
		return false
	})
}

//loops through the listener getting all matches and returning the appended matches object
func (lm *listenersMap) GetMatches(matches []*listenerRecord, user *id.User, messageType int32) ([]*listenerRecord) {
	listeners, ok := lm.GetListenerRecords(user, messageType)
	if ok {
		for _, listener := range listeners {
			matches = appendIfUnique(matches, listener)
		}
	}
	return matches
}

func appendIfUnique(matches []*listenerRecord, newListener *listenerRecord) []*listenerRecord {
	// Search for the listener ID
	found := false
	for _, l := range matches {
		found = found || (l.id == newListener.id)
	}
	if !found {
		// If we didn't find it, it's OK to append it to the slice
		return append(matches, newListener)
	} else {
		// We already matched this listener, and shouldn't append it
		return matches
	}
}

// Delete a listenerRecord from the slice at a specific location in a slice maintaining its order.
func deleteSliceElem(loc int, records []*listenerRecord) []*listenerRecord {
	copy(records[loc:], records[loc+1:])        // Shift a[i+1:] left one index.
	records[len(records)-1] = &listenerRecord{} // Erase last element (write zero value).
	records = records[:len(records)-1]          // Truncate slice.

	return records
}
