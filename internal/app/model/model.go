package model

import "time"

type Transaction struct {
	Merchant string    `json:"merchant"`
	Amount   int       `json:"amount"`
	Time     time.Time `json:"time"`
}

type Account struct {
	Id             int  `json:"-"`
	ActiveCard     bool `json:"activeCard"`
	AvailableLimit int  `json:"availableLimit"`
}
