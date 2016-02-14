package main
import (
	"sync"
	"fmt"
	"reflect"
	"time"
	"github.com/botchniaque/eventsourcing-cqrs-go/eventsourcing"
)

var store = eventsourcing.NewStore()

func main() {
	wg := sync.WaitGroup{}
	wg.Add(1)

	as := eventsourcing.NewAccountService(store)
	mt := eventsourcing.NewMoneyTransferService(store)

	accChan := as.CommandChannel()
	transChan := mt.CommandChannel()

	go handler(accChan, transChan)
	go as.CommandHandler()
	go mt.CommandHandler()

	acc1 := as.OpenAccount(100)
	acc2 := as.OpenAccount(10)
	mt.Transfer(10, acc1, acc2)

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
