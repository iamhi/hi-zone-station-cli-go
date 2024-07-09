package main

import "github.com/iamhi/station-cli-go/internal/services/executionservice"

func main() {
	// TODO: Need to refresh refresh token
	// stationclient.Login("admin1", "admin1")

	executionservice.ExecuteCommand()
}
