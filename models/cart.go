package models

type Cart struct {
	Model
	BuyerID string `json:"buyer_id"`
}
