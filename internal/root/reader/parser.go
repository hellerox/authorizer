package reader

import (
	"authorizer/internal/app/service"
	"encoding/json"

	log "github.com/sirupsen/logrus"
)

// defaultID is the Id used when we don't get an input ID
const defaultID = 1

// ReadCreateAccount gets the struct from the text line received
func ReadCreateAccount(s string) *service.CreateAccount {
	createAccount := &service.CreateAccount{}

	if err := json.Unmarshal([]byte(s), createAccount); err != nil {
		log.Errorf("error unmarshaling request: %+v", err)

		return nil
	}

	if createAccount.Account.Id == 0 {
		createAccount.Account.Id = defaultID
	}

	return createAccount
}

// ReadProcessTransaction gets the struct from the text line received
func ReadProcessTransaction(s string) *service.ProcessTransaction {
	processTransaction := &service.ProcessTransaction{}

	if err := json.Unmarshal([]byte(s), processTransaction); err != nil {
		log.Errorf("error unmarshaling request: %+v", err)

		return nil
	}

	if processTransaction.AccountID == 0 {
		processTransaction.AccountID = defaultID
	}

	return processTransaction
}
