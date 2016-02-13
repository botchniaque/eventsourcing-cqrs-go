package eventsourcing
import (
	"testing"
	"github.com/twinj/uuid"
	"github.com/stretchr/testify/assert"
)

//Example event used for testing
type MyTestEvent struct {
	value int
	BaseEvent
}


func TestSaveEvent(t *testing.T)  {
	var storeUT = NewStore()
	var e1 = MyTestEvent{value:1}

	storeUT.Save([]Event{&e1})

	assert.Len(t, storeUT.events, 1)
	assert.Len(t, storeUT.store, 1)
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
	assert.Equal(t, 2, version)
	assert.Len(t, events, 2)
	assert.Equal(t, 3, events[0].(*MyTestEvent).value)
	assert.Equal(t, 4, events[1].(*MyTestEvent).value)
}

func TestUpdateEvent(t *testing.T) {
	var storeUT = NewStore()
	guid := storeUT.Save([]Event{
		&MyTestEvent{value:1},
		&MyTestEvent{value:2},
		&MyTestEvent{value:3},
	})

	err := storeUT.Update(guid, 3, []Event{
		&MyTestEvent{value:4},
		&MyTestEvent{value:5},
	})
	assert.False(t, err)

	e, v := storeUT.Find(guid);

	assert.Len(t, e, 5)
	assert.Equal(t, 5, v)

}

func TestUpdateEvent_WrongVersion(t *testing.T) {
	var storeUT = NewStore()
	guid := storeUT.Save([]Event{
		&MyTestEvent{value:1},
		&MyTestEvent{value:2},
		&MyTestEvent{value:3},
	})

	err := storeUT.Update(guid, 4, []Event{})
	assert.True(t, err)
}

func TestUpdateEvent_UseAsSave(t *testing.T) {
	var storeUT = NewStore()

	guid := uuid.NewV4().String()
	err := storeUT.Update(guid, 0, []Event{
		&MyTestEvent{value:1},
		&MyTestEvent{value:2},
	})
	assert.False(t, err)
	e, _ := storeUT.Find(guid);
	assert.Len(t, e, 2)
}

func TestFindEvents(t *testing.T) {
	var storeUT = NewStore()
	guid1 := storeUT.Save([]Event{
		&MyTestEvent{value:1},
		&MyTestEvent{value:2},
		&MyTestEvent{value:3},
	})
	guid2 := storeUT.Save([]Event{
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