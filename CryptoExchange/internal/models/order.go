package models

type Order struct {
	ID       int     `json:"id"`
	UserID   int     `json:"user_id"`
	PairID   int     `json:"pair_id"`
	Quantity float64 `json:"quantity"`
	Price    float64 `json:"price"`
	Type     string  `json:"type"`
	Closed   string  `json:"closed"`
}
