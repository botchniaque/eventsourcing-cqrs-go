package eventsourcing

// Common interface for all events
type Event interface {
	addGuid(Guid)
	Guid() Guid
}

// Base implementation for all events
type BaseEvent struct {
	Event
	guid Guid
}

func (e *BaseEvent) addGuid(g Guid) {
	e.guid = g
}
func (e *BaseEvent) Guid() Guid {
	return e.guid
}
