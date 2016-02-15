package main
import (
	"sync"
	"fmt"
	"time"
	"eventsourcing"
)

// Runs example service calls
func main() {
	var store = eventsourcing.NewInMemStore()
	wg := sync.WaitGroup{}
	wg.Add(1)

	as := eventsourcing.NewAccountService(store)
	mt := eventsourcing.NewMoneyTransferService(store)
	eh := eventsourcing.NewEventHandler(store, as.CommandChannel(), mt.CommandChannel())

	go eh.HandleEvents()
	go as.HandleCommands()
	go mt.HandleCommands()

	fmt.Printf("- Open account 1 with balance 10:\tOK\n")
	acc1 := as.OpenAccount(10) // acc1: balance=10
	fmt.Printf("- Open account 2 with balance 10:\tOK\n")
	acc2 := as.OpenAccount(10) // acc2: balance=10
	fmt.Printf("- Credit account 1 with amount 190:\tOK\n")
	as.CreditAccount(acc1, 190) // acc1: balance=200
	fmt.Printf("- Debit account 1 with amount 100:\tOK\n")
	as.DebitAccount(acc1, 100) // acc1: balance=100
	fmt.Printf("- Debit account 2 with amount 100:\tFAIL\n")
	as.DebitAccount(acc2, 100) // Will fail -> no change
	fmt.Printf("- Transfer 10 from account 1 to account 2:\tOK\n")
	trans1 := mt.Transfer(10, acc1, acc2)
	fmt.Printf("- Transfer 100 from account 2 to account 1:\tFAIL\n")
	trans2 := mt.Transfer(100, acc2, acc1)

	//wait and print
	go func() {
		time.Sleep(200*time.Millisecond)
		printEvents(store.GetEvents(0, 100))
		fmt.Printf("-----------------\nAggregates:\n\n")
		fmt.Printf("%v\n------------------\n", eventsourcing.RestoreAccount(acc1, store))
		fmt.Printf("%v\n------------------\n", eventsourcing.RestoreAccount(acc2, store))
		fmt.Printf("%v\n------------------\n", eventsourcing.RestoreMoneyTransfer(trans1, store))
		fmt.Printf("%v\n------------------\n", eventsourcing.RestoreMoneyTransfer(trans2, store))
		wg.Done()
	}()

	wg.Wait()

}

func printEvents(events []eventsourcing.Event) {
	fmt.Printf("-----------------\nEvents after all operations:\n\n")
	for i, e := range events {
		fmt.Printf("%v: %T\n", i, e)
	}
}

