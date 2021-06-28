package service

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"authorizer/internal/app/model"
)

func TestNew(t *testing.T) {
	type args struct {
		storage Storage
	}

	tests := []struct {
		name string
		args args
		want *Service
	}{
		{"success",
			args{
				storage: &mockStorage{},
			},
			&Service{
				storage: &mockStorage{},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := New(tt.args.storage)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestService_CreateAccount(t *testing.T) {
	type fields struct {
		storage Storage
	}

	type args struct {
		ca CreateAccount
	}

	tests := []struct {
		name         string
		fields       fields
		args         args
		wantResponse TransactionResponse
		wantErr      error
	}{
		{"success",
			fields{
				storage: &mockStorage{},
			},
			args{
				CreateAccount{
					Account: model.Account{
						Id:             1,
						ActiveCard:     true,
						AvailableLimit: 10,
					},
				},
			},
			TransactionResponse{
				Account: model.Account{
					Id:             1,
					ActiveCard:     true,
					AvailableLimit: 10,
				},
				Violations: []string{},
			},
			nil,
		},
		{"alreadyExists",
			fields{
				storage: &mockStorage{},
			},
			args{
				CreateAccount{
					Account: model.Account{
						Id:             2,
						ActiveCard:     true,
						AvailableLimit: 110,
					},
				},
			},
			TransactionResponse{
				Account: model.Account{
					Id:             2,
					ActiveCard:     true,
					AvailableLimit: 110,
				},
				Violations: []string{"account-already-initialized"},
			},
			nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Service{
				storage: tt.fields.storage,
			}

			gotResponse, err := s.CreateAccount(tt.args.ca)
			assert.Equal(t, tt.wantResponse, gotResponse)
			assert.Equal(t, tt.wantErr, err)
		})
	}
}

func TestService_ProcessTransaction(t *testing.T) {
	type fields struct {
		storage Storage
	}

	type args struct {
		tx ProcessTransaction
	}

	currentTime := time.Now()

	tests := []struct {
		name         string
		fields       fields
		args         args
		wantResponse TransactionResponse
		wantErr      error
	}{
		{"success", fields{
			storage: &mockStorage{},
		},
			args{
				tx: ProcessTransaction{
					Transaction: model.Transaction{
						Merchant: "uno",
						Amount:   10,
						Time:     currentTime,
					},
					AccountID: 2,
				},
			},
			TransactionResponse{
				Account: model.Account{
					Id:             2,
					ActiveCard:     true,
					AvailableLimit: 110,
				},
				Violations: []string{},
			},
			nil,
		},
		{"card-not-active", fields{
			storage: &mockStorage{},
		},
			args{
				tx: ProcessTransaction{
					Transaction: model.Transaction{
						Merchant: "uno",
						Amount:   10,
						Time:     currentTime,
					},
					AccountID: 1,
				},
			},
			TransactionResponse{
				Account:    model.Account{},
				Violations: []string{"card-not-active"},
			},
			nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Service{
				storage: tt.fields.storage,
			}

			gotResponse, err := s.ProcessTransaction(tt.args.tx)
			assert.Equal(t, tt.wantResponse, gotResponse)
			assert.Equal(t, tt.wantErr, err)
		})
	}
}

type mockStorage struct{}

func (m *mockStorage) GetTransactions(accountID int) []model.Transaction {
	return []model.Transaction{}
}

func (m *mockStorage) CreateAccount(a model.Account) error {
	return nil
}

func (m *mockStorage) GetAccount(aID int) model.Account {
	if aID == 2 {
		return model.Account{
			Id:             2,
			ActiveCard:     true,
			AvailableLimit: 110,
		}
	}

	return model.Account{}
}

func (m *mockStorage) Close() error {
	return nil
}

func (m *mockStorage) ExecuteTransaction(a model.Account, t model.Transaction) (model.Account, error) {
	if a.Id == 2 {
		account := model.Account{
			Id:             2,
			ActiveCard:     true,
			AvailableLimit: 110,
		}

		return account, nil
	}

	return model.Account{}, nil
}
