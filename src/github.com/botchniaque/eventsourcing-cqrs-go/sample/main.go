package main
import (
	"sync"
	"fmt"
	"reflect"
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
	mt.Transfer(10, acc1, acc2)

	//wait and print
	go func() {
		time.Sleep(200*time.Millisecond)
		printEvents(store.GetEvents(0, 100))
		a1 := eventsourcing.NewAccount()
		eventsourcing.RestoreAggregate(acc1, a1, store)
		a2 := eventsourcing.NewAccount()
		eventsourcing.RestoreAggregate(acc2, a2, store)

		fmt.Printf("acc1: %+v\n", a1)
		fmt.Printf("acc2: %+v\n", a2)
		wg.Done()
	}()

	wg.Wait()

}

func printEvents(events []eventsourcing.Event) {
	for i, e := range events {
		fmt.Printf("%v: %#v\n", i, e)
	}
}

func HandleEvents(accComm chan<- eventsourcing.Guider, transComm chan<- eventsourcing.Guider)  {
	h := eventsourcing.Handler{Store:store, AccChan:accComm, TransChan:transComm}
	eventChan := store.GetEventChan()
	for {
		event := <-eventChan
		fmt.Printf("Got event %v\n", reflect.TypeOf(event))
		h.HandleEvent(event)

	}
}
