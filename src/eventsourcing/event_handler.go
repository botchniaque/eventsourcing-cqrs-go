package eventsourcing

type Handler struct {
	Store EventStore
	AccChan chan<- Command
	TransChan chan<- Command
}

func (this *Handler) HandleEvent(event Event) {
	switch e := event.(type){
	case *MoneyTransferCreatedEvent:
		this.AccChan <- &DebitAccountBecauseOfMoneyTransferCommand{amount:e.amount, from:e.from, to:e.to, BaseCommand:BaseCommand{guid:e.from}}
	}

}