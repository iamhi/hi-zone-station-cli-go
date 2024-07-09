package stationclient

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/iamhi/station-cli-go/internal/config/stationapiconfig"
	"github.com/iamhi/station-cli-go/internal/storage/tokenstorage"
)

const application_json_media_type = "application/json"

func Login(username string, password string) error {
	login_response, err := loginAttempt(LoginRequest{
		username,
		password,
	})

	if err == nil {
		tokenstorage.StoreTokens(login_response.AccessToken, login_response.RefreshToken)

		return nil
	}

	return err
}

func loginAttempt(login_request LoginRequest) (LoginResponse, error) {
	login_request_json, marshal_error := json.Marshal(login_request)

	if marshal_error != nil {
		return LoginResponse{}, fmt.Errorf("Error marshaling the request")
	}

	login_response, login_error := http.Post(
		stationapiconfig.GetStationApiUrl()+"/token/generate",
		application_json_media_type,
		bytes.NewBuffer(login_request_json))

	if login_error != nil {
		return LoginResponse{}, fmt.Errorf("Error connectin to authentication service")
	}

	defer login_response.Body.Close()

	var login_result LoginResponse

	decoding_error := json.NewDecoder(login_response.Body).Decode(&login_result)

	if decoding_error != nil {
		return LoginResponse{}, fmt.Errorf("Error decoing the response from authentication service")
	}

	return login_result, nil
}

func GetLatestTransactions() ([]TransactionResponse, error) {
	req, err := http.NewRequest("GET", stationapiconfig.GetStationApiUrl()+"/transaction", nil)

	if err != nil {
		return []TransactionResponse{}, nil
	}

	access_token, err := tokenstorage.ReadAccessToken()

	if err != nil {
		return []TransactionResponse{}, err
	}

	req.Header.Set("Authorization", access_token)

	client := &http.Client{}

	transaction_response, err := client.Do(req)

	if err != nil {
		return []TransactionResponse{}, err
	}

	defer transaction_response.Body.Close()

	if transaction_response.StatusCode == http.StatusUnauthorized {
		access_token, access_token_err := refreshAccessToken()

		if access_token_err != nil {
			return []TransactionResponse{}, access_token_err
		}

		req.Header.Set("Authorization", access_token)

		transaction_response, err = client.Do(req)
	}

	if err != nil {
		return []TransactionResponse{}, err
	}

	if transaction_response.StatusCode == http.StatusUnauthorized {
		return []TransactionResponse{}, fmt.Errorf("Unable to rewnew access token")
	}

	var transaction_list_response []TransactionResponse

	decoding_error := json.NewDecoder(transaction_response.Body).Decode(&transaction_list_response)

	if decoding_error != nil {
		return []TransactionResponse{}, decoding_error
	}

	return transaction_list_response, nil
}

func CreateTransaction(
	description string,
	category string,
	value float32,
) (TransactionResponse, error) {
	transaction_create_request, marshal_error := json.Marshal(TransactionCreateRequest{
		description,
		category,
		value,
	})

	if marshal_error != nil {
		return TransactionResponse{}, marshal_error
	}

	req, err := http.NewRequest(
		"POST",
		stationapiconfig.GetStationApiUrl()+"/transaction",
		bytes.NewBuffer(transaction_create_request))

	if err != nil {
		return TransactionResponse{}, nil
	}

	access_token, err := tokenstorage.ReadAccessToken()

	if err != nil {
		return TransactionResponse{}, err
	}

	req.Header.Set("Authorization", access_token)
	req.Header.Set("Content-Type", application_json_media_type)

	client := &http.Client{}

	transaction_create_response, err := client.Do(req)

	if err != nil {
		return TransactionResponse{}, err
	}

	defer transaction_create_response.Body.Close()

	if transaction_create_response.StatusCode == http.StatusUnauthorized {
		access_token, access_token_err := refreshAccessToken()

		if access_token_err != nil {
			return TransactionResponse{}, access_token_err
		}

		req.Header.Set("Authorization", access_token)

		transaction_create_response, err = client.Do(req)
	}

	if err != nil {
		return TransactionResponse{}, err
	}

	if transaction_create_response.StatusCode == http.StatusUnauthorized {
		return TransactionResponse{}, fmt.Errorf("Unable to rewnew access token")
	}

	var transaction_response TransactionResponse

	decoding_error := json.NewDecoder(transaction_create_response.Body).Decode(&transaction_response)

	if decoding_error != nil {
		return TransactionResponse{}, decoding_error
	}

	return transaction_response, nil
}

func RefreshTokens() {
	refreshRefreshToken()
	refreshAccessToken()
}

func refreshRefreshToken() (string, error) {
	refresh_token, err := tokenstorage.ReadRefreshToken()

	if err != nil {
		return "", err
	}

	renew_request_json, err := json.Marshal(RenewTokenRequest{refresh_token})

	renew_response, err := http.Post(
		stationapiconfig.GetStationApiUrl()+"/token/renew",
		application_json_media_type,
		bytes.NewBuffer(renew_request_json))

	if err != nil {
		return "", err
	}

	defer renew_response.Body.Close()

	var renew_token_response RenewTokenResponse

	decoding_error := json.NewDecoder(renew_response.Body).Decode(&renew_token_response)

	if decoding_error != nil {
		// TODO: Log the error message or the response body as string
		return "", fmt.Errorf("Error decoing the response from renewing access token")
	}

	tokenstorage.StoreRefreshToken(renew_token_response.Token)

	return renew_token_response.Token, nil
}

func refreshAccessToken() (string, error) {
	refresh_token, err := tokenstorage.ReadRefreshToken()

	if err != nil {
		return "", err
	}

	renew_request_json, err := json.Marshal(RenewTokenRequest{refresh_token})

	renew_response, err := http.Post(
		stationapiconfig.GetStationApiUrl()+"/token/access",
		application_json_media_type,
		bytes.NewBuffer(renew_request_json))

	if err != nil {
		return "", err
	}

	defer renew_response.Body.Close()

	var renew_token_response RenewTokenResponse

	decoding_error := json.NewDecoder(renew_response.Body).Decode(&renew_token_response)

	if decoding_error != nil {
		// TODO: Log the error message or the response body as string
		return "", fmt.Errorf("Error decoing the response from renewing access token")
	}

	tokenstorage.StoreAccessToken(renew_token_response.Token)

	return renew_token_response.Token, nil
}
