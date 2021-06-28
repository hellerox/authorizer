package rules

import (
	"authorizer/internal/app/model"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func Test_isActive(t *testing.T) {
	type args struct {
		input BusinessRule
	}

	tests := []struct {
		name  string
		args  args
		want  bool
		want1 string
	}{
		{"success",
			args{
				BusinessRule{
					Transaction:      model.Transaction{},
					PastTransactions: nil,
					Account: model.Account{
						Id:             1,
						ActiveCard:     true,
						AvailableLimit: 0,
					},
				},
			},
			true,
			"",
		},
		{"notActive",
			args{
				BusinessRule{
					Transaction:      model.Transaction{},
					PastTransactions: nil,
					Account: model.Account{
						Id:             1,
						ActiveCard:     false,
						AvailableLimit: 0,
					},
				},
			},
			false,
			"card-not-active",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1 := tt.args.input.isActive()
			assert.Equal(t, tt.want, got)
			assert.Equal(t, tt.want1, got1)
		})
	}
}

func Test_sufficientLimit(t *testing.T) {
	type args struct {
		input BusinessRule
	}

	currentTime := time.Now()

	tx := model.Transaction{
		Merchant: "Merchant",
		Amount:   100,
		Time:     currentTime,
	}

	tests := []struct {
		name  string
		args  args
		want  bool
		want1 string
	}{
		{"success", args{
			BusinessRule{
				Transaction:      tx,
				PastTransactions: nil,
				Account: model.Account{
					Id:             1,
					ActiveCard:     true,
					AvailableLimit: 1000,
				},
			},
		},
			true,
			"",
		},
		{"insufficient", args{
			BusinessRule{
				Transaction:      tx,
				PastTransactions: nil,
				Account: model.Account{
					Id:             1,
					ActiveCard:     true,
					AvailableLimit: 10,
				},
			},
		},
			false,
			"insufficient-limit",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1 := tt.args.input.sufficientLimit()
			assert.Equal(t, tt.want, got)
			assert.Equal(t, tt.want1, got1)
		})
	}
}

func Test_doubleTransaction(t *testing.T) {
	type args struct {
		input BusinessRule
	}

	currentTime := time.Now()
	tx := model.Transaction{
		Merchant: "uno",
		Amount:   10,
		Time:     currentTime,
	}

	tests := []struct {
		name  string
		args  args
		want  bool
		want1 string
	}{
		{"success",
			args{
				BusinessRule{
					Transaction: tx,
					PastTransactions: []model.Transaction{
						{
							Merchant: "cero",
							Amount:   100,
							Time:     currentTime.Add(-5 * time.Minute),
						},
					},
					Account: model.Account{
						Id:             1,
						ActiveCard:     true,
						AvailableLimit: 110,
					},
				},
			},
			true,
			"",
		},
		{"doubleTransaction",
			args{
				BusinessRule{
					Transaction: tx,
					PastTransactions: []model.Transaction{
						{
							Merchant: "uno",
							Amount:   10,
							Time:     currentTime.Add(-90 * time.Second),
						},
					},
					Account: model.Account{
						Id:             1,
						ActiveCard:     true,
						AvailableLimit: 110,
					},
				},
			},
			false,
			"doubled-transaction",
		},
		{"disorderedTransactions",
			args{
				BusinessRule{
					Transaction: tx,
					PastTransactions: []model.Transaction{
						{
							Merchant: "uno",
							Amount:   10,
							Time:     currentTime.Add(-90 * time.Second),
						},
						{
							Merchant: "uno",
							Amount:   10,
							Time:     currentTime.Add(-90 * time.Minute),
						},
					},
					Account: model.Account{
						Id:             1,
						ActiveCard:     true,
						AvailableLimit: 110,
					},
				},
			},
			false,
			"doubled-transaction",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1 := tt.args.input.doubleTransaction()
			assert.Equal(t, tt.want, got)
			assert.Equal(t, tt.want1, got1)
		})
	}
}

