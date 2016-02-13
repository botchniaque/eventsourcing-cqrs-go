package eventsourcing
import (
	"testing"
	"github.com/stretchr/testify/assert"
)

func TestOkTransfer(t *testing.T) {
	trans := new(MoneyTransfer)
	trans.applyEvents([]Event{
		MoneyTransferCreatedEvent{amount:100, to:NewGuid(), from:NewGuid()},
		MoneyTransferDebitedEvent{},
		MoneyTransferCompletedEvent{},
	})

	assert.Equal(t, Completed, trans.state)
}

func TestFailedTransfer(t *testing.T) {
	trans := new(MoneyTransfer)
	trans.applyEvents([]Event{
		MoneyTransferCreatedEvent{amount:100, to:NewGuid(), from:NewGuid()},
		MoneyTransferFailedDueToLackOfFundsEvent{},
	})

	assert.Equal(t, Failed, trans.state)
}

func TestEventsOutOfOrder(t *testing.T) {
	trans := new(MoneyTransfer)
	trans.applyEvents([]Event{
		MoneyTransferCreatedEvent{amount:100, to:NewGuid(), from:NewGuid()},
		MoneyTransferCompletedEvent{},
		MoneyTransferDebitedEvent{},
		MoneyTransferDebitedEvent{},
		MoneyTransferDebitedEvent{},
	})
	assert.Equal(t, Debited, trans.state)

	trans.applyEvents([]Event{
		MoneyTransferCompletedEvent{},
		MoneyTransferFailedDueToLackOfFundsEvent{},
	})
	assert.Equal(t, Completed, trans.state)
}

func TestFullTransfer(t *testing.T) {
	acc1 := new(Account)
	acc1.guid = NewGuid()
	acc1.applyEvents(acc1.processCommand(OpenAccountCommand{initialBalance:100}))

	acc2 := new(Account)
	acc2.guid = NewGuid()
	acc2.applyEvents(acc1.processCommand(OpenAccountCommand{initialBalance:10}))

	trans := new(MoneyTransfer)
	mtCreated := trans.processCommand(CreateMoneyTransferCommand{from:acc1.guid, to:acc2.guid, amount:67})

	trans.applyEvents(mtCreated)
	assert.Equal(t, trans.state, Created)

	a1Debited := acc1.processCommand(DebitAccountBecauseOfMoneyTransferCommand{amount:67, from:acc1.guid, to:acc2.guid})
	acc1.applyEvents(a1Debited)

	trans.applyEvents(trans.processCommand(DebitMoneyTransferCommand{}))

	a2Credit := acc2.processCommand(CreditAccountBecauseOfMoneyTransferCommand{amount:67, from:acc1.guid, to:acc2.guid})
	acc2.applyEvents(a2Credit)

	trans.applyEvents(trans.processCommand(CompleteMoneyTransferCommand{}))

	assert.Equal(t, 33, acc1.balance)
	assert.Equal(t, 77, acc2.balance)
	assert.Equal(t, Completed, trans.state)

}