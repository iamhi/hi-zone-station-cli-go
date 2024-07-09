package tokenservice

import (
	"time"

	"github.com/iamhi/station-cli-go/internal/external/stationclient"
	"github.com/iamhi/station-cli-go/internal/storage/tokenstorage"
)

func CheckTokens() {
	token_bundle, err := tokenstorage.ReadTokenBundle()

	if err != nil {
		return
	}

	refreshed_time, err := time.Parse(time.RFC3339, token_bundle.RefreshedTime)

	if err != nil {
		if refreshed_time.Add(time.Hour * 48).Before(time.Now()) {
			stationclient.RefreshTokens()
		}
	}

}
