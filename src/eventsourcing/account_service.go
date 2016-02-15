package eventsourcing

// common properties for all customer facing services
type Service struct {
	commandChannel chan Command
	store EventStore
}

// Getter for command channel - will allow others to post commands
func (a Service) CommandChannel() chan<- Command {
	return a.commandChannel
}

// Account Service - allows simple account management (open, credit, debit)
type AccountService struct {
	Service
}

// Create account service - initialize command channel, and let it use passed event store
func NewAccountService(store EventStore) *AccountService{
	acc := &AccountService{
		Service:Service{
			commandChannel:make(chan Command),
			store:store,
		},
	}
	return acc
}

// Reads from command channel,
// restores an aggregate,
// processes the command and
// persists received events.
// This method *blocks* until command is available,
// therefore should run in a goroutine
func (a *AccountService) HandleCommands() {
	for {
		c := <- a.commandChannel
		acc := RestoreAccount(c.GetGuid(), a.store)
		a.store.Update(c.GetGuid(), acc.Version, acc.processCommand(c))

	}
}

// Open a new account
// Returns accounts GUID
func (a AccountService) OpenAccount(balance int) guid {
	guid := newGuid()
	c := &OpenAccountCommand{
		InitialBalance:balance,
		withGuid:withGuid{Guid:guid},
	}
	a.commandChannel <- c
	return guid
}

// Creding an existing account
func (a AccountService) CreditAccount(guid guid, amount int) {
	c := &CreditAccountCommand{
		Amount:amount,
		withGuid:withGuid{Guid:guid},
	}
	a.commandChannel <- c
}

//debit an existing account
func (a AccountService) DebitAccount(guid guid, amount int) {
	c := &DebitAccountCommand{
		Amount:amount,
		withGuid:withGuid{Guid:guid},
	}
	a.commandChannel <- c
}
