package eventsourcing
import (
	"testing"
)

//Example event used for testing
type MyTestEvent struct {
	BaseEvent
	value int
}


func TestSaveEvent(t *testing.T)  {
	var storeUT = NewStore()
	var e1 = MyTestEvent{value:1}

	storeUT.Save([]Event{&e1})

	if (len(storeUT.events) != 1) {
		t.Errorf("Expected 1 event but got %s", len(storeUT.events))
	}
	if (len(storeUT.store) != 1) {
		t.Errorf("Expected 1 enity in store but got %s", len(storeUT.store))
	}
}

func TestFindEvent(t *testing.T) {
	var storeUT = NewStore()
	storeUT.Save([]Event{
		&MyTestEvent{value:1},
		&MyTestEvent{value:2},
	})
	var guid = storeUT.Save([]Event{
		&MyTestEvent{value:3},
		&MyTestEvent{value:4},
	})

	var events, version = storeUT.Find(guid)
	if version != 2 {
		t.Errorf("Expected version to be %v but got %v", 2, version)
		t.Fail()
	}

	if len(events) != 2 {
		t.Errorf("Expected to get %v event, but got %v", 2, len(events))
		t.Fail()
	}

	if events[0].(*MyTestEvent).value != 3 {
		t.Errorf("Expected different event, but got %v", events[0])
		t.Fail()
	}
	if events[1].(*MyTestEvent).value != 4 {
		t.Errorf("Expected different event, but got %v", events[1])
		t.Fail()
	}
}