package main

import (
	"github.com/robertvitoriano/penguin-server/models"
)

type UpdateOtherPlayerPositionEvent struct {
	Event    string   `json:"event"`
	ID       string   `json:"id"`
	Position Position `json:"position"`
}

type SetInitialPlayersPositionEvent struct {
	Event   string           `json:"event"`
	Players []*models.Player `json:"players"`
}
