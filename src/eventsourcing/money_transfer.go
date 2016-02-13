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
	guid Guid
	from Guid
	to Guid
	amount int
	state State
}

type MoneyTransferCreatedEvent struct {
	Event
	from Guid
	to Guid
	amount int
}

type MoneyTransferDebitedEvent struct {
	Event
}

type MoneyTransferCompletedEvent struct {
	Event
}

type MoneyTransferFailedDueToLackOfFundsEvent struct {
	Event
}

func (t *MoneyTransfer) applyEvents(events []Event) {
	for _, e := range events {
		switch event := e.(type){
		case MoneyTransferCreatedEvent:
			t.amount = event.amount
			t.from = event.from
			t.to = event.to
			t.state = Created
		case MoneyTransferDebitedEvent:
			if t.state == Created {
				t.state = Debited
			}
		case MoneyTransferCompletedEvent:
			if t.state == Debited {
				t.state = Completed
			}
		case MoneyTransferFailedDueToLackOfFundsEvent:
			if t.state == Created {
				t.state = Failed
			}

		default:
			panic(fmt.Sprintf("Unknown event %#v", event))
		}
	}
}

type CreateMoneyTransferCommand struct {
	Command
	from, to Guid
	amount int
}
type DebitMoneyTransferCommand struct {}

type CompleteMoneyTransferCommand struct {}

type FailMoneyTransferCommand struct {}


func (t MoneyTransfer) processCommand(command Command) []Event {
	switch comm := command.(type){
	case CreateMoneyTransferCommand: return []Event{
		MoneyTransferCreatedEvent{amount:comm.amount, from:comm.from, to:comm.to},
	}
	case DebitMoneyTransferCommand: return []Event{
		MoneyTransferDebitedEvent{},
	}
	case CompleteMoneyTransferCommand: return []Event{
		MoneyTransferCompletedEvent{},
	}
	case FailMoneyTransferCommand: return []Event{
		MoneyTransferFailedDueToLackOfFundsEvent{},
	}
	default:
		panic(fmt.Sprintf("Unknown command %#v", comm))
	}
}