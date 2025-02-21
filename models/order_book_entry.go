package models

type OrderBookEntry struct {
	Price float64 `json:"price"`
	Liquidity float64 `json:"liquidity"`
	Type OrderType `json:"type"`
}