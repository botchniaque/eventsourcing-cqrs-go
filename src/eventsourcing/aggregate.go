package eventsourcing

type Aggregate interface {
//	restore(Guid, EventStore)
	applyEvents([]Event)
	processCommand(Command) []Event
}

type BaseAggregate struct {
	Aggregate
	version int
}

func restore(guid Guid, a Aggregate, store EventStore)  {
	events, _ := store.Find(guid)
	a.applyEvents(events)
//	a.version = version
//	return a
}

