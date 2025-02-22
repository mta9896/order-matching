package models

type OrderType string

const Buy OrderType = "BUY"
const Sell OrderType = "SELL"

type Order struct {
	ID string `json:"uuid" binding:"required,uuid4" example:"550e8400-e29b-41d4-a716-646655440000"`
	Action OrderType `json:"action" binding:"required,oneof=BUY SELL"`
	Price float64 `json:"price" binding:"required" example:"100.0"`
	Amount float64 `json:"amount" binding:"required" example:"10.0"`
}
