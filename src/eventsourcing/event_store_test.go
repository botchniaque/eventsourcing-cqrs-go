package eventsourcing
import (
	"testing"
	"github.com/stretchr/testify/assert"
)

//Example event used for testing
type MyTestEvent struct {
	value int
	withGuid
}


func TestFindEvent(t *testing.T) {
	var storeUT = NewInMemStore()
	storeUT.Update(newGuid(), 0, []Event{
		&MyTestEvent{value:1},
		&MyTestEvent{value:2},
	})
	var guid = newGuid()
	storeUT.Update(guid, 0, []Event{
		&MyTestEvent{value:3},
		&MyTestEvent{value:4},
	})

	var events, version = storeUT.Find(guid)
	assert.Equal(t, 2, version)
	assert.Len(t, events, 2)
	assert.Equal(t, 3, events[0].(*MyTestEvent).value)
	assert.Equal(t, 4, events[1].(*MyTestEvent).value)
}

func TestUpdateEvent(t *testing.T) {
	var storeUT = NewInMemStore()
	guid := newGuid()
	storeUT.Update(guid, 0, []Event{
		&MyTestEvent{value:1},
		&MyTestEvent{value:2},
		&MyTestEvent{value:3},
	})

	err := storeUT.Update(guid, 3, []Event{
		&MyTestEvent{value:4},
		&MyTestEvent{value:5},
	})
	assert.Nil(t, err)

	e, v := storeUT.Find(guid);

	assert.Len(t, e, 5)
	assert.Equal(t, 5, v)

}

func TestUpdateEvent_WrongVersion(t *testing.T) {
	var storeUT = NewInMemStore()
	guid := newGuid()
	storeUT.Update(guid, 0, []Event{
		&MyTestEvent{value:1},
		&MyTestEvent{value:2},
		&MyTestEvent{value:3},
	})

	err := storeUT.Update(guid, 4, []Event{})
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "Optimistic locking")
}

func TestUpdateEvent_UseAsSave(t *testing.T) {
	var storeUT = NewInMemStore()

	guid := newGuid()
	err := storeUT.Update(guid, 0, []Event{
		&MyTestEvent{value:1},
		&MyTestEvent{value:2},
	})
	assert.Nil(t, err)
	e, _ := storeUT.Find(guid);
	assert.Len(t, e, 2)
}

func TestFindEvents(t *testing.T) {
	var storeUT = NewInMemStore()
	guid1 := newGuid()
	storeUT.Update(guid1, 0, []Event{
		&MyTestEvent{value:1},
		&MyTestEvent{value:2},
		&MyTestEvent{value:3},
	})
	guid2 := newGuid()
	storeUT.Update(guid2, 0, []Event{
		&MyTestEvent{value:4},
		&MyTestEvent{value:5},
		&MyTestEvent{value:6},
	})
	storeUT.Update(guid1, 3, []Event{
		&MyTestEvent{value:7},
		&MyTestEvent{value:8},
	})
	storeUT.Update(guid2, 3, []Event{
		&MyTestEvent{value:9},
		&MyTestEvent{value:10},
	})


	e := storeUT.GetEvents(0, 10)
	assert.Len(t, e, 10)
	e = storeUT.GetEvents(0, 100)
	assert.Len(t, e, 10)
}