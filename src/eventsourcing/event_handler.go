package eventsourcing

// Event Handler subscribes 'persisted events channel' (provided by event store)
// and reacts with commands for some of them
type EventHandler struct {
	store     *MemEventStore
	accChan   chan<- Command
	transChan chan<- Command
}


// Handles events from 'persisted events channel' in an endless loop
// This method blocks as it listens on a channel in a loop
// therefore should run in a goroutine
func (this EventHandler) HandleEvents() {
	eventChan := this.store.GetEventChan()
	for {
		event := <-eventChan
		this.handleEvent(event)
	}
}


// Handle event logic
func (this *EventHandler) handleEvent(event Event) {
	switch e := event.(type){
	// Account Commands
	case *MoneyTransferCreatedEvent:
		this.accChan <- &DebitAccountBecauseOfMoneyTransferCommand{
			mTDetails:e.mTDetails,
			withGuid:withGuid{Guid:e.From},
		}
	case *MoneyTransferDebitedEvent:
		this.accChan <- &CreditAccountBecauseOfMoneyTransferCommand{
			mTDetails:e.mTDetails,
			withGuid:withGuid{Guid:e.To},
		}

	// Money Transfer Commands
	case *AccountDebitedBecauseOfMoneyTransferEvent:
		this.transChan <- &DebitMoneyTransferCommand{
			mTDetails:e.mTDetails,
			withGuid:withGuid{e.Transaction},
		}
	case *AccountDebitBecauseOfMoneyTransferFailedEvent:
		this.transChan <- &FailMoneyTransferCommand{
			mTDetails:e.mTDetails,
			withGuid:withGuid{e.Transaction},
		}
	case *AccountCreditedBecauseOfMoneyTransferEvent:
		this.transChan <- &CompleteMoneyTransferCommand{
			mTDetails:e.mTDetails,
			withGuid:withGuid{e.Transaction},
		}
	}

}


// Constructor for new Event handler.
func NewEventHandler(store *MemEventStore, accChan chan<- Command, transChan chan<- Command) *EventHandler {
	return &EventHandler{store:store, accChan:accChan, transChan:transChan}
}

