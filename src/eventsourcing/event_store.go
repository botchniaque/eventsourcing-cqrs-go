package eventsourcing
import (
	"github.com/twinj/uuid"
	"math"
)

type Guid string

type EventStore interface {
	Save(events []Event) Guid
	Find(guid Guid) (events []Event, version int)
	Update(guid Guid, version int, events []Event) bool
	GetEvents(offset int, batchSize int) []Event
}

//in-memory event store
type MemEventStore struct {
	store map[Guid][]Event
	events []Event
	eventChan chan Event
}

func (es *MemEventStore) Save(events []Event) Guid {
	guid := NewGuid()
//	for _, event := range events {
//		event.addGuid(guid)
//	}
	es.appendEvents(events)
	es.store[guid] = events
	return guid
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
func (es *MemEventStore) Update(guid Guid, version int, events []Event) (err bool){
	changes := es.store[guid]
	if len(changes) == version {
		err = false
		for _, event := range events {
			event.addGuid(guid)
		}
		es.appendEvents(events)
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
func NewStore() *MemEventStore {
	return &MemEventStore{
		store:map[Guid][]Event{},
		events:make([]Event, 0),
		eventChan:make(chan Event, 100),
	}
}

func NewGuid() Guid {
	return Guid(uuid.NewV4().String())
}