package redis

import (
	"time"

	"github.com/google/uuid"
	"github.com/robertvitoriano/penguin-server/internal/domain/entities"
)

var ChatMessagesByID = make(map[string][]*entities.ChatMessage)

func GetChatMessages(playerId string) []*entities.ChatMessage {
	return ChatMessagesByID[playerId]
}

func SaveChatMessage(playerId string, newChatMessage string) {
	chatMessage := &entities.ChatMessage{
		ID:         uuid.New().String(),
		SenderId:   playerId,
		ReceiverId: uuid.Nil.String(),
		Message:    newChatMessage,
		Timestamp:  time.Now().Local().GoString(),
	}
	ChatMessagesByID[playerId] = append(ChatMessagesByID[playerId], chatMessage)
}
