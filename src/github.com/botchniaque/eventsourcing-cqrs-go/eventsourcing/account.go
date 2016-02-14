package eventsourcing
import "fmt"

type Account struct {
	BaseAggregate
	Balance int
	guid    Guid
}

func NewAccount() *Account {
	acc := Account{}
	return &acc
}

type AccountOpenedEvent struct {
	WithGuid
	initialBalance int
}
type AccountCreditedEvent struct {
	WithGuid
	amount int
}
type AccountDebitedEvent struct {
	WithGuid
	amount int
}
type AccountDebitFailedEvent struct {
	WithGuid
}

type AccountDebitedBecauseOfMoneyTransferEvent struct {
	WithGuid
	from, to Guid
	amount int
}
type AccountCreditedBecauseOfMoneyTransferEvent struct {
	WithGuid
	from, to Guid
	amount int
}

type AccountDebitBecauseOfMoneyTransferFailedEvent struct {
	WithGuid
}

func (a *Account) ApplyEvents(events []Event) {
	for _, e := range events {
		switch event := e.(type){
		case *AccountOpenedEvent: a.Balance = event.initialBalance;
		case *AccountCreditedEvent: a.Balance += event.amount
		case *AccountCreditedBecauseOfMoneyTransferEvent: a.Balance += event.amount
		case *AccountDebitedEvent: a.Balance -= event.amount
		case *AccountDebitedBecauseOfMoneyTransferEvent: a.Balance -= event.amount
		case *AccountDebitFailedEvent: //do nothing
		case *AccountDebitBecauseOfMoneyTransferFailedEvent: //do nothing
		default:
			panic(fmt.Sprintf("Unknown event %#v", event))
		}
	}
}


type DebitAccountCommand struct {
	WithGuid
	Amount int
}

type CreditAccountCommand struct {
	WithGuid
	Amount int
}

type OpenAccountCommand struct {
	WithGuid
	InitialBalance int
}

type DebitAccountBecauseOfMoneyTransferCommand struct {
	WithGuid
	amount int
	from, to Guid
	transaction int
}

type CreditAccountBecauseOfMoneyTransferCommand struct {
	WithGuid
	amount int
	from, to Guid
	transaction int
}

func (a Account) ProcessCommand(command Guider) []Event {
	var event Event
	switch c := command.(type){
	case *OpenAccountCommand:
		event = &AccountOpenedEvent{initialBalance:c.InitialBalance}
	case *DebitAccountCommand:
		if a.Balance < c.Amount {
			event = &AccountDebitFailedEvent{}
		} else {
			event = &AccountDebitedEvent{amount:c.Amount}
		}
	case *CreditAccountCommand:
		event = &AccountCreditedEvent{amount:c.Amount}
	case *CreditAccountBecauseOfMoneyTransferCommand:
		event = &AccountCreditedBecauseOfMoneyTransferEvent{amount:c.amount, from:c.from, to:c.to}
	case *DebitAccountBecauseOfMoneyTransferCommand:
		if a.Balance < c.amount {
			event = &AccountDebitBecauseOfMoneyTransferFailedEvent{}
		} else {
			event = &AccountDebitedBecauseOfMoneyTransferEvent{amount:c.amount, from:c.from, to:c.to}
		}
	default:
		panic(fmt.Sprintf("Unknown command %#v", c))
	}
	event.SetGuid(command.Guid())
	return []Event{event}
}