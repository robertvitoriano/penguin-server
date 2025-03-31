package controllers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
	"github.com/robertvitoriano/penguin-server/auth"
	"github.com/robertvitoriano/penguin-server/repositories"
)

type Websocket struct {
	Connections map[*websocket.Conn]bool
	mu          sync.Mutex
}

type BaseMessage struct {
	Event string `json:"event"`
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func (ws *Websocket) ServeWebsocket(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}

	// Add the new connection
	ws.AddConnection(conn)

	// Remove when done
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

		ws.handleIncomingMessage(conn, GameReceiveEvent(message.Event), data)
	}
}
func (ws *Websocket) AddConnection(conn *websocket.Conn) {
	ws.mu.Lock()
	defer ws.mu.Unlock()
	ws.Connections[conn] = true
}

func (ws *Websocket) RemoveConnection(conn *websocket.Conn) {
	ws.mu.Lock()
	defer ws.mu.Unlock()
	delete(ws.Connections, conn)
}

func (ws *Websocket) Broadcast(message []byte) {
	ws.mu.Lock()
	defer ws.mu.Unlock()

	for conn := range ws.Connections {
		if err := conn.WriteMessage(websocket.TextMessage, message); err != nil {
			log.Println("Write error:", err)
			conn.Close()
			delete(ws.Connections, conn)
		}
	}
}
func (ws *Websocket) handleIncomingMessage(currentConn *websocket.Conn, eventType GameReceiveEvent, data []byte) {
	switch eventType {
	case StartGame:
		{
			var eventPayload StartGameEvent
			if err := json.Unmarshal(data, &eventPayload); err != nil {
				fmt.Println("error parsing StartGame event")
				return
			}

			claims, err := auth.ParseToken(eventPayload.Token)
			if err != nil {
				fmt.Println("Error parsing token")
				return
			}

			for _, player := range repositories.Players {
				if player.ID == claims["id"] {
					player.Position.X = eventPayload.Position.X
					player.Position.Y = eventPayload.Position.Y

					var emitEventPayload SetInitialPlayersPositionEvent
					emitEventPayload.Event = "set_initial_players_position"
					emitEventPayload.Players = repositories.Players

					emitEventPayloadJSON, err := json.Marshal(emitEventPayload)
					if err != nil {
						fmt.Println("Error converting players to json")
						return
					}

					ws.broadcastMessage(emitEventPayloadJSON)
					break
				}
			}
		}
	case PlayerMoved:
		{
			var eventPayload PlayerMovedEvent
			if err := json.Unmarshal(data, &eventPayload); err != nil {
				fmt.Println("error parsing PlayerMoved event")
				return
			}

			claims, err := auth.ParseToken(eventPayload.Token)
			if err != nil {
				fmt.Println("Error parsing token")
				return
			}

			for _, player := range repositories.Players {
				if player.ID == claims["id"] {
					player.Position.X = eventPayload.Position.X
					player.Position.Y = eventPayload.Position.Y

					var emitEventPayload UpdateOtherPlayerPositionEvent

					emitEventPayload.ID = player.ID
					emitEventPayload.Position = eventPayload.Position
					emitEventPayload.Event = "update_player_position"
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
	}
}
func (ws *Websocket) broadcastMessage(message []byte) {
	ws.mu.Lock()
	defer ws.mu.Unlock()

	for conn := range ws.Connections {
		if err := conn.WriteMessage(websocket.TextMessage, message); err != nil {
			fmt.Println("Error broadcasting message:", err)
			conn.Close()
			delete(ws.Connections, conn)
		}
	}
}

func (ws *Websocket) broadcastMessageExcept(message []byte, excludeConn *websocket.Conn) {
	ws.mu.Lock()
	defer ws.mu.Unlock()

	for conn := range ws.Connections {
		if conn != excludeConn {
			if err := conn.WriteMessage(websocket.TextMessage, message); err != nil {
				fmt.Println("Error broadcasting message:", err)
				conn.Close()
				delete(ws.Connections, conn)
			}
		}
	}
}
