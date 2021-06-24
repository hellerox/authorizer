package cmd

import (
	"authorizer/internal/app/errors"
	"authorizer/internal/app/model"
	"authorizer/internal/app/service"
	"bytes"
	"io"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

type MockAuthorizer struct {
	countExecCreate  int
	countExecProcess int
}

func (m *MockAuthorizer) CreateAccount(ca service.CreateAccount) (response service.TransactionResponse, err error) {
	if m.countExecCreate == 1 {
		violations := []string{errors.ViolationAccountAlreadyExists}
		return service.TransactionResponse{Violations: violations}, nil
	}

	account := model.Account{
		Id:             1,
		ActiveCard:     true,
		Active:         true,
		AvailableLimit: 10,
	}

	m.countExecCreate++

	return service.TransactionResponse{Account: account, Violations: []string{}}, nil
}

func (m *MockAuthorizer) ProcessTransaction(pt service.ProcessTransaction) (
	response service.TransactionResponse,
	err error,
) {
	if m.countExecCreate == 1 {
		account := model.Account{
			Id:             1,
			ActiveCard:     true,
			Active:         true,
			AvailableLimit: 100,
		}

		m.countExecProcess++

		return service.TransactionResponse{Account: account, Violations: []string{}}, nil
	}

	account := model.Account{
		Id:             1,
		ActiveCard:     true,
		Active:         true,
		AvailableLimit: 10,
	}

	m.countExecProcess++

	return service.TransactionResponse{Account: account, Violations: []string{}}, nil
}

func TestExecute(t *testing.T) {
	type args struct {
		auth   Authorizer
		reader io.Reader
	}

	tests := []struct {
		name           string
		buf            *bytes.Buffer
		args           args
		expectedOutput string
	}{
		{"unknownCommand",
			new(bytes.Buffer),
			args{
				auth:   &MockAuthorizer{},
				reader: strings.NewReader("abcde"),
			},
			"unknown-command \n"},
		{"createAccount",
			new(bytes.Buffer),
			args{
				auth:   &MockAuthorizer{},
				reader: strings.NewReader("{\"account\": { \"activeCard\": true, \"availableLimit\": 10 } }"),
			},
			"{\"account\":{\"activeCard\":true,\"availableLimit\":10},\"violations\":[]} \n"},
		{"doubleCreateAccount",
			new(bytes.Buffer),
			args{
				auth: &MockAuthorizer{},
				reader: strings.NewReader("{\"account\": { \"activeCard\": true, \"availableLimit\": 10 } } \n" +
					"						{\"account\": { \"activeCard\": true, \"availableLimit\": 10 } }"),
			},
			"{\"account\":{\"activeCard\":true,\"availableLimit\":10},\"violations\":[]} \n" +
				"{\"account\":{\"activeCard\":false,\"availableLimit\":0},\"violations\":[\"account-already-initialized\"]} \n"},
		{"createAndProcess",
			new(bytes.Buffer),
			args{
				auth: &MockAuthorizer{},
				reader: strings.NewReader("{\"account\": { \"activeCard\": true, \"availableLimit\": 10 } } \n" +
					"{\"account\": { \"activeCard\": true, \"availableLimit\": 10 } } \n" +
					"{ \"transaction\": { \"merchant\": \"Habbib's\", \"amount\": 90, \"time\": \"2019-02-13T11:00:00.000Z\" } } \n"),
			},
			"{\"account\":{\"activeCard\":true,\"availableLimit\":10},\"violations\":[]} \n" +
				"{\"account\":{\"activeCard\":false,\"availableLimit\":0},\"violations\":[\"account-already-initialized\"]} \n" +
				"{\"account\":{\"activeCard\":true,\"availableLimit\":100},\"violations\":[]} \n"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			Execute(tt.args.auth, tt.args.reader, tt.buf)
			assert.Equal(t, tt.expectedOutput, tt.buf.String())
		})
	}
}
