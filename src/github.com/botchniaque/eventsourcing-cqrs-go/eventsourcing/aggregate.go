package eventsourcing
import "fmt"

type Aggregate interface {
	Guider
	ApplyEvents([]Event)
	ProcessCommand(Guider) []Event
	setVersion(int)
	Version() int
}

type BaseAggregate struct {
	WithGuid
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
	a.SetGuid(guid)
}

