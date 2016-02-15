package eventsourcing
import (
	"fmt"
	"gopkg.in/yaml.v2"
)

// An aggregate implementation representing a bank account
type Account struct {
	baseAggregate
	Balance int
}

// Make sure it implements Aggregate
var _ Aggregate = (*Account)(nil)

// @see Aggregate.applyEvents
func (a *Account) applyEvents(events []Event) {
	for _, e := range events {
		switch event := e.(type){
		case *AccountOpenedEvent: a.Balance = event.initialBalance;
		case *AccountCreditedEvent: a.Balance += event.amount
		case *AccountCreditedBecauseOfMoneyTransferEvent: a.Balance += event.Amount
		case *AccountDebitedEvent: a.Balance -= event.amount
		case *AccountDebitedBecauseOfMoneyTransferEvent: a.Balance -= event.Amount
		case *AccountDebitFailedEvent: //do nothing
		case *AccountDebitBecauseOfMoneyTransferFailedEvent: //do nothing
		default:
			panic(fmt.Sprintf("Unknown event %#v", event))
		}
	}
	a.Version = len(events)
}

// @see Aggregate.processCommand
func (a Account) processCommand(command Command) []Event {
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
		event = &AccountCreditedBecauseOfMoneyTransferEvent{mTDetails:c.mTDetails}
	case *DebitAccountBecauseOfMoneyTransferCommand:
		if a.Balance < c.Amount {
			event = &AccountDebitBecauseOfMoneyTransferFailedEvent{mTDetails:c.mTDetails}
		} else {
			event = &AccountDebitedBecauseOfMoneyTransferEvent{mTDetails:c.mTDetails}
		}
	default:
		panic(fmt.Sprintf("Unknown command %#v", c))
	}
	event.SetGuid(command.GetGuid())
	return []Event{event}
}


// Helper function to restore account according to persisted state in event store
func RestoreAccount(guid guid, store EventStore) *Account {
	a:= NewAccount()
	RestoreAggregate(guid, a, store)
	return a
}

// create new account in an initial state
func NewAccount() *Account {
	return &Account{}
}

// pretty print in YAML
func (a Account) String() string {
	yaml, _ := yaml.Marshal(&a)
	return fmt.Sprintf("Account:\n%v", string(yaml))
}

