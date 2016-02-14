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

type MoneyTransfer struct {
	baseAggregate
	mTDetails
	State State
}

type mTDetails struct {
	From        guid
	To          guid
	Amount      int
	Transaction guid
}

type MoneyTransferCreatedEvent struct {
	withGuid
	mTDetails
}

type MoneyTransferDebitedEvent struct {
	withGuid
	mTDetails
}

type MoneyTransferCompletedEvent struct {
	withGuid
	mTDetails
}

type MoneyTransferFailedDueToLackOfFundsEvent struct {
	withGuid
	mTDetails
}

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
}

type CreateMoneyTransferCommand struct {
	withGuid
	mTDetails
}
type DebitMoneyTransferCommand struct {
	withGuid
	mTDetails
}

type CompleteMoneyTransferCommand struct {
	withGuid
	mTDetails
}

type FailMoneyTransferCommand struct {
	withGuid
	mTDetails
}


func (t MoneyTransfer) String() string {
	yaml, _ := yaml.Marshal(&t)
	return fmt.Sprintf("MoneyTransfer:\n%v", string(yaml))
}

func (t MoneyTransfer) ProcessCommand(command Guider) []Event {
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


func RestoreMoneyTransfer(guid guid, store EventStore) *MoneyTransfer {
	t := new(MoneyTransfer)
	RestoreAggregate(guid, t, store)
	return t
}