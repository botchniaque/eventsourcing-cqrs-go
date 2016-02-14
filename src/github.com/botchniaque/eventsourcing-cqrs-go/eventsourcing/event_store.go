package eventsourcing
import (
	"math"
	"errors"
	"fmt"
)

type EventStore interface {
	Find(guid Guid) (events []Event, version int)
	Update(guid Guid, version int, events []Event) error
	GetEvents(offset int, batchSize int) []Event
}

//in-memory event store
type MemEventStore struct {
	store map[Guid][]Event
	events []Event
	eventChan chan Event
}

func (es *MemEventStore) Find(guid Guid) ([]Event, int) {
	events := es.store[guid]
	return events, len(events)
}

func (es *MemEventStore) GetEventChan() <-chan Event {
	return es.eventChan
}

func (es *MemEventStore) appendEvents(events []Event) {
	es.events = append(es.events, events...)
	for _, e := range events {
		es.eventChan <- e
	}
}

// Update aggregate with events. Returns true if version did not match
func (es *MemEventStore) Update(guid Guid, version int, events []Event) error{
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

func (es *MemEventStore) GetEvents(offset int, batchSize int) []Event {
	until := int(math.Min(float64(offset + batchSize), float64(len(es.events))))
	return es.events[offset:until]
}

// initializer for event store
func NewStore() *MemEventStore {
	return &MemEventStore{
		store:map[Guid][]Event{},
		events:make([]Event, 0),
		eventChan:make(chan Event, 100),
	}
}
