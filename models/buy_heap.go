package models

// We use a max heap for buy order prices
type BuyHeap []float64

func (bh BuyHeap) Len() int {
	return len(bh)
} 

func (bh BuyHeap) Less(i int, j int) bool {
	return bh[i] > bh[j] // We want the highest buy price to be on top
}

func (bh BuyHeap) Swap(i int, j int) {
	bh[i], bh[j] = bh[j], bh[i]
}

func (bh *BuyHeap) Push(element any) {
	price := element.(float64)
	*bh = append(*bh, price)
}

func (bh *BuyHeap) Pop() (any) {
	old := *bh
	n := len(old)
	item := old[n - 1]
	*bh = old[0 : n-1]

	return item
} 