func TestBusinessRule_highFrequency(t *testing.T) {
	type fields struct {
		Transaction      model.Transaction
		PastTransactions []model.Transaction
		Account          model.Account
	}

	currentTime := time.Now()
	tx := model.Transaction{
		Merchant: "uno",
		Amount:   100,
		Time:     currentTime,
	}

	tests := []struct {
		name   string
		fields fields
		want   bool
		want1  string
	}{
		{"success",
			fields{
				Transaction:      tx,
				PastTransactions: nil,
				Account:          model.Account{},
			},
			true,
			"",
		},
		{"highFrequency",
			fields{
				Transaction: tx,
				PastTransactions: []model.Transaction{
					{
						Merchant: "dos",
						Amount:   10,
						Time:     currentTime.Add(1 * time.Second),
					},
					{
						Merchant: "tres",
						Amount:   20,
						Time:     currentTime.Add(1 * time.Minute),
					},
					{
						Merchant: "cuatro",
						Amount:   30,
						Time:     currentTime.Add(14 * time.Second),
					},
				},
				Account: model.Account{},
			},
			false,
			"high-frequency-small-interval",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			br := &BusinessRule{
				Transaction:      tt.fields.Transaction,
				PastTransactions: tt.fields.PastTransactions,
				Account:          tt.fields.Account,
			}

			got, got1 := br.highFrequency()
			assert.Equal(t, tt.want, got)
			assert.Equal(t, tt.want1, got1)
		})
	}
}

func TestBusinessRule_ExecuteRules(t *testing.T) {
	type fields struct {
		Transaction      model.Transaction
		PastTransactions []model.Transaction
		Account          model.Account
	}

	currentTime := time.Now()
	tests := []struct {
		name   string
		fields fields
		want   bool
		want1  string
	}{
		{"success",
			fields{
				Transaction: model.Transaction{
					Merchant: "uno",
					Amount:   11,
					Time:     time.Now(),
				},
				PastTransactions: []model.Transaction{},
				Account: model.Account{
					Id:             1,
					ActiveCard:     true,
					AvailableLimit: 100,
				},
			},
			true,
			"",
		},
		{"cardNotActive",
			fields{
				Transaction: model.Transaction{
					Merchant: "uno",
					Amount:   11,
					Time:     time.Now(),
				},
				PastTransactions: []model.Transaction{},
				Account: model.Account{
					Id:             1,
					ActiveCard:     false,
					AvailableLimit: 100,
				},
			},
			false,
			"card-not-active",
		},
		{"doubleTransaction",
			fields{
				Transaction: model.Transaction{
					Merchant: "uno",
					Amount:   11,
					Time:     currentTime,
				},
				PastTransactions: []model.Transaction{
					{
						Merchant: "uno",
						Amount:   11,
						Time:     currentTime.Add(1 * time.Second),
					},
				},
				Account: model.Account{
					Id:             1,
					ActiveCard:     false,
					AvailableLimit: 100,
				},
			},
			false,
			"card-not-active",
		},
		{"highFrequency",
			fields{
				Transaction: model.Transaction{
					Merchant: "uno",
					Amount:   11,
					Time:     currentTime,
				},
				PastTransactions: []model.Transaction{
					{
						Merchant: "uno2",
						Amount:   111,
						Time:     currentTime.Add(1 * time.Second),
					},
					{
						Merchant: "uno2",
						Amount:   112,
						Time:     currentTime.Add(2 * time.Second),
					},
				},
				Account: model.Account{
					Id:             1,
					ActiveCard:     false,
					AvailableLimit: 100,
				},
			},
			false,
			"card-not-active",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			br := &BusinessRule{
				Transaction:      tt.fields.Transaction,
				PastTransactions: tt.fields.PastTransactions,
				Account:          tt.fields.Account,
			}

			got, got1 := br.ExecuteRules()
			assert.Equal(t, tt.want, got)
			assert.Equal(t, tt.want1, got1)
		})
	}
}
