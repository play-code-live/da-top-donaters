package main

import (
	"container/heap"
	"fmt"

	donationClient "He110/donation-report-manager/src/donation-client"
	nameAdapter "He110/donation-report-manager/src/name-adapter"
)

const topAmount = 10

func main() {
	client, err := donationClient.NewClient(
		"10386",
		"x9Auz25j1PULNJXl4FScvSnnEKzJIf95oXXYPgvq",
	)
	if err != nil {
		panic(err)
	}
	donations, err := client.GetAllDonations()
	if err != nil {
		panic(err)
	}
	adapter := nameAdapter.NewNameAdapter()

	top := map[string]float64{}
	for _, donation := range donations {
		name := adapter.Perform(donation.Username)
		if adapter.ShouldBeSkipped(name) {
			continue
		}
		top[name] += donation.Amount
	}

	h := MaxHeap{}
	for name, amount := range top {
		heap.Push(&h, TopItem{
			Amount: amount,
			Name:   name,
		})
	}

	finalText := ""
	for i := 0; i < topAmount && h.Len() > 0; i++ {
		item := heap.Pop(&h).(TopItem)
		finalText += fmt.Sprintf("%d. **%s**: %d₽\n", i+1, item.Name, int(item.Amount))
	}

	fmt.Println(finalText)
}

type TopItem struct {
	Amount float64
	Name   string
}

type MaxHeap []TopItem

func (h *MaxHeap) Less(i, j int) bool { return (*h)[i].Amount > (*h)[j].Amount }
func (h *MaxHeap) Swap(i, j int)      { (*h)[i], (*h)[j] = (*h)[j], (*h)[i] }
func (h *MaxHeap) Len() int           { return len(*h) }
func (h *MaxHeap) Push(v interface{}) { *h = append(*h, v.(TopItem)) }
func (h *MaxHeap) Pop() (v interface{}) {
	*h, v = (*h)[:h.Len()-1], (*h)[h.Len()-1]
	return v
}
