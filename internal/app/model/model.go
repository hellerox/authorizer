package model

import "time"

// Transaction is the object that represents the operation
// executed on the AvailableLimit of the account
type Transaction struct {
	Merchant string    `json:"merchant"`
	Amount   int       `json:"amount"`
	Time     time.Time `json:"time"`
}

// Account is the object that represents the account of a person
// from which we want to subtract balance with each transaction
type Account struct {
	Id             int  `json:"-"`
	ActiveCard     bool `json:"activeCard"`
	AvailableLimit int  `json:"availableLimit"`
}
