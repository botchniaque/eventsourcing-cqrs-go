package eventsourcing
import (
	"fmt"
	"reflect"
)

type AccountService struct {
	commandChannel chan Guider
	store EventStore
}

func NewAccountService(store EventStore) *AccountService{
	acc := &AccountService{
		commandChannel:make(chan Guider),
		store:store,
	}
	return acc
}

func (a *AccountService) HandleCommands() {
	for {
		c := <- a.commandChannel
		fmt.Printf("Got command %v\n", reflect.TypeOf(c))
		acc := NewAccount()
		RestoreAggregate(c.Guid(), acc, a.store)
		fmt.Printf("Account %v balance %v\n", acc.Guid(), acc.Balance)
		a.store.Update(c.Guid(), acc.Version(), acc.ProcessCommand(c))

	}
}

func (a AccountService) CommandChannel() chan<- Guider {
	return a.commandChannel
}

func (a AccountService) OpenAccount(balance int) Guid {
	c := &OpenAccountCommand{InitialBalance:balance}
	guid := NewGuid()
	c.SetGuid(guid)
	a.commandChannel <- c
	return guid
}

type MoneyTransferService struct {
	commandChannel chan Guider
	store EventStore
}

func NewMoneyTransferService(store EventStore) *MoneyTransferService{
	mt := &MoneyTransferService{
		commandChannel:make(chan Guider),
		store:store,
	}
	return mt
}

func (a *MoneyTransferService) HandleCommands() {
	for {
		c := <- a.commandChannel
		fmt.Printf("Got command %v\n", reflect.TypeOf(c))
		mt := new (MoneyTransfer)
		RestoreAggregate(c.Guid(), mt, a.store)
		a.store.Update(c.Guid(), mt.Version(), mt.ProcessCommand(c))

	}
}


func (a MoneyTransferService) Transfer(amount int, from Guid, to Guid) Guid {
	c := &CreateMoneyTransferCommand{
		From:from,
		To:to,
		Amount:amount,
	}

	guid := NewGuid()
	c.SetGuid(guid)
	a.commandChannel <- c
	return guid
}

func (a MoneyTransferService) CommandChannel() chan<- Guider {
	return a.commandChannel
}