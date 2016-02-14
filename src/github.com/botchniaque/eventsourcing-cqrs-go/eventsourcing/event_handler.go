package eventsourcing

type Handler struct {
	Store EventStore
	AccChan chan<- Guider
	TransChan chan<- Guider
}

func (this *Handler) HandleEvent(event Event) {
	switch e := event.(type){
	case *MoneyTransferCreatedEvent:
		this.AccChan <- &DebitAccountBecauseOfMoneyTransferCommand{amount:e.amount, from:e.from, to:e.to, WithGuid:WithGuid{guid:e.from}}
	}

}