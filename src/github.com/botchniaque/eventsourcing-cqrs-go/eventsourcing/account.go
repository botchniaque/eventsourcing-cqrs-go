package eventsourcing
import (
	"fmt"
	"gopkg.in/yaml.v2"
)

type Account struct {
	baseAggregate
	Balance int
}

func NewAccount() *Account {
	return &Account{}
}

func (a Account) String() string {
	yaml, _ := yaml.Marshal(&a)
	return fmt.Sprintf("Account:\n%v", string(yaml))
}

type AccountOpenedEvent struct {
	withGuid
	initialBalance int
}
type AccountCreditedEvent struct {
	withGuid
	amount int
}
type AccountDebitedEvent struct {
	withGuid
	amount int
}
type AccountDebitFailedEvent struct {
	withGuid
}

type AccountDebitedBecauseOfMoneyTransferEvent struct {
	withGuid
	mTDetails
}
type AccountCreditedBecauseOfMoneyTransferEvent struct {
	withGuid
	mTDetails
}

type AccountDebitBecauseOfMoneyTransferFailedEvent struct {
	withGuid
	mTDetails
}

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
}


type DebitAccountCommand struct {
	withGuid
	Amount int
}

type CreditAccountCommand struct {
	withGuid
	Amount int
}

type OpenAccountCommand struct {
	withGuid
	InitialBalance int
}

type DebitAccountBecauseOfMoneyTransferCommand struct {
	withGuid
	mTDetails
}

type CreditAccountBecauseOfMoneyTransferCommand struct {
	withGuid
	mTDetails
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

func RestoreAccount(guid guid, store EventStore) *Account {
	a:= NewAccount()
	RestoreAggregate(guid, a, store)
	return a
}