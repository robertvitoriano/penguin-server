package handler

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/robertvitoriano/penguin-server/internal/domain/entities"
	"github.com/robertvitoriano/penguin-server/internal/domain/events"
	"github.com/robertvitoriano/penguin-server/internal/domain/payloads"
	"github.com/robertvitoriano/penguin-server/internal/domain/repository"
	"github.com/robertvitoriano/penguin-server/internal/infra/repository/redis"
	"github.com/robertvitoriano/penguin-server/internal/utils"
)

type Websocket struct {
	Connections                 map[*websocket.Conn]bool
	addConnection               chan *websocket.Conn
	removeConnection            chan *websocket.Conn
	broadcast                   chan broadcastMessage
	playerLiveDataRepository    repository.PlayerRepository
	playerPersistencyRepository repository.PlayerRepository
}

type BaseMessage struct {
	Event string `json:"event"`
}

type broadcastMessage struct {
	message   []byte
	exclude   *websocket.Conn
	errorChan chan error
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func NewWebsocket(playerLiveDataRepository repository.PlayerRepository, playerPersistencyRepository repository.PlayerRepository) *Websocket {
	ws := &Websocket{
		Connections:                 make(map[*websocket.Conn]bool),
		addConnection:               make(chan *websocket.Conn),
		removeConnection:            make(chan *websocket.Conn),
		broadcast:                   make(chan broadcastMessage),
		playerLiveDataRepository:    playerLiveDataRepository,
		playerPersistencyRepository: playerPersistencyRepository,
	}

	go ws.hub()
	return ws
}

func (ws *Websocket) ServeWebsocket(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}

	if !ws.Connections[conn] {
		ws.AddConnection(conn)
	}

	defer ws.RemoveConnection(conn)
	defer conn.Close()

	for {
		_, data, err := conn.ReadMessage()
		if err != nil {
			log.Println("read error:", err)
			break
		}

		var message BaseMessage
		json.Unmarshal(data, &message)
		fmt.Printf("RECEIVED WEBSOCKET_MESSAGE:%v\n", message)
		ws.handleIncomingMessage(conn, events.GameReceiveEvent(message.Event), data)
	}
}
func (ws *Websocket) AddConnection(conn *websocket.Conn) {
	ws.addConnection <- conn
}

func (ws *Websocket) RemoveConnection(conn *websocket.Conn) {

	ws.removeConnection <- conn
}

