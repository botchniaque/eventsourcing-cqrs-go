package main
import (
	"testing"
	"github.com/stretchr/testify/assert"
)

func TestAccountRestore(t *testing.T) {
	acc := new(Account)
	acc.applyEvents([]Event{
		AccountOpenedEvent{initialBalance:100},
		AccountCreditedEvent{amount:100},
		AccountDebitedEvent{amount:50},
		AccountDebitFailedEvent{},
	})
	assert.Equal(t, 150, acc.balance)
}

func TestAccountCommand(t *testing.T) {
	acc := new(Account)
	e := acc.processCommand(OpenAccountCommand{initialBalance: 100})
	assert.Equal(t, []Event{AccountOpenedEvent{initialBalance:100}}, e)
}