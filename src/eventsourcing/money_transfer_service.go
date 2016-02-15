package eventsourcing

// Money Transfer Service - allows creating new transfers
type MoneyTransferService struct {
	Service
}

// Create money transfer service - initialize command channel, and let it use passed event store
func NewMoneyTransferService(store EventStore) *MoneyTransferService{
	mt := &MoneyTransferService{
		Service:Service{
			commandChannel:make(chan Command),
			store:store,
		},
	}
	return mt
}

// Reads from command channel,
// restores an aggregate,
// processes the command and
// persists received events.
// This method *blocks* until command is available,
// therefore should run in a goroutine
func (a *MoneyTransferService) HandleCommands() {
	for {
		c := <- a.commandChannel
		mt := RestoreMoneyTransfer(c.GetGuid(), a.store)
		a.store.Update(c.GetGuid(), mt.Version, mt.processCommand(c))
	}
}

// Create a new money transfer between 2 existing accounts
func (a MoneyTransferService) Transfer(amount int, from guid, to guid) guid {
	guid := newGuid()
	c := &CreateMoneyTransferCommand{
		mTDetails:mTDetails{
			Amount:amount,
			From:from,
			To:to,
			Transaction:guid,
		},
		withGuid:withGuid{Guid:guid},
	}
	a.commandChannel <- c
	return guid
}
