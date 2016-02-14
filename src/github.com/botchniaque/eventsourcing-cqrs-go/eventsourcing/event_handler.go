package eventsourcing

type Handler struct {
	Store EventStore
	AccChan chan<- Guider
	TransChan chan<- Guider
}

func (this *Handler) HandleEvent(event Event) {
	switch e := event.(type){
	// Account Commands
	case *MoneyTransferCreatedEvent:
		this.AccChan <- &DebitAccountBecauseOfMoneyTransferCommand{
			mTDetails:e.mTDetails,
			withGuid:withGuid{Guid:e.From},
		}
	case *MoneyTransferDebitedEvent:
		this.AccChan <- &CreditAccountBecauseOfMoneyTransferCommand{
			mTDetails:e.mTDetails,
			withGuid:withGuid{Guid:e.To},
		}

	// Money Transfer Commands
	case *AccountDebitedBecauseOfMoneyTransferEvent:
		this.TransChan <- &DebitMoneyTransferCommand{
			mTDetails:e.mTDetails,
			withGuid:withGuid{e.Transaction},
		}
	case *AccountDebitBecauseOfMoneyTransferFailedEvent:
		this.TransChan <- &FailMoneyTransferCommand{
			mTDetails:e.mTDetails,
			withGuid:withGuid{e.Transaction},
		}
	case *AccountCreditedBecauseOfMoneyTransferEvent:
		this.TransChan <- &CompleteMoneyTransferCommand{
			mTDetails:e.mTDetails,
			withGuid:withGuid{e.Transaction},
		}
	}

}