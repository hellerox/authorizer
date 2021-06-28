package rules

import (
	"math"

	log "github.com/sirupsen/logrus"

	"authorizer/internal/app/errors"
	"authorizer/internal/app/model"
)

// BusinessRule contains the list of fields needed for business rules to take a decision
type BusinessRule struct {
	Transaction      model.Transaction
	PastTransactions []model.Transaction
	Account          model.Account
}

// ExecuteRules lists and executes all the business rules
func (br *BusinessRule) ExecuteRules() (bool, string) {
	response, violation := br.isActive()
	if !response {
		return response, violation
	}

	response, violation = br.sufficientLimit()
	if !response {
		return response, violation
	}

	response, violation = br.doubleTransaction()
	if !response {
		return response, violation
	}

	response, violation = br.highFrequency()
	if !response {
		return response, violation
	}

	return true, ""
}

// isActive verifies that your account has an active card
func (br *BusinessRule) isActive() (bool, string) {
	if !br.Account.ActiveCard {
		log.Errorf("violation:%s id:%d", errors.ViolationCardNotActive, br.Account.Id)
		return false, errors.ViolationCardNotActive
	}

	return true, ""
}

// sufficientLimit verifies that your account has enough available limit
// to execute the transaction
func (br *BusinessRule) sufficientLimit() (bool, string) {
	if (br.Account.AvailableLimit - br.Transaction.Amount) < 0 {
		log.Errorf("violation:%s id:%d", errors.ViolationInsufficientLimit, br.Account.Id)

		return false, errors.ViolationInsufficientLimit
	}

	return true, ""
}

// doubleTransaction compares current transaction time against every other transaction trying to find one transactions
// within 2 minutes and with the same amount and merchant
func (br *BusinessRule) doubleTransaction() (bool, string) {
	if len(br.PastTransactions) > 0 {
		for _, pastTx := range br.PastTransactions {
			if br.Transaction.Amount == pastTx.Amount &&
				br.Transaction.Merchant == pastTx.Merchant &&
				math.Abs(br.Transaction.Time.Sub(pastTx.Time).Minutes()) < 2 {
				log.Errorf("violation:%s id:%d", errors.ViolationDoubledTransaction, br.Account.Id)

				return false, errors.ViolationDoubledTransaction
			}
		}
	}

	return true, ""
}

// highFrequency compares current transaction time against every other transaction trying to find
// two other transactions within 2 minutes
func (br *BusinessRule) highFrequency() (bool, string) {
	countInPeriod := 0

	if len(br.PastTransactions) > 0 {
		for _, pastTx := range br.PastTransactions {
			if math.Abs(br.Transaction.Time.Sub(pastTx.Time).Minutes()) < 2 {
				countInPeriod++
				if countInPeriod == 2 {
					log.Errorf("violation:%s id:%d", errors.ViolationHighFrequencySmallInterval, br.Account.Id)

					return false, errors.ViolationHighFrequencySmallInterval
				}
			}
		}
	}

	return true, ""
}
