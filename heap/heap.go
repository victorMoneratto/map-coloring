package heap

import (
	"container/heap"

	"github.com/victormoneratto/map-coloring/graph"
)

type Heap struct {
	Items     []Heaper
	UseDegree bool
}

type Heaper interface {
	Less(h Heap, other Heaper) bool
}

type NodeItem struct {
	Node *graph.Node
}

func NewNodeItem(n *graph.Node) NodeItem {
	return NodeItem{Node: n}
}

func (this NodeItem) Less(h Heap, other Heaper) bool {
	b := other.(NodeItem)
	switch {
	case this.Node.NumTaken > b.Node.NumTaken:
		return true
	case this.Node.NumTaken < b.Node.NumTaken:
		return false
	default:
		return h.UseDegree && this.Node.Degree > b.Node.Degree
	}
}

func NewNodeHeap(initialCap int, useDegree bool) Heap {
	return Heap{
		Items:     make([]Heaper, 0, initialCap),
		UseDegree: useDegree,
	}
}

type ColorItem struct {
	Color    graph.Color
	Priority int
}

func (this ColorItem) Less(h Heap, other Heaper) bool {
	b := other.(ColorItem)
	return this.Priority > b.Priority
}

func NewColorHeap(src *graph.Node) Heap {
	items := make([]int, graph.NumColors-1)
	for _, adj := range src.Adj {
		if adj.C == graph.Blank {
			for i := graph.Blank + 1; i < graph.NumColors; i++ {
				items[i-1] += adj.TakenColors[i]
			}
		}
	}

	h := Heap{Items: make([]Heaper, graph.NumColors-1)}
	for i := int(graph.Blank + 1); i < graph.NumColors; i++ {
		h.Items[i-1] = ColorItem{
			Color:    graph.Color(i),
			Priority: items[i-1],
		}
	}

	heap.Init(&h)
	return h
}

func (h Heap) Len() int {
	return len(h.Items)
}

func (h Heap) Less(i, j int) bool {
	a, b := h.Items[i], h.Items[j]
	return a.Less(h, b)
}

func (h Heap) Swap(i, j int) {
	h.Items[i], h.Items[j] = h.Items[j], h.Items[i]
}

func (h *Heap) Push(x interface{}) {
	item := x.(Heaper)
	h.Items = append(h.Items, item)
}

func (h *Heap) Pop() interface{} {
	old := h.Items
	n := len(old)
	x := old[n-1]
	h.Items = old[:n-1]

	return x
}
