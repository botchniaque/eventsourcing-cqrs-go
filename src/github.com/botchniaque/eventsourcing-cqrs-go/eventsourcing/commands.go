package eventsourcing

type Command interface {
	Guid() Guid
	SetGuid(Guid)
}

// Base implementation for all events
type BaseCommand struct {
	guid Guid
}

func (e *BaseCommand) SetGuid(g Guid) {
	e.guid = g
}
func (e *BaseCommand) Guid() Guid {
	return e.guid
}

