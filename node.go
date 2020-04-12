package main

import (
	"fmt"
	"sync"
	"time"
)

type State string

const (
	SINK     State = "SINK"
	INTERNAL State = "INTERNAL"
	SOURCE   State = "SOURCE"
	ASLEEP   State = "ASLEEP"
)

type Node struct {
	id             int
	state          State
	buffer         int
	receivedWeight []int
	mux            sync.Mutex
	wg             sync.WaitGroup
}

func NewNode(id int) *Node {
	node := new(Node)
	node.id = id
	node.state = ASLEEP
	node.buffer = 0
	return node
}

func (u *Node) compare(other *Node) bool {
	return u.id > other.id
}

func (u *Node) handleYoDownMsg(msg Message, g *Graph) {
	u.mux.Lock()
	u.buffer++
	u.receivedWeight = append(u.receivedWeight, msg.candidate)
	u.mux.Unlock()

	min := 100000000000
	if u.buffer == len(g.edges[u.id]) {
		fmt.Printf("[%d] has received all nodes result", u.id)
		// when received all weights and start reply accordingly
		for _, v := range u.receivedWeight {
			if v < min {
				min = v
			}
		}
		fmt.Printf("And The minimal node is %v\n", min)
		for _, v := range u.receivedWeight {
			if v == min {
				g.links[edge{u.id, v}] <- Message{YES, v, u.id}
				continue
			}
			fmt.Printf("Before fliping %v \n", g.dAG)
			flipEdge(g, v, u.id)
			fmt.Printf("After fliping %v \n", g.dAG)
			if len(g.dAG[v]) == 0 {
				fmt.Printf("[%d] become a SINK Node ending pruning current Node\n", v)
				g.links[edge{u.id, v}] <- Message{PRUNE, min, u.id}
				continue
			} else {
				g.links[edge{u.id, v}] <- Message{NO, min, u.id}
			}
		}
		u.mux.Lock()
		u.receivedWeight = []int{}
		u.mux.Unlock()
	}
}

func (u *Node) handleYoUpMessage(msg Message, g *Graph) {

}

func (u *Node) handleMessage(g *Graph) {
	for _, v := range g.edges[u.id] {
		go func(receiver chan Message, v int) {
			for {
				select {
				case msg, ok := <-receiver:
					if !(ok) {
						fmt.Printf("[%d, %d] pipe is closing ", u.id, v)
						break
					}
					switch msg.messagetype {
					case YoDown:
						fmt.Printf("[%d] Received YoDown message %v\n", u.id, msg)
						u.handleYoDownMsg(msg, g)
					case YES:
						fmt.Printf("[%d] Received YES message %v\n", u.id, msg)
					case NO:
						fmt.Printf("[%d] Received NO message %v\n", u.id, msg)
					case PRUNE:
						fmt.Printf("[%d] Received PRUNE message %v\n", u.id, msg)
						g.links[edge{u.id, msg.sender}] <- Message{YoUp, msg.candidate, u.id}
						break
					case YoUp:
						fmt.Printf("[%d] Received PRUNE message %v\n", u.id, msg)
						u.handleYoUpMessage(msg, g)
					}
				default:
					continue
				}
			}
		}(g.links[edge{v.id, u.id}], v.id)

		// go func(receiver chan Message) {
		// 	for {
		// 		select {
		// 		case msg := <-receiver:
		// 			fmt.Printf("[%d] Received Yes No message %v %s\n", u.id, msg, time.Now())
		// 			continue
		// 		default:
		// 			continue
		// 		}
		// 	}
		// }(g.links[edge{u.id, v.id}])
	}
}

func (u *Node) SinkYoDOWN(msg TypesOfMessage, graph *Graph) {
	switch msg {
	case YoDown:
		for _, v := range graph.dAG[u.id] {
			fmt.Printf("Sending from  [%d] to %d %v \n", u.id, v.id, time.Now())
			graph.links[edge{u.id, v.id}] <- Message{YoDown, u.id, u.id}
		}
	}
}
