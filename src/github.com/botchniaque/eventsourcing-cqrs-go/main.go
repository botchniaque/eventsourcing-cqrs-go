package main
import (
	"sync"
	"fmt"
	"eventsourcing"
	"time"
)


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

	acc1 := as.OpenAccount(10) // acc1: balance=10
	acc2 := as.OpenAccount(10) // acc2: balance=10
	as.CreditAccount(acc1, 190) // acc1: balance=200
	as.DebitAccount(acc1, 100) // acc1: balance=100
	as.DebitAccount(acc2, 100) // Will fail -> no change
	trans1 := mt.Transfer(10, acc1, acc2)
	trans2 := mt.Transfer(100, acc2, acc1)

	//wait and print
	go func() {
		time.Sleep(200*time.Millisecond)
		printEvents(store.GetEvents(0, 100))
		fmt.Printf("%v\n------------------\n", eventsourcing.RestoreAccount(acc1, store))
		fmt.Printf("%v\n------------------\n", eventsourcing.RestoreAccount(acc2, store))
		fmt.Printf("%v\n------------------\n", eventsourcing.RestoreMoneyTransfer(trans1, store))
		fmt.Printf("%v\n------------------\n", eventsourcing.RestoreMoneyTransfer(trans2, store))
		wg.Done()
	}()

	wg.Wait()

}

func printEvents(events []eventsourcing.Event) {
	for i, e := range events {
		fmt.Printf("%v: %#+v\n", i, e)
	}
}

