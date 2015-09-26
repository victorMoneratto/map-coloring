package heap

import (
	"container/heap"

	"github.com/victormoneratto/map-coloring/graph"
)

// Type for heap for Heaper interface
type Heap struct {
	Items     []Heaper
	UseDegree bool
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

// Heaper interface
type Heaper interface {
	Less(h Heap, other Heaper) bool
}

// Node Item for heap (implements Heaper)
type NodeItem struct {
	*graph.Node
}

func NewNodeHeap(initialCap int, useDegree bool) Heap {
	return Heap{
		Items:     make([]Heaper, 0, initialCap),
		UseDegree: useDegree,
	}
}

func NewNodeItem(n *graph.Node) NodeItem {
	return NodeItem{n}
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

// Node Item for Color
type ColorItem struct {
	Color    graph.NodeColor
	Priority int
}

func NewColorHeap(src *graph.Node) Heap {
	h := Heap{Items: make([]Heaper, graph.NumColors-1)}
	for i := int(graph.Blank + 1); i < graph.NumColors; i++ {
		h.Items[i-1] = ColorItem{
			Color: graph.NodeColor(i),
		}
	}

	for _, adj := range src.Adj {
		if adj.Color == graph.Blank {
			for i := graph.Blank + 1; i < graph.NumColors; i++ {
				priority := h.Items[i-1].(ColorItem)
				priority.Priority += adj.TakenColors[i]
				h.Items[i-1] = priority
			}
		}
	}

	heap.Init(&h)
	return h
}

func (this ColorItem) Less(h Heap, other Heaper) bool {
	b := other.(ColorItem)
	return this.Priority > b.Priority
}
