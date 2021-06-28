package service

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"authorizer/internal/app/model"
)

func TestReadCreateAccount(t *testing.T) {
	type args struct {
		s string
	}

	account := model.Account{
		Id:             1,
		ActiveCard:     true,
		AvailableLimit: 1010,
	}

	successCreateAccount := CreateAccount{Account: account}

	emptyAccount := model.Account{
		Id:             1,
		ActiveCard:     false,
		AvailableLimit: 0,
	}

	otherStructure := CreateAccount{Account: emptyAccount}

	tests := []struct {
		name string
		args args
		want *CreateAccount
	}{
		{"successCase",
			args{s: "{\"account\": { \"activeCard\": true, \"availableLimit\": 1010} }"},
			&successCreateAccount,
		},
		{"OtherStructure",
			args{s: "{\"accounts\": { \"none\": true, \"availableLimit\": 1010} }"},
			&otherStructure,
		},
		{"invalidString",
			args{s: "---"},
			nil,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			got := ReadCreateAccount(tt.args.s)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestReadProcessTransaction(t *testing.T) {
	type args struct {
		s string
	}

	txTime := time.Date(
		2019, 02, 13, 11, 00, 00, 0, time.UTC)
	tx := model.Transaction{
		Merchant: "Habbib's",
		Amount:   90,
		Time:     txTime,
	}

	successProcessTx := ProcessTransaction{
		Transaction: tx,
		AccountID:   defaultID,
	}

	emptyProcessTx := ProcessTransaction{AccountID: 1}

	tests := []struct {
		name string
		args args
		want *ProcessTransaction
	}{
		{"successCase",
			args{s: "{ \"transaction\": { \"merchant\": \"Habbib's\", \"amount\": 90," +
				" \"time\": \"2019-02-13T11:00:00.000Z\" } }"},
			&successProcessTx,
		},
		{"OtherStructure",
			args{s: "{ \"tx\": { \"merchant\": \"Habbib's\", \"amount\": 90, \"time\": \"2019-02-13T11:00:00.000Z\" } }"},
			&emptyProcessTx,
		},
		{"invalidString",
			args{s: "---"},
			nil,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			got := ReadProcessTransaction(tt.args.s)
			assert.Equal(t, tt.want, got)
		})
	}
}
