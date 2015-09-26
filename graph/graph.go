package graph

// Type for graph
type Graph []*Node

func NewGraph(initialCount int) Graph {
	return make([]*Node, 0, initialCount)
}

// Type for node in graph
type Node struct {
	name        string
	Color       NodeColor
	Adj         []*Node
	TakenColors []int
	NumTaken    int
	Degree      int
}

func NewNode(name string) *Node {
	return &Node{name: name, TakenColors: make([]int, NumColors)}
}

func (n Node) Name() string {
	return n.name
}

// Connect nodes in graph
func (n *Node) Connect(dest *Node) {
	const initialCap int = 3
	if n.Adj == nil {
		n.Adj = make([]*Node, 0, initialCap)
	}

	n.Adj = append(n.Adj, dest)
	n.Degree++
}

// Type for colors a node can be
type NodeColor int

const (
	Blank = iota
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

const NumColors = len(ColorNames)

// Implement fmt.Stringer interface
func (c NodeColor) String() string {
	return ColorNames[c]
}
