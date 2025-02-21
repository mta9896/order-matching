package services

import (
	"container/heap"
	"order-matching/models"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPlaceSellOrder_WhenBuyHeapIsEmpty(t *testing.T) {
	t.Parallel()
	ob := NewOrderBook()

	sellOrder := models.Order{
		ID:     "550e8400-e29b-41d4-a716-446655440000",
        Action: models.Sell,
        Price:  100.0,
        Amount: 2.0,
	}

	matchedOrder := ob.PlaceOrder(&sellOrder)

	assert.Equal(t, 0, len(matchedOrder))
	assert.Equal(t, 1, len(ob.SellOrders[sellOrder.Price]))
	assert.Equal(t, 1, len(ob.SellPricesHeap))
}

func TestPlaceSellOrder_WhenNoBuyOrderIsMatched(t *testing.T) {
	t.Parallel()
	ob := NewOrderBook()

	ob.BuyOrders = map[float64][]models.Order{
		80.0: {
			{
				ID:     "550e8400-e29b-41d4-a716-666655442000",
        		Action: models.Buy,
       			Price:  80.0,
       			Amount: 2.0,
			},
		},
	}

	heap.Push(&ob.BuyPricesHeap, 80.0)

	sellOrder := models.Order{
		ID:     "550e8400-e29b-41d4-a716-446655440000",
        Action: models.Sell,
        Price:  100.0,
        Amount: 2.0,
	}

	matchedOrder := ob.PlaceOrder(&sellOrder)

	assert.Equal(t, 0, len(matchedOrder))
	assert.Equal(t, 1, len(ob.SellOrders[sellOrder.Price]))
	assert.Equal(t, 1, len(ob.SellPricesHeap))
}

func TestPlaceSellOrder_WhenOneBuyOrderIsMatched(t *testing.T) {
	t.Parallel()
	ob := NewOrderBook()

	ob.BuyOrders = map[float64][]models.Order{
		80.0: {
			{
				ID:     "550e8400-e29b-41d4-a716-666655442000",
        		Action: models.Buy,
       			Price:  80.0,
       			Amount: 2.0,
			},
		},
		100.0: {
			{
				ID:     "550e8400-e29b-41d4-a716-77755442000",
        		Action: models.Buy,
       			Price:  100.0,
       			Amount: 2.0,
			},
		},
	}

	heap.Push(&ob.BuyPricesHeap, 80.0)
	heap.Push(&ob.BuyPricesHeap, 100.0)

	sellOrder := models.Order{
		ID:     "550e8400-e29b-41d4-a716-446655440000",
        Action: models.Sell,
        Price:  100.0,
        Amount: 2.0,
	}

	matchedOrders := ob.PlaceOrder(&sellOrder)
	
	assert.Equal(t, 1, len(matchedOrders))
	expectedMatchedOrders := []models.Order {
		{
			ID:     "550e8400-e29b-41d4-a716-77755442000",
			Action: models.Buy,
			   Price:  100.0,
			   Amount: 2.0,
		},
	}
	assert.Equal(t, expectedMatchedOrders, matchedOrders)

	assert.Equal(t, 0, len(ob.SellOrders[sellOrder.Price]))
	assert.Equal(t, 0, len(ob.SellPricesHeap))

	assert.Equal(t, 0, len(ob.BuyOrders[sellOrder.Price]))
	assert.Equal(t, 1, len(ob.BuyPricesHeap))
}

func TestPlaceSellOrder_WhenTwoBuyOrdersAreMatched(t *testing.T) {
	t.Parallel()
	ob := NewOrderBook()

	ob.BuyOrders = map[float64][]models.Order{
		80.0: {
			{
				ID:     "550e8400-e29b-41d4-a716-666655442000",
        		Action: models.Buy,
       			Price:  80.0,
       			Amount: 2.0,
			},
		},
		100.0: {
			{
				ID:     "550e8400-e29b-41d4-a716-77755442000",
        		Action: models.Buy,
       			Price:  100.0,
       			Amount: 2.0,
			},
			{
				ID:     "550e8400-e29b-41d4-a716-77755442001",
        		Action: models.Buy,
       			Price:  100.0,
       			Amount: 2.0,
			},
			{
				ID:     "550e8400-e29b-41d4-a716-77755442002",
        		Action: models.Buy,
       			Price:  100.0,
       			Amount: 3.0,
			},
		},
	}

	heap.Push(&ob.BuyPricesHeap, 80.0)
	heap.Push(&ob.BuyPricesHeap, 100.0)

	sellOrder := models.Order{
		ID:     "550e8400-e29b-41d4-a716-446655440000",
        Action: models.Sell,
        Price:  100.0,
        Amount: 2.0,
	}

	matchedOrders := ob.PlaceOrder(&sellOrder)
	
	assert.Equal(t, 2, len(matchedOrders))
	expectedMatchedOrders := []models.Order {
		{
			ID:     "550e8400-e29b-41d4-a716-77755442000",
			Action: models.Buy,
			   Price:  100.0,
			   Amount: 2.0,
		},
		{
			ID:     "550e8400-e29b-41d4-a716-77755442001",
			Action: models.Buy,
			   Price:  100.0,
			   Amount: 2.0,
		},
	}
	assert.Equal(t, expectedMatchedOrders, matchedOrders)

	assert.Equal(t, 0, len(ob.SellOrders[sellOrder.Price]))
	assert.Equal(t, 0, len(ob.SellPricesHeap))
}

func TestPlaceBuyOrder_WhenSellHeapIsEmpty(t *testing.T) {
	t.Parallel()
	ob := NewOrderBook()

	buyOrder := models.Order{
		ID:     "550e8400-e29b-41d4-a716-446655440000",
        Action: models.Buy,
        Price:  100.0,
        Amount: 2.0,
	}

	matchedOrder := ob.PlaceOrder(&buyOrder)

	assert.Equal(t, 0, len(matchedOrder))
	assert.Equal(t, 1, len(ob.BuyOrders[buyOrder.Price]))
	assert.Equal(t, 1, len(ob.BuyOrders))
}

func TestPlaceBuyOrder_WhenNoSellOrderIsMatched(t *testing.T) {
	t.Parallel()
	ob := NewOrderBook()

	ob.SellOrders = map[float64][]models.Order{
		120.0: {
			{
				ID:     "550e8400-e29b-41d4-a716-666655442000",
        		Action: models.Sell,
       			Price:  120.0,
       			Amount: 2.0,
			},
		},
	}

	heap.Push(&ob.SellPricesHeap, 120.0)

	buyOrder := models.Order{
		ID:     "550e8400-e29b-41d4-a716-446655440000",
        Action: models.Buy,
        Price:  100.0,
        Amount: 2.0,
	}

	matchedOrder := ob.PlaceOrder(&buyOrder)

	assert.Equal(t, 0, len(matchedOrder))
	assert.Equal(t, 1, len(ob.BuyOrders[buyOrder.Price]))
	assert.Equal(t, 1, len(ob.BuyPricesHeap))
}

func TestPlaceBuyOrder_WhenOneSellOrderIsMatched(t *testing.T) {
	t.Parallel()
	ob := NewOrderBook()

	ob.SellOrders = map[float64][]models.Order{
		120.0: {
			{
				ID:     "550e8400-e29b-41d4-a716-666655442000",
        		Action: models.Sell,
       			Price:  120.0,
       			Amount: 2.0,
			},
		},
		100.0: {
			{
				ID:     "550e8400-e29b-41d4-a716-77755442000",
        		Action: models.Sell,
       			Price:  100.0,
       			Amount: 2.0,
			},
		},
	}

	heap.Push(&ob.SellPricesHeap, 120.0)
	heap.Push(&ob.SellPricesHeap, 100.0)

	buyOrder := models.Order{
		ID:     "550e8400-e29b-41d4-a716-446655440000",
        Action: models.Buy,
        Price:  100.0,
        Amount: 2.0,
	}

	matchedOrders := ob.PlaceOrder(&buyOrder)
	
	assert.Equal(t, 1, len(matchedOrders))
	expectedMatchedOrders := []models.Order {
		{
			ID:     "550e8400-e29b-41d4-a716-77755442000",
			Action: models.Sell,
			Price:  100.0,
			Amount: 2.0,
		},
	}
	assert.Equal(t, expectedMatchedOrders, matchedOrders)

	assert.Equal(t, 0, len(ob.BuyOrders[buyOrder.Price]))
	assert.Equal(t, 0, len(ob.BuyPricesHeap))

	assert.Equal(t, 0, len(ob.SellOrders[buyOrder.Price]))
	assert.Equal(t, 1, len(ob.SellPricesHeap))
}

func TestPlaceBuyOrder_WhenTwoSellOrdersAreMatched(t *testing.T) {
	t.Parallel()
	ob := NewOrderBook()

	ob.SellOrders = map[float64][]models.Order{
		120.0: {
			{
				ID:     "550e8400-e29b-41d4-a716-666655442000",
        		Action: models.Sell,
       			Price:  120.0,
       			Amount: 2.0,
			},
		},
		100.0: {
			{
				ID:     "550e8400-e29b-41d4-a716-77755442000",
        		Action: models.Sell,
       			Price:  100.0,
       			Amount: 2.0,
			},
			{
				ID:     "550e8400-e29b-41d4-a716-77755442001",
        		Action: models.Sell,
       			Price:  100.0,
       			Amount: 2.0,
			},
			{
				ID:     "550e8400-e29b-41d4-a716-77755442002",
        		Action: models.Sell,
       			Price:  100.0,
       			Amount: 3.0,
			},
		},
	}

	heap.Push(&ob.SellPricesHeap, 120.0)
	heap.Push(&ob.SellPricesHeap, 100.0)

	buyOrder := models.Order{
		ID:     "550e8400-e29b-41d4-a716-446655440000",
        Action: models.Buy,
        Price:  100.0,
        Amount: 2.0,
	}

	matchedOrders := ob.PlaceOrder(&buyOrder)
	
	assert.Equal(t, 2, len(matchedOrders))
	expectedMatchedOrders := []models.Order {
		{
			ID:     "550e8400-e29b-41d4-a716-77755442000",
			Action: models.Sell,
			   Price:  100.0,
			   Amount: 2.0,
		},
		{
			ID:     "550e8400-e29b-41d4-a716-77755442001",
			Action: models.Sell,
			   Price:  100.0,
			   Amount: 2.0,
		},
	}
	assert.Equal(t, expectedMatchedOrders, matchedOrders)

	assert.Equal(t, 0, len(ob.BuyOrders[buyOrder.Price]))
	assert.Equal(t, 0, len(ob.BuyPricesHeap))
}

func TestGetOrderBook(t *testing.T) {
	t.Parallel()
	ob := NewOrderBook()
	ob.SellOrders = map[float64][]models.Order{
		120.0: {
			{
				ID:     "550e8400-e29b-41d4-a716-666655442000",
        		Action: models.Sell,
       			Price:  120.0,
       			Amount: 2.0,
			},
		},
		100.0: {
			{
				ID:     "550e8400-e29b-41d4-a716-77755442000",
        		Action: models.Sell,
       			Price:  100.0,
       			Amount: 2.0,
			},
			{
				ID:     "550e8400-e29b-41d4-a716-77755442001",
        		Action: models.Sell,
       			Price:  100.0,
       			Amount: 2.0,
			},
		},
	}

	heap.Push(&ob.BuyPricesHeap, 80.0)
	heap.Push(&ob.BuyPricesHeap, 100.0)

	ob.BuyOrders = map[float64][]models.Order{
		80.0: {
			{
				ID:     "550e8400-e29b-41d4-a716-666655442000",
        		Action: models.Buy,
       			Price:  80.0,
       			Amount: 2.0,
			},
		},
		100.0: {
			{
				ID:     "550e8400-e29b-41d4-a716-77755442000",
        		Action: models.Buy,
       			Price:  100.0,
       			Amount: 2.0,
			},
		},
	}

	heap.Push(&ob.SellPricesHeap, 120.0)
	heap.Push(&ob.SellPricesHeap, 100.0)

	orderbook := ob.GetOrderBook(2)

	expected := []models.OrderBookEntry{
		{
			Price: 120.0,
			Type: models.Sell,
			Liquidity: 2.0,
		},
		{
			Price: 100.0,
			Type: models.Sell,
			Liquidity: 4.0,
		},
		{
			Price: 100.0,
			Type: models.Buy,
			Liquidity: 2.0,
		},
		{
			Price: 80.0,
			Type: models.Buy,
			Liquidity: 2.0,
		},
	}

	assert.Equal(t, expected, orderbook)
}