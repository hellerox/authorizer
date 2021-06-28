package service

import (
	"authorizer/internal/app/service/rules"

	log "github.com/sirupsen/logrus"

	"authorizer/internal/app/model"
	"authorizer/internal/app/violations"
)

// Service contains the logic to execute the commands
type Service struct {
	storage Storage
}

// Storage interface used in service to execute or simulate an storage
type Storage interface {
	CreateAccount(a model.Account) error
	GetAccount(aID int) model.Account
	ExecuteTransaction(a model.Account, t model.Transaction) (model.Account, error)
	GetTransactions(accountID int) []model.Transaction
	Close() error
}

// CreateAccount is the input of the createAccount operation
type CreateAccount struct {
	Account model.Account `json:"account"`
}

// TransactionResponse is the response for any operation
type TransactionResponse struct {
	Account    model.Account `json:"account"`
	Violations []string      `json:"violations"`
}

// ProcessTransaction is the input of the transaction operation
type ProcessTransaction struct {
	Transaction model.Transaction `json:"transaction"`
	AccountID   int               `json:"-"`
}

// New creates a new service instance
func New(storage Storage) *Service {
	return &Service{
		storage: storage,
	}
}

// CreateAccount contains the logic to create a new account (for this example only account 1 is created)
// 1.- Verify if the account was already created,
//	if it was already created return the violation ViolationAccountAlreadyExists
// 2.- If it wasn't created before, create a new account in storage
func (s *Service) CreateAccount(ca CreateAccount) (response TransactionResponse, err error) {
	response.Account = ca.Account

	account := s.storage.GetAccount(ca.Account.Id)
	if account.ActiveCard {
		log.Errorf("error:%s id:%d", violations.ViolationAccountAlreadyExists, ca.Account.Id)

		response.Account = account
		response.Violations = append(response.Violations, violations.ViolationAccountAlreadyExists)

		return response, nil
	}

	if err = s.storage.CreateAccount(ca.Account); err != nil {
		response.Violations = append(response.Violations, err.Error())

		return response, err
	}

	response.Violations = []string{}

	return response, nil
}

// ProcessTransaction processes the transaction received in the input json
// 1.- Get account information based on the accountID (always 1 in this example)
// 2.- Get all the transactions executed by this account (info used by the business rules)
// 3.- Execute all the business rules, the rules are functions with the same input and outputs
//      If one of them fail, the response contains the violation
// 4.- If transaction passed all the business rules, then we execute the transaction on the storage
//      updating the availableLimit and registering the new transaction in the history
func (s *Service) ProcessTransaction(tx ProcessTransaction) (response TransactionResponse, err error) {
	accountFound := s.storage.GetAccount(tx.AccountID)
	response.Account = accountFound

	pastTransactions := s.storage.GetTransactions(tx.AccountID)

	br := rules.BusinessRule{
		Transaction:      tx.Transaction,
		PastTransactions: pastTransactions,
		Account:          accountFound,
	}

	isValid, violation := br.ExecuteRules()
	if !isValid {
		response.Violations = []string{violation}
		return response, nil
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
