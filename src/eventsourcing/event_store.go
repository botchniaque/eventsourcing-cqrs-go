package eventsourcing
import (
	"fmt"
	"github.com/twinj/uuid"
)

// Common interface for all events
type Event interface {
	addGuid(string)
	Guid() string
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


type EventStore interface {
	Save(events []Event) string
	Find(guid string) (events []Event, version int)
	Update(guid string, version int, events []Event)
	GetEvents(offset int, batchSize int)
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
		fmt.Printf("%#v\n", event)
		fmt.Printf("%v\n", event.Guid())
	}
	es.events = append(es.events, events...)
	es.store[guid] = events
	fmt.Printf("%v\n", len(es.events))
	fmt.Printf("%#v\n", es.store)
	return guid
}

func (es *MemEventStore) Find(guid string) ([]Event, int) {
	var events = es.store[guid]
	return events, len(events)
}

// initializer for event store
func NewStore() MemEventStore {
	return MemEventStore{store:map[string][]Event{}, events:make([]Event, 0)}
}