package main

type TypesOfMessage string

const (
	YoDown TypesOfMessage = "DOWN"
	YoUp   TypesOfMessage = "UP"

	PRUNE TypesOfMessage = "PRUNE"
	YES   TypesOfMessage = "YES"
	NO    TypesOfMessage = "NO"
)

type Message struct {
	messagetype TypesOfMessage
	candidate   int
	sender      int
}
