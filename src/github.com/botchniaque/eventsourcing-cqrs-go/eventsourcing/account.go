package eventsourcing
import "fmt"

type Account struct {
	BaseAggregate
	Balance int
	guid    Guid
}

func NewAccount(store EventStore) Account {
	acc := Account{BaseAggregate:BaseAggregate{store:store}}
	return acc
}

type AccountOpenedEvent struct {
	BaseEvent
	initialBalance int
}
type AccountCreditedEvent struct {
	BaseEvent
	amount int
}
type AccountDebitedEvent struct {
	BaseEvent
	amount int
}
type AccountDebitFailedEvent struct {
	BaseEvent
}

type AccountDebitedBecauseOfMoneyTransferEvent struct {
	BaseEvent
	from, to Guid
	amount int
}
type AccountCreditedBecauseOfMoneyTransferEvent struct {
	BaseEvent
	from, to Guid
	amount int
}

type AccountDebitBecauseOfMoneyTransferFailedEvent struct {
	BaseEvent
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
	BaseCommand
	Amount int
}

type CreditAccountCommand struct {
	BaseCommand
	Amount int
}

type OpenAccountCommand struct {
	BaseCommand
	InitialBalance int
}

type DebitAccountBecauseOfMoneyTransferCommand struct {
	BaseCommand
	amount int
	from, to Guid
	transaction int
}

type CreditAccountBecauseOfMoneyTransferCommand struct {
	BaseCommand
	amount int
	from, to Guid
	transaction int
}

func (a Account) ProcessCommand(command Command) []Event {
	switch c := command.(type){
	case *OpenAccountCommand:
		return []Event{&AccountOpenedEvent{initialBalance:c.InitialBalance}}
	case *DebitAccountCommand:
		if a.Balance < c.Amount {
			return []Event{&AccountDebitFailedEvent{}}
		} else {
			return []Event{&AccountDebitedEvent{amount:c.Amount}}
		}
	case *CreditAccountCommand:
		return []Event{&AccountCreditedEvent{amount:c.Amount}}
	case *CreditAccountBecauseOfMoneyTransferCommand:
		return []Event{&AccountCreditedBecauseOfMoneyTransferEvent{amount:c.amount, from:c.from, to:c.to}}
	case *DebitAccountBecauseOfMoneyTransferCommand:
		if a.Balance < c.amount {
			return []Event{&AccountDebitBecauseOfMoneyTransferFailedEvent{}}
		} else {
			return []Event{&AccountDebitedBecauseOfMoneyTransferEvent{amount:c.amount, from:c.from, to:c.to}}
		}
	default:
		panic(fmt.Sprintf("Unknown command %#v", c))
	}
}