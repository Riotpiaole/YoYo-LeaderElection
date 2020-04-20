package main

import (
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
}
