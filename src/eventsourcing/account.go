package eventsourcing
import "fmt"

type Account struct {
	balance int
}


func (a *Account) applyEvents(events []Event) {
	for _, e := range events {
		switch event := e.(type){
		case AccountOpenedEvent: a.balance = event.initialBalance;
		case AccountCreditedEvent: a.balance += event.amount
		case AccountDebitedEvent: a.balance -= event.amount
		case AccountDebitFailedEvent: //do nothing
		default:
			panic(fmt.Sprintf("Unknown event %#v", event))
		}
	}
}

func (a Account) processCommand(command Command) []Event {
	switch comm := command.(type){
	case OpenAccountCommand: return []Event{AccountOpenedEvent{initialBalance:comm.initialBalance}}
	case DebitAccountCommand:
		if a.balance < comm.amount {
			return []Event{AccountDebitFailedEvent{}}
		} else {
			return []Event{AccountDebitedEvent{amount:comm.amount}}
		}
	case CreditAccountCommand: return []Event{AccountCreditedEvent{amount:comm.amount}}
	default:
		panic(fmt.Sprintf("Unknown command %#v", comm))
	}
}