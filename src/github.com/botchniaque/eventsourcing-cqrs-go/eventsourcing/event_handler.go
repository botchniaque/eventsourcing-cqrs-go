package eventsourcing

type EventHandler struct {
	store     *MemEventStore
	accChan   chan<- Command
	transChan chan<- Command
}

func (this EventHandler) HandleEvents() {
	eventChan := this.store.GetEventChan()
	for {
		event := <-eventChan
		this.handleEvent(event)
	}
}

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
func NewEventHandler(store *MemEventStore, accChan chan<- Command, transChan chan<- Command) *EventHandler {
	return &EventHandler{store:store, accChan:accChan, transChan:transChan}
}

