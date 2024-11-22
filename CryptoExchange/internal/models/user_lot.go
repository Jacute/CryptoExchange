package models

type UserLot struct {
	LotID    int     `json:"lot_id"`
	Quantity float64 `json:"quantity"`
}
