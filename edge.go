package main

import "fmt"

type Edge []*Node

type Link chan Message

func removeEdge(e *Edge, val int) Edge {
	var index int
	for i, v := range *e {
		if v.id == val {
			index = i
			break
		}
		if i == len(*e)-1 {
			fmt.Printf("Given Key is not found in the edge\n")
			fmt.Printf("%v\n", *e)
			return *e
		}
	}
	return append((*e)[:index], (*e)[index+1:]...)
}
