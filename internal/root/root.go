package cmd

import (
	reader3 "authorizer/internal/root/reader"
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"strings"

	log "github.com/sirupsen/logrus"

	"authorizer/internal/app/service"
)

const createAccount = "account"
const processTransaction = "transaction"

// Authorizer is the interface of the service with the basic operations createAccount and processTransaction
type Authorizer interface {
	CreateAccount(ca service.CreateAccount) (response service.TransactionResponse, err error)
	ProcessTransaction(pt service.ProcessTransaction) (response service.TransactionResponse, err error)
}

// Execute is the function that controls the flow of the application getting the lines from the stdin
// and executing the operation related to the json received
func Execute(auth Authorizer, reader io.Reader, writer io.Writer) {
	scanner := bufio.NewScanner(reader)

	for scanner.Scan() {
		var response []byte

		line := scanner.Text()

		switch {
		case strings.Contains(line, createAccount):
			createAccount := reader3.ReadCreateAccount(line)

			createAccountResponse, err := auth.CreateAccount(*createAccount)
			if err != nil {
				log.Errorf("error creating account: %+v", err)
			}

			response, err = json.Marshal(createAccountResponse)
			if err != nil {
				log.Fatalf("error marshaling response: %+v", err)

				continue
			}

		case strings.Contains(line, processTransaction):
			processTransaction := reader3.ReadProcessTransaction(line)

			responseTransaction, err := auth.ProcessTransaction(*processTransaction)
			if err != nil {
				log.Errorf("error processing transaction: %+v", err)
			}

			response, err = json.Marshal(responseTransaction)
			if err != nil {
				log.Fatalf("error marshaling response: %+v", err)

				continue
			}

		default:
			response = []byte("unknown-command")
		}

		fmt.Fprintf(writer, "%s\n", string(response))
	}
}
