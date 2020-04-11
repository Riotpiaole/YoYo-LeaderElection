package main

import (
	"fmt"
)

func main() {
	fmt.Println("Hello World")
	graph := CreateTestGraph()
	for i := 1; i <= 9; i++ {
		fmt.Printf(
			"graph Node %v", graph.nodes[i],
		)
		fmt.Printf(
			"graph Node %v\n", graph.edges[i],
		)
	}
}
