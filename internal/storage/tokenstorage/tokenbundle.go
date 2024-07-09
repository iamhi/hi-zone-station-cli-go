package tokenstorage

type tokenBundle struct {
	AccessToken string `json:"accessToken"`

	RefreshToken string `json:"refreshToken"`

	RefreshedTime string `json:"refreshedTime"`
}
