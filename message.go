package main

type TypesOfMessage string

const (
	YoDown TypesOfMessage = "DOWN"
	YoUp   TypesOfMessage = "Up"
	PRUNE  TypesOfMessage = "PRUNE"
)

type Message struct {
}
