package service

import (
	log "github.com/sirupsen/logrus"

	"authorizer/internal/app/errors"
	"authorizer/internal/app/model"
)

type Service struct {
	storage Storage
	rules   []func(ruleInput) (bool, string)
}

type Storage interface {
	CreateAccount(a model.Account) error
	GetAccount(aID int) model.Account
	ExecuteTransaction(a model.Account, t model.Transaction) (model.Account, error)
	GetTransactions(accountID int) []model.Transaction
}

type CreateAccount struct {
	Account model.Account `json:"account"`
}

type TransactionResponse struct {
	Account    model.Account `json:"account"`
	Violations []string      `json:"violations"`
}

type ProcessTransaction struct {
	Transaction model.Transaction `json:"transaction"`
	AccountID   int               `json:"-"`
}

func New(storage Storage) *Service {
	rules := []func(ruleInput) (bool, string){
		isActive,
		sufficientLimit,
		doubleTransaction,
		highFrequency,
	}

	return &Service{
		storage: storage,
		rules:   rules,
	}
}

func (s *Service) CreateAccount(ca CreateAccount) (response TransactionResponse, err error) {
	response.Account = ca.Account

	account := s.storage.GetAccount(ca.Account.Id)
	if account.ActiveCard {
		log.Errorf("error:%s id:%d", errors.ViolationAccountAlreadyExists, ca.Account.Id)
		response.Violations = append(response.Violations, errors.ViolationAccountAlreadyExists)

		return response, nil
	}

	if err = s.storage.CreateAccount(ca.Account); err != nil {
		response.Violations = append(response.Violations, err.Error())

		return response, err
	}

	response.Violations = []string{}

	return response, nil
}

func (s *Service) ProcessTransaction(tx ProcessTransaction) (response TransactionResponse, err error) {
	accountFound := s.storage.GetAccount(tx.AccountID)
	response.Account = accountFound

	pastTransactions := s.storage.GetTransactions(tx.AccountID)

	ruleInput := ruleInput{
		Transaction:      tx.Transaction,
		PastTransactions: pastTransactions,
		Account:          accountFound,
	}

	for _, ruleFunc := range s.rules {
		isValid, violation := ruleFunc(ruleInput)
		if !isValid {
			response.Violations = []string{violation}
			return response, nil
		}
	}

	account, err := s.storage.ExecuteTransaction(accountFound, tx.Transaction)
	if err != nil {
		log.Errorf("error:%s id:%d", err, tx.AccountID)

		return response, err
	}

	response.Account = account
	response.Violations = []string{}

	return response, nil
}
