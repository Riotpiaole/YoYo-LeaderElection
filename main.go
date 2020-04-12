package main

import (
	"fmt"
	"time"
)

func main() {
	fmt.Println("Hello World")
	graph := NewGraph(2)
	graph.nodes[31] = NewNode(31)
	graph.addEdge(1, []int{31})
	graph.addEdge(2, []int{31})

	graph.addEdge(31, []int{1, 2})
	graph.Initalizes()
	for k, v := range graph.dAG {
		for _, i := range v {
			fmt.Printf("src %v  dest %v ", k, i.id)
		}
		fmt.Print("\n")
	}
	fmt.Printf("%v\n", graph.dAG)
	activateAllNodesHandler(graph)
	// // graph.nodes[2].handleMessage(graph)
	// // graph.nodes[1].handleMessage(graph)
	time.Sleep(1 * time.Second)
	graph.nodes[1].SinkYoDOWN(
		YoDown,
		graph,
	)
	// time.Sleep(1 * time.Second)
	graph.nodes[2].SinkYoDOWN(
		YoDown,
		graph,
	)

	time.Sleep(1 * time.Second)
	graph.closeAllChannel()
}
