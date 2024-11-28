package models

type Pair struct {
	ID        int `json:"id"`
	BuyLotID  int `json:"buy_lot_id"`
	SellLotID int `json:"sell_lot_id"`
}
