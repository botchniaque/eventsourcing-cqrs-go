package eventsourcing


// Common interface for all commands
type Command interface {
	Guider
}

//----------------
//Account Commands
//----------------

// Open a new account with specified initial balance
type OpenAccountCommand struct {
	withGuid
	InitialBalance int
}

// Debit an account with specified amount
type DebitAccountCommand struct {
	withGuid
	Amount int
}

// Credit an account with specified amount
type CreditAccountCommand struct {
	withGuid
	Amount int
}

//-----------------------
//Money Transfer Commands
//-----------------------

// Start a new money transfer between 2 accounts
type CreateMoneyTransferCommand struct {
	withGuid
	mTDetails
}

// Debit an account in the process of money transfer
type DebitAccountBecauseOfMoneyTransferCommand struct {
	withGuid
	mTDetails
}

// Credit an account in the process of money transfer
type CreditAccountBecauseOfMoneyTransferCommand struct {
	withGuid
	mTDetails
}

// Mark Money Transfer as 'Debit successful'
type DebitMoneyTransferCommand struct {
	withGuid
	mTDetails
}

// Mark Money Transfer as 'Credit successful'
type CompleteMoneyTransferCommand struct {
	withGuid
	mTDetails
}

// Mark Money Transfer as 'Debit failed'
type FailMoneyTransferCommand struct {
	withGuid
	mTDetails
}

