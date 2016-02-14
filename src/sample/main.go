package main
import (
	"eventsourcing"
	"sync"
	"fmt"
	"reflect"
	"time"
)

var store = eventsourcing.NewStore()

func main() {
	wg := sync.WaitGroup{}
	wg.Add(1)

	accChan := make(chan eventsourcing.Command)
	transChan := make(chan eventsourcing.Command)

	go handler(accChan, transChan)
	go AccoutCommandHandler(accChan)
	go TransferCommandHandler(transChan)

	create1 := &eventsourcing.OpenAccountCommand{InitialBalance:100}
	create1.SetGuid(eventsourcing.NewGuid())
	accChan <- create1
	create2 := &eventsourcing.OpenAccountCommand{InitialBalance:10}
	create2.SetGuid(eventsourcing.NewGuid())
	accChan <- create2
	transChan <- &eventsourcing.CreateMoneyTransferCommand{
				From:create1.Guid(),
				To:create2.Guid(),
				Amount:10,
			}

	time.Sleep(time.Duration(1000000000))
	for i, e := range store.GetEvents(0, 100) {
		fmt.Printf("%v: %#v\n", i, e)

	}

	wg.Wait()

}

func handler(accComm chan<- eventsourcing.Command, transComm chan<- eventsourcing.Command)  {
	h := eventsourcing.Handler{Store:store, AccChan:accComm, TransChan:transComm}
	eventChan := store.GetEventChan()
	for {
		event := <-eventChan
		fmt.Printf("Got event %v\n", reflect.TypeOf(event))
		h.HandleEvent(event)

	}
}

func AccoutCommandHandler(commChan <-chan eventsourcing.Command) {
	for {
		c := <-commChan
		fmt.Printf("Got command %v\n", reflect.TypeOf(c))
		acc := eventsourcing.NewAccount(store)
		eventsourcing.RestoreAggregate(c.Guid(), &acc, store)
		fmt.Printf("Account %v balance %v\n", acc.Guid(), acc.Balance)
		store.Save(acc.ProcessCommand(c))

	}
}

func TransferCommandHandler(commChan <-chan eventsourcing.Command) {
	for {
		c := <-commChan
		fmt.Printf("Got command %v\n", reflect.TypeOf(c))
		t := new (eventsourcing.MoneyTransfer)
		eventsourcing.RestoreAggregate(c.Guid(), t, store)
		store.Save(t.ProcessCommand(c))

	}
}
