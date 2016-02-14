package eventsourcing

type AccountService struct {
	commandChannel chan Guider
	store EventStore
}

func NewAccountService(store EventStore) *AccountService{
	acc := &AccountService{
		commandChannel:make(chan Guider),
		store:store,
	}
	return acc
}

func (a *AccountService) HandleCommands() {
	for {
		c := <- a.commandChannel
		acc := RestoreAccount(c.GetGuid(), a.store)
		a.store.Update(c.GetGuid(), acc.Version, acc.ProcessCommand(c))

	}
}

func (a AccountService) CommandChannel() chan<- Guider {
	return a.commandChannel
}

func (a AccountService) OpenAccount(balance int) guid {
	c := &OpenAccountCommand{InitialBalance:balance}
	guid := newGuid()
	c.SetGuid(guid)
	a.commandChannel <- c
	return guid
}

type MoneyTransferService struct {
	commandChannel chan Guider
	store EventStore
}

func NewMoneyTransferService(store EventStore) *MoneyTransferService{
	mt := &MoneyTransferService{
		commandChannel:make(chan Guider),
		store:store,
	}
	return mt
}

func (a *MoneyTransferService) HandleCommands() {
	for {
		c := <- a.commandChannel
		mt := RestoreMoneyTransfer(c.GetGuid(), a.store)
		a.store.Update(c.GetGuid(), mt.Version, mt.ProcessCommand(c))
	}
}


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

func (a MoneyTransferService) CommandChannel() chan<- Guider {
	return a.commandChannel
}