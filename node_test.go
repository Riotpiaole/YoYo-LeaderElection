package main

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNodeCompare(t *testing.T) {
	node1 := NewNode(1)
	node2 := NewNode(2)
	if node1.compare(node2.id) {
		t.Errorf("Error on compare bteween two nodes ")
	}
}

func TestFindLocalLeader(t *testing.T) {
	msg := []Message{Message{YoDown, 3, 3}, Message{YoDown, 4, 4}}
	node1, node2 := NewNode(3), NewNode(4)
	pickedIndex, min, cnt := node1.findLocalLeader(msg)
	assert.Equal(t, pickedIndex, 0)
	assert.Equal(t, min, 3)
	assert.Equal(t, cnt, 1)

	msg = []Message{Message{YoDown, 3, 3}, Message{YoDown, 4, 4}, Message{YoDown, 3, 5}}
	pickedIndex, min, cnt = node2.findLocalLeader(msg)
	assert.Equal(t, pickedIndex, 0)
	assert.Equal(t, min, 3)
	assert.Equal(t, cnt, 2)
}

func TestFindLocalLeaderEnhance(t *testing.T) {
	// msg := []Message{Message{YoDown, 3, 3}, Message{YoDown, 4, 4}}
	// node1, node2 := NewNode(3), NewNode(4)
}

func TestLecureExample(t *testing.T) {
	// when sink received all of its messages is about to forward upward
	graph := CreateTestGraph()

	graph.Initalizes()

	activateAllNodesHandler(graph)
	graph.YoDown()
	assert.Equal(
		t, LEADER, graph.nodes[1].state,
	)
}

func TestLecureFullYOYOExample(t *testing.T) {
	graph := CreateLectureGraph()
	graph.Initalizes()
	activateAllNodesHandler(graph)
	graph.YoDown()
	assert.Equal(
		t, LEADER, graph.nodes[2].state,
	)
	// fmt.Println(graph.stats.visualizesResult("LectureExample", "Random", ","))
	graph.stats.exportCSV("LectureFullYoyo", "Rand", "./lectureExample.csv")
}

func TestHyperCube(t *testing.T) {
	g := HyperCube(7)
	g.Initalizes()
	g.PrintGraph(1, 2)
	activateAllNodesHandler(g)

	g.YoDown()
	assert.Equal(
		t, LEADER, g.nodes[1].state,
	)
}

func TestCompleteGraph(t *testing.T) {
	g := CreateCompleteTopology(10)
	fmt.Printf("edges %v\n", g.edges)
	g.Initalizes()
	g.PrintGraph(1, 2)
	activateAllNodesHandler(g)

	g.YoDown()
	assert.Equal(
		t, LEADER, g.nodes[1].state,
	)
}

func TestRingTopology(t *testing.T) {
	g := CreateRingTopology(5)
	fmt.Printf("Edge %v, Nodes %v \n", g.edges, g.nodes)
	g.Initalizes()
	g.PrintGraph(1, 2)
	activateAllNodesHandler(g)
	g.YoDown()
}

func TestCircularLadder(t *testing.T) {
	g := CreateLatterTopology(5)
	fmt.Printf("Edge %v, Nodes %v \n", g.edges, g.nodes)
	g.Initalizes()
	g.PrintGraph(1, 2)
	activateAllNodesHandler(g)
	g.YoDown()
}
