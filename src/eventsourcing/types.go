package eventsourcing
import "github.com/twinj/uuid"

// Common interface for all events
type Command interface {
	Guider
}

// Common interface for all events
type Event interface {
	Guider
}

// Base implementation for all events
type withGuid struct {
	Guid guid
}

func (e *withGuid) SetGuid(g guid) {
	e.Guid = g
}
func (e *withGuid) GetGuid() guid {
	return e.Guid
}

type guid string


func newGuid() guid {
	return guid(uuid.NewV4().String())
}

type Guider interface {
	GetGuid() guid
	SetGuid(guid)
}
