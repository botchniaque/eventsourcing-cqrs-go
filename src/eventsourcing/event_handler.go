package eventsourcing

type Handler struct {
	store EventStore
}

func (this *Handler) handleEvents(events []Event) {
	for _, event := range events {
		switch e := event.(type){
		case MoneyTransferCreatedEvent:
			acc := NewAccount(this.store)
			restore(e.from, acc, this.store)
			newEvents := acc.processCommand(DebitAccountBecauseOfMoneyTransferCommand{amount:e.amount, from:e.from, to:e.to})
			this.store.Update(acc.guid, acc.version, newEvents)
		}
	}
}