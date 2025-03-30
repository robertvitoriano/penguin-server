package main

type GameEmitEvent string

const (
	PlayerJoined GameEmitEvent = "PlayerJoined"
)

var emitEventDefinitions = map[GameEmitEvent]interface{}{
	PlayerJoined: struct {
		PlayerID string `json:"playerId"`
		Username string `json:"username"`
	}{},
}
