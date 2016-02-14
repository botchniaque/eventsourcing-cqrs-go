package eventsourcing
import "github.com/twinj/uuid"

// Common interface for all events
type Event interface {
	Guider
}

// Base implementation for all events
type WithGuid struct {
	guid Guid
}

func (e *WithGuid) SetGuid(g Guid) {
	e.guid = g
}
func (e *WithGuid) Guid() Guid {
	return e.guid
}

type Guid string


func NewGuid() Guid {
	return Guid(uuid.NewV4().String())
}

type Guider interface {
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

