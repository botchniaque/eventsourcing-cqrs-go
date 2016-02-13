package eventsourcing

type Aggregate interface {
//	restore(Guid, EventStore)
	applyEvents([]Event)
	processCommand(Command) []Event
	setVersion(int)
}

type BaseAggregate struct {
	Aggregate
	store EventStore
	version int
}

func (a *BaseAggregate) setVersion(ver int) {
	a.version = ver
}

func restore(guid Guid, a Aggregate, store EventStore)  {
	events, version := store.Find(guid)
	a.applyEvents(events)
	a.setVersion(version)
}

