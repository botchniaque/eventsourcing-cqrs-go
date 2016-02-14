package eventsourcing
import (
	"testing"
	"github.com/stretchr/testify/assert"
)

func TestAccountRestore(t *testing.T) {
	acc := new(Account)
	acc.applyEvents([]Event{
		&AccountOpenedEvent{initialBalance:100},
		&AccountCreditedEvent{amount:100},
		&AccountDebitedEvent{amount:50},
		&AccountDebitFailedEvent{},
	})
	assert.Equal(t, 150, acc.Balance)
}

func TestAccountCommand(t *testing.T) {
	acc := NewAccount()
	e := acc.processCommand(&OpenAccountCommand{InitialBalance: 100})
	assert.Equal(t, []Event{&AccountOpenedEvent{initialBalance:100}}, e)
}