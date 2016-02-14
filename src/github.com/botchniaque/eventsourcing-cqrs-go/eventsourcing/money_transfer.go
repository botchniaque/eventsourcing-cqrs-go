package eventsourcing
import "fmt"

type State int

const (
	Created = State(iota)
	Debited = State(iota)
	Completed = State(iota)
	Failed = State(iota)
)

type MoneyTransfer struct {
	BaseAggregate
	guid Guid
	from Guid
	to Guid
	amount int
	state State
}

type MoneyTransferCreatedEvent struct {
	WithGuid
	from Guid
	to Guid
	amount int
}

type MoneyTransferDebitedEvent struct {
	WithGuid
}

type MoneyTransferCompletedEvent struct {
	WithGuid
}

type MoneyTransferFailedDueToLackOfFundsEvent struct {
	WithGuid
}

func (t *MoneyTransfer) ApplyEvents(events []Event) {
	for _, e := range events {
		switch event := e.(type){
		case *MoneyTransferCreatedEvent:
			t.amount = event.amount
			t.from = event.from
			t.to = event.to
			t.state = Created
		case *MoneyTransferDebitedEvent:
			if t.state == Created {
				t.state = Debited
			}
		case *MoneyTransferCompletedEvent:
			if t.state == Debited {
				t.state = Completed
			}
		case *MoneyTransferFailedDueToLackOfFundsEvent:
			if t.state == Created {
				t.state = Failed
			}

		default:
			panic(fmt.Sprintf("Unknown event %#v", event))
		}
	}
}

type CreateMoneyTransferCommand struct {
	WithGuid
	From, To Guid
	Amount   int
}
type DebitMoneyTransferCommand struct {
	WithGuid
}

type CompleteMoneyTransferCommand struct {
	WithGuid
}

type FailMoneyTransferCommand struct {
	WithGuid
}


func (t MoneyTransfer) ProcessCommand(command Guider) []Event {
	switch comm := command.(type){
	case *CreateMoneyTransferCommand: return []Event{
		&MoneyTransferCreatedEvent{amount:comm.Amount, from:comm.From, to:comm.To},
	}
	case *DebitMoneyTransferCommand: return []Event{
		&MoneyTransferDebitedEvent{},
	}
	case *CompleteMoneyTransferCommand: return []Event{
		&MoneyTransferCompletedEvent{},
	}
	case *FailMoneyTransferCommand: return []Event{
		&MoneyTransferFailedDueToLackOfFundsEvent{},
	}
	default:
		panic(fmt.Sprintf("Unknown command %#v", comm))
	}
}