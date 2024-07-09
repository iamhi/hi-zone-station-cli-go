package stationclient

type TransactionCreateRequest struct {
	Description string `json:"description"`

	Category string `json:"category"`

	Value float32 `json:"value"`
}
