package eventsourcing
import "fmt"

type Aggregate interface {
	ApplyEvents([]Event)
	ProcessCommand(Command) []Event
	setVersion(int)
	Version() int
	setGuid(Guid)
	Guid() Guid
}

type BaseAggregate struct {
	store EventStore
	version int
	guid Guid
}

func (a *BaseAggregate) setVersion(ver int) {
	a.version = ver
}

func (a BaseAggregate) Version() int {
	return a.version
}

func (a *BaseAggregate) setGuid(g Guid) {
	a.guid = g
}

func (a BaseAggregate) Guid() Guid {
	return a.guid
}

func RestoreAggregate(guid Guid, a Aggregate, store EventStore)  {
	events, version := store.Find(guid)
	fmt.Printf("Restoring %v from events %v\n", guid, events)
	a.ApplyEvents(events)
	a.setVersion(version)
	a.setGuid(guid)
}

func PersistResult(a Aggregate, c Command, s EventStore) Guid {
	events := a.ProcessCommand(c)
	if a.Guid() == "" {
		return s.Save(events)
	} else {
		s.Update(a.Guid(), a.Version(), events)
		return a.Guid()
	}
}