package main

type Edge []int

type Link chan Message

func removeEdge(e Edge, val int) Edge {
	if len(e) == 0 {
		return e
	}
	var index int
	for i, v := range e {
		if v == val {
			index = i
			break
		}
		if i == len(e)-1 {
			return e
		}
	}
	return append(e[:index], e[index+1:]...)
}

func removeInt(slice []int, index int) []int {
	return append(slice[:index], slice[index+1:]...)
}

func removeByVal(slice []int, val int) []int {
	for i, v := range slice {
		if val == v {
			return append(slice[:i], slice[i+1:]...)
		}
	}
	return slice
}

func removeMsgQueue(slice []Message, val int) []Message {
	for i, v := range slice {
		if v.sender == val {
			return append(slice[:i], slice[i+1:]...)
		}
	}
	return slice
}

func appendEdge(slice Edge, val int) Edge {
	for _, v := range slice {
		if v == val {
			return slice
		}
	}
	return append(slice, val)
}
