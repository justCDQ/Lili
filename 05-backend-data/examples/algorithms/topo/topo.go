package topo

import (
	"container/heap"
	"errors"
	"fmt"
)

type stringHeap []string

func (h stringHeap) Len() int           { return len(h) }
func (h stringHeap) Less(i, j int) bool { return h[i] < h[j] }
func (h stringHeap) Swap(i, j int)      { h[i], h[j] = h[j], h[i] }
func (h *stringHeap) Push(x any)        { *h = append(*h, x.(string)) }
func (h *stringHeap) Pop() any {
	old := *h
	x := old[len(old)-1]
	*h = old[:len(old)-1]
	return x
}

type Edge struct{ Before, After string }

func Sort(nodes []string, edges []Edge) ([]string, error) {
	next := make(map[string][]string, len(nodes))
	indegree := make(map[string]int, len(nodes))
	for _, node := range nodes {
		if node == "" {
			return nil, errors.New("node id is empty")
		}
		if _, exists := next[node]; exists {
			return nil, fmt.Errorf("duplicate node %q", node)
		}
		next[node] = nil
		indegree[node] = 0
	}
	seenEdge := make(map[Edge]struct{}, len(edges))
	for _, edge := range edges {
		if _, ok := next[edge.Before]; !ok {
			return nil, fmt.Errorf("unknown node %q", edge.Before)
		}
		if _, ok := next[edge.After]; !ok {
			return nil, fmt.Errorf("unknown node %q", edge.After)
		}
		if _, duplicate := seenEdge[edge]; duplicate {
			return nil, fmt.Errorf("duplicate edge %+v", edge)
		}
		seenEdge[edge] = struct{}{}
		next[edge.Before] = append(next[edge.Before], edge.After)
		indegree[edge.After]++
	}
	ready := &stringHeap{}
	heap.Init(ready)
	for node, degree := range indegree {
		if degree == 0 {
			heap.Push(ready, node)
		}
	}
	order := make([]string, 0, len(nodes))
	for ready.Len() > 0 {
		node := heap.Pop(ready).(string)
		order = append(order, node)
		for _, successor := range next[node] {
			indegree[successor]--
			if indegree[successor] == 0 {
				heap.Push(ready, successor)
			}
		}
	}
	if len(order) != len(nodes) {
		return nil, errors.New("dependency cycle")
	}
	return order, nil
}
