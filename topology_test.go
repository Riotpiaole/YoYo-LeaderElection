package main

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestAddNewEdge(t *testing.T) {
	graph := NewGraph(2)
	graph.addEdge(1, []int{2})
	graph.addEdge(2, []int{1})
	if graph.edges[1][0] != 2 {
		t.Errorf("Edge not properly initalizes")
	}
}

func TestInitalizes(t *testing.T) {
	graph := CreateTestGraph()
	graph.Initalizes()
	expected := []int{6, 7, 9}
	compare := []int{}
	for _, v := range graph.dAG[5] {
		compare = append(compare, v)
	}
	assert.Equal(t, graph.nodes[1].state, SOURCE)
	assert.Equal(t, graph.nodes[8].state, SINK)
	assert.Equal(t, graph.nodes[6].state, INTERNAL)
	assert.Equal(t, compare, expected, "The outgoing of 5 should be this")
	assert.Equal(t, len(graph.inComing[6]), 2)
	assert.Equal(t, len(graph.inComing[7]), 2)
	assert.Equal(t, len(graph.inComing[4]), 3)
}

// func TestChannelCheck(t *testing.T) {
// 	graph := NewGraph(2)
// 	graph.addEdge(1, []int{2})
// 	graph.addEdge(2, []int{1})
// 	graph.Initalizes()
// 	// fmt.Printf("%v \n", graph.dAG[1][0])
// 	// fmt.Printf("%v\n", graph.links[edge{1, 2}])

// 	go graph.nodes[2].receiveMessage(graph)
// 	graph.nodes[1].sendMessage(
// 		YoDown,
// 		graph,
// 	)
// }

func TestEdgeFLip(t *testing.T) {
	graph := CreateTestGraph()
	graph.Initalizes()

	// assert.Equal(t, len(graph.inComing[1]), 0)

	// flipEdge(graph, 1, 2)

	// // 2 -> 1
	// assert.Equal(t, len(graph.inComing[1]), 1)
	// assert.Equal(t, graph.inComing[1][0], 2)
	// incoming of 1 should be 2
}

func TestPruneNode(t *testing.T) {
	graph := NewGraph(2)

	graph.dAG[1] = append(graph.dAG[1], 2)
	graph.inComing[2] = append(graph.inComing[2], 1)

	graph.pruneNode(2)
	assert.Equal(t, len(graph.dAG[1]), 0)
	assert.Equal(t, len(graph.dAG[2]), 0)
	assert.Equal(t, len(graph.inComing[2]), 0)
}

func TesthandlePruneSinkNode(t *testing.T) {
	graph := NewGraph(2)
	graph.nodes[31] = NewNode(31)
	graph.addEdge(1, []int{31})
	graph.addEdge(2, []int{31})

	graph.addEdge(31, []int{1, 2})
	graph.Initalizes()
	activateAllNodesHandler(graph)

	time.Sleep(1 * time.Second)
	graph.nodes[1].SinkYoDOWN(graph)
	// time.Sleep(1 * time.Second)
	graph.nodes[2].SinkYoDOWN(graph)

	time.Sleep(1 * time.Second)
	fmt.Printf("OutGoing %v\n", graph.dAG)
	fmt.Printf("Incoming %v\n", graph.inComing)
	graph.closeAllChannel()
	assert.Equal(t, 0, -1)
	assert.Equal(t, graph.dAG[1][0], 31)
	assert.Equal(t, len(graph.dAG[2]), 0)
	assert.Equal(t, graph.inComing[31][0], 1)
}

func TesthandlePruneEdge(t *testing.T) {
	graph := NewGraph(2)
	graph.nodes[31] = NewNode(31)
	graph.addEdge(1, []int{31})
	graph.addEdge(2, []int{31})

	graph.addEdge(31, []int{1, 2})

	// cause the grpah contains 1, 2 has 1 as min point toward 31
	graph.nodes[2].min = 1

	graph.Initalizes()
	activateAllNodesHandler(graph)

	time.Sleep(1 * time.Second)
	graph.nodes[1].SinkYoDOWN(graph)
	// time.Sleep(1 * time.Second)
	graph.nodes[2].SinkYoDOWN(graph)

	time.Sleep(1 * time.Second)
	fmt.Printf("OutGoing %v\n", graph.dAG)
	fmt.Printf("Incoming %v\n", graph.inComing)
	graph.closeAllChannel()
	assert.Equal(t, graph.dAG[1][0], 31)
	assert.Equal(t, len(graph.dAG[2]), 0)
	assert.Equal(t, graph.inComing[31][0], 1)
}

func TestFlipEdge(t *testing.T) {
	graph := CreateTestGraph()
	graph.Initalizes()
	// graph.PrintGraph(5, 9)
	flipEdge(graph, 5, 9)
	// graph.PrintGraph(5, 9)
	assert.Equal(t, graph.inComing[9][0], 6)
	assert.Equal(t, graph.dAG[9][0], 5)

	assert.Equal(t, graph.dAG[5], Edge([]int{6, 7}))
	assert.Equal(t, graph.inComing[5][0], 9)
	graph = CreateLectureGraph()
	graph.Initalizes()
	graph.PrintGraph(14, 4)
	flipEdge(graph, 5, 11)
	graph.PrintGraph(14, 4)
}
