package tokenstorage

import (
	"encoding/json"
	"log"
	"os"
	"time"
)

const TOKEN_FILE_NAME = "tokenbundle.json"

func StoreTokens(access_token string, refresh_token string) {
	token_bundle := tokenBundle{
		access_token,
		refresh_token,
		time.Now().Format(time.RFC3339),
	}

	bytes, err := json.Marshal(token_bundle)

	if err == nil {
		file, err := os.Create(TOKEN_FILE_NAME)

		defer file.Close()

		if err == nil {
			file.Write(bytes)
		}
	}
}

func StoreRefreshToken(refresh_token string) error {
	token_bundle, err := ReadTokenBundle()

	if err != nil {
		return err
	}

	StoreTokens(token_bundle.AccessToken, refresh_token)

	return nil
}

func StoreAccessToken(access_token string) error {
	token_bundle, err := ReadTokenBundle()

	if err != nil {
		return err
	}

	StoreTokens(access_token, token_bundle.RefreshToken)

	return nil
}

func ReadAccessToken() (string, error) {
	token_bundle, err := ReadTokenBundle()

	if err != nil {
		return "", err
	}

	return token_bundle.AccessToken, nil
}

func ReadRefreshToken() (string, error) {
	token_bundle, err := ReadTokenBundle()

	if err != nil {
		return "", err
	}

	return token_bundle.RefreshToken, nil
}

func ReadTokenBundle() (tokenBundle, error) {
	content, err := os.ReadFile(TOKEN_FILE_NAME)

	if err != nil {
		log.Fatal("Error while reading token bundle")

		return tokenBundle{}, err
	}

	var token_bundle tokenBundle

	err = json.Unmarshal(content, &token_bundle)

	if err != nil {
		log.Fatal("Error while unmarshaling token bundle content")

		return tokenBundle{}, err
	}

	return token_bundle, nil
}
