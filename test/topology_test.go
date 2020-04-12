package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCreateTopology(t *testing.T) {
	graph := CreateTestGraph()
	for i := 1; i <= 9; i++ {

		if graph.nodes[i].id != i && graph.nodes[i].state != ASLEEP {
			t.Errorf("Node initalizes error")
		}
		for _, neighbour := range graph.edges[i] {
			if neighbour == nil {
				t.Errorf("Edge not properly initalizes")
			}
		}

	}
}

func TestAddNewEdge(t *testing.T) {
	graph := NewGraph(2)
	graph.addEdge(1, []int{2})
	graph.addEdge(2, []int{1})
	if graph.edges[1][0].id != 2 {
		t.Errorf("Edge not properly initalizes")
	}
}

func TestInitalizes(t *testing.T) {
	graph := CreateTestGraph()
	graph.Initalizes()
	expected := []int{6, 7, 9}
	compare := []int{}
	for _, v := range graph.dAG[5] {
		compare = append(compare, v.id)
	}
	assert.Equal(t, graph.nodes[1].state, SOURCE)
	assert.Equal(t, graph.nodes[8].state, SINK)
	assert.Equal(t, graph.nodes[6].state, INTERNAL)
	assert.Equal(t, compare, expected, "The outgoing of 5 should be this")
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
	graph := NewGraph(2)
	graph.dAG[1] = append(graph.dAG[1], graph.nodes[2])
	flipEdge(graph, 1, 2)
	assert.Equal(t, graph.dAG[2][0].id, 1)
	assert.Equal(t, len(graph.dAG[1]), 0)
}
