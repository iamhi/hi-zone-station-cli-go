package stationclient

type TransactionResponse struct {
	Uuid string `json:"uuid"`

	Description string `json:"description"`

	Category string `json:"category"`

	Value float32 `json:"value"`

	CreatedAt string `json:"createdAt"`
}
