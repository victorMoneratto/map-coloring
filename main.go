package main

import (
	"bufio"
	"errors"
	"flag"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/victormoneratto/map-coloring/graph"
)

//Heuristic flag enum
const (
	NoHeuristic = 'a' + iota
	MRVOnly
	MRVandDegree
	AllHeuristics
)

func main() {
	// profiler settings
	// cfg := profile.Config{ProfilePath: ".",
	// CPUProfile: true,
	// MemProfile: true,
	// }
	// defer profile.Start(&cfg).Stop()

	//heuristic level flag
	heuristicLevel := rune((*flag.String("heuristic", " ", "a, b, c or d"))[0])
	//input file flag
	inputFile := flag.String("file", "", "filename of input file, standard input for default")
	flag.Parse()

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
	graphDescription := parseInputFile(file, &heuristicLevel)
	file.Close()

	//populate graph
	g := graph.NewGraph(len(graphDescription))
	populateGraph(g, graphDescription)

	// perform backtracking
	if colorMap(g) {
		for _, node := range g.Nodes {
			fmt.Println(node.Name()+":", node.C)
		}
	} else {
		fmt.Println("Impossível")
	}
}

func populateGraph(g *graph.Graph, graphDescription [][]string) {
	nodesMap := make(map[string]int)

	// create all nodes (each is src once)
	for _, nodeDescription := range graphDescription {
		srcName := nodeDescription[0]
		node := graph.NewNode(srcName)
		g.Nodes = append(g.Nodes, node)
		nodesMap[srcName] = len(g.Nodes) - 1
	}

	// connect srcs and dests
	for _, nodeDescription := range graphDescription {
		src := nodesMap[nodeDescription[0]]
		for _, destName := range nodeDescription[1:] {
			dest := nodesMap[destName]
			g.Nodes[src].Connect(dest)
		}
	}
}

// cache for backtrack
type BacktrackCache struct {
	TakenColors [][]int
}

// arrange map coloring backtrack algorithm
func colorMap(g *graph.Graph) bool {
	cache := &BacktrackCache{TakenColors: make([][]int, len(g.Nodes))}
	for i := 0; i < len(g.Nodes); i++ {
		cache.TakenColors[i] = make([]int, len(graph.ColorNames))
	}
	return colorMapBacktrack(g, cache, 0)
}

// recursive map coloring (backtrack)
func colorMapBacktrack(g *graph.Graph, cache *BacktrackCache, height int) bool {

	// have all nodes been colored?
	if height == len(g.Nodes) {
		return true
	}

	// for each color
	for color := graph.Color(graph.Blank + 1); int(color) < len(graph.ColorNames); color++ {
		//can we use it?
		if cache.TakenColors[height][color] == 0 {
			updateColorAndAdj(g, cache, height, color)
			//proceed with backtrack recursion
			if colorMapBacktrack(g, cache, height+1) {
				return true
			}
		}
	}

	// reset color and ajacents
	updateColorAndAdj(g, cache, height, graph.Blank)
	return false
}

// update node color and cache values for it's adjacents
func updateColorAndAdj(g *graph.Graph, cache *BacktrackCache, nodeIndex int, newColor graph.Color) {
	prevColor := g.Nodes[nodeIndex].C
	for _, value := range g.Nodes[nodeIndex].Adj {
		cache.TakenColors[value][newColor]++
		cache.TakenColors[value][prevColor]--
	}

	g.Nodes[nodeIndex].C = newColor
}

// Parse input file to [][]string of the format:
// [
// [src1 dst1 dst2... dstn][value]z
// ...
// [src2 dst4 dst1... dstm]
// ]
// where each node is src once and many dst can connect to it
func parseInputFile(file *os.File, heuristic *rune) [][]string {
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

	if len(fields) > 1 && *heuristic == ' ' {
		*heuristic = rune(fields[1][0])
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
