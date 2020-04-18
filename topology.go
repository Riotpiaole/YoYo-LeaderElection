package main

import (
	"fmt"
	"sync"
	"time"
)

type edge struct {
	from int
	to   int
}

type Graph struct {
	nodes map[int]*Node

	links map[edge]Link

	edges    map[int]Edge
	dAG      map[int]Edge
	inComing map[int]Edge
	muxDAG   sync.Mutex

	source []int
	muxSRC sync.Mutex

	sourcewg sync.WaitGroup
}

func (g *Graph) addEdge(id int, links []int) {
	for _, link := range links {
		g.edges[id] = append(g.edges[id], link)
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
	graph.inComing = make(map[int]Edge)

	graph.links = make(map[edge]Link)

	for i := 1; i <= N; i++ {
		graph.nodes[i] = NewNode(i)
		graph.edges[i] = Edge([]int{})
		graph.dAG[i] = Edge([]int{})
	}
	return graph
}

func CreateSmallerTestGraph() *Graph {
	graph := NewGraph(4)
	graph.addEdge(1, []int{2, 4})
	graph.addEdge(2, []int{1, 4})
	graph.addEdge(3, []int{4})
	graph.addEdge(4, []int{1, 2, 3})
	return graph
}

func CreateTestGraph() *Graph {
	graph := NewGraph(9)
	graph.addEdge(1, []int{2, 4})
	graph.addEdge(2, []int{1, 4, 7})
	graph.addEdge(3, []int{4, 6})
	graph.addEdge(4, []int{1, 2, 3})
	graph.addEdge(5, []int{6, 7, 9})
	graph.addEdge(6, []int{3, 5, 8, 9})
	graph.addEdge(7, []int{2, 5, 8})
	graph.addEdge(8, []int{6, 7})
	graph.addEdge(9, []int{5, 6})
	return graph
}

func DefaultNodesGraph(nodes []int) *Graph {
	graph := new(Graph)

	graph.nodes = make(map[int]*Node)

	graph.edges = make(map[int]Edge)
	graph.dAG = make(map[int]Edge)
	graph.inComing = make(map[int]Edge)

	graph.links = make(map[edge]Link)
	for _, node := range nodes {
		graph.nodes[node] = NewNode(node)
		graph.edges[node] = Edge([]int{})
		graph.dAG[node] = Edge([]int{})
	}
	return graph
}

func CreateLectureGraph() *Graph {
	graph := DefaultNodesGraph(
		[]int{2, 3, 4, 5, 7, 11, 12, 14, 20, 31, 41})
	graph.addEdge(2, []int{31})
	graph.addEdge(3, []int{11, 12, 14})
	graph.addEdge(4, []int{14, 20})
	graph.addEdge(5, []int{11, 12, 20})
	graph.addEdge(7, []int{20, 31})
	graph.addEdge(11, []int{3, 5, 12})
	graph.addEdge(12, []int{3, 5, 11, 20})
	graph.addEdge(14, []int{3, 4})
	graph.addEdge(20, []int{4, 5, 7, 12, 41})
	graph.addEdge(31, []int{2, 7})
	graph.addEdge(41, []int{20})
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
			// fmt.Printf("%d checking  %d : compare result [%t] ", v.id, neighbour.id, neighbour.compare(v))
			if !v.compare(neighbour) {
				g.dAG[v.id] = append(g.dAG[v.id], neighbour)
				g.inComing[neighbour] = append(g.inComing[neighbour], v.id)
				cnt++
			}
		}
		if cnt == len(g.edges[v.id]) {
			v.state = SOURCE
			g.source = append(g.source, v.id)
		} else if cnt == 0 {
			v.state = SINK
		} else {
			v.state = INTERNAL
		}
		cnt = 0
	}
}

func flipEdge(g *Graph, u int, v int) {
	// u = 1  v = 2  1 -> 2
	// outgoing  2 -> 1
	// incoming  1 <- 2 2 remove incoming 1

	// fmt.Printf("Fliping edges between [%d] [%d]\n", u, v)
	// fmt.Printf("Before fliping \n\tincoming[%v]\n\toutgoing[%v]\n",
	// 	g.inComing, g.dAG)
	fmt.Printf("[%d] is fliping edge {%d , %v}\n", u, u, v)
	g.muxDAG.Lock()

	g.inComing[u] = appendEdge(g.inComing[u], v)
	g.dAG[u] = removeEdge(g.dAG[u], v)

	g.inComing[v] = removeEdge(g.inComing[v], u)
	g.dAG[v] = appendEdge(g.dAG[v], u)

	g.muxDAG.Unlock()
	// fmt.Printf("After fliping \n\tincoming[%v]\n\toutgoing[%v]\n",
	// 	g.inComing, g.dAG)
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
	// removed current Node u
	fmt.Printf("Pruning anything related to %d", u)
	g.muxDAG.Lock()
	for k := range g.dAG {
		g.dAG[k] = removeEdge(g.dAG[k], u)
	}
	g.inComing[u] = Edge([]int{})
	g.muxDAG.Unlock()
}

func (g *Graph) pruneEdge(u int, edges []int) {
	if len(edges) == 0 {
		return
	}
	fmt.Printf("[%d] is pruning all outgoing in %d and incoming of %v\n", u, u, edges)
	fmt.Printf("[%d] is checking Grpah \n\toutGoing: %v \n\tinComing: %v\n", u, g.dAG, g.inComing)
	g.muxDAG.Lock()
	for _, edge := range edges {
		g.inComing[edge] = removeEdge(g.inComing[edge], u)
		g.inComing[u] = removeEdge(g.inComing[u], edge)
		g.dAG[u] = removeEdge(g.dAG[u], edge)
		g.dAG[edge] = removeEdge(g.dAG[edge], u)
	}
	g.muxDAG.Unlock()
	fmt.Printf("After PRUNE Edge")
	fmt.Printf("\n[%d] is checking Grpah \n\toutGoing: %v \n\tinComing: %v\n\n", u, g.dAG, g.inComing)
}

func YoYo(g *Graph) {
	fmt.Printf("")
	activateAllNodesHandler(g)
	time.Sleep(1 * time.Second)
	for _, sourceNode := range g.source {
		go g.nodes[sourceNode].SinkYoDOWN(g)
	}
	g.closeAllChannel()
}

func (g *Graph) PrintGraph(u int, v int) {
	fmt.Println("\n\n========================================================================================")
	fmt.Printf("[%d , %d] is checking Grpah \n\toutGoing: %v \n\tinComing: %v\n", u, v, g.dAG, g.inComing)
	fmt.Printf(time.Now().Format("2006-01-02 15:04:05.000000\n"))
	fmt.Printf("[%d] source are %v\n", u, g.source)
	for _, node := range g.nodes {
		node.printNode(g, Message{CHECK, 1000, 1000}, 1000)
	}
	fmt.Printf("========================================================================================\n\n")
}

func (g *Graph) GlobalUpdate() {
	for _, node := range g.nodes {
		node.updateState(g)
	}
}

func (g *Graph) YoDown() {

	for len(g.source) != 0 {
		for _, sourceNode := range g.source {
			g.nodes[sourceNode].SinkYoDOWN(g)
			g.sourcewg.Add(1)
		}
		// fmt.Printf("Check nusm  %d", cnt)
		g.sourcewg.Wait()
	}
}
