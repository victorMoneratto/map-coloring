package graph

import "testing"

func TestAdd(t *testing.T) {
	var graph Graph
	node1, err := graph.Add("node1")
	if err != nil {
		t.Error("Couldn't add first node")
	}

	node2, err := graph.Add("node2")
	if err != nil {
		t.Error("Couldn't add first node")
	}

	if node1.Name() == node2.Name() {
		t.Error("Different nodes have the same name")
	}
}

func TestConnect(t *testing.T) {
	var graph Graph

	node1, _ := graph.Add("node1")
	node2, _ := graph.Add("node2")

	err := node1.Connect(node2, false)
	if err != nil {
		t.Error(err)
	}
	err = node2.Connect(node1, false)
	if err != nil {
		t.Error(err)
	}

	if node1.Adj["node2"] != node2 {
		t.Error("Value in map is different than set (undirected)")
	}

	node3, _ := graph.Add("node3")

	err = node3.Connect(node2, true)
	if err != nil {
		t.Error(err)
	}

	if node3.Adj[node2.Name()] != node2 {
		t.Error("Value in map is different than set (bidirected)")
	}
	if node2.Adj[node3.Name()] != node3 {
		t.Error("Value in map is different than set (bidirected)")
	}
	if _, exists := node1.Adj[node3.Name()]; exists {
		t.Error("node1 shouldn't be connected to node3")
	}
}
