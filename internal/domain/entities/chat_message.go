package entities

type ChatMessage struct {
	ID         string `json:"id"`
	SenderId   string `json:"senderId"`
	ReceiverId string `json:"receiverId"`
	Message    string `json:"message"`
	Timestamp  string `json:"timestamp"`
}
