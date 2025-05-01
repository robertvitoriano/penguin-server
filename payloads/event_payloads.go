package payloads

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

type WebrtcAnswerReceivedEvent struct {
	Event  string `json:"event"`
	Answer string `json:"answer"`
}
type RTCIceCandidate struct {
	Address          string  `json:"address"`
	Candidate        string  `json:"candidate"`
	Component        string  `json:"component"`
	Foundation       string  `json:"foundation"`
	Port             int     `json:"port"`
	Priority         uint32  `json:"priority"`
	Protocol         string  `json:"protocol"`
	RelatedAddress   *string `json:"relatedAddress"`
	RelatedPort      *int    `json:"relatedPort"`
	RelayProtocol    *string `json:"relayProtocol"`
	SdpMLineIndex    int     `json:"sdpMLineIndex"`
	SdpMid           string  `json:"sdpMid"`
	TcpType          *string `json:"tcpType"`
	Type             string  `json:"type"`
	URL              *string `json:"url"`
	UsernameFragment string  `json:"usernameFragment"`
}
type WebrtcCandidateReceidEvent struct {
	Event     string          `json:"event"`
	Candidate RTCIceCandidate `json:"candidate"`
}
type OfferAnswer struct {
	SDP  string `json:"sdp"`
	Type string `json:"type"`
}
type WebRTCOfferReceivedEvent struct {
	Event string      `json:"event"`
	Offer OfferAnswer `json:"offer"`
}
type WebRTCAnswerReceivedEvent struct {
	Event  string      `json:"event"`
	Answer OfferAnswer `json:"answer"`
}

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
type AudioChunkSentEvent struct {
	Chunk   []int64 `json:"chunk"`
	Message string  `json:"message"`
	Token   string  `json:"token"`
}

type WebRTCCandidateSentEvent struct {
	Candidate RTCIceCandidate `json:"candidate"`
	Token     string          `json:"token"`
}
type WebRTCOfferSentEvent struct {
	Offer OfferAnswer `json:"offer"`
	Token string      `json:"token"`
}
type WebRTCAnswerSentEvent struct {
	Answer OfferAnswer `json:"answer"`
	Token  string      `json:"token"`
}
