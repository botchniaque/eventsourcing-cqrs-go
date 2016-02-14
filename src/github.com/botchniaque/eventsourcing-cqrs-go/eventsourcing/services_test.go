package eventsourcing
import (
	"testing"
	"time"
	"sync"
	"github.com/stretchr/testify/assert"
)

func TestServicesScenario(t *testing.T) {
	store := NewInMemStore()
	wg := sync.WaitGroup{}
	wg.Add(1)

	as := NewAccountService(store)
	mt := NewMoneyTransferService(store)
	eh := NewEventHandler(store, as.CommandChannel(), mt.CommandChannel())

	go eh.HandleEvents()
	go as.HandleCommands()
	go mt.HandleCommands()

	accGuid1 := as.OpenAccount(10) // acc1: balance=10
	accGuid2 := as.OpenAccount(10) // acc2: balance=10
	as.CreditAccount(accGuid1, 190) // acc1: balance=200
	as.DebitAccount(accGuid1, 100) // acc1: balance=100
	as.DebitAccount(accGuid2, 100) // Will fail -> no change
	mtGuid1 := mt.Transfer(10, accGuid1, accGuid2) // acc1: balance 90, acc2: balance 20
	mtGuid2 := mt.Transfer(100, accGuid2, accGuid1) // Will fail -> no change

	//wait and print
	go func() {
		time.Sleep(200*time.Millisecond)
		assertAccount(t, store, accGuid1, 4, 90)
		assertAccount(t, store, accGuid2, 4, 20)
		assertMoneyTransfer(t, store, mtGuid1, 3, 10, accGuid1, accGuid2, mtGuid1, Completed)
		assertMoneyTransfer(t, store, mtGuid2, 2, 100, accGuid2, accGuid1, mtGuid2, Failed)
		wg.Done()
	}()

	wg.Wait()
}

func assertAccount(t *testing.T, s EventStore, guid guid, version int, balance int) {
	acc := RestoreAccount(guid, s)
	assert.Equal(t, version, acc.Version, "Wrong version")
	assert.Equal(t, balance, acc.Balance, "Wrong balance")
}

func assertMoneyTransfer(t *testing.T, s EventStore, guid guid, version int, amount int,
from guid, to guid, trans guid, state State) {
	mt := RestoreMoneyTransfer(guid, s)
	assert.Equal(t, version, mt.Version, "Wrong version")
	assert.Equal(t, amount, mt.Amount, "Wrong amount")
	assert.Equal(t, from, mt.From, "Wrong from")
	assert.Equal(t, to, mt.To, "Wrong to")
	assert.Equal(t, trans, mt.Transaction, "Wrong transaction id")
	assert.Equal(t, state, mt.State, "Wrong state")
}