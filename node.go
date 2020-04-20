package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

type State string

const (
	SINK     State = "SINK"
	INTERNAL State = "INTERNAL"
	SOURCE   State = "SOURCE"
	ASLEEP   State = "ASLEEP"
	LEADER   State = "LEADER"
)

type Node struct {
	id    int
	state State

	upwardMsgs  []Message
	muxInComing sync.Mutex

	downWardMsgs []Message
	mux          sync.Mutex //prevent change downWardMsgs

	min int // represent min Canadiate
}

func (u *Node) printNode(g *Graph, msg Message, chanDest int) {
	u.mux.Lock()
	u.muxInComing.Lock()
	fmt.Println()
	fmt.Printf("[%d, %s, min %d] Received %s message [ %v ], through chan [%d, %d]\n", u.id, u.state, u.min, msg.messagetype, msg, u.id, chanDest)
	fmt.Printf("[%d, %s] \treceived UpwardMsg: %v\tDownWardMsg: %v\n", u.id, u.state, u.upwardMsgs, u.downWardMsgs)
	fmt.Printf("[%d, %s] \tStateCheck OutGoing: %v\tIncoming: %v\n", u.id, u.state, g.dAG[u.id], g.inComing[u.id])
	fmt.Println()
	fmt.Println()
	u.mux.Unlock()
	u.muxInComing.Unlock()
}

func NewNode(id int) *Node {
	node := new(Node)
	node.id = id
	node.state = ASLEEP
	node.min = id
	node.downWardMsgs = []Message{}
	node.upwardMsgs = []Message{}
	return node
}

func (u *Node) compare(other int) bool {
	return u.id > other
}

func (u *Node) findLocalLeader(downWardMsgs []Message) (int, int, int) {
	min, cnt, pickedIndex := u.min, 0, -1
	yesMessage := []Message{}
	for _, receivedMsg := range downWardMsgs {
		if receivedMsg.candidate < min {

			u.min = receivedMsg.candidate
			min = receivedMsg.candidate
		}
	}
	for _, downWardMsgs := range downWardMsgs {
		if downWardMsgs.candidate == min {
			cnt++
		}
	}

	for i, msg := range downWardMsgs {
		if msg.candidate == min {
			if msg.sender == min {
				pickedIndex = i
			}
			yesMessage = append(yesMessage, msg)
		}
	}
	if cnt > 1 {
		if pickedIndex == -1 { // haven't found a sender is the Min
			fmt.Printf("[%d , %s] Checking repeatedYesMessages %v and cnt %d \n", u.id, u.state, len(yesMessage), cnt)
			if len(downWardMsgs) == 0 {
				pickedIndex = -1
			} else {
				pickedIndex = rand.Intn(len(yesMessage))
				for i, msg := range downWardMsgs {
					if compareMessage(msg, yesMessage[pickedIndex]) {
						pickedIndex = i
					}
				}
			}
		}
	} else if cnt == 1 {
		for i, msg := range downWardMsgs {
			if msg.candidate == min {
				pickedIndex = i
			}
		}
	}

	return pickedIndex, min, cnt
}

func (u *Node) updateState(g *Graph) State {
	g.muxDAG.Lock()
	u.mux.Lock()

	previousState := u.state
	numsIncoming, numsOutGoing := len(g.inComing[u.id]), len(g.dAG[u.id])

	if numsOutGoing == 0 {
		u.state = SINK
	} else if numsOutGoing == 0 && numsIncoming == 0 {
		u.state = ASLEEP
	} else if numsOutGoing > 0 && numsIncoming > 0 {
		u.state = INTERNAL
	}
	u.mux.Unlock()

	if previousState != u.state && previousState == SOURCE {
		g.muxSRC.Lock()
		fmt.Printf("[%d] Remove %d from source %v\n", u.id, u.id, g.source)
		g.source = removeByVal(g.source, u.id)
		fmt.Printf("How many time does remove source is ran\n")
		g.muxSRC.Unlock()
		// fmt.Printf("How many time does remove source is ran")
		g.sourcewg.Done()
	}
	g.muxDAG.Unlock()

	return previousState
}

