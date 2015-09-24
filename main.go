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

//Heuristic flag enum
const (
	NoHeuristic = 'a' + iota
	MRVOnly
	MRVandDegree
	MRVDegreeAndLCV
)

var heuristic rune

func main() {
	// profiler settings
	// cfg := profile.Config{ProfilePath: ".",
	// CPUProfile: true,
	// MemProfile: true,
	// }
	// defer profile.Start(&cfg).Stop()

	//heuristic level flag
	heuristicString := flag.String("heuristic", "a", "a, b, c or d")
	//input file flag
	inputFile := flag.String("file", "", "filename of input file, standard input for default")
	flag.Parse()

	heuristic = rune((*heuristicString)[0])

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
	graphDescription := parseInputFile(file)
	file.Close()

	//populate graph
	g := graph.NewGraph(len(graphDescription))
	populateGraph(&g, graphDescription)
	// perform backtracking
	if colorMap(g) {
		for _, node := range g {
			fmt.Println(node.Name()+":", node.C)
		}
	} else {
		fmt.Println("Impossível")
	}
}

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

// arrange map coloring backtrack algorithm
func colorMap(g graph.Graph) bool {
	var h nh.Heap
	if heuristic >= MRVOnly {
		h = nh.NewNodeHeap(len(g), heuristic >= MRVandDegree)
		for _, n := range g {
			h.Items = append(h.Items, nh.NewNodeItem(n))
		}
		heap.Init(&h)
	}
	return colorMapBacktrack(g, 0, &h)
}

// recursive map coloring (backtrack)
func colorMapBacktrack(g graph.Graph, height int, h *nh.Heap) bool {
	// have all nodes been colored?
	if height == len(g) {
		return true
	}

	// current node
	var node *graph.Node
	if heuristic >= MRVOnly {
		node = heap.Pop(h).(nh.NodeItem).Node
	} else {
		node = g[height]
	}

	colors := make([]graph.Color, 0, graph.NumColors-1)
	if heuristic >= MRVDegreeAndLCV {
		colorHeap := nh.NewColorHeap(node)
		for colorHeap.Len() > 0 {
			color := heap.Pop(&colorHeap).(nh.ColorItem).Color
			colors = append(colors, color)
		}
	} else {
		for i := graph.Blank + 1; i < graph.NumColors; i++ {
			colors = append(colors, graph.Color(i))
		}
	}

	// for each color
	for _, color := range colors {
		//can we use it?
		if node.TakenColors[color] == 0 {
			//update colors and adj
			updateColorAndAdj(node, color, h)
			//proceed with backtrack recursion
			if colorMapBacktrack(g, height+1, h) {
				return true
			}
		}
	}

	// reset color and ajacents
	if heuristic >= MRVOnly {
		heap.Push(h, node)
	}
	updateColorAndAdj(node, graph.Blank, h)
	return false
}

// update node color and cache values for it's adjacents
func updateColorAndAdj(node *graph.Node, newColor graph.Color, h *nh.Heap) {
	prevColor := node.C

	for _, adj := range node.Adj {
		if heuristic >= MRVOnly {
			updateTakenColors(adj, prevColor, newColor)
		}
		adj.TakenColors[newColor]++
		adj.TakenColors[prevColor]--
	}

	node.C = newColor

	if heuristic >= MRVOnly {
		heap.Init(h)
	}
}

func updateTakenColors(adj *graph.Node, prevColor, newColor graph.Color) {
	//if node was blank
	if prevColor == graph.Blank {
		//node is the only neighbor this color
		if adj.TakenColors[newColor] == 0 {
			adj.NumTaken++
		}
		//inc taken color for new color

		// if node is losing color
	} else if newColor == graph.Blank {
		// node was the only neighbor that color
		if adj.TakenColors[prevColor] == 1 {
			adj.NumTaken--
		}
	}
}

// Parse input file to [][]string of the format:
// [
// [src1 dst1 dst2... dstn][value]z
// ...
// [src2 dst4 dst1... dstm]
// ]
// where each node is src once and many dst can connect to it
func parseInputFile(file *os.File) [][]string {
	buffer := bufio.NewReader(file)

	//read first line
	line, err := buffer.ReadString('\n')
	if err != nil {
		panic(errors.New("File out of format: " + err.Error()))
	}

	//first line is numberOfLines OptionalHeuristicOverride
	fields := strings.Fields(line)
	nodesCount, err := strconv.Atoi(fields[0])
	if err != nil {
		panic(errors.New("File out of format: " + err.Error()))
	}

	if len(fields) > 1 {
		heuristic = rune(fields[1][0])
	}

	ret := make([][]string, nodesCount)

	//for each line as predefined
	for i := 0; i < nodesCount; i++ {
		// read line and split by ': or ',' or '.'
		line, _ = buffer.ReadString('\n')
		fields := strings.FieldsFunc(line, func(r rune) bool {
			return r == ',' || r == ':' || r == '.' || r == '\n'
		})

		// allocate string slice and store elements from the line
		ret[i] = make([]string, 0, len(fields))
		for _, field := range fields {
			ret[i] = append(ret[i], strings.TrimSpace(field))
		}
	}

	return ret
}
