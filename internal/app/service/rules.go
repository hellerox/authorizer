package service

import (
	"math"

	log "github.com/sirupsen/logrus"

	"authorizer/internal/app/errors"
	"authorizer/internal/app/model"
)

type ruleInput struct {
	Transaction      model.Transaction
	PastTransactions []model.Transaction
	Account          model.Account
}

func isActive(input ruleInput) (bool, string) {
	if !input.Account.ActiveCard {
		log.Errorf("violation:%s id:%d", errors.ViolationCardNotActive, input.Account.Id)
		return false, errors.ViolationCardNotActive
	}

	return true, ""
}

func sufficientLimit(input ruleInput) (bool, string) {
	if (input.Account.AvailableLimit + input.Transaction.Amount) < 0 {
		log.Errorf("violation:%s id:%d", errors.ViolationInsufficientLimit, input.Account.Id)

		return false, errors.ViolationInsufficientLimit
	}

	return true, ""
}

func doubleTransaction(input ruleInput) (bool, string) {
	if len(input.PastTransactions) > 0 {
		if math.Abs(input.Transaction.Time.Sub(input.PastTransactions[len(input.PastTransactions)-1].Time).Minutes()) < 2 &&
			input.Transaction.Amount == input.PastTransactions[len(input.PastTransactions)-1].Amount &&
			input.Transaction.Merchant == input.PastTransactions[len(input.PastTransactions)-1].Merchant {
			log.Errorf("violation:%s id:%d", errors.ViolationDoubledTransaction, input.Account.Id)

			return false, errors.ViolationDoubledTransaction
		}
	}

	return true, ""
}

func highFrequency(input ruleInput) (bool, string) {
	return true, ""
}
