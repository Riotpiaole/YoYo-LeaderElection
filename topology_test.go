package main

import (
	"testing"
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