func (u *Node) handleMessage(g *Graph) {
	// fmt.Printf("OutGoing %v \nInComing %v", )
	for _, v := range g.edges[u.id] {
		fmt.Printf("[%d] Channel Listening from [%d, %d]\n", u.id, v, u.id)
		go func(receiver chan Message, v int) {
			for {
				select {
				case msg, ok := <-receiver:
					if !(ok) {
						fmt.Printf("[%d, %d] pipe is closing \n", u.id, v)
						fmt.Printf("checking current \n\toutgoing %v\t\n incoming %v",
							g.dAG[u.id], g.inComing[u.id])
						return
					}

					switch msg.messagetype {

					case YoDown:
						u.handleYoDownMsg(msg, g)
						// u.printNode(g, msg, v)

					case YES:
						u.handleUpwardMessages(msg, g)
						// u.printNode(g, msg, v)
						//continue send the YoDown Message to all is outGoing Link

					case NO:
						u.handleUpwardMessages(msg, g)
						// u.printNode(g, msg, v)
						//update the min and change of the state from source to internal
						// or sink

					case YESPRUNE:
						// pruning unwanted repeated Link
						u.handleUpwardMessages(msg, g)
						fmt.Printf("[%d, %d] stops listening \n", u.id, v)
						break
					}
				default:
					continue
				}
			}
		}(g.links[edge{v, u.id}], v)
	}
}

func (u *Node) SinkYoDOWN(graph *Graph) {
	if u.state != SOURCE {
		fmt.Printf("[%d] Should be sending a YODON from a non source node \n", u.id)
		return
	}
	for _, v := range graph.dAG[u.id] {
		fmt.Printf("Sending from  [%d] to %d %v\n", u.id, v, time.Now())
		graph.links[edge{u.id, v}] <- Message{YoDown, u.min, u.id}
	}
}

func (u *Node) handleSinkUpward(g *Graph) {
	// only handling SINK
	u.updateState(g)
	pickedIndex, min, _ := u.findLocalLeader(u.downWardMsgs)
	if min < u.min {
		u.min = min
	}
	fmt.Printf("[%s,%d, min:%d] checking for msg %v and index %d\n", u.state, u.id, u.min, u.downWardMsgs, pickedIndex)

	for index, msg := range u.downWardMsgs {
		if index != pickedIndex && u.min == msg.candidate {
			// synchronization mechanics
			u.forwardMessage(msg.sender, YESPRUNE, u.min, g)
			u.downWardMsgs = removeMsgQueue(u.downWardMsgs, msg.sender)
		} else if msg.candidate != u.min {
			u.forwardMessage(msg.sender, NO, u.min, g)
			u.downWardMsgs = removeMsgQueue(u.downWardMsgs, msg.sender)
		}
	}
	fmt.Printf("Checking From SinkDown AWAIT\n")
	u.printNode(g, Message{CHECK, 300, 300}, 30)
	u.updateState(g)
	if len(u.downWardMsgs) == 1 {
		msg := u.downWardMsgs[0]
		if u.state == SINK && len(g.inComing[u.id]) == 1 {
			fmt.Printf("[%d, %s] has become a sink with only one outGoing Edges sending YESPRUNE\n",
				u.id, u.state)
			u.forwardMessage(msg.sender, YESPRUNE, u.min, g)
			u.downWardMsgs = removeMsgQueue(u.downWardMsgs, msg.sender)
		} else {
			u.forwardMessage(msg.sender, YES, u.min, g)
			u.downWardMsgs = removeMsgQueue(u.downWardMsgs, msg.sender)
		}
	}
	fmt.Printf("Checking From SinkDown AWAIT SHOULD HITS HERE\n")

	u.mux.Lock()
	u.upwardMsgs = []Message{}
	u.downWardMsgs = []Message{}
	u.mux.Unlock()
}

