package topk

import (
	"container/heap"
	"errors"
	"sort"
)

type Item struct {
	ID    string
	Score int
}
type minHeap []Item

func (h minHeap) Len() int { return len(h) }
func (h minHeap) Less(i, j int) bool {
	if h[i].Score != h[j].Score {
		return h[i].Score < h[j].Score
	}
	return h[i].ID > h[j].ID
}
func (h minHeap) Swap(i, j int)   { h[i], h[j] = h[j], h[i] }
func (h *minHeap) Push(value any) { *h = append(*h, value.(Item)) }
func (h *minHeap) Pop() any {
	old := *h
	last := old[len(old)-1]
	old[len(old)-1] = Item{}
	*h = old[:len(old)-1]
	return last
}
func better(a, b Item) bool {
	return a.Score > b.Score || (a.Score == b.Score && a.ID < b.ID)
}
func Select(items []Item, k int) ([]Item, error) {
	if k < 0 {
		return nil, errors.New("k must be non-negative")
	}
	if k == 0 {
		return []Item{}, nil
	}
	selected := make(minHeap, 0, min(k, len(items)))
	heap.Init(&selected)
	for _, item := range items {
		if item.ID == "" {
			return nil, errors.New("item id is empty")
		}
		if selected.Len() < k {
			heap.Push(&selected, item)
		} else if better(item, selected[0]) {
			selected[0] = item
			heap.Fix(&selected, 0)
		}
	}
	result := append([]Item(nil), selected...)
	sort.Slice(result, func(i, j int) bool { return better(result[i], result[j]) })
	return result, nil
}
