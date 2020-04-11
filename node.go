package main

type State string

const (
	SINK     State = "SINK"
	INTERNAL State = "INTERNAL"
	SOURCE   State = "SOURCE"
	ASLEEP   State = "ASLEEP"
)

type Node struct {
	id    int
	state State
}

func NewNode(id int) *Node {
	node := new(Node)
	node.id = id
	node.state = ASLEEP
	return node
}

func (n *Node) compare(other *Node) bool {
	return n.id > other.id
}

func (n *Node) handle(msg Message) {

}
