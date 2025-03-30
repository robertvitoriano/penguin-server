package main

type GameReceiveEvent string

const (
	StartGame   GameReceiveEvent = "start_game"
	PlayerMoved GameReceiveEvent = "player_moved"
)

var receiveEventDefinitions = map[GameReceiveEvent]interface{}{
	StartGame: struct {
		Event string `json:"event"`
		Token string `json:"token"`
	}{},
	PlayerMoved: struct {
		Token string  `json:"token"`
		X     float64 `json:"x"`
		Y     float64 `json:"y"`
	}{},
}
