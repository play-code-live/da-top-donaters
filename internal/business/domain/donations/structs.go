package donations

type MaxDonationHeap []DonationItem

func (h *MaxDonationHeap) Less(i, j int) bool { return (*h)[i].Amount > (*h)[j].Amount }
func (h *MaxDonationHeap) Swap(i, j int)      { (*h)[i], (*h)[j] = (*h)[j], (*h)[i] }
func (h *MaxDonationHeap) Len() int           { return len(*h) }
func (h *MaxDonationHeap) Push(v interface{}) { *h = append(*h, v.(DonationItem)) }
func (h *MaxDonationHeap) Pop() (v interface{}) {
	*h, v = (*h)[:h.Len()-1], (*h)[h.Len()-1]
	return v
}
