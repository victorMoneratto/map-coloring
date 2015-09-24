package heap

import (
	"container/heap"

	"github.com/victormoneratto/map-coloring/graph"

	"testing"
)

func TestHeap(t *testing.T) {
	h := NewNodeHeap(10, false)
	for i := 0; i < 10; i++ {
		h.Items = append(h.Items, &NodeItem{
			Node: &graph.Node{NumTaken: i},
		})
	}
	heap.Init(&h)

	for h.Len() > 2 {
		a, b := heap.Pop(&h), heap.Pop(&h)
		if a.(*NodeItem).Node.NumTaken < b.(*NodeItem).Node.NumTaken {
			t.Error("a should be > b")
		}
	}
}