func (u *Node) handleUpwardMessages(msg Message, g *Graph) {
	u.muxInComing.Lock()
	u.upwardMsgs = append(u.upwardMsgs, msg)
	u.muxInComing.Unlock()

	if len(u.upwardMsgs) == len(g.dAG[u.id]) {

		fmt.Printf("[%d, %s] has received all upward message handling\n", u.id, u.state)
		pickedIndex, min, cnt := u.findLocalLeader(u.upwardMsgs)
		if min < u.min {
			u.min = min
		}

		for _, msg := range u.upwardMsgs {
			if msg.messagetype == YESPRUNE {
				g.pruneEdge(u.id, []int{msg.sender})
				cnt--
			} else if msg.messagetype == NO {
				flipEdge(g, u.id, msg.sender)
			}
		}
		if u.leaderChecking(g) {
			g.source = removeByVal(g.source, u.id)
			g.sourcewg.Done()
			u.state = LEADER
			fmt.Printf("[%d, %s] has elected as Leader\n\n", u.id, u.state)
			fmt.Printf("How many time does remove source is ran\n")
			return
		}

		u.updateState(g)
		fmt.Printf("[%d, %s] pickedIndex %d , min %d, cnt %d\n", u.id, u.state, pickedIndex, min, cnt)
		u.printNode(g, msg, msg.sender)

		switch u.state {
		case SOURCE:
			for _, upwardMsg := range u.upwardMsgs {
				if upwardMsg.messagetype == YESPRUNE {
					fmt.Printf("[%d , %s] pruning the edge according to messages %v\n", u.id, u.state, upwardMsg)
					g.pruneEdge(u.id, []int{upwardMsg.sender})
				}
			}
			// update state and flip all the no edge

			for _, upwardMsg := range u.upwardMsgs {
				if upwardMsg.messagetype == NO {
					u.muxInComing.Lock()
					flipEdge(g, u.id, upwardMsg.sender)
					u.min = upwardMsg.candidate
					u.upwardMsgs = removeMsgQueue(u.upwardMsgs, upwardMsg.sender)
					u.muxInComing.Unlock()
				}
			}
			u.updateState(g)
			/**
			* Recursive Loop perform additional YoDown
			 */
			fmt.Printf("[%d, %s, min:%d] has Received All upward incoming candiates\n", u.id, u.state, u.min)
			fmt.Printf("[%d, %s, min:%d] upward %v downward %v\n", u.id, u.state, u.min, u.upwardMsgs, u.downWardMsgs)
			fmt.Printf("[%d, %s, min:%d] should prun Edges? %d and index is %d min %d\n",
				u.id, u.state, u.min, cnt, pickedIndex, min)
			if u.state == SOURCE {
				fmt.Printf("How many time does remove source is ran\n")
				g.sourcewg.Done()
			}
		case INTERNAL:
			u.pruneRepeatedCanadiateEdge(g)
			for index, inComingMsg := range u.downWardMsgs {
				if inComingMsg.candidate != u.min {
					u.forwardMessage(inComingMsg.sender, NO, u.min, g)
				} else if inComingMsg.candidate == u.min && cnt <= 1 {
					u.forwardMessage(inComingMsg.sender, YES, u.min, g)
				} else if cnt > 1 && pickedIndex != -1 {
					if inComingMsg.candidate == u.min && index != pickedIndex {
						u.forwardMessage(inComingMsg.sender, YESPRUNE, u.min, g)
					}
				}
			}
			u.updateState(g)
		case SINK:
			u.pruneRepeatedCanadiateEdge(g)
			if len(g.inComing[u.id]) == 1 {
				u.forwardMessage(g.inComing[u.id][0], YESPRUNE, u.min, g)
			} else {
				for _, inComingMsg := range u.downWardMsgs {
					if inComingMsg.candidate != u.min {
						u.forwardMessage(inComingMsg.sender, NO, u.min, g)
					} else {
						u.forwardMessage(inComingMsg.sender, YES, u.min, g)
					}
				}
			}
			u.updateState(g)
		}
		u.muxInComing.Lock()
		u.mux.Lock()
		u.downWardMsgs = []Message{}
		u.upwardMsgs = []Message{}
		u.mux.Unlock()
		u.muxInComing.Unlock()
	}
	g.GlobalUpdate()
}

