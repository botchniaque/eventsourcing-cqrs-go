package eventsourcing

// Common interface for all event-sourced aggregates
type Aggregate interface {
	Guider
	// apply a list of events to restore actual state (core event sourcing)
	applyEvents([]Event)

	// Process a command according to own actual state (eg. debit account checks account.balance)
	// Produce proper state-changing events
	processCommand(Command) []Event
}

// base implementation for all aggregates - with GUID and Version
type baseAggregate struct {
	withGuid
	Version int
}

// restores given empty aggregate from a state stored in event store
func RestoreAggregate(guid guid, a Aggregate, store EventStore)  {
	events, _ := store.Find(guid)
	a.applyEvents(events)
	a.SetGuid(guid)
}

