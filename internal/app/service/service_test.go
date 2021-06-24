package service

import (
	"authorizer/internal/app/model"
	"reflect"
	"testing"
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
			args{storage: &mockStorage{}},
			&Service{storage: &mockStorage{}},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := New(tt.args.storage); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("New() = %v, want %v", got, tt.want)
			}
		})
	}
}

type mockStorage struct{}

func (m *mockStorage) GetTransactions(accountID int) []model.Transaction {
	panic("implement me")
}

func (m *mockStorage) CreateAccount(a model.Account) error {
	return nil
}

func (m *mockStorage) GetAccount(aID int) model.Account {
	return model.Account{}
}

func (m *mockStorage) ExecuteTransaction(a model.Account, t model.Transaction) (model.Account, error) {
	return model.Account{}, nil
}
