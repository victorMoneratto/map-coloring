package graph

import (
	"errors"
)

type Graph struct {
	Map map[string]*Node
}

type Node struct {
	name string
	C    Color
	Adj  map[string]*Node
}

type Color int

const (
	Blank Color = iota
	Blue
	Yellow
	Red
	Green
)

var ColorNames = [...]string{
	"Sem cor",
	"Azul",
	"Amarelo",
	"Vermelho",
	"Verde",
}

func (c Color) String() string {
	return ColorNames[c]
}

func NewNode(name string) *Node {
	return &Node{name: name}
}

func (n Node) Name() string {
	return n.name
}

func (g *Graph) GetByName(name string) *Node {
	return g.Map[name]
}

func (n *Node) Connect(dest *Node, bidirectional bool) error {
	if dest == nil {
		return errors.New("Can't connect to nil node " + dest.Name())
	}

	if n.Adj == nil {
		n.Adj = make(map[string]*Node)
	}

	n.Adj[dest.Name()] = dest
	if bidirectional {
		dest.Connect(n, false)
	}
	return nil
}

// Add node to graph by name
// n is number of neighbors to pre allocate
func (g *Graph) Add(name string) (*Node, error) {

	if g.Map == nil {
		g.Map = make(map[string]*Node)
	}

	node, present := g.Map[name]
	if present {
		return node, &DuplicateNodeError{}
	}

	node = NewNode(name)
	g.Map[name] = node
	return node, nil
}

type DuplicateNodeError struct {
}

func (e *DuplicateNodeError) Error() string {
	return "Node was already present"
}
