package eventsourcing

type Handler struct {
	store EventStore
}

func (h *Handler) handleEvents(events []Event) {
	for _, event := range events {
		switch e := event.(type){
		case MoneyTransferCreatedEvent:
			acc := NewAccount()
			restore(e.from, acc, h.store)
			newEvents := acc.processCommand(DebitAccountBecauseOfMoneyTransferCommand{amount:e.amount, from:e.from, to:e.to})
			h.store.Update(acc.guid, acc.version, newEvents)
		}
	}
}