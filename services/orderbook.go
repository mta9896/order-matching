package services

import (
	"container/heap"
	"order-matching/models"
)

type OrderBook struct {
	BuyPricesHeap models.BuyHeap
	SellPricesHeap models.SellHeap
	BuyOrders map[float64][]models.Order
	SellOrders map[float64][]models.Order
}

func NewOrderBook() *OrderBook {
	orderBook := &OrderBook{
		BuyPricesHeap: models.BuyHeap{},
		SellPricesHeap: models.SellHeap{},
		BuyOrders: make(map[float64][]models.Order),
		SellOrders: make(map[float64][]models.Order),
	}

	heap.Init(&orderBook.BuyPricesHeap)
	heap.Init(&orderBook.SellPricesHeap)

	return orderBook
}

func (ob *OrderBook) PlaceOrder(order *models.Order) (matchedOrders []models.Order){
	if order.Action == models.Buy {
		matchedOrders = ob.handleBuyAction(order)
	} else { // sell action
		matchedOrders = ob.handleSellAction(order)
	}

	return matchedOrders
}

func (ob *OrderBook) GetOrderBook(limit int) []models.OrderBookEntry {
	var sellOrders, buyOrders []models.OrderBookEntry

	startIndex := limit
	if limit >= ob.SellPricesHeap.Len() {
		startIndex = ob.SellPricesHeap.Len() - 1
	}

	for i := startIndex; i >= 0; i-- {
		price := ob.SellPricesHeap[i]
		liquidity := 0.0
		for _, order := range ob.SellOrders[price] {
			liquidity += order.Amount
		}
		
		sellOrders = append(sellOrders, models.OrderBookEntry{
			Price: price,
			Liquidity: liquidity,
			Type: models.Sell,
		})
	}

	for i := 0; i < ob.BuyPricesHeap.Len() && i < limit; i++ {
		price := ob.BuyPricesHeap[i]
		liquidity := 0.0
		for _, order := range ob.BuyOrders[price] {
			liquidity += order.Amount
		}
		
		buyOrders = append(buyOrders, models.OrderBookEntry{
			Price: price,
			Liquidity: liquidity,
			Type: models.Buy,
		})
	}

	return append(sellOrders, buyOrders...)
}

func (ob *OrderBook) GetOrderList(page int, pageSize int) []models.Order {
	var allOrders []models.Order
	for _, orders := range ob.BuyOrders {
		allOrders = append(allOrders, orders...)
	}
	for _, orders := range ob.SellOrders {
		allOrders = append(allOrders, orders...)
	}

	totalOrders := len(allOrders)
	
	start := (page - 1) * pageSize
	end := start + pageSize
	if end > totalOrders {
		end = totalOrders
	}

	return allOrders[start:end]
}

func (ob *OrderBook) handleBuyAction(order *models.Order) (matchedOrders []models.Order) {
	if ob.SellPricesHeap.Len() > 0 {
		cheapestSell := ob.SellPricesHeap[0]

		if order.Price >= cheapestSell {
			//check the map
			if sellOrders, exists := ob.SellOrders[order.Price]; exists {
				matchedOrders = findMatchingOrdersByAmount(sellOrders, order)

				for i, sellOrder := range sellOrders {
					if sellOrder.Amount == order.Amount {
						ob.SellOrders[sellOrder.Price] = append(sellOrders[:i], sellOrders[i+1:]...)
						if len(ob.SellOrders[order.Price]) == 0 {
							delete(ob.SellOrders, order.Price)
							heap.Pop(&ob.SellPricesHeap)
						}
						break
					}
				}
			}
			return
		} 
	}

	// no match found or the SellPricesHeap were empty
	heap.Push(&ob.BuyPricesHeap, order.Price)
	// insert the order in map
	if orders, exists := ob.BuyOrders[order.Price]; exists {
		ob.BuyOrders[order.Price] = append(orders, *order) // makes a copy of the order and puts it in the heap
	} else {
		ob.BuyOrders[order.Price] = []models.Order{*order}
	}

	return
}

func (ob *OrderBook) handleSellAction(order *models.Order) (matchedOrders []models.Order) {
	if ob.BuyPricesHeap.Len() > 0 {
		highestBid := ob.BuyPricesHeap[0] 
		if highestBid >= order.Price {
			//check the map
			if buyOrders, exists := ob.BuyOrders[order.Price]; exists {
				matchedOrders = findMatchingOrdersByAmount(buyOrders, order)

				for i, buyOrder := range buyOrders {
					if buyOrder.Amount == order.Amount {
						ob.BuyOrders[buyOrder.Price] = append(buyOrders[:i], buyOrders[i+1:]...)
						if len(ob.BuyOrders[order.Price]) == 0 {
							delete(ob.BuyOrders, order.Price)
							heap.Pop(&ob.BuyPricesHeap)
						}
						break
					}
				}
			}
			return
		}
	}
	// no match found or BuyPricesHeap were empty
	heap.Push(&ob.SellPricesHeap, order.Price)
	if orders, exists := ob.SellOrders[order.Price]; exists {
		ob.SellOrders[order.Price] = append(orders, *order) // makes a copy of the order and puts it in the heap
	} else {
		ob.SellOrders[order.Price] = []models.Order{*order}
	}

	return
}

func findMatchingOrdersByAmount(orders []models.Order, order *models.Order) []models.Order {
	matchedOrders := []models.Order{}
	for _, o := range orders {
		if o.Amount == order.Amount {
			matchedOrders = append(matchedOrders, o)
		}
	}

	return matchedOrders
}

