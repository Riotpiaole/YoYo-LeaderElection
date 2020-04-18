package main

import (
	"testing"
	"time"

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
	assert.Equal(t, pickedIndex, -1)
	assert.Equal(t, min, 3)
	assert.Equal(t, cnt, 1)

	msg = []Message{Message{YoDown, 3, 3}, Message{YoDown, 4, 4}, Message{YoDown, 3, 5}}
	pickedIndex, min, cnt = node2.findLocalLeader(msg)
	assert.Equal(t, pickedIndex, 0)
	assert.Equal(t, min, 3)
	assert.Equal(t, cnt, 2)
}

func TestHandleSinkUpward(t *testing.T) {
	// when sink received all of its messages is about to forward upward
	graph := NewGraph(3)
	graph.addEdge(1, []int{3})
	graph.addEdge(2, []int{3})
	graph.addEdge(3, []int{1, 2})
	graph.Initalizes()

	activateAllNodesHandler(graph)
	graph.nodes[3].downWardMsgs =
		[]Message{Message{YoDown, 1, 2}, Message{YoDown, 1, 1}}

	graph.nodes[3].handleSinkUpward(graph)
	graph.PrintGraph(1, 2)
	time.Sleep(1 * time.Second)
	assert.Equal(t, Edge([]int{}), graph.dAG[2])
	assert.Equal(t, Edge([]int{}), graph.dAG[1])
	assert.Equal(t, Edge([]int{}), graph.inComing[3])
	assert.Equal(t, LEADER, graph.nodes[1].state)
}

func TestHandleSinkUpwardNormalCase(t *testing.T) {
	// when sink received all of its messages is about to forward upward
	graph := NewGraph(3)
	graph.addEdge(1, []int{3})
	graph.addEdge(2, []int{3})
	graph.addEdge(3, []int{1, 2})
	graph.Initalizes()

	activateAllNodesHandler(graph)
	graph.nodes[3].downWardMsgs =
		[]Message{Message{YoDown, 2, 2}, Message{YoDown, 1, 1}}

	graph.nodes[3].handleSinkUpward(graph)
	graph.PrintGraph(1, 2)
	time.Sleep(1 * time.Second)
	assert.Equal(t, graph.dAG[2], Edge([]int{}))
	assert.Equal(t, graph.dAG[1], Edge([]int{3}))
	assert.Equal(t, graph.dAG[3], Edge([]int{2}))
	assert.Equal(t, graph.inComing[3], Edge([]int{1}))
}

func TestHandleSinkUpwardNormal(t *testing.T) {
	// when sink received all of its messages is about to forward upward
	graph := NewGraph(3)
	graph.addEdge(3, []int{1, 2})
	graph.addEdge(2, []int{3})
	graph.addEdge(3, []int{1, 2})
	graph.Initalizes()

	activateAllNodesHandler(graph)
	graph.nodes[3].downWardMsgs =
		[]Message{Message{YoDown, 2, 2}, Message{YoDown, 1, 1}}

	graph.nodes[3].handleSinkUpward(graph)
	graph.PrintGraph(1, 2)
	time.Sleep(1 * time.Second)
	assert.Equal(t, graph.dAG[2], Edge([]int{}))
	assert.Equal(t, graph.dAG[1], Edge([]int{}))
	assert.Equal(t, graph.dAG[3], Edge([]int{}))
	assert.Equal(t, graph.inComing[3], Edge([]int{}))
}

func TestHandleSinkUpwardPruneOnly(t *testing.T) {
	// when sink received all of its messages is about to forward upward
	graph := NewGraph(3)
	graph.addEdge(1, []int{2})
	graph.addEdge(2, []int{1})
	graph.Initalizes()

	activateAllNodesHandler(graph)
	graph.nodes[2].downWardMsgs =
		[]Message{Message{YoDown, 1, 1}}

	graph.nodes[2].handleSinkUpward(graph)
	time.Sleep(1 * time.Second)
	assert.Equal(t, Edge([]int{}), graph.dAG[2])
	assert.Equal(t, Edge([]int{}), graph.dAG[1])
	assert.Equal(t, Edge([]int{}), graph.inComing[2])
	assert.Equal(t, Edge([]int{}), graph.inComing[2])
	assert.Equal(t, LEADER, graph.nodes[1].state)
}

func TestPruneMultipleRepeatedEdges(t *testing.T) {
	// when sink received all of its messages is about to forward upward
	graph := NewGraph(4)
	graph.addEdge(1, []int{3})
	graph.addEdge(2, []int{3})
	graph.addEdge(4, []int{3})
	graph.addEdge(3, []int{1, 2, 4})

	graph.Initalizes()
	for _, sourceNode := range graph.source {
		graph.nodes[sourceNode].SinkYoDOWN(graph)
	}
	activateAllNodesHandler(graph)
	time.Sleep(1 * time.Second)
	assert.Equal(t, LEADER, graph.nodes[1].state)
}

func TestLecureExample(t *testing.T) {
	// when sink received all of its messages is about to forward upward
	// graph := CreateTestGraph()

	// graph.Initalizes()

	// activateAllNodesHandler(graph)

	// for _, sourceNode := range graph.source {
	// 	graph.nodes[sourceNode].SinkYoDOWN(graph)
	// }
	// time.Sleep(2 * time.Second)
	// // c := exec.Command("clear")
	// // c.Stdout = os.Stdout
	// // c.Run()
	// graph.PrintGraph(1, 2)
	// for _, sourceNode := range graph.source {
	// 	graph.nodes[sourceNode].SinkYoDOWN(graph)
	// }
	// time.Sleep(2 * time.Second)
	// c = exec.Command("clear")
	// c.Stdout = os.Stdout
	// c.Run()
	// for _, sourceNode := range graph.source {
	// 	graph.nodes[sourceNode].SinkYoDOWN(graph)
	// }
	// time.Sleep(2 * time.Second)
	graph := CreateLectureGraph()
	graph.Initalizes()

	activateAllNodesHandler(graph)
	graph.YoDown()
	graph.PrintGraph(0, 122)

}

func TestLecureFullYOYOExample(t *testing.T) {

	// assert.Equal(
	// 	t, graph.nodes[2].state, LEADER,
	// )
}
