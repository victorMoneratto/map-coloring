package main

import (
	"bufio"
	"log"
	"os"
	"regexp"
	"strconv"

	"github.com/victormoneratto/map-coloring/graph"
)

func main() {
	var graph graph.Graph
	graphDescription := parseInputFile("input/brasil.in")

	// create all nodes (each is src once)
	for _, nodeDescription := range graphDescription {
		src := nodeDescription[0]
		_, err := graph.Add(src)
		if err != nil {
			panic(err)
		}
	}
	// connect srcs and dests
	for _, nodeDescription := range graphDescription {
		src := graph.GetByName(nodeDescription[0])
		for _, destName := range nodeDescription[1:] {
			dest := graph.GetByName(destName)
			if err := src.Connect(dest, false); err != nil {
				panic(err)
			}
		}
		log.Println(src)
	}
}

// Parse input file to [][]string of the format:
// [
//  [src1 dst1 dst2... dstn]
//  [src2 dst4 dst1... dstm]
// ]
// where each node is src once and many dst can connect to it
func parseInputFile(filename string) [][]string {
	file, err := os.Open(filename)
	defer file.Close()
	if err != nil {
		panic(err)
	}

	buffer := bufio.NewReader(file)
	firstLine, err := buffer.ReadString('\n')
	count, err := strconv.Atoi(firstLine[:len(firstLine)-1])

	if err != nil {
		panic(err)
	}

	ret := make([][]string, count)

	toExp := regexp.MustCompile("(?:([^,:]*))(?:[,|.|:]) *")

	for i := 0; i < count; i++ {
		line, _ := buffer.ReadString('\n')
		for _, name := range toExp.FindAllStringSubmatch(line, count-1) {
			ret[i] = append(ret[i], name[1])
		}
	}

	return ret
}
