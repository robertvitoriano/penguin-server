package mysql

import (
	"time"

	"github.com/google/uuid"
	"github.com/robertvitoriano/penguin-server/internal/models"
	"gorm.io/gorm"
)

type ChatMessagesMysqlRepository struct {
	Db *gorm.DB
}

func NewChatMessagesRepository(db *gorm.DB) *ChatMessagesMysqlRepository {

	return &ChatMessagesMysqlRepository{
		Db: db,
	}
}

var ChatMessagesByID = make(map[string][]*models.ChatMessage)

func GetChatMessages(playerId string) []*models.ChatMessage {
	return ChatMessagesByID[playerId]
}

func SaveChatMessage(playerId string, newChatMessage string) {
	chatMessage := &models.ChatMessage{
		ID:         uuid.New().String(),
		SenderId:   playerId,
		ReceiverId: uuid.Nil.String(),
		Message:    newChatMessage,
		Timestamp:  time.Now().Local().GoString(),
	}
	ChatMessagesByID[playerId] = append(ChatMessagesByID[playerId], chatMessage)
}
