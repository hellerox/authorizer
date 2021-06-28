package storage

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"

	"authorizer/internal/app/model"
)

func TestInMemory_GenerateAccountID(t *testing.T) {
	tests := []struct {
		name string
		want int
	}{
		{"defaultID", 1},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			im := InMemory{}

			got := im.GenerateAccountID()

			assert.Equal(t, tt.want, got)
		})
	}
}

func TestInMemory_CreateAccount(t *testing.T) {
	type fields struct {
		History map[int][]Transaction
		Account map[int]Account
	}

	type args struct {
		a model.Account
	}

	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{"successfulCreation",
			fields{},
			args{
				a: model.Account{
					Id:             1,
					ActiveCard:     true,
					AvailableLimit: 1,
				}},
			false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			im := &InMemory{
				History: tt.fields.History,
				Account: tt.fields.Account,
			}

			if err := im.CreateAccount(tt.args.a); (err != nil) != tt.wantErr {
				t.Errorf("CreateAccount() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestInMemory_GetTransactions(t *testing.T) {
	type fields struct {
		History map[int][]Transaction
		Account map[int]Account
	}

	type args struct {
		accountID int
	}

	currentTime := time.Now()

	tests := []struct {
		name   string
		fields fields
		args   args
		want   []model.Transaction
	}{
		{"success",
			fields{
				History: map[int][]Transaction{
					1: {
						{
							Id:       uuid.MustParse("5171e74b-93dc-4198-8d14-f8b4731fa9c0"),
							Merchant: "Uno",
							Amount:   100,
							Time:     currentTime,
						},
						{
							Id:       uuid.MustParse("5171e74b-93dc-4191-8d14-f8b4731fa9c0"),
							Merchant: "dos",
							Amount:   101,
							Time:     currentTime.Add(2 * time.Hour),
						},
					},
				},
				Account: nil,
			},
			args{accountID: 1},
			[]model.Transaction{
				{
					Merchant: "Uno",
					Amount:   100,
					Time:     currentTime,
				},
				{
					Merchant: "dos",
					Amount:   101,
					Time:     currentTime.Add(2 * time.Hour),
				},
			},
		},
		{"empty",
			fields{
				History: map[int][]Transaction{},
				Account: nil,
			},
			args{accountID: 1},
			[]model.Transaction{},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			im := &InMemory{
				History: tt.fields.History,
				Account: tt.fields.Account,
			}

			got := im.GetTransactions(tt.args.accountID)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestInMemory_GetAccount(t *testing.T) {
	type fields struct {
		History map[int][]Transaction
		Account map[int]Account
	}

	type args struct {
		accountID int
	}

	tests := []struct {
		name   string
		fields fields
		args   args
		want   model.Account
	}{
		{"success",
			fields{
				History: nil,
				Account: map[int]Account{
					1: {
						Id:             1,
						ActiveCard:     true,
						AvailableLimit: 11,
					},
				},
			},
			args{accountID: 1},
			model.Account{
				Id:             1,
				ActiveCard:     true,
				AvailableLimit: 11,
			},
		},
		{"empty",
			fields{
				History: nil,
				Account: map[int]Account{},
			},
			args{accountID: 1},
			model.Account{
				Id:             1,
				ActiveCard:     false,
				AvailableLimit: 0,
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			im := &InMemory{
				History: tt.fields.History,
				Account: tt.fields.Account,
			}

			got := im.GetAccount(tt.args.accountID)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestInMemory_Close(t *testing.T) {
	tests := []struct {
		name    string
		wantErr bool
	}{
		{"success", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			im := InMemory{}
			if err := im.Close(); (err != nil) != tt.wantErr {
				t.Errorf("Close() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
