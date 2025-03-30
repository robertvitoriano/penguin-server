package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/robertvitoriano/penguin-server/auth"
	"github.com/robertvitoriano/penguin-server/repositories"
)

type Websocket struct {
	connection *websocket.Conn
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

func (ws *Websocket) serveWebsocket(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)

	ws.connection = conn

	if err != nil {
		log.Println(err)
		return
	}

	fmt.Println("trying to connect...")

	for {
		_, data, err := ws.connection.ReadMessage()

		if err != nil {
			fmt.Println("There was an error")
		}

		var receiveMessage BaseMessage

		err = json.Unmarshal(data, &receiveMessage)

		if err != nil {
			fmt.Println("Error parsing json")
		}

		ws.handleIncomingMessage(GameReceiveEvent(receiveMessage.Event), data)

	}

}

func (ws *Websocket) handleIncomingMessage(eventType GameReceiveEvent, data []byte) {

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
			}

			for _, player := range repositories.Players {
				if player.ID == claims["id"] {

					player.Position.X = eventPayload.Position.X
					player.Position.Y = eventPayload.Position.Y

					playersJSON, err := json.Marshal(repositories.Players)

					if err != nil {
						fmt.Println("Error conveting players to json")
					}

					ws.connection.WriteMessage(websocket.TextMessage, playersJSON)

					break
				}
			}
		}
	case PlayerMoved:
		{
			var eventPayload PlayerMovedEvent

			if err := json.Unmarshal(data, &eventPayload); err != nil {
				fmt.Println("error parsing StartGame event")
				return
			}

			claims, err := auth.ParseToken(eventPayload.Token)

			if err != nil {
				fmt.Println("Error parsing token")
			}

			for _, player := range repositories.Players {
				if player.ID == claims["id"] {

					player.Position.X = eventPayload.Position.X
					player.Position.Y = eventPayload.Position.Y

					playersJSON, err := json.Marshal(repositories.Players)

					if err != nil {
						fmt.Println("Error conveting players to json")
					}

					ws.connection.WriteMessage(websocket.TextMessage, playersJSON)

					break
				}
			}
		}

	}

}

func (ws *Websocket) handleEmitMessage(eventType GameReceiveEvent, data []byte) {

}
