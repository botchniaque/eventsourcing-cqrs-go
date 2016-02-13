package eventsourcing

type Command interface  {

}

type DebitAccountCommand struct {
	Command
	amount int
}

type CreditAccountCommand struct {
	Command
	amount int
}

type OpenAccountCommand struct {
	Command
	initialBalance int
}