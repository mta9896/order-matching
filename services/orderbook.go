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
		if ob.SellPricesHeap.Len() > 0 {
			cheapestSell := ob.SellPricesHeap[0]

			if order.Price >= cheapestSell {
				//check the map
				if sellOrders, exists := ob.SellOrders[order.Price]; exists {
					for _, sellOrder := range sellOrders {
						if sellOrder.Amount == order.Amount {
							matchedOrders = append(matchedOrders, sellOrder)
						}
					}

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
				// delete(ob.orderMap, bestSell.ID) delete from the existing order uuids map?
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

	} else { // sell action
		if ob.BuyPricesHeap.Len() > 0 {
			highestBid := ob.BuyPricesHeap[0] 
			if highestBid >= order.Price {
				//check the map
				if buyOrders, exists := ob.BuyOrders[order.Price]; exists {
					for _, buyOrder := range buyOrders {
						if buyOrder.Amount == order.Amount {
							matchedOrders = append(matchedOrders, buyOrder)
						}
					}

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
				// delete(ob.orderMap, bestSell.ID) delete from the existing order uuids map?
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
	}

	return matchedOrders
}

