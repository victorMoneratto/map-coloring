package graph

type Graph struct {
	Nodes []Node
}

type Node struct {
	name string
	C    Color
	Adj  []int
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

func NewNode(name string) Node {
	return Node{name: name}
}

func (n Node) Name() string {
	return n.name
}

func (n *Node) Connect(dest int) {
	const initialCap int = 3
	if n.Adj == nil {
		n.Adj = make([]int, 0, initialCap)
	}

	n.Adj = append(n.Adj, dest)
}

func NewGraph(initialCount int) *Graph {
	return &Graph{Nodes: make([]Node, 0, initialCount)}
}
