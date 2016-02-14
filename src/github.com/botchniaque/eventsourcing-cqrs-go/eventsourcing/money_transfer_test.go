package eventsourcing
import (
	"testing"
	"github.com/stretchr/testify/assert"
)

func TestOkTransfer(t *testing.T) {
	trans := new(MoneyTransfer)
	trans.applyEvents([]Event{
		&MoneyTransferCreatedEvent{
			mTDetails:mTDetails{
				Amount:100,
				To:newGuid(),
				From:newGuid(),
			}},
		&MoneyTransferDebitedEvent{},
		&MoneyTransferCompletedEvent{},
	})

	assert.Equal(t, Completed, trans.State)
}

func TestFailedTransfer(t *testing.T) {
	trans := new(MoneyTransfer)
	trans.applyEvents([]Event{
		&MoneyTransferCreatedEvent{
			mTDetails:mTDetails{
				Amount:100,
				To:newGuid(),
				From:newGuid(),
			}},
		&MoneyTransferFailedDueToLackOfFundsEvent{},
	})

	assert.Equal(t, Failed, trans.State)
}

func TestEventsOutOfOrder(t *testing.T) {
	trans := new(MoneyTransfer)
	trans.applyEvents([]Event{
		&MoneyTransferCreatedEvent{
			mTDetails:mTDetails{
				Amount:100,
				To:newGuid(),
				From:newGuid(),
			}},
		&MoneyTransferCompletedEvent{},
		&MoneyTransferDebitedEvent{},
		&MoneyTransferDebitedEvent{},
		&MoneyTransferDebitedEvent{},
	})
	assert.Equal(t, Debited, trans.State)

	trans.applyEvents([]Event{
		&MoneyTransferCompletedEvent{},
		&MoneyTransferFailedDueToLackOfFundsEvent{},
	})
	assert.Equal(t, Completed, trans.State)
}

func TestFullTransfer(t *testing.T) {
	//create account with 100
	acc1 := new(Account)
	acc1.Guid = newGuid()
	acc1.applyEvents(acc1.ProcessCommand(&OpenAccountCommand{InitialBalance:100}))

	//create account with 10
	acc2 := new(Account)
	acc2.Guid = newGuid()
	acc2.applyEvents(acc1.ProcessCommand(&OpenAccountCommand{InitialBalance:10}))

	//transfer 67 form 1 to 2
	trans := new(MoneyTransfer)
	fromAcc := acc1.Guid
	toAcc := acc2.Guid
	mtCreated := trans.ProcessCommand(&CreateMoneyTransferCommand{
		mTDetails:mTDetails{
			Amount:67,
			To:toAcc,
			From:fromAcc,
		}})

	//mock handler logic
	trans.applyEvents(mtCreated)
	assert.Equal(t, trans.State, Created)

	a1Debited := acc1.ProcessCommand(&DebitAccountBecauseOfMoneyTransferCommand{
		mTDetails:mTDetails{
			Amount:67,
			To:toAcc,
			From:fromAcc,
		}})
	assert.IsType(t, &AccountDebitedBecauseOfMoneyTransferEvent{}, a1Debited[0])
	acc1.applyEvents(a1Debited)

	trans.applyEvents(trans.ProcessCommand(&DebitMoneyTransferCommand{}))

	a2Credit := acc2.ProcessCommand(&CreditAccountBecauseOfMoneyTransferCommand{
		mTDetails:mTDetails{
			Amount:67,
			To:toAcc,
			From:fromAcc,
		}})
	acc2.applyEvents(a2Credit)

	trans.applyEvents(trans.ProcessCommand(&CompleteMoneyTransferCommand{}))

	//assert final state
	assert.Equal(t, 33, acc1.Balance)
	assert.Equal(t, 77, acc2.Balance)
	assert.Equal(t, Completed, trans.State)

}

func TestFullTransfer_Failed(t *testing.T) {
	//create account with 100
	acc1 := new(Account)
	acc1.Guid = newGuid()
	acc1.applyEvents(acc1.ProcessCommand(&OpenAccountCommand{InitialBalance:100}))

	//create account with 10
	acc2 := new(Account)
	acc2.Guid = newGuid()
	acc2.applyEvents(acc1.ProcessCommand(&OpenAccountCommand{InitialBalance:10}))

	//transfer 67 form 2 to 1 (should fail)
	trans := new(MoneyTransfer)
	fromAcc := acc2.Guid
	toAcc := acc1.Guid
	mtCreated := trans.ProcessCommand(&CreateMoneyTransferCommand{
		mTDetails:mTDetails{
			Amount:67,
			To:toAcc,
			From:fromAcc,
		}})

	//mock handler logic
	trans.applyEvents(mtCreated)
	assert.Equal(t, trans.State, Created)

	a1NotDebited := acc2.ProcessCommand(&DebitAccountBecauseOfMoneyTransferCommand{
		mTDetails:mTDetails{
			Amount:67,
			To:toAcc,
			From:fromAcc,
		}})
	assert.IsType(t, &AccountDebitBecauseOfMoneyTransferFailedEvent{}, a1NotDebited[0])
	acc1.applyEvents(a1NotDebited)

	trans.applyEvents(trans.ProcessCommand(&FailMoneyTransferCommand{}))

	//assert final state
	assert.Equal(t, 100, acc1.Balance)
	assert.Equal(t, 10, acc2.Balance)
	assert.Equal(t, Failed, trans.State)


}