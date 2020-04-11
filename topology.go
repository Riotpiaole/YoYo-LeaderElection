package main

type Graph struct {
	nodes map[int]*Node
	edges map[int]Edge
	dAG   map[int]Edge
}

func (g *Graph) addEdge(id int, links []int) {
	for _, link := range links {
		g.edges[id] = append(g.edges[id], g.nodes[link])
	}
}

func NewGraph(N int) *Graph {
	graph := new(Graph)
	graph.nodes = make(map[int]*Node)
	graph.edges = make(map[int]Edge)
	for i := 1; i <= N; i++ {
		graph.nodes[i] = NewNode(i)
		graph.edges[i] = Edge([]*Node{})
	}

	return graph
}

func CreateTestGraph() *Graph {
	graph := NewGraph(9)
	graph.addEdge(1, []int{2, 4})
	graph.addEdge(2, []int{1, 4, 7})
	graph.addEdge(3, []int{4, 6})
	graph.addEdge(4, []int{1, 2, 3})
	graph.addEdge(5, []int{6, 7, 9})
	graph.addEdge(6, []int{3, 8, 9})
	graph.addEdge(7, []int{2, 5, 8})
	graph.addEdge(8, []int{6, 7})
	graph.addEdge(9, []int{5, 6})
	return graph
}


func Convert