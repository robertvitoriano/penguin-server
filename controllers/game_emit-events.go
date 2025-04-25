package controllers

import "github.com/robertvitoriano/penguin-server/models"

type UpdateOtherPlayerPositionEvent struct {
	Event        string   `json:"event"`
	ID           string   `json:"id"`
	Position     Position `json:"position"`
	CurrentState string   `json:"currentState"`
	IsFlipped    bool     `json:"isFlipped"`
}

type MessageReceivedEvent struct {
	Event    string `json:"event"`
	SenderID string `json:"senderId"`
	Message  string `json:"message"`
}

type AudioChuckReceivedEvent struct {
	Event    string  `json:"event"`
	SenderID string  `json:"senderId"`
	Chunk    []int64 `json:"chunk"`
}

type PlayerWithMessages struct {
	ID           string                `json:"id"`
	Username     string                `json:"username"`
	Color        string                `json:"color"`
	Position     Position              `json:"position"`
	ChatMessages []*models.ChatMessage `json:"chatMessages"`
}

type SetInitialPlayersPositionEvent struct {
	Event   string               `json:"event"`
	Players []PlayerWithMessages `json:"players"`
}
