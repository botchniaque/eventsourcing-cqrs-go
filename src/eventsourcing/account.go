package eventsourcing
import "fmt"

type Account struct {
	balance int
	guid Guid
}

type AccountOpenedEvent struct {
	Event
	initialBalance int
}
type AccountCreditedEvent struct {
	Event
	amount int
}
type AccountDebitedEvent struct {
	Event
	amount int
}
type AccountDebitFailedEvent struct {
	Event
}

type AccountDebitedBecauseOfMoneyTransferEvent struct {
	Event
	from, to Guid
	amount int
}
type AccountCreditedBecauseOfMoneyTransferEvent struct {
	Event
	from, to Guid
	amount int
}

type AccountDebitBecauseOfMoneyTransferFailedEvent struct {
	Event
}

func (a *Account) applyEvents(events []Event) {
	for _, e := range events {
		switch event := e.(type){
		case AccountOpenedEvent: a.balance = event.initialBalance;
		case AccountCreditedEvent: a.balance += event.amount
		case AccountCreditedBecauseOfMoneyTransferEvent: a.balance += event.amount
		case AccountDebitedEvent: a.balance -= event.amount
		case AccountDebitedBecauseOfMoneyTransferEvent: a.balance -= event.amount
		case AccountDebitFailedEvent: //do nothing
		default:
			panic(fmt.Sprintf("Unknown event %#v", event))
		}
	}
}


type DebitAccountCommand struct {
	amount int
}

type CreditAccountCommand struct {
	amount int
}

type OpenAccountCommand struct {
	initialBalance int
}

type DebitAccountBecauseOfMoneyTransferCommand struct {
	amount int
	from, to Guid
	transaction int
}

type CreditAccountBecauseOfMoneyTransferCommand struct {
	amount int
	from, to Guid
	transaction int
}

func (a Account) processCommand(command Command) []Event {
	switch c := command.(type){
	case OpenAccountCommand:
		return []Event{AccountOpenedEvent{initialBalance:c.initialBalance}}
	case DebitAccountCommand:
		if a.balance < c.amount {
			return []Event{AccountDebitFailedEvent{}}
		} else {
			return []Event{AccountDebitedEvent{amount:c.amount}}
		}
	case CreditAccountCommand:
		return []Event{AccountCreditedEvent{amount:c.amount}}
	case CreditAccountBecauseOfMoneyTransferCommand:
		return []Event{AccountCreditedBecauseOfMoneyTransferEvent{amount:c.amount, from:c.from, to:c.to}}
	case DebitAccountBecauseOfMoneyTransferCommand:
		if a.balance < c.amount {
			return []Event{AccountDebitBecauseOfMoneyTransferFailedEvent{}}
		} else {
			return []Event{AccountDebitedBecauseOfMoneyTransferEvent{amount:c.amount, from:c.from, to:c.to}}
		}
	default:
		panic(fmt.Sprintf("Unknown command %#v", c))
	}
}