package graph

type Graph []*Node

type Node struct {
	name        string
	C           Color
	Adj         []*Node
	TakenColors []int
	NumTaken    int
	Degree      int
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

const NumColors = len(ColorNames)

func (c Color) String() string {
	return ColorNames[c]
}

func NewNode(name string) *Node {
	return &Node{name: name, TakenColors: make([]int, NumColors)}
}

func (n Node) Name() string {
	return n.name
}

func (n *Node) Connect(dest *Node) {
	const initialCap int = 3
	if n.Adj == nil {
		n.Adj = make([]*Node, 0, initialCap)
	}

	n.Adj = append(n.Adj, dest)
	n.Degree++
}

func NewGraph(initialCount int) Graph {
	return make([]*Node, 0, initialCount)
}
