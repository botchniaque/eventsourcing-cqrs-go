package eventsourcing
import "github.com/twinj/uuid"

// An item having a GUID
type Guider interface {
	GetGuid() guid
	SetGuid(guid)
}

// Base implementation for all Guiders
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


// Create a new GUID - use UUID v4
func newGuid() guid {
	return guid(uuid.NewV4().String())
}
