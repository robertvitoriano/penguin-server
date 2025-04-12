package controllers

type GameReceiveEvent string

const (
	StartGame   GameReceiveEvent = "start_game"
	PlayerMoved GameReceiveEvent = "player_moved"
	MessageSent GameReceiveEvent = "message_sent"
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
	IsFlipped    bool     `json:"isFlipped"`
	CurrentState string   `json:"currentState"`
	Token        string   `json:"token"`
	Position     Position `json:"position"`
}

type MessageSentEvent struct {
	Token   string `json:"token"`
	Message string `json:"message"`
}

var receiveEventDefinitions = map[GameReceiveEvent]interface{}{
	StartGame:   StartGameEvent{},
	PlayerMoved: PlayerMovedEvent{},
}
