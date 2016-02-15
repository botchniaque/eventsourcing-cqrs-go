package eventsourcing

// Common interface for all events
type Event interface {
	Guider
}

// --------------
// Account Events
// --------------

// An account was opened with given initial balance
type AccountOpenedEvent struct {
	withGuid
	initialBalance int
}

// An account was credited specified amount
type AccountCreditedEvent struct {
	withGuid
	amount int
}

// An account was debited specified amount
type AccountDebitedEvent struct {
	withGuid
	amount int
}

// An account was not debited due to insufficient funds.
type AccountDebitFailedEvent struct {
	withGuid
}

// An account was debited in the process of money transfer
type AccountDebitedBecauseOfMoneyTransferEvent struct {
	withGuid
	mTDetails
}

// An account was credited in the process of money transfer
type AccountCreditedBecauseOfMoneyTransferEvent struct {
	withGuid
	mTDetails
}

// An account was not debited in the process of money transfer due to insufficient funds
type AccountDebitBecauseOfMoneyTransferFailedEvent struct {
	withGuid
	mTDetails
}


// --------------------
//Money Transfer Events
// --------------------


// money transfer was initiated (state Created)
type MoneyTransferCreatedEvent struct {
	withGuid
	mTDetails
}

// 'From account' was debited in the process of money transfer (state Debited)
type MoneyTransferDebitedEvent struct {
	withGuid
	mTDetails
}

// 'To account' was credited in the process of money transfer (state Completed)
type MoneyTransferCompletedEvent struct {
	withGuid
	mTDetails
}

// 'From account' was not debited in the process of money transfer (state Failed)
type MoneyTransferFailedDueToLackOfFundsEvent struct {
	withGuid
	mTDetails
}
