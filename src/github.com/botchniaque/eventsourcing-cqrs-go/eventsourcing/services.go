package eventsourcing

type Service struct {
	commandChannel chan Command
	store EventStore
}


func (a Service) CommandChannel() chan<- Command {
	return a.commandChannel
}

type AccountService struct {
	Service
}

func NewAccountService(store EventStore) *AccountService{
	acc := &AccountService{
		Service:Service{
			commandChannel:make(chan Command),
			store:store,
		},
	}
	return acc
}

func (a *AccountService) HandleCommands() {
	for {
		c := <- a.commandChannel
		acc := RestoreAccount(c.GetGuid(), a.store)
		a.store.Update(c.GetGuid(), acc.Version, acc.processCommand(c))

	}
}

func (a AccountService) OpenAccount(balance int) guid {
	guid := newGuid()
	c := &OpenAccountCommand{
		InitialBalance:balance,
		withGuid:withGuid{Guid:guid},
	}
	a.commandChannel <- c
	return guid
}

func (a AccountService) CreditAccount(guid guid, amount int) {
	c := &CreditAccountCommand{
		Amount:amount,
		withGuid:withGuid{Guid:guid},
	}
	a.commandChannel <- c
}

func (a AccountService) DebitAccount(guid guid, amount int) {
	c := &DebitAccountCommand{
		Amount:amount,
		withGuid:withGuid{Guid:guid},
	}
	a.commandChannel <- c
}

type MoneyTransferService struct {
	Service
}

func NewMoneyTransferService(store EventStore) *MoneyTransferService{
	mt := &MoneyTransferService{
		Service:Service{
			commandChannel:make(chan Command),
			store:store,
		},
	}
	return mt
}

func (a *MoneyTransferService) HandleCommands() {
	for {
		c := <- a.commandChannel
		mt := RestoreMoneyTransfer(c.GetGuid(), a.store)
		a.store.Update(c.GetGuid(), mt.Version, mt.processCommand(c))
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
