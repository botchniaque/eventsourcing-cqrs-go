package eventsourcing
import (
	"fmt"
	"gopkg.in/yaml.v2"
)

type State string

const (
	Created = State("Created")
	Debited = State("Debited")
	Completed = State("Completed")
	Failed = State("Failed")
)

// An aggregate implementation representing a Money Transfer
type MoneyTransfer struct {
	baseAggregate
	mTDetails
	State State
}

// Make sure it implements Aggregate
var _ Aggregate = (*MoneyTransfer)(nil)

// Details of money transfer
type mTDetails struct {
	From        guid
	To          guid
	Amount      int
	Transaction guid
}

// @see Aggregate.applyEvents
func (t *MoneyTransfer) applyEvents(events []Event) {
	for _, e := range events {
		switch event := e.(type){
		case *MoneyTransferCreatedEvent:
			t.Amount = event.Amount
			t.From = event.From
			t.To = event.To
			t.State = Created
			t.Transaction = event.Transaction
		case *MoneyTransferDebitedEvent:
			if t.State == Created {
				t.State = Debited
			}
		case *MoneyTransferCompletedEvent:
			if t.State == Debited {
				t.State = Completed
			}
		case *MoneyTransferFailedDueToLackOfFundsEvent:
			if t.State == Created {
				t.State = Failed
			}

		default:
			panic(fmt.Sprintf("Unknown event %#v", event))
		}
	}
	t.Version = len(events)

}

// @see Aggregate.processCommand
func (t MoneyTransfer) processCommand(command Command) []Event {
	switch c := command.(type){
	case *CreateMoneyTransferCommand: return []Event{
		&MoneyTransferCreatedEvent{
			mTDetails:c.mTDetails,
			withGuid:c.withGuid,
		},
	}
	case *DebitMoneyTransferCommand: return []Event{
		&MoneyTransferDebitedEvent{
			mTDetails:c.mTDetails,
			withGuid:c.withGuid,
		},
	}
	case *CompleteMoneyTransferCommand: return []Event{
		&MoneyTransferCompletedEvent{
			mTDetails:c.mTDetails,
			withGuid:c.withGuid,
		},
	}
	case *FailMoneyTransferCommand: return []Event{
		&MoneyTransferFailedDueToLackOfFundsEvent{
			mTDetails:c.mTDetails,
			withGuid:c.withGuid,
		},
	}
	default:
		panic(fmt.Sprintf("Unknown command %#v", c))
	}
}

// Helper function to restore money transfer according to persisted state in event store
func RestoreMoneyTransfer(guid guid, store EventStore) *MoneyTransfer {
	t := NewMoneyTransfer()
	RestoreAggregate(guid, t, store)
	return t
}

// Creator function for new money transfers in initial state
func NewMoneyTransfer() *MoneyTransfer {
	return &MoneyTransfer{}
}


// Pretty print in YAML
func (t MoneyTransfer) String() string {
	yaml, _ := yaml.Marshal(&t)
	return fmt.Sprintf("MoneyTransfer:\n%v", string(yaml))
}
