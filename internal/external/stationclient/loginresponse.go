package stationclient

type LoginResponse struct {
	AccessToken string `json:"accessToken"`

	RefreshToken string `json:"refreshToken"`
}
