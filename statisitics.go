package main

import (
	"fmt"
	"os"
	"sync"
	"time"
)

var msgSequence = []TypesOfMessage{
	YoDown,
	YES,
	YESPRUNE,
	NO,
}

type MessageCnter struct {
	count map[TypesOfMessage]int
}

type Statisics struct {
	mux           *sync.Mutex
	yoDownTime    map[int]time.Duration
	totalDuration time.Duration
	numsYoYo      int
	yoStages      map[int]*MessageCnter
	totlaMsgCnter *MessageCnter
	numsNode      int
	numsEdge      int
}

func NewMessageCnter() *MessageCnter {
	msgCnter := new(MessageCnter)
	msgCnter.count = map[TypesOfMessage]int{
		YoDown:   0,
		YES:      0,
		NO:       0,
		YESPRUNE: 0,
	}
	return msgCnter
}

func (msgcnt *MessageCnter) UnpackMsgCnter(sep string) string {
	result := ""
	for _, msgType := range msgSequence {
		if sep != "," && msgType == NO {
			result += "       "
		}
		result += fmt.Sprintf("%s%d", sep, msgcnt.count[msgType])
	}
	return result
}

func NewStatistiscs(g *Graph) *Statisics {
	stats := new(Statisics)
	stats.numsYoYo = 1 // indicate yoyo stages it
	stats.yoStages = make(map[int]*MessageCnter)
	stats.yoDownTime = map[int]time.Duration{
		0: time.Duration(0),
	}

	stats.totlaMsgCnter = NewMessageCnter()
	stats.mux = new(sync.Mutex)
	stats.yoStages[stats.numsYoYo] = NewMessageCnter()
	fmt.Printf("%v\n", stats)
	return stats
}

func (stats *Statisics) logTimeTaken(timeDiff time.Duration) {
	stats.yoDownTime[stats.numsYoYo] = timeDiff
}

func (stats *Statisics) increment() {
	stats.mux.Lock()
	stats.numsYoYo++
	stats.yoStages[stats.numsYoYo] = NewMessageCnter()
	stats.yoDownTime[stats.numsYoYo] = time.Duration(0)
	stats.mux.Unlock()

}

func (stats *Statisics) addMessage(msg Message) {
	stats.mux.Lock()
	stats.yoStages[stats.numsYoYo].count[msg.messagetype]++
	stats.totlaMsgCnter.count[msg.messagetype]++
	stats.mux.Unlock()
}

func (stats *Statisics) parseGraphProperties(g *Graph) {
	stats.numsNode = len(g.nodes)
	stats.numsEdge = g.numsEdges
}

func (stats *Statisics) visualizesResult(
	title string,
	typesOfTopology string,
	sep string,
	filename string) string {

	result := ""
	_, err := os.Stat(filename)
	if err != nil || filename == "" || sep != "," {
		fmt.Printf("%s not exists creating one...\n", filename)
		result += title + fmt.Sprintf("\nDescription%s#YoDown%s#Yes%s#YesPrune%s#No%sDuration%s#numsNode%s#nusmEdges\n",
			sep, sep, sep, sep, sep, sep, sep)
	}

	for stages, msgCnter := range stats.yoStages {
		result += fmt.Sprintf("%s_%d", typesOfTopology, stages) + msgCnter.UnpackMsgCnter(sep) + sep
		result += fmt.Sprintf("%s\n", stats.yoDownTime[stages])
	}
	result += fmt.Sprintf("%s", typesOfTopology) + stats.totlaMsgCnter.UnpackMsgCnter(sep) + sep
	result += fmt.Sprintf("%s%s", stats.totalDuration, sep)
	result += fmt.Sprintf("%d%s", stats.numsNode, sep)
	if sep != "," {
		result += "        "
	}
	result += fmt.Sprintf("%d\n", stats.numsEdge)
	return result
}

func (stats *Statisics) exportCSV(title string, typesOfTopology string, filename string) {
	experiement := stats.visualizesResult(title, typesOfTopology, ",", filename)
	f, err := os.OpenFile(filename, os.O_APPEND|os.O_WRONLY, 0600)
	if err != nil {
		f, err = os.Create(filename)
	}

	if _, err = f.WriteString(experiement); err != nil {
		panic(err)
	}
	defer f.Close()

}
