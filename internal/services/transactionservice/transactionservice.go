package transactionservice

import (
	"log"

	"github.com/iamhi/station-cli-go/internal/external/stationclient"
	"github.com/iamhi/station-cli-go/internal/services/tokenservice"
)

func GetLastestTransactions() []TransactionDto {
	tokenservice.CheckTokens()
	transaction_response, err := stationclient.GetLatestTransactions()

	if err != nil {
		log.Println("Error getting the transactions " + err.Error())
	}

	return mapTransactionResponseListToDto(transaction_response)
}

func CreateTransaction(
	description string,
	category string,
	value float32,
) (TransactionDto, error) {
	transaction_response, err := stationclient.CreateTransaction(description, category, value)

	if err != nil {
		return TransactionDto{}, err
	}

	return mapTransactionResponseToDto(transaction_response), nil
}

func mapTransactionResponseToDto(transaction_response stationclient.TransactionResponse) TransactionDto {
	return TransactionDto{
		transaction_response.Uuid,
		transaction_response.Description,
		transaction_response.Category,
		transaction_response.Value,
		transaction_response.CreatedAt,
	}

}

func mapTransactionResponseListToDto(transaction_response_list []stationclient.TransactionResponse) []TransactionDto {
	transaction_dto_list := make([]TransactionDto, len(transaction_response_list))

	for i := range transaction_response_list {
		transaction_dto_list[i] = mapTransactionResponseToDto(transaction_response_list[i])
	}

	return transaction_dto_list
}
