package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/robertvitoriano/penguin-server/auth"
	"github.com/robertvitoriano/penguin-server/models"
	"github.com/robertvitoriano/penguin-server/repositories"
	receiveEvents "github.com/robertvitoriano/penguin-server/reveiveEvents"
)

type Websocket struct {
	connection *websocket.Conn
}

type ReceiveMessage struct {
	Event    receiveEvents.ReceiveEvent `json:"event"`
	Token    string                     `json:"token"`
	Position models.Position            `json:"position"`
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

		var receiveMessage ReceiveMessage

		err = json.Unmarshal(data, &receiveMessage)

		if err != nil {
			fmt.Println("Error parsing json")
		}

		if receiveMessage.Event == receiveEvents.Close {
			break
		}

		ws.handleIncomingMessage(receiveMessage)

	}

	ws.connection.Close()

}

func (ws *Websocket) handleIncomingMessage(receiveMessage ReceiveMessage) {

	switch receiveMessage.Event {
	case receiveEvents.StartGame:
		{
			claims, err := auth.ParseToken(receiveMessage.Token)

			for _, player := range repositories.Players {
				if player.ID == claims["id"] {

					player.Position.X = receiveMessage.Position.X
					player.Position.Y = receiveMessage.Position.Y

					playersJSON, err := json.Marshal(repositories.Players)

					if err != nil {
						fmt.Println("Error conveting players to json")
					}

					ws.connection.WriteMessage(websocket.BinaryMessage, playersJSON)

					break
				}
			}

			if err != nil {
				fmt.Println("Error parsing token")
			}
		}

	}

}
