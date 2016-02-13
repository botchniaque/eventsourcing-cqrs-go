package eventsourcing

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

type AccountOpenedEvent struct {
	Event
	initialBalance int
}
type AccountCreditedEvent struct {
	Event
	amount int
}
type AccountDebitedEvent struct {
	Event
	amount int
}
type AccountDebitFailedEvent struct {
	Event
}
