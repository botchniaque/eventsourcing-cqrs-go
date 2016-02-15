package eventsourcing
import (
	"math"
	"errors"
	"fmt"
)


// Event store's common interface
type EventStore interface {
	// find all events for given ID (aggregate).
	// returns event list as well as aggregate version
	Find(guid guid) (events []Event, version int)

	// Update an aggregate with new events. If the version specified
	// does not match with the version in the Event Store, an error is returned
	Update(guid guid, version int, events []Event) error

	// Get events from Event Store.
	// Supports pagination with use of offset and batchsize.
	GetEvents(offset int, batchSize int) []Event
}

//in-memory event store. Uses slice for 'complete events catalogue'
// and a map for 'per aggregate' events
type MemEventStore struct {
	store map[guid][]Event
	events []Event
	eventChan chan Event
}

// @see EventStore.GetEvents
func (es *MemEventStore) Find(guid guid) ([]Event, int) {
	events := es.store[guid]
	return events, len(events)
}


// @see EventStore.Update
func (es *MemEventStore) Update(guid guid, version int, events []Event) error{
	changes, ok := es.store[guid]
	if !ok {
		// initialize if not exists
		changes = []Event{}
	}
	if len(changes) == version {
		for _, event := range events {
			event.SetGuid(guid)
		}
		es.appendEvents(events)
		es.store[guid] = append(changes, events...)
	} else {
		return errors.New(
			fmt.Sprintf("Optimistic locking exeption - client has version %v, but store %v", version, len(changes)))
	}
	return nil
}

// @see EventStore.GetEvents
func (es *MemEventStore) GetEvents(offset int, batchSize int) []Event {
	until := int(math.Min(float64(offset + batchSize), float64(len(es.events))))
	return es.events[offset:until]
}

// initializer for event store
func NewInMemStore() *MemEventStore {
	return &MemEventStore{
		store:map[guid][]Event{},
		events:make([]Event, 0),
		eventChan:make(chan Event, 100),
	}
}

// Get persisted events channel -
// channel notifies of any change persisted int the event store
func (es *MemEventStore) GetEventChan() <-chan Event {
	return es.eventChan
}

// Add events to the store and send them down the channel
func (es *MemEventStore) appendEvents(events []Event) {
	es.events = append(es.events, events...)
	for _, e := range events {
		es.eventChan <- e
	}
}
