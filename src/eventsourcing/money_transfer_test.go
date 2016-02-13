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
	//create account with 100
	acc1 := new(Account)
	acc1.guid = NewGuid()
	acc1.applyEvents(acc1.processCommand(OpenAccountCommand{initialBalance:100}))

	//create account with 10
	acc2 := new(Account)
	acc2.guid = NewGuid()
	acc2.applyEvents(acc1.processCommand(OpenAccountCommand{initialBalance:10}))

	//transfer 67 form 1 to 2
	trans := new(MoneyTransfer)
	fromAcc := acc1.guid
	toAcc := acc2.guid
	mtCreated := trans.processCommand(CreateMoneyTransferCommand{from:fromAcc, to:toAcc, amount:67})

	//mock handler logic
	trans.applyEvents(mtCreated)
	assert.Equal(t, trans.state, Created)

	a1Debited := acc1.processCommand(DebitAccountBecauseOfMoneyTransferCommand{amount:67, from:fromAcc, to:toAcc})
	assert.IsType(t, &AccountDebitedBecauseOfMoneyTransferEvent{}, a1Debited[0])
	acc1.applyEvents(a1Debited)

	trans.applyEvents(trans.processCommand(DebitMoneyTransferCommand{}))

	a2Credit := acc2.processCommand(CreditAccountBecauseOfMoneyTransferCommand{amount:67, from:fromAcc, to:toAcc})
	acc2.applyEvents(a2Credit)

	trans.applyEvents(trans.processCommand(CompleteMoneyTransferCommand{}))

	//assert final state
	assert.Equal(t, 33, acc1.balance)
	assert.Equal(t, 77, acc2.balance)
	assert.Equal(t, Completed, trans.state)

}

func TestFullTransfer_Failed(t *testing.T) {
	//create account with 100
	acc1 := new(Account)
	acc1.guid = NewGuid()
	acc1.applyEvents(acc1.processCommand(OpenAccountCommand{initialBalance:100}))

	//create account with 10
	acc2 := new(Account)
	acc2.guid = NewGuid()
	acc2.applyEvents(acc1.processCommand(OpenAccountCommand{initialBalance:10}))

	//transfer 67 form 2 to 1 (should fail)
	trans := new(MoneyTransfer)
	fromAcc := acc2.guid
	toAcc := acc1.guid
	mtCreated := trans.processCommand(CreateMoneyTransferCommand{from:fromAcc, to:toAcc, amount:67})

	//mock handler logic
	trans.applyEvents(mtCreated)
	assert.Equal(t, trans.state, Created)

	a1NotDebited := acc2.processCommand(DebitAccountBecauseOfMoneyTransferCommand{amount:67, from:fromAcc, to:toAcc})
	assert.IsType(t, &AccountDebitBecauseOfMoneyTransferFailedEvent{}, a1NotDebited[0])
	acc1.applyEvents(a1NotDebited)

	trans.applyEvents(trans.processCommand(FailMoneyTransferCommand{}))

	//assert final state
	assert.Equal(t, 100, acc1.balance)
	assert.Equal(t, 10, acc2.balance)
	assert.Equal(t, Failed, trans.state)


}