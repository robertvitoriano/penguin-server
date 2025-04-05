package controllers

import (
	"github.com/robertvitoriano/penguin-server/models"
)

type UpdateOtherPlayerPositionEvent struct {
	Event        string   `json:"event"`
	ID           string   `json:"id"`
	Position     Position `json:"position"`
	CurrentState string   `json:"currentState"`
	IsFlipped    bool     `json:"isFlipped"`
}

type SetInitialPlayersPositionEvent struct {
	Event   string           `json:"event"`
	Players []*models.Player `json:"players"`
}
