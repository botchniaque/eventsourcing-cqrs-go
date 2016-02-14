package eventsourcing
import (
	"testing"
	"github.com/stretchr/testify/assert"
)

func TestOkTransfer(t *testing.T) {
	trans := new(MoneyTransfer)
	trans.ApplyEvents([]Event{
		&MoneyTransferCreatedEvent{amount:100, to:NewGuid(), from:NewGuid()},
		&MoneyTransferDebitedEvent{},
		&MoneyTransferCompletedEvent{},
	})

	assert.Equal(t, Completed, trans.state)
}

func TestFailedTransfer(t *testing.T) {
	trans := new(MoneyTransfer)
	trans.ApplyEvents([]Event{
		&MoneyTransferCreatedEvent{amount:100, to:NewGuid(), from:NewGuid()},
		&MoneyTransferFailedDueToLackOfFundsEvent{},
	})

	assert.Equal(t, Failed, trans.state)
}

func TestEventsOutOfOrder(t *testing.T) {
	trans := new(MoneyTransfer)
	trans.ApplyEvents([]Event{
		&MoneyTransferCreatedEvent{amount:100, to:NewGuid(), from:NewGuid()},
		&MoneyTransferCompletedEvent{},
		&MoneyTransferDebitedEvent{},
		&MoneyTransferDebitedEvent{},
		&MoneyTransferDebitedEvent{},
	})
	assert.Equal(t, Debited, trans.state)

	trans.ApplyEvents([]Event{
		&MoneyTransferCompletedEvent{},
		&MoneyTransferFailedDueToLackOfFundsEvent{},
	})
	assert.Equal(t, Completed, trans.state)
}

func TestFullTransfer(t *testing.T) {
	//create account with 100
	acc1 := new(Account)
	acc1.guid = NewGuid()
	acc1.ApplyEvents(acc1.ProcessCommand(&OpenAccountCommand{InitialBalance:100}))

	//create account with 10
	acc2 := new(Account)
	acc2.guid = NewGuid()
	acc2.ApplyEvents(acc1.ProcessCommand(&OpenAccountCommand{InitialBalance:10}))

	//transfer 67 form 1 to 2
	trans := new(MoneyTransfer)
	fromAcc := acc1.guid
	toAcc := acc2.guid
	mtCreated := trans.ProcessCommand(&CreateMoneyTransferCommand{From:fromAcc, To:toAcc, Amount:67})

	//mock handler logic
	trans.ApplyEvents(mtCreated)
	assert.Equal(t, trans.state, Created)

	a1Debited := acc1.ProcessCommand(&DebitAccountBecauseOfMoneyTransferCommand{amount:67, from:fromAcc, to:toAcc})
	assert.IsType(t, &AccountDebitedBecauseOfMoneyTransferEvent{}, a1Debited[0])
	acc1.ApplyEvents(a1Debited)

	trans.ApplyEvents(trans.ProcessCommand(&DebitMoneyTransferCommand{}))

	a2Credit := acc2.ProcessCommand(&CreditAccountBecauseOfMoneyTransferCommand{amount:67, from:fromAcc, to:toAcc})
	acc2.ApplyEvents(a2Credit)

	trans.ApplyEvents(trans.ProcessCommand(&CompleteMoneyTransferCommand{}))

	//assert final state
	assert.Equal(t, 33, acc1.Balance)
	assert.Equal(t, 77, acc2.Balance)
	assert.Equal(t, Completed, trans.state)

}

func TestFullTransfer_Failed(t *testing.T) {
	//create account with 100
	acc1 := new(Account)
	acc1.guid = NewGuid()
	acc1.ApplyEvents(acc1.ProcessCommand(&OpenAccountCommand{InitialBalance:100}))

	//create account with 10
	acc2 := new(Account)
	acc2.guid = NewGuid()
	acc2.ApplyEvents(acc1.ProcessCommand(&OpenAccountCommand{InitialBalance:10}))

	//transfer 67 form 2 to 1 (should fail)
	trans := new(MoneyTransfer)
	fromAcc := acc2.guid
	toAcc := acc1.guid
	mtCreated := trans.ProcessCommand(&CreateMoneyTransferCommand{From:fromAcc, To:toAcc, Amount:67})

	//mock handler logic
	trans.ApplyEvents(mtCreated)
	assert.Equal(t, trans.state, Created)

	a1NotDebited := acc2.ProcessCommand(&DebitAccountBecauseOfMoneyTransferCommand{amount:67, from:fromAcc, to:toAcc})
	assert.IsType(t, &AccountDebitBecauseOfMoneyTransferFailedEvent{}, a1NotDebited[0])
	acc1.ApplyEvents(a1NotDebited)

	trans.ApplyEvents(trans.ProcessCommand(&FailMoneyTransferCommand{}))

	//assert final state
	assert.Equal(t, 100, acc1.Balance)
	assert.Equal(t, 10, acc2.Balance)
	assert.Equal(t, Failed, trans.state)


}