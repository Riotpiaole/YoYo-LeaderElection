package main

import "testing"

func TestNodeCompare(t *testing.T) {
	node1 := NewNode(1)
	node2 := NewNode(2)
	if node1.compare(node2) {
		t.Errorf("Error on compare bteween two nodes ")
	}
}
