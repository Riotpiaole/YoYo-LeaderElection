package main

// package main

// import (
// 	"fmt"
// 	"time"
// )

// func main() {

// 	// wg.Add(1)
// 	// graph := NewGraph(2)
// 	graph := CreateLectureGraph()
// 	// graph.nodes[31] = NewNode(31)
// 	// graph.addEdge(1, []int{31})
// 	// graph.addEdge(2, []int{31})

// 	// graph.addEdge(31, []int{1, 2})
// 	graph.Initalizes()
// 	// graph.nodes[2].min = 1

// 	// for k, v := range graph.dAG {
// 	// 	for _, i := range v {
// 	// 		fmt.Printf("src %v  dest %v ", k, i)
// 	// 	}
// 	// 	fmt.Print("\n")
// 	// }
// 	// fmt.Printf("OutGoing %v\n", graph.dAG)
// 	// fmt.Printf("Incoming %v\n", graph.inComing)
// 	activateAllNodesHandler(graph)
// 	for _, sourceNode := range graph.source {
// 		graph.nodes[sourceNode].SinkYoDOWN(graph)
// 	}
// 	// time.Sleep(1 * time.Second)
// 	// graph.nodes[1].SinkYoDOWN(graph)
// 	// // time.Sleep(1 * time.Second)
// 	// graph.nodes[2].SinkYoDOWN(
// 	// 	graph,
// 	// )
// 	time.Sleep(5 * time.Second)
// 	fmt.Println("chekcing Final Graph")
// 	fmt.Printf("Incoming %v\n", graph.inComing)
// 	fmt.Printf("OutGoing %v\n", graph.dAG)
// 	fmt.Printf("Remaing Source is %v\n", graph.source)
// 	// graph.closeAllChannel()
// }

func main() {
	graph := HyperCube(2)

}
