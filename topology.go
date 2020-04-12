package main

import (
	"sync"
)

type edge struct {
	from int
	to   int
}

type Graph struct {
	nodes map[int]*Node

	links map[edge]Link

	edges  map[int]Edge
	dAG    map[int]Edge
	muxDAG sync.Mutex

	source []int
	muxSRC sync.Mutex

	sink    []int
	muxSINK sync.Mutex

	internal    []int
	muxINTERNAL sync.Mutex
}

func (g *Graph) addEdge(id int, links []int) {
	for _, link := range links {
		g.edges[id] = append(g.edges[id], g.nodes[link])
		// fmt.Printf("%v is expecting size %d chan\n", edge{id, link}, len(links))
		g.links[edge{id, link}] = make(chan Message, len(links))
		g.links[edge{link, id}] = make(chan Message, len(links))
	}
}

func NewGraph(N int) *Graph {
	graph := new(Graph)
	graph.nodes = make(map[int]*Node)
	graph.edges = make(map[int]Edge)
	graph.dAG = make(map[int]Edge)

	graph.links = make(map[edge]Link)

	for i := 1; i <= N; i++ {
		graph.nodes[i] = NewNode(i)
		graph.edges[i] = Edge([]*Node{})
		graph.dAG[i] = Edge([]*Node{})
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

func (g *Graph) Initalizes() {
	visited := make(map[int]bool)
	for _, v := range g.nodes {
		visited[v.id] = false
	}
	cnt := 0
	for _, v := range g.nodes {
		for _, neighbour := range g.edges[v.id] {
			if neighbour.compare(v) {
				g.dAG[v.id] = append(g.dAG[v.id], neighbour)
				cnt++
			}
		}
		if cnt == len(g.edges[v.id]) {
			v.state = SOURCE
			g.source = append(g.source, v.id)
		} else if cnt == 0 {
			v.state = SINK
			g.sink = append(g.sink, v.id)
		} else {
			v.state = INTERNAL
			g.internal = append(g.internal, v.id)
		}
		cnt = 0
	}

}

func flipEdge(g *Graph, u int, v int) {
	g.muxDAG.Lock()
	e := g.dAG[u]
	g.dAG[u] = removeEdge(&e, v)
	g.dAG[v] = append(g.dAG[v], g.nodes[u])
	g.muxDAG.Unlock()
}

func activateAllNodesHandler(g *Graph) {
	for _, v := range g.nodes {
		v.handleMessage(g)
	}
}

func (g *Graph) closeAllChannel() {
	for _, v := range g.links {
		close(v)
	}
}

func (g *Graph) pruneNode(u int) {

}
