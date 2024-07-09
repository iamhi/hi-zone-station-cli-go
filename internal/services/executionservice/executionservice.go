package executionservice

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/iamhi/station-cli-go/internal/external/stationclient"
	"github.com/iamhi/station-cli-go/internal/services/transactionservice"
)

const command_login = "login"
const command_transaction_list = "transaction_list"
const command_create_transaction = "transaction_create"

func ExecuteCommand() {
	command := os.Args[1]

	switch command {

	case command_login:
		executeLoginCommand()

	case command_transaction_list:
		executeTransactionList()

	case command_create_transaction:
		executeTransactionCreate()

	default:
		fmt.Println("Command not recognized" + command)
	}
}

func executeTransactionList() {
	transactions := transactionservice.GetLastestTransactions()

	for _, transaction := range transactions {
		fmt.Println(transaction)
	}
}

func executeLoginCommand() {
	username := os.Args[2]
	password := os.Args[3]

	err := stationclient.Login(username, password)

	if err == nil {
		fmt.Println("Login successful")
	} else {
		fmt.Println("Unable to login")
	}
}

func executeTransactionCreate() {
	description := os.Args[2]
	category := os.Args[3]
	value, err := strconv.ParseFloat(strings.TrimSpace(os.Args[4]), 32)

	if err == nil {
		created_transaction, err := transactionservice.CreateTransaction(description, category, float32(value))

		if err == nil {
			fmt.Println("Created transaction: " + created_transaction.Description + " - " + created_transaction.Category + " - " + fmt.Sprint(created_transaction.Value))
		} else {
			fmt.Println(err.Error())
		}
	} else {
		fmt.Println(err.Error())
	}
}
