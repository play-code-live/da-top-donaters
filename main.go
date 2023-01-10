package main

import (
	"container/heap"
	"fmt"

	donationClient "He110/donation-report-manager/src/donation-client"
	nameAdapter "He110/donation-report-manager/src/name-adapter"
)

const topAmount = 10

func main() {
	client := donationClient.NewClient(
		"10386",
		"x9Auz25j1PULNJXl4FScvSnnEKzJIf95oXXYPgvq",
		"eyJ0eXAiOiJKV1QiLCJhbGciOiJSUzI1NiJ9.eyJhdWQiOiIxMDM4NiIsImp0aSI6Ijk2NmY5MjM0NDc1YTZmN2E3Zjc3MDEwODVhM2RiMGQ4YmJlNDE5YmRiOTZmZjc5ZjcwNTNjMzk1YTY1YWE4ODliZTFlY2E5MWE5NzQ1MDdhIiwiaWF0IjoxNjczMjg0MTE5LjMzNDMsIm5iZiI6MTY3MzI4NDExOS4zMzQzLCJleHAiOjIzMDQ0MzYxMTkuMjc4Mywic3ViIjoiNTM4NTIzOCIsInNjb3BlcyI6WyJvYXV0aC1kb25hdGlvbi1pbmRleCJdfQ.Ee8pYXqiSCB1Kz32jEnYnvdyGZFZ5cQVirvApuJC-l5k9afzeHY7p7spAqFSnFISSk3v8HiukB2qRl9kRzi49cs4NtV_XDJej3tL1q1NHWxUvOh1hDli_WWCn6SA0q6CNlmFMgqle7Ir2igamkQ1yW7OcHzuzNIJ9uKDx3-zQqxQuNPyPhlWEuOf9_i37MgNnB7mlalaICvP4_CPEXVWiynHjYgnNZ8jsvpBQs8GZoKq1nwVpfCx58DGt-7cnlL9UaGvYInW5UrEzfkiYTfdvKkrJBOR35qhQkjmhLpfzrv9ycebV1uNvh6PvhbSUXHaSxEootitpXd9LsbII2i4geOobGUP9zWmmcjoqGaht16eWC2kyJy-JxiwkdRCQ_pIxe6FdaMlDZ800LYCFUA7RZ1gX7Lc8xkrLIaIN_Ir8GNmKsRtooDGZ3mlrPFuJDtooi0zTLkNK5ap2ZfpoH1NMS3Y1E2WVnqUlzTUfWxW00bfqIDtZaukcUCsTrBrCLKKuan43g2tC2mKXeSRMBwzmxp3Q8nGfYk8B6yUlAXgOwWfZ3ijjMm4eEluTLNnhvisx2DEI2epB_iqjm1hRUNY38L-goOjOeUWLM90y1Z719Qlaz9aQ-k5u-5JQw3GyDv4ZI4twxbhRspvibykc2FuRCrbTIfkgX9Kwqh_K4ypFAs",
	)
	donations, _ := client.GetAllDonations()
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
