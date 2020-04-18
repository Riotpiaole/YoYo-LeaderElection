package main

type TypesOfMessage string

const (
	YoDown TypesOfMessage = "DOWN"
	YoUp   TypesOfMessage = "UP"
	CHECK  TypesOfMessage = "CHECK"

	PRUNE    TypesOfMessage = "PRUNE"
	YES      TypesOfMessage = "YES"
	YESPRUNE TypesOfMessage = "YESPRUNE"
	NO       TypesOfMessage = "NO"
)

type Message struct {
	messagetype TypesOfMessage
	candidate   int
	sender      int
}

func compareMessage(m1 Message, m2 Message) bool {
	return m1.messagetype == m2.messagetype && m1.candidate == m2.candidate && m1.sender == m2.sender
}

func findMinAmongMessage(msgs []Message, min int) int {
	for _, receivedMsg := range msgs {
		if receivedMsg.candidate < min {
			min = receivedMsg.candidate
		}
	}
	return min
}
