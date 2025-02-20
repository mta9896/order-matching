package models

// We use a min heap for sell order prices
type SellHeap []float64

func (sh SellHeap) Len() int {
	return len(sh)
} 

func (sh SellHeap) Less(i int, j int) bool {
	return sh[i] < sh[j] // We want the lowest sell price to be on top
}

func (sh SellHeap) Swap(i int, j int) {
	sh[i], sh[j] = sh[j], sh[i]
}

func (sh *SellHeap) Push(element any) {
	price := element.(float64)
	*sh = append(*sh, price)
}

func (sh *SellHeap) Pop() (any) {
	old := *sh
	n := len(old)
	item := old[n - 1]
	*sh = old[0 : n-1]

	return item
} 