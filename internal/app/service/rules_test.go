package service

import (
	"authorizer/internal/app/model"
	"testing"
)

func Test_isActive(t *testing.T) {
	type args struct {
		input ruleInput
	}

	tests := []struct {
		name  string
		args  args
		want  bool
		want1 string
	}{
		{"success",
			args{
				ruleInput{
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
				ruleInput{
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
			got, got1 := isActive(tt.args.input)
			if got != tt.want {
				t.Errorf("isActive() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("isActive() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}