func (u *Node) forwardMessage(sender int, types TypesOfMessage, candidate int, g *Graph) {
	msg := Message{types, candidate, u.id}
	fmt.Printf("[%d] sending Message %v through channel [%d , %d]\n",
		u.id, msg, u.id, sender)
	// sending from u to sender with given message type
	g.links[edge{u.id, sender}] <- msg
}

func (u *Node) replyMessage(sender int, types TypesOfMessage, candidate int, g *Graph) {
	msg := Message{types, candidate, u.id}
	fmt.Printf("[%d] replying Message %v through channel [%d , %d]\n",
		u.id, msg, sender, u.id)
	g.links[edge{sender, u.id}] <- msg
}

func (u *Node) pruneRepeatedCanadiateEdge(g *Graph) {
	// double check
	//================================================================
	// handling reply when the node is still an internal Node
	//================================================================

	pickedIndex, min, cnt := u.findLocalLeader(u.downWardMsgs)
	u.min = min

	/**
	* Random select a edge from a edge sent the same weight
	* For all edge prune the unNeeded Channel
	 */
	if cnt > 1 {
		for index, msg := range u.downWardMsgs {
			if index != pickedIndex && msg.candidate == u.min {
				u.forwardMessage(msg.sender, YESPRUNE, u.min, g)
				u.downWardMsgs = removeMsgQueue(u.downWardMsgs, msg.sender)
			}
		}
	}
}

func (u *Node) handleYoDownMsg(msg Message, g *Graph) {
	u.mux.Lock()
	// u.downWardMsgs = append(u.downWardMsgs, msg.candidate)
	u.downWardMsgs = append(u.downWardMsgs, msg)
	u.mux.Unlock()
	/*
	 * If received all of the incoming message
	 */

	if len(u.downWardMsgs) == len(g.inComing[u.id]) {
		/**
		 * Prune yourself if yourself is the only sink Node
		 * Update your current state
		 */

		_, min, _ := u.findLocalLeader(u.downWardMsgs)
		if min < u.min {
			u.min = min
		}
		fmt.Printf("[%d, %s] has Received All downward incoming candiates and the min is %d \n", u.id, u.state, min)
		switch u.state {

		case INTERNAL:

			fmt.Printf("[%d] And The minimal node is %v\n", u.id, min)
			/*****
			* Forward your current min to all outGoing Edges with YoDown Message
			* Update your current state
			 */
			for _, outGoingLink := range g.dAG[u.id] {
				fmt.Printf("[%d] sending to %d with msg %v\t", u.id, outGoingLink, Message{YoDown, u.min, u.id})
				u.forwardMessage(outGoingLink, YoDown, u.min, g)
			}
			fmt.Println()

		case SINK:
			// what if different one
			// when i received multiple same candidate
			fmt.Printf("[%d, %s] handling YODOWN message\n", u.id, u.state)
			fmt.Printf("[%d, %s]\treceived \n\tUpwardMsg: %v\n\tDownWardMsg: %v \n", u.id, u.state, u.upwardMsgs, u.downWardMsgs)
			u.handleSinkUpward(g)
			u.mux.Lock()
			u.muxInComing.Lock()
			u.downWardMsgs = []Message{}
			u.upwardMsgs = []Message{}
			u.mux.Unlock()
			u.muxInComing.Unlock()
		}

	}
	g.GlobalUpdate()
}

func (u *Node) leaderChecking(g *Graph) bool {
	fmt.Printf("[%d, %s] min : %d is doing a leader checking \n", u.id, u.state, u.min)
	if u.min != u.id {
		return false
	}
	for _, v := range g.dAG {
		if len(v) != 0 {
			return false
		}
	}
	for _, v := range g.inComing {
		if len(v) != 0 {
			return false
		}
	}
	return true
}
