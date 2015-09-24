package heap

import "github.com/victormoneratto/map-coloring/graph"

type NodeHeap struct {
	Items     []*NodeItem
	UseDegree bool
}

func NewNodeHeap(initialCap int, useDegree bool) NodeHeap {
	return NodeHeap{
		Items:     make([]*NodeItem, 0, initialCap),
		UseDegree: useDegree,
	}
}

type NodeItem struct {
	Node  *graph.Node
	Index int
}

func NewNodeItem(n *graph.Node) *NodeItem {
	return &NodeItem{Node: n}
}

func (h NodeHeap) Len() int {
	return len(h.Items)
}

func (h NodeHeap) Less(i, j int) bool {
	a, b := h.Items[i].Node, h.Items[j].Node
	switch {
	case a.NumTaken > b.NumTaken:
		return true
	case a.NumTaken < b.NumTaken:
		return false
	default:
		return h.UseDegree && a.Degree > b.Degree
	}
}

func (h NodeHeap) Swap(i, j int) {
	h.Items[i], h.Items[j] = h.Items[j], h.Items[i]
	h.Items[i].Index = i
	h.Items[j].Index = j
}

func (h *NodeHeap) Push(x interface{}) {
	item := x.(*NodeItem)
	item.Index = h.Len()
	h.Items = append(h.Items, item)
}

func (h *NodeHeap) Pop() interface{} {
	old := h.Items
	n := len(old)
	x := old[n-1]
	h.Items = old[:n-1]

	return x
}
