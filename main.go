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

func main() {
	heuristics := rune((*flag.String("heuristics", "a", "a, b, c or d"))[0])
	flag.Parse()

	var g graph.Graph
	graphDescription := parseInputFile("input/brasil.in", &heuristics)

	// create all nodes (each is src once)
	for _, nodeDescription := range graphDescription {
		src := nodeDescription[0]
		_, err := g.Add(src)
		if err != nil {
			panic(err)
		}
	}
	// connect srcs and dests
	for _, nodeDescription := range graphDescription {
		src := g.GetByName(nodeDescription[0])
		for _, destName := range nodeDescription[1:] {
			dest := g.GetByName(destName)
			if err := src.Connect(dest, false); err != nil {
				panic(err)
			}
		}
	}

	if backtracking(&g) {
		for _, node := range g.Map {
			fmt.Println(node.Name()+":", node.C)
		}
	} else {
		fmt.Println("Impossível")
	}
}

// Parse input file to [][]string of the format:
// [
//  [src1 dst1 dst2... dstn]
// 	...
//  [src2 dst4 dst1... dstm]
// ]
// where each node is src once and many dst can connect to it
func parseInputFile(filename string, level *rune) [][]string {
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
	count, err := strconv.Atoi(fields[0])
	if err != nil {
		panic(err)
	}
	if len(fields) > 1 {
		*level = rune(fields[1][0])
	}

	ret := make([][]string, count)

	statesExp := regexp.MustCompile("(?:([^,:]*))(?:[,|.|:]) *")

	for i := 0; i < count; i++ {
		line, _ := buffer.ReadString('\n')
		for _, name := range statesExp.FindAllStringSubmatch(line, count-1) {
			ret[i] = append(ret[i], name[1])
		}
	}

	return ret
}

func backtracking(g *graph.Graph) bool {
	nodes := make([]*graph.Node, 0, len(g.Map))

	for _, value := range g.Map {
		nodes = append(nodes, value)
	}

	return backtrackingRecursive(nodes, 0)
}

func backtrackingRecursive(nodes []*graph.Node, count int) bool {
	if count == len(nodes) {
		return true
	}

	curr := nodes[count]
	takenColors := selectColor(curr)

	for color := graph.Color(1); int(color) < len(takenColors); color++ {
		if !takenColors[color] {
			curr.C = color
			if backtrackingRecursive(nodes, count+1) {
				return true
			}
		}
	}
	curr.C = graph.Blank
	return false
}

func selectColor(node *graph.Node) []bool {
	takenColors := make([]bool, len(graph.ColorNames))
	for _, value := range node.Adj {
		takenColors[value.C] = true
	}

	return takenColors
}
