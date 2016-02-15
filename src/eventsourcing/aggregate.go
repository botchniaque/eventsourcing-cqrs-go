package eventsourcing

type Aggregate interface {
	Guider
	applyEvents([]Event)
	processCommand(Command) []Event
	setVersion(int)
	getVersion() int
	String() string
}

type baseAggregate struct {
	withGuid
	Version int
}

func (a *baseAggregate) setVersion(ver int) {
	a.Version = ver
}

func (a baseAggregate) getVersion() int {
	return a.Version
}

func (a *baseAggregate) setGuid(g guid) {
	a.Guid = g
}

func (a baseAggregate) GetGuid() guid {
	return a.Guid
}

func RestoreAggregate(guid guid, a Aggregate, store EventStore)  {
	events, version := store.Find(guid)
	a.applyEvents(events)
	a.setVersion(version)
	a.SetGuid(guid)
}

