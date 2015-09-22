package main

import (
	"os"
	"testing"

	"github.com/victormoneratto/map-coloring/graph"
)

func TestParseInputFile(t *testing.T) {
	heuristic := 'a'
	file, err := os.Open("input/brasil.in")
	if err != nil {
		t.Error(err)
	}
	descs := parseInputFile(file, &heuristic)
	if descs[0][0] != "Acre" {
		t.Error("Should be Acre, was", descs[0][0])
	}

	if descs[0][1] != "Amazonas" {
		t.Error("Should be Amazonas, was", descs[0][1])
	}
	if descs[0][2] != "Rondônia" {
		t.Error("Should be Rondônia, was", descs[0][2])
	}

	if descs[26][0] != "Tocantins" {
		t.Error("Should be Tocantins, was", descs[26][0])
	}
}

func TestColorMap(t *testing.T) {
	heuristic := 'a'
	file, err := os.Open("input/brasil.in")
	if err != nil {
		t.Error(err)
	}
	descs := parseInputFile(file, &heuristic)
	g := graph.NewGraph(len(descs))
	populateGraph(g, descs)

	if colorMap(g) {
		for _, node := range g.Nodes {
			for _, adj := range node.Adj {
				if node.C == g.Nodes[adj].C {
					t.Error("Adjacent nodes have the same color:", node, g.Nodes[adj])
				}
			}
		}
	} else {
		t.Error("Graph coloring for should be possible")
	}
}
