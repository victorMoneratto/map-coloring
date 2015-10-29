package main

import (
	"bufio"
	"container/heap"
	"errors"
	"flag"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/victormoneratto/map-coloring/graph"
	nh "github.com/victormoneratto/map-coloring/heap"
)

// Heuristic flag enum
const (
	NoHeuristic = 'a' + iota
	FCOnly
	FCandMRV
	FCandMRVandDegree
)

var Heuristic struct {
	MRV, FC, Degree bool
}

func main() {
	// Heuristic level flag
	heuristicString := flag.String("heuristic", "a", "a, b, c or d")
	// Input file flag
	inputFile := flag.String("file", "input/usa.in", "filename of input file, standard input for default")
	flag.Parse()

	// Open file if specified, use stdin otherwise
	var file *os.File
	if *inputFile != "" {
		var err error
		file, err = os.Open(*inputFile)
		if err != nil {
			panic(err)
		}
	} else {
		file = os.Stdin
	}

	//parse input file
	graphDescription, fileHeuristic := parseInputFile(file)
	file.Close()
	if len(fileHeuristic) > 0 {
		heuristicString = &fileHeuristic
	}

	// Initialize Heuristic global
	parseHeuristic(*heuristicString)

	// Populate graph
	g := graph.NewGraph(len(graphDescription))
	populateGraph(&g, graphDescription)

	// start := time.Now()
	if colorMap(g) {
		// elapsed := time.Now().Sub(start)
		// fmt.Println(elapsed)
		for _, node := range g {
			fmt.Println(node.Name()+":", node.Color.String()+".")
		}
	} else {
		fmt.Println("Impossível")
	}
}

// Parse heuristc flag and init Heuristc global
func parseHeuristic(input string) {
	firstChar := input[0]
	switch {
	case firstChar == FCandMRVandDegree:
		Heuristic.Degree = true
		fallthrough
	case firstChar == FCandMRV:
		Heuristic.MRV = true
		fallthrough
	case firstChar == FCOnly:
		Heuristic.FC = true
	}
}

// Create graph based on description
func populateGraph(g *graph.Graph, graphDescription [][]string) {
	nodesMap := make(map[string]*graph.Node)

	// create all nodes (each is src once)
	for _, nodeDescription := range graphDescription {
		srcName := nodeDescription[0]
		node := graph.NewNode(srcName)
		*g = append(*g, node)
		nodesMap[srcName] = node
	}

	// connect srcs and dests
	for _, nodeDescription := range graphDescription {
		src := nodesMap[nodeDescription[0]]
		for _, destName := range nodeDescription[1:] {
			dest := nodesMap[destName]
			src.Connect(dest)
		}
	}
}

// Setup for graph coloring backtracking algorithm,
// search heuristics are set based on heuristc flag
func colorMap(g graph.Graph) bool {

	// For MRV, we'll use a heap
	var h nh.Heap
	if Heuristic.MRV {
		h = nh.NewNodeHeap(len(g), Heuristic.Degree)
		for _, n := range g {
			h.Items = append(h.Items, nh.NewNodeItem(n))
		}
		heap.Init(&h)
	}

	// Init colors cache
	ColorsCache.Colors = make([][]graph.NodeColor, len(g))
	for i := range g {
		ColorsCache.Colors[i] = make([]graph.NodeColor, 0, graph.NumColors-1)
		for color := graph.Blank + 1; color < graph.NumColors; color++ {
			ColorsCache.Colors[i] = append(ColorsCache.Colors[i], graph.NodeColor(color))
		}
	}

	return colorMapBacktrack(g, 0, &h)
}

// Cache for colors
var ColorsCache struct {
	Colors [][]graph.NodeColor
}

var assignments int

// Graph coloring backtrack,
// search heuristics are set based on heuristc flag
func colorMapBacktrack(g graph.Graph, height int, h *nh.Heap) bool {
	// Have all nodes been colored?
	if height == len(g) {
		return true
	}

	// Assign current node by MRV heap or sequentally
	var currNode *graph.Node
	if Heuristic.MRV {
		// Top from the heap has the Minimum Remaining Values
		currNode = heap.Pop(h).(nh.NodeItem).Node
	} else {
		currNode = g[height]
	}

	colors := ColorsCache.Colors[height]

	for _, color := range colors {
		// Is this color avaliable for use?
		if currNode.TakenColors[color] == 0 {
			// Update currNode color and warn adjacents
			impossible := updateColorAndAdj(currNode, color, h)

			// Forward Checking
			if Heuristic.FC {
				if impossible {
					resetColorAndAdj(currNode)
					continue
				}
			}

			// Proceed with recursion
			if colorMapBacktrack(g, height+1, h) {
				return true
			}
			resetColorAndAdj(currNode)
		}
	}

	// Need to backtrack
	if Heuristic.MRV {
		heap.Push(h, nh.NewNodeItem(currNode))
	}
	return false
}

func resetColorAndAdj(node *graph.Node) {
	for _, adj := range node.Adj {
		adj.NeighborLostColor(node.Color)
	}
	node.Color = graph.Blank
}

// Update node color and udpate colors info on adjacents
func updateColorAndAdj(node *graph.Node, newColor graph.NodeColor, h *nh.Heap) bool {
	anyZeroed := false
	for _, adj := range node.Adj {
		adj.NeighborGotColor(newColor)
		anyZeroed = anyZeroed || adj.NumTaken == graph.NumColors-1
	}
	if Heuristic.MRV {
		heap.Init(h)
	}
	node.Color = newColor

	return anyZeroed
}

// Parse input file to [][]string of the format:
// [
// [src1 dst1 dst2... dstn][value]z
// ...
// [src2 dst4 dst1... dstm]
// ]
// where each node is src once and many dst can connect to it
func parseInputFile(file *os.File) ([][]string, string) {
	buffer := bufio.NewReader(file)

	// Read first line
	line, err := buffer.ReadString('\n')
	if err != nil {
		panic(errors.New("File out of format: " + err.Error()))
	}

	// First line is numberOfLines OptionalHeuristicOverride
	fields := strings.Fields(line)
	nodesCount, err := strconv.Atoi(fields[0])
	if err != nil {
		panic(errors.New("File out of format: " + err.Error()))
	}

	// Read heuristic flag
	var heuristic string
	if len(fields) >= 2 {
		heuristic = fields[1]
	}

	ret := make([][]string, nodesCount)

	// For each line as predefined
	for i := 0; i < nodesCount; i++ {
		// Read line and split by ': or ',' or '.'
		line, _ = buffer.ReadString('\n')
		fields := strings.FieldsFunc(line, func(r rune) bool {
			return r == ',' || r == ':' || r == '.' || r == '\n' || r == '\r'
		})
		// fmt.Println(fields)

		// Allocate string slice and store elements from the line
		ret[i] = make([]string, 0, len(fields))
		for _, field := range fields {
			ret[i] = append(ret[i], strings.TrimSpace(field))
		}
	}

	return ret, heuristic
}
