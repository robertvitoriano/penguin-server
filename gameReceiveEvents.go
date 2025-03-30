package main

type GameReceiveEvent string

const (
	StartGame   GameReceiveEvent = "start_game"
	PlayerMoved GameReceiveEvent = "player_moved"
)

type Position struct {
	X float64 `json:"x"`
	Y float64 `json:"y"`
}

type StartGameEvent struct {
	Token    string   `json:"token"`
	Position Position `json:"position"`
}

type PlayerMovedEvent struct {
	Token    string   `json:"token"`
	Position Position `json:"position"`
}

var receiveEventDefinitions = map[GameReceiveEvent]interface{}{
	StartGame:   StartGameEvent{},
	PlayerMoved: PlayerMovedEvent{},
}
