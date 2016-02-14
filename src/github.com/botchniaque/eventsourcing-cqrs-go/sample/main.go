package main
import (
	"sync"
	"fmt"
	"github.com/botchniaque/eventsourcing-cqrs-go/eventsourcing"
	"time"
)

var store = eventsourcing.NewStore()

func main() {
	wg := sync.WaitGroup{}
	wg.Add(1)

	as := eventsourcing.NewAccountService(store)
	mt := eventsourcing.NewMoneyTransferService(store)

	go HandleEvents(as.CommandChannel(), mt.CommandChannel())
	go as.HandleCommands()
	go mt.HandleCommands()

	acc1 := as.OpenAccount(100)
	acc2 := as.OpenAccount(10)
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

func HandleEvents(accComm chan<- eventsourcing.Guider, transComm chan<- eventsourcing.Guider)  {
	h := eventsourcing.Handler{Store:store, AccChan:accComm, TransChan:transComm}
	eventChan := store.GetEventChan()
	for {
		event := <-eventChan
		h.HandleEvent(event)

	}
}
