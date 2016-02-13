package eventsourcing

type Aggregate interface {
	applyEvents([]Event)
	processCommand(Command) []Event
}

