package cmd

import (
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

type Authorizer interface {
	CreateAccount(ca service.CreateAccount) (response service.TransactionResponse, err error)
	ProcessTransaction(pt service.ProcessTransaction) (response service.TransactionResponse, err error)
}

func Execute(auth Authorizer, reader io.Reader, writer io.Writer) {
	scanner := bufio.NewScanner(reader)

	for scanner.Scan() {
		var response []byte

		line := scanner.Text()

		switch {
		case strings.Contains(line, createAccount):
			createAccount := service.ReadCreateAccount(line)

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
			processTransaction := service.ReadProcessTransaction(line)

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