func (ws *Websocket) Broadcast(message []byte) {

	for conn := range ws.Connections {
		if err := conn.WriteMessage(websocket.TextMessage, message); err != nil {
			log.Println("Write error:", err)
			conn.Close()
			delete(ws.Connections, conn)
		}
	}
}
func (ws *Websocket) handleIncomingMessage(currentConn *websocket.Conn, eventType events.GameReceiveEvent, data []byte) {
	switch eventType {
	case events.StartGame:
		{
			var eventPayload payloads.StartGameEvent
			if err := json.Unmarshal(data, &eventPayload); err != nil {
				fmt.Println("error parsing StartGame event")
				return
			}
			claims, err := utils.ParseToken(eventPayload.Token)
			if err != nil {
				fmt.Println("Error parsing token")
				return
			}

			var existingPlayer *entities.Player

			players, err := ws.playerLiveDataRepository.List()

			if err != nil {
				log.Println(err.Error())
			}

			for _, player := range players {
				if player.ID == claims["id"] {

					existingPlayer = player

					if existingPlayer.Position == nil {
						player.Position = &entities.Position{
							X: &eventPayload.Position.X,
							Y: &eventPayload.Position.Y,
						}
					}
					ws.playerLiveDataRepository.Save(player)

					break
				}
			}

			if existingPlayer == nil {
				playerNotFoundPayload := struct {
					Event events.GameEmitEvent `json:"event"`
				}{
					Event: events.PlayerNotFound,
				}

				playerNotFoundPayloadJSON, err := json.Marshal(playerNotFoundPayload)
				if err != nil {
					fmt.Println("Error converting player not found payload to JSON:", err)
					return
				}

				ws.broadcastMessage(playerNotFoundPayloadJSON)
			}

			var emitEventPayload payloads.SetInitialPlayersPositionEvent
			emitEventPayload.Event = events.SetInitialPlayersPosition

			for _, player := range players {
				emitEventPayload.Players = append(emitEventPayload.Players, payloads.PlayerWithMessages{
					ID:       player.ID,
					Username: player.Username,
					Color:    player.Color,
					Position: payloads.Position{
						X: *player.Position.X,
						Y: *player.Position.Y,
					},
					ChatMessages: redis.GetChatMessages(player.ID),
				})
			}

			emitEventPayloadJSON, err := json.Marshal(emitEventPayload)

			if err != nil {
				fmt.Println("Error converting players to json")
				return
			}

			ws.broadcastMessage(emitEventPayloadJSON)
		}
	case events.PlayerMoved:
		{
			var eventPayload payloads.PlayerMovedEvent

			if err := json.Unmarshal(data, &eventPayload); err != nil {
				fmt.Println("error parsing PlayerMoved event")
				return
			}

			claims, err := utils.ParseToken(eventPayload.Token)

			if err != nil {
				fmt.Println("Error parsing token")
				return
			}

			players, err := ws.playerLiveDataRepository.List()

			if err != nil {
				fmt.Println(err.Error())
			}
			for _, player := range players {
				if player.ID == claims["id"] {
					player.Position.X = &eventPayload.Position.X
					player.Position.Y = &eventPayload.Position.Y
					ws.playerLiveDataRepository.Save(player)

					var emitEventPayload payloads.UpdateOtherPlayerPositionEvent

					emitEventPayload.ID = player.ID
					emitEventPayload.Position = eventPayload.Position
					emitEventPayload.Event = string(events.UpdatePlayerPosition)
					emitEventPayload.CurrentState = eventPayload.CurrentState
					emitEventPayload.IsFlipped = eventPayload.IsFlipped
					emitPayLoadJSON, err := json.Marshal(emitEventPayload)

					if err != nil {
						fmt.Println("Error converting players to json")
						return
					}

					ws.broadcastMessageExcept(emitPayLoadJSON, currentConn)
					break
				}
			}
		}
	case events.MessageSent:
		{
			var eventPayload payloads.MessageSentEvent
			if err := json.Unmarshal(data, &eventPayload); err != nil {
				fmt.Println("error parsing PlayerMoved event")
				return
			}

			claims, err := utils.ParseToken(eventPayload.Token)
			if err != nil {
				fmt.Println("Error parsing token")
				return
			}
			emitEventPayload := payloads.MessageReceivedEvent{
				Event:    string(events.MessageReceived),
				SenderID: claims["id"].(string),
				Message:  eventPayload.Message,
			}
			redis.SaveChatMessage(claims["id"].(string), eventPayload.Message)

			emitPayLoadJSON, err := json.Marshal(emitEventPayload)

			if err != nil {
				fmt.Println("Error converting players to json")
				return
			}

			ws.broadcastMessageExcept(emitPayLoadJSON, currentConn)

		}

	case events.AudioChuckSent:
		{
			var eventPayload payloads.AudioChunkSentEvent
			if err := json.Unmarshal(data, &eventPayload); err != nil {
				fmt.Println("error parsing PlayerMoved event")
				return
			}

			claims, err := utils.ParseToken(eventPayload.Token)
			if err != nil {
				fmt.Println("Error parsing token")
				return
			}
			emitEventPayload := payloads.AudioChuckReceivedEvent{
				Event:    string(events.AudioChunkReceived),
				SenderID: claims["id"].(string),
				Chunk:    eventPayload.Chunk,
			}

			emitPayLoadJSON, err := json.Marshal(emitEventPayload)

			if err != nil {
				fmt.Println("Error converting audio chunk message to json")
				return
			}

			ws.broadcastMessageExcept(emitPayLoadJSON, currentConn)
		}
	case events.WebRTCOfferSent:
		{
			var eventPayload payloads.WebRTCOfferSentEvent
			if err := json.Unmarshal(data, &eventPayload); err != nil {
				fmt.Println("error parsing PlayerMoved event")
				return
			}

			_, err := utils.ParseToken(eventPayload.Token)
			if err != nil {
				fmt.Println("Error parsing token")
				return
			}
			emitEventPayload := payloads.WebRTCOfferReceivedEvent{
				Event: string(events.WebRTCOfferReceived),
				Offer: eventPayload.Offer,
			}

			emitPayLoadJSON, err := json.Marshal(emitEventPayload)

			if err != nil {
				fmt.Println("Error converting audio chunk message to json")
				return
			}

			ws.broadcastMessageExcept(emitPayLoadJSON, currentConn)
		}
	case events.WebRTCCandidateSent:
		{
			var eventPayload payloads.WebRTCCandidateSentEvent
			if err := json.Unmarshal(data, &eventPayload); err != nil {
				fmt.Println("error parsing PlayerMoved event")
				return
			}

			_, err := utils.ParseToken(eventPayload.Token)
			if err != nil {
				fmt.Println("Error parsing token")
				return
			}
			emitEventPayload := payloads.WebrtcCandidateReceidEvent{
				Event:     string(events.WebRTCCandidateReceived),
				Candidate: eventPayload.Candidate,
			}

			emitPayLoadJSON, err := json.Marshal(emitEventPayload)

			if err != nil {
				fmt.Println("Error converting audio chunk message to json")
				return
			}

			ws.broadcastMessageExcept(emitPayLoadJSON, currentConn)
		}
	case events.WebRTCAnswerSent:
		{
			var eventPayload payloads.WebRTCAnswerSentEvent
			if err := json.Unmarshal(data, &eventPayload); err != nil {
				fmt.Println("error parsing PlayerMoved event")
				return
			}

			_, err := utils.ParseToken(eventPayload.Token)
			if err != nil {
				fmt.Println("Error parsing token")
				return
			}
			emitEventPayload := payloads.WebRTCAnswerReceivedEvent{
				Event:  string(events.WebRTCAnswerReceived),
				Answer: eventPayload.Answer,
			}

			emitPayLoadJSON, err := json.Marshal(emitEventPayload)

			if err != nil {
				fmt.Println("Error converting audio chunk message to json")
				return
			}

			ws.broadcastMessageExcept(emitPayLoadJSON, currentConn)
		}
	}

}
func (ws *Websocket) broadcastMessage(message []byte) {

	ws.broadcast <- broadcastMessage{
		message: message,
	}
}

func (ws *Websocket) broadcastMessageExcept(message []byte, excludeConn *websocket.Conn) {

	ws.broadcast <- broadcastMessage{
		message: message,
		exclude: excludeConn,
	}
}

func (ws *Websocket) hub() {
	for {
		select {
		case conn := <-ws.addConnection:
			ws.Connections[conn] = true

		case conn := <-ws.removeConnection:
			if _, ok := ws.Connections[conn]; ok {
				conn.Close()
				delete(ws.Connections, conn)
			}

		case msg := <-ws.broadcast:
			var err error
			for conn := range ws.Connections {
				if msg.exclude != nil && conn == msg.exclude {
					continue
				}
				if writeErr := conn.WriteMessage(websocket.TextMessage, msg.message); writeErr != nil {
					log.Println("Write error:", writeErr)
					ws.removeConnection <- conn
					msg.errorChan <- writeErr
				}
			}
			if msg.errorChan != nil {
				msg.errorChan <- err
			}
		}
	}
}
