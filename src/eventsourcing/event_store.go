package main
import (
	"github.com/twinj/uuid"
	"math"
)

// Common interface for all events
type Event interface {
	addGuid(string)
	Guid() string
}


type EventStore interface {
	Save(events []Event) string
	Find(guid string) (events []Event, version int)
	Update(guid string, version int, events []Event)
	GetEvents(offset int, batchSize int) []Event
}

// Base implementation for all events
type BaseEvent struct {
	Event
	guid string
}

func (e *BaseEvent) addGuid(g string) {
	e.guid = g
}
func (e *BaseEvent) Guid()string {
	return e.guid
}

//in-memory event store
type MemEventStore struct {
	store map[string][]Event
	events []Event
}

func (es *MemEventStore) Save(events []Event) string {
	guid := uuid.NewV4().String()
	for _, event := range events {
		event.addGuid(guid)
	}
	es.events = append(es.events, events...)
	es.store[guid] = events
	return guid
}

func (es *MemEventStore) Find(guid string) ([]Event, int) {
	events := es.store[guid]
	return events, len(events)
}

// Update aggregate with events. Returns true if version did not match
func (es *MemEventStore) Update(guid string, version int, events []Event) (err bool){
	changes := es.store[guid]
	if len(changes) == version {
		err = false
		for _, event := range events {
			event.addGuid(guid)
		}
		es.events = append(es.events, events...)
		es.store[guid] = append(changes, events...)
	} else {
		err = true
	}
	return err
}

func (es *MemEventStore) GetEvents(offset int, batchSize int) []Event {
	until := int(math.Min(float64(offset + batchSize), float64(len(es.events))))
	return es.events[offset:until]
}

// initializer for event store
func NewStore() MemEventStore {
	return MemEventStore{store:map[string][]Event{}, events:make([]Event, 0)}
}