package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"

	"github.com/victormoneratto/map-coloring/graph"
)

const (
	NoHeuristic = 'a' + iota
)

func main() {
	// cfg := profile.Config{MemProfile: true, ProfilePath: "."}
	// defer profile.Start(&cfg).Stop()

	//heuristic leve flag
	heuristicLevel := rune((*flag.String("heuristic", " ", "a, b, c or d"))[0])
	flag.Parse()

	//read file and populate graph
	graphDescription := parseInputFile("input/usa.in", &heuristicLevel)
	g := graph.NewGraph(len(graphDescription))
	populateGraph(g, graphDescription)

	// perform backtracking
	colorMap(g)
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
	TakenColors []bool
}

// arrange map coloring backtrack algorithm
func colorMap(g *graph.Graph) bool {
	cache := make([]BacktrackCache, len(g.Nodes))
	for i := 0; i < len(g.Nodes); i++ {
		cache[i].TakenColors = make([]bool, len(graph.ColorNames))
	}
	return colorMapBacktrack(g, cache, 0)
}

// recursive map coloring (backtrack)
func colorMapBacktrack(g *graph.Graph, cache []BacktrackCache, height int) bool {

	// have all nodes been colored?
	if height == len(g.Nodes) {
		return true
	}

	takenColors := takenColors(g, cache, height)

	// for each color this node can take
	for color := graph.Color(graph.Blank + 1); int(color) < len(takenColors); color++ {
		if !takenColors[color] {
			g.Nodes[height].C = color
			if colorMapBacktrack(g, cache, height+1) {
				return true
			}
		}
	}

	g.Nodes[height].C = graph.Blank
	return false
}

// returns array where true represents a taken color, false otherwise
func takenColors(g *graph.Graph, cache []BacktrackCache, nodeIndex int) []bool {
	for i := range graph.ColorNames {
		cache[nodeIndex].TakenColors[i] = false
	}

	for _, destIndex := range g.Nodes[nodeIndex].Adj {
		dest := &g.Nodes[destIndex]
		cache[nodeIndex].TakenColors[dest.C] = true
	}

	return cache[nodeIndex].TakenColors
}

// Parse input file to [][]string of the format:
// [
// [src1 dst1 dst2... dstn]
// ...
// [src2 dst4 dst1... dstm]
// ]
// where each node is src once and many dst can connect to it
func parseInputFile(filename string, heuristic *rune) [][]string {
	file, err := os.Open(filename)
	defer file.Close()
	if err != nil {
		panic(err)
	}

	buffer := bufio.NewReader(file)
	firstLine, err := buffer.ReadString('\n')
	if err != nil {
		panic(err)
	}
	fields := strings.Fields(firstLine)
	nodesCount, err := strconv.Atoi(fields[0])
	if err != nil {
		panic(err)
	}
	if len(fields) > 1 && *heuristic == ' ' {
		*heuristic = rune(fields[1][0])
	}

	ret := make([][]string, nodesCount)

	statesExp := regexp.MustCompile("(?:([^,:]*))(?:[,|.|:]) *")

	for i := 0; i < nodesCount; i++ {
		line, _ := buffer.ReadString('\n')
		submatches := statesExp.FindAllStringSubmatch(line, nodesCount-1)
		ret[i] = make([]string, 0, len(submatches))
		for _, name := range submatches {
			ret[i] = append(ret[i], name[1])
		}
	}

	return ret
}
