package models

import (
)

type OrderType string

const Buy OrderType = "BUY"
const Sell OrderType = "SELL"

type Order struct {
	ID string `json:"uuid" binding:"required,uuid4"`
	Action OrderType `json:"action" binding:"required,oneof=BUY SELL"`
	Price float64 `json:"price" binding:"required"`
	Amount float64 `json:"amount" binding:"required"`
	Index int `json:"index"`
}
