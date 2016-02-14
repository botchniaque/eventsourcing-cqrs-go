package eventsourcing

type Aggregate interface {
	Guider
	applyEvents([]Event)
	ProcessCommand(Guider) []Event
	setVersion(int)
	GetVersion() int
	String() string
}

type baseAggregate struct {
	withGuid
	Version int
}

func (a *baseAggregate) setVersion(ver int) {
	a.Version = ver
}

func (a baseAggregate) GetVersion() int {
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

