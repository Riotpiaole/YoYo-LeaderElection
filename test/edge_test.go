package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRemoveEdge(t *testing.T) {
	edge := Edge([]*Node{NewNode(1), NewNode(2)})
	edge = removeEdge(&(edge), 1)
	assert.Equal(t, 1, len(edge))
	edge = removeEdge(&(edge), 2)
	assert.Equal(t, 0, len(edge))
	edge = append(edge, NewNode(5))
	assert.Equal(t, 1, len(edge))
}